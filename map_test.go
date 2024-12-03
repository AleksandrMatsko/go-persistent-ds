package main

import (
	"errors"
	"go-persistent-ds/internal"
	"testing"
)

func versionOk(t *testing.T, got, expected uint64) {
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

func TestMap_GetSet(t *testing.T) {
	t.Run("Set and Get works expected", func(t *testing.T) {
		t.Parallel()

		key := "a"
		val := "a"

		m, initialVersion := NewMap[string, string]()
		versionOk(t, initialVersion, 0)

		gotVal, found := m.Get(initialVersion, key)
		isTrue(t, !found)
		isTrue(t, gotVal == "")

		gotVersion, err := m.Set(initialVersion, key, val)
		errIsNil(t, err)
		versionOk(t, gotVersion, 1)

		gotVal, found = m.Get(initialVersion, key)
		isTrue(t, !found)
		isTrue(t, gotVal == "")

		gotVal, found = m.Get(gotVersion, key)
		isTrue(t, found)
		isTrue(t, gotVal == val)
	})

	t.Run("Set for not existing version", func(t *testing.T) {
		t.Parallel()

		m, initialVersion := NewMap[string, string]()
		versionOk(t, initialVersion, 0)

		newVersion, err := m.Set(2, "a", "a")
		errShouldBe(t, err, internal.ErrVersionNotFound)
		versionOk(t, newVersion, 0)
	})

	t.Run("Get for not existing version", func(t *testing.T) {
		t.Parallel()

		m, initialVersion := NewMap[string, string]()
		versionOk(t, initialVersion, 0)

		_, err := m.Set(0, "a", "1")
		errIsNil(t, err)

		val, found := m.Get(2, "a")
		isTrue(t, !found)
		isTrue(t, val == "")
	})

	t.Run("Set and Get with different branches", func(t *testing.T) {
		t.Parallel()

		m, initialVersion := NewMap[string, string]()
		versionOk(t, initialVersion, 0)

		v, err := m.Set(0, "a", "0")
		errIsNil(t, err)
		versionOk(t, v, 1)

		v, err = m.Set(1, "b", "1")
		errIsNil(t, err)
		versionOk(t, v, 2)

		v, err = m.Set(1, "c", "1")
		errIsNil(t, err)
		versionOk(t, v, 3)

		v, err = m.Set(2, "c", "2")
		errIsNil(t, err)
		versionOk(t, v, 4)

		v, err = m.Set(3, "b", "2")
		errIsNil(t, err)
		versionOk(t, v, 5)

		type testCase struct {
			givenVersion  uint64
			givenKey      string
			valShouldBe   string
			foundShouldBe bool
		}

		testCases := make([]testCase, 0, 6)
		testCases = append(testCases,
			testCase{
				givenVersion:  0,
				givenKey:      "a",
				valShouldBe:   "",
				foundShouldBe: false,
			})

		for v = range 5 {
			testCases = append(testCases, testCase{
				givenVersion:  v + 1,
				givenKey:      "a",
				valShouldBe:   "0",
				foundShouldBe: true,
			})
		}

		otherTestCases := []testCase{
			{
				givenVersion:  0,
				givenKey:      "b",
				valShouldBe:   "",
				foundShouldBe: false,
			},
			{
				givenVersion:  0,
				givenKey:      "c",
				valShouldBe:   "",
				foundShouldBe: false,
			},
			{
				givenVersion:  1,
				givenKey:      "b",
				valShouldBe:   "",
				foundShouldBe: false,
			},
			{
				givenVersion:  1,
				givenKey:      "c",
				valShouldBe:   "",
				foundShouldBe: false,
			},
			{
				givenVersion:  2,
				givenKey:      "b",
				valShouldBe:   "1",
				foundShouldBe: true,
			},
			{
				givenVersion:  2,
				givenKey:      "c",
				valShouldBe:   "",
				foundShouldBe: false,
			},
			{
				givenVersion:  3,
				givenKey:      "b",
				valShouldBe:   "",
				foundShouldBe: false,
			},
			{
				givenVersion:  3,
				givenKey:      "c",
				valShouldBe:   "1",
				foundShouldBe: true,
			},
			{
				givenVersion:  4,
				givenKey:      "b",
				valShouldBe:   "1",
				foundShouldBe: true,
			},
			{
				givenVersion:  4,
				givenKey:      "c",
				valShouldBe:   "2",
				foundShouldBe: true,
			},
			{
				givenVersion:  5,
				givenKey:      "b",
				valShouldBe:   "2",
				foundShouldBe: true,
			},
			{
				givenVersion:  5,
				givenKey:      "c",
				valShouldBe:   "1",
				foundShouldBe: true,
			},
		}

		testCases = append(testCases, otherTestCases...)

		for i, singleCase := range testCases {
			t.Logf("case %v: m.Get(%v, %v) -> %v, %v",
				i+1,
				singleCase.givenVersion,
				singleCase.givenKey,
				singleCase.valShouldBe,
				singleCase.foundShouldBe)
			gotVal, found := m.Get(singleCase.givenVersion, singleCase.givenKey)
			isTrue(t, found == singleCase.foundShouldBe)
			isTrue(t, gotVal == singleCase.valShouldBe)
		}
	})
}
