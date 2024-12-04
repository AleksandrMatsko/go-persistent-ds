package main

import (
	"errors"
	"go-persistent-ds/internal"
	"maps"
	"testing"
)

func versionShouldBe(t *testing.T, got, expected uint64) {
	if got != expected {
		t.Errorf("expected version: %v, got: %v", expected, got)
	}
}

func isTrue(t *testing.T, flag bool) {
	if !flag {
		t.Errorf("expected true, got false")
	}
}

func errIsNil(t *testing.T, err error) {
	if err != nil {
		t.Errorf("error should be nil, but got %T: %s", err, err)
	}
}

func errShouldBe(t *testing.T, gotErr error, expectedErr error) {
	if !errors.Is(gotErr, expectedErr) {
		t.Errorf("expected err: (%T: %s), got: (%T: %s)", expectedErr, expectedErr, gotErr, gotErr)
	}
}

func getBranchedMap(t *testing.T) *Map[string, string] {
	m, initialVersion := NewMap[string, string]()
	versionShouldBe(t, initialVersion, 0)

	v, err := m.Set(0, "a", "0")
	errIsNil(t, err)
	versionShouldBe(t, v, 1)

	v, err = m.Set(1, "b", "1")
	errIsNil(t, err)
	versionShouldBe(t, v, 2)

	v, err = m.Set(1, "c", "1")
	errIsNil(t, err)
	versionShouldBe(t, v, 3)

	v, err = m.Set(2, "c", "2")
	errIsNil(t, err)
	versionShouldBe(t, v, 4)

	v, err = m.Set(3, "b", "2")
	errIsNil(t, err)
	versionShouldBe(t, v, 5)

	return m
}

func TestMap_with_GetSet(t *testing.T) {
	t.Run("Set and Get works expected", func(t *testing.T) {
		t.Parallel()

		key := "a"
		val := "a"

		m, initialVersion := NewMap[string, string]()
		versionShouldBe(t, initialVersion, 0)

		gotVal, err := m.Get(initialVersion, key)
		errShouldBe(t, err, ErrNotFound)
		isTrue(t, gotVal == "")

		gotVersion, err := m.Set(initialVersion, key, val)
		errIsNil(t, err)
		versionShouldBe(t, gotVersion, 1)

		gotVal, err = m.Get(initialVersion, key)
		errShouldBe(t, err, ErrNotFound)
		isTrue(t, gotVal == "")

		gotVal, err = m.Get(gotVersion, key)
		errIsNil(t, err)
		isTrue(t, gotVal == val)
	})

	t.Run("Set updates existed key", func(t *testing.T) {
		t.Parallel()

		m, initialVersion := NewMap[string, string]()
		versionShouldBe(t, initialVersion, 0)

		v, err := m.Set(0, "a", "1")
		errIsNil(t, err)
		versionShouldBe(t, v, 1)

		v, err = m.Set(1, "a", "2")
		errIsNil(t, err)
		versionShouldBe(t, v, 2)

		val, err := m.Get(1, "a")
		errIsNil(t, err)
		isTrue(t, val == "1")

		val, err = m.Get(2, "a")
		errIsNil(t, err)
		isTrue(t, val == "2")
	})

	t.Run("Set for not existing version", func(t *testing.T) {
		t.Parallel()

		m, initialVersion := NewMap[string, string]()
		versionShouldBe(t, initialVersion, 0)

		newVersion, err := m.Set(2, "a", "a")
		errShouldBe(t, err, internal.ErrVersionNotFound)
		versionShouldBe(t, newVersion, 0)
	})

	t.Run("Get for not existing version", func(t *testing.T) {
		t.Parallel()

		m, initialVersion := NewMap[string, string]()
		versionShouldBe(t, initialVersion, 0)

		_, err := m.Set(0, "a", "1")
		errIsNil(t, err)

		val, err := m.Get(2, "a")
		errShouldBe(t, err, ErrNotFound)
		isTrue(t, val == "")
	})

	t.Run("Set and Get with different branches", func(t *testing.T) {
		t.Parallel()

		m := getBranchedMap(t)

		type testCase struct {
			givenVersion uint64
			givenKey     string
			valShouldBe  string
			errShouldBe  error
		}

		testCases := make([]testCase, 0, 6)
		testCases = append(testCases,
			testCase{
				givenVersion: 0,
				givenKey:     "a",
				valShouldBe:  "",
				errShouldBe:  ErrNotFound,
			})

		for v := range 5 {
			testCases = append(testCases, testCase{
				givenVersion: uint64(v + 1),
				givenKey:     "a",
				valShouldBe:  "0",
				errShouldBe:  nil,
			})
		}

		otherTestCases := []testCase{
			{
				givenVersion: 0,
				givenKey:     "b",
				valShouldBe:  "",
				errShouldBe:  ErrNotFound,
			},
			{
				givenVersion: 0,
				givenKey:     "c",
				valShouldBe:  "",
				errShouldBe:  ErrNotFound,
			},
			{
				givenVersion: 1,
				givenKey:     "b",
				valShouldBe:  "",
				errShouldBe:  ErrNotFound,
			},
			{
				givenVersion: 1,
				givenKey:     "c",
				valShouldBe:  "",
				errShouldBe:  ErrNotFound,
			},
			{
				givenVersion: 2,
				givenKey:     "b",
				valShouldBe:  "1",
				errShouldBe:  nil,
			},
			{
				givenVersion: 2,
				givenKey:     "c",
				valShouldBe:  "",
				errShouldBe:  ErrNotFound,
			},
			{
				givenVersion: 3,
				givenKey:     "b",
				valShouldBe:  "",
				errShouldBe:  ErrNotFound,
			},
			{
				givenVersion: 3,
				givenKey:     "c",
				valShouldBe:  "1",
				errShouldBe:  nil,
			},
			{
				givenVersion: 4,
				givenKey:     "b",
				valShouldBe:  "1",
				errShouldBe:  nil,
			},
			{
				givenVersion: 4,
				givenKey:     "c",
				valShouldBe:  "2",
				errShouldBe:  nil,
			},
			{
				givenVersion: 5,
				givenKey:     "b",
				valShouldBe:  "2",
				errShouldBe:  nil,
			},
			{
				givenVersion: 5,
				givenKey:     "c",
				valShouldBe:  "1",
				errShouldBe:  nil,
			},
		}

		testCases = append(testCases, otherTestCases...)

		for i, singleCase := range testCases {
			t.Logf("case %v: m.Get(%v, %v) -> %v, %v",
				i+1,
				singleCase.givenVersion,
				singleCase.givenKey,
				singleCase.valShouldBe,
				singleCase.errShouldBe)
			gotVal, err := m.Get(singleCase.givenVersion, singleCase.givenKey)
			errShouldBe(t, err, singleCase.errShouldBe)
			isTrue(t, gotVal == singleCase.valShouldBe)
		}
	})
}

func TestMap_with_GetSetDelete(t *testing.T) {
	t.Run("Get after Delete", func(t *testing.T) {
		t.Parallel()

		m, initialVersion := NewMap[string, string]()
		versionShouldBe(t, initialVersion, 0)

		v, err := m.Set(0, "a", "1")
		errIsNil(t, err)
		versionShouldBe(t, v, 1)

		v, err = m.Delete(1, "a")
		errIsNil(t, err)
		versionShouldBe(t, v, 2)

		val, err := m.Get(2, "a")
		errShouldBe(t, err, ErrNotFound)
		isTrue(t, val == "")
	})

	t.Run("Delete not existing element", func(t *testing.T) {
		t.Parallel()

		m, initialVersion := NewMap[string, string]()
		versionShouldBe(t, initialVersion, 0)

		v, err := m.Set(0, "a", "1")
		errIsNil(t, err)
		versionShouldBe(t, v, 1)

		v, err = m.Delete(1, "b")
		errShouldBe(t, err, ErrNotFound)
		versionShouldBe(t, v, 0)
	})

	t.Run("Delete with not existing version", func(t *testing.T) {
		t.Parallel()

		m, initialVersion := NewMap[string, string]()
		versionShouldBe(t, initialVersion, 0)

		v, err := m.Set(0, "a", "1")
		errIsNil(t, err)
		versionShouldBe(t, v, 1)

		v, err = m.Delete(2, "a")
		errShouldBe(t, err, ErrNotFound)
		versionShouldBe(t, v, 0)
	})
}

func TestMap_Len(t *testing.T) {
	t.Run("With empty map", func(t *testing.T) {
		t.Parallel()

		m, initialVersion := NewMap[string, string]()
		versionShouldBe(t, initialVersion, 0)

		size, err := m.Len(0)
		errIsNil(t, err)
		isTrue(t, size == 0)
	})

	t.Run("With not existing version", func(t *testing.T) {
		t.Parallel()

		m, initialVersion := NewMap[string, string]()
		versionShouldBe(t, initialVersion, 0)

		v, err := m.Set(0, "a", "1")
		errIsNil(t, err)
		versionShouldBe(t, v, 1)

		size, err := m.Len(2)
		errShouldBe(t, err, internal.ErrVersionNotFound)
		isTrue(t, size == 0)
	})

	t.Run("After Set", func(t *testing.T) {
		t.Parallel()

		m, initialVersion := NewMap[string, string]()
		versionShouldBe(t, initialVersion, 0)

		v, err := m.Set(0, "a", "1")
		errIsNil(t, err)
		versionShouldBe(t, v, 1)

		size, err := m.Len(1)
		errIsNil(t, err)
		isTrue(t, size == 1)
	})

	t.Run("After Set and then update key", func(t *testing.T) {
		t.Parallel()

		m, initialVersion := NewMap[string, string]()
		versionShouldBe(t, initialVersion, 0)

		v, err := m.Set(0, "a", "1")
		errIsNil(t, err)
		versionShouldBe(t, v, 1)

		v, err = m.Set(1, "a", "2")
		errIsNil(t, err)
		versionShouldBe(t, v, 2)

		size, err := m.Len(1)
		errIsNil(t, err)
		isTrue(t, size == 1)

		size, err = m.Len(2)
		errIsNil(t, err)
		isTrue(t, size == 1)
	})

	t.Run("After Set and then Delete", func(t *testing.T) {
		t.Parallel()

		m, initialVersion := NewMap[string, string]()
		versionShouldBe(t, initialVersion, 0)

		v, err := m.Set(0, "a", "1")
		errIsNil(t, err)
		versionShouldBe(t, v, 1)

		v, err = m.Delete(1, "a")
		errIsNil(t, err)
		versionShouldBe(t, v, 2)

		size, err := m.Len(1)
		errIsNil(t, err)
		isTrue(t, size == 1)

		size, err = m.Len(2)
		errIsNil(t, err)
		isTrue(t, size == 0)
	})

	t.Run("With version branches", func(t *testing.T) {
		t.Parallel()

		m := getBranchedMap(t)

		type testCase struct {
			givenVersion uint64
			expectedLen  int
		}

		testCases := []testCase{
			{
				givenVersion: 0,
				expectedLen:  0,
			},
			{
				givenVersion: 1,
				expectedLen:  1,
			},
			{
				givenVersion: 2,
				expectedLen:  2,
			},
			{
				givenVersion: 3,
				expectedLen:  2,
			},
			{
				givenVersion: 4,
				expectedLen:  3,
			},
			{
				givenVersion: 5,
				expectedLen:  3,
			},
		}

		for i, c := range testCases {
			t.Logf("case %v: Len(%v) -> %v", i+1, c.givenVersion, c.expectedLen)
			gotLen, err := m.Len(c.givenVersion)
			errIsNil(t, err)
			isTrue(t, gotLen == c.expectedLen)
		}
	})
}

func TestMap_ToGoMap(t *testing.T) {
	t.Run("With empty map", func(t *testing.T) {
		t.Parallel()

		m, initialVersion := NewMap[string, string]()
		versionShouldBe(t, initialVersion, 0)

		gotMap, err := m.ToGoMap(0)
		errIsNil(t, err)
		isTrue(t, maps.Equal(gotMap, map[string]string{}))
	})

	t.Run("With not existing version", func(t *testing.T) {
		t.Parallel()

		m, initialVersion := NewMap[string, string]()
		versionShouldBe(t, initialVersion, 0)

		v, err := m.Set(0, "a", "1")
		errIsNil(t, err)
		versionShouldBe(t, v, 1)

		gotMap, err := m.ToGoMap(2)
		errShouldBe(t, err, internal.ErrVersionNotFound)
		isTrue(t, gotMap == nil)
	})

	t.Run("After Set", func(t *testing.T) {
		t.Parallel()

		m, initialVersion := NewMap[string, string]()
		versionShouldBe(t, initialVersion, 0)

		v, err := m.Set(0, "a", "1")
		errIsNil(t, err)
		versionShouldBe(t, v, 1)

		gotMap, err := m.ToGoMap(1)
		errIsNil(t, err)
		isTrue(t, maps.Equal(
			gotMap,
			map[string]string{
				"a": "1",
			}))
	})

	t.Run("After Set and then update key", func(t *testing.T) {
		t.Parallel()

		m, initialVersion := NewMap[string, string]()
		versionShouldBe(t, initialVersion, 0)

		v, err := m.Set(0, "a", "1")
		errIsNil(t, err)
		versionShouldBe(t, v, 1)

		v, err = m.Set(1, "a", "2")
		errIsNil(t, err)
		versionShouldBe(t, v, 2)

		gotMap, err := m.ToGoMap(1)
		errIsNil(t, err)
		isTrue(t, maps.Equal(
			gotMap,
			map[string]string{
				"a": "1",
			}))

		gotMap, err = m.ToGoMap(2)
		errIsNil(t, err)
		isTrue(t, maps.Equal(
			gotMap,
			map[string]string{
				"a": "2",
			}))
	})

	t.Run("After Set and then Delete", func(t *testing.T) {
		t.Parallel()

		m, initialVersion := NewMap[string, string]()
		versionShouldBe(t, initialVersion, 0)

		v, err := m.Set(0, "a", "1")
		errIsNil(t, err)
		versionShouldBe(t, v, 1)

		v, err = m.Delete(1, "a")
		errIsNil(t, err)
		versionShouldBe(t, v, 2)

		gotMap, err := m.ToGoMap(1)
		errIsNil(t, err)
		isTrue(t, maps.Equal(
			gotMap,
			map[string]string{
				"a": "1",
			}))

		gotMap, err = m.ToGoMap(2)
		errIsNil(t, err)
		isTrue(t, maps.Equal(gotMap, map[string]string{}))
	})

	t.Run("With version branches", func(t *testing.T) {
		t.Parallel()

		m := getBranchedMap(t)

		type testCase struct {
			givenVersion uint64
			expectedMap  map[string]string
		}

		testCases := []testCase{
			{
				givenVersion: 0,
				expectedMap:  map[string]string{},
			},
			{
				givenVersion: 1,
				expectedMap: map[string]string{
					"a": "0",
				},
			},
			{
				givenVersion: 2,
				expectedMap: map[string]string{
					"a": "0",
					"b": "1",
				},
			},
			{
				givenVersion: 3,
				expectedMap: map[string]string{
					"a": "0",
					"c": "1",
				},
			},
			{
				givenVersion: 4,
				expectedMap: map[string]string{
					"a": "0",
					"b": "1",
					"c": "2",
				},
			},
			{
				givenVersion: 5,
				expectedMap: map[string]string{
					"a": "0",
					"b": "2",
					"c": "1",
				},
			},
		}

		for i, c := range testCases {
			t.Logf("case %v: version %v", i+1, c.givenVersion)
			gotMap, err := m.ToGoMap(c.givenVersion)
			errIsNil(t, err)
			isTrue(t, maps.Equal(gotMap, c.expectedMap))
		}
	})
}
