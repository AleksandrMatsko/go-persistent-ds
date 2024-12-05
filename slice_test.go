package main

import (
	"errors"
	"go-persistent-ds/internal"
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

func getBranchedSlice(t *testing.T) *Slice[string] {
	s, initialVersion := NewSlice[string]()
	versionShouldBe(t, initialVersion, 0)

	v, err := s.Append(0, "a")
	errIsNil(t, err)
	versionShouldBe(t, v, 1)

	v, err = s.Append(1, "b")
	errIsNil(t, err)
	versionShouldBe(t, v, 2)

	v, err = s.Append(1, "c")
	errIsNil(t, err)
	versionShouldBe(t, v, 3)

	v, err = s.Append(2, "c")
	errIsNil(t, err)
	versionShouldBe(t, v, 4)

	v, err = s.Append(3, "b")
	errIsNil(t, err)
	versionShouldBe(t, v, 5)

	return s
}

func TestSlice_AppendAndGet(t *testing.T) {
	t.Run("Append and Get ok", func(t *testing.T) {
		t.Parallel()

		s, initialVersion := NewSlice[string]()
		versionShouldBe(t, initialVersion, 0)

		v, err := s.Append(0, "a")
		errIsNil(t, err)
		versionShouldBe(t, v, 1)

		val, err := s.Get(0, 0)
		errShouldBe(t, err, ErrIndexOutOfRange)
		isTrue(t, val == "")

		val, err = s.Get(1, 0)
		errIsNil(t, err)
		isTrue(t, val == "a")
	})

	t.Run("Append to non existing version", func(t *testing.T) {
		t.Parallel()

		s, initialVersion := NewSlice[string]()
		versionShouldBe(t, initialVersion, 0)

		v, err := s.Append(1, "a")
		errShouldBe(t, err, internal.ErrVersionNotFound)
		versionShouldBe(t, v, 0)
	})

	t.Run("Branched Append and Get", func(t *testing.T) {
		t.Parallel()

		s := getBranchedSlice(t)

		type testCase struct {
			givenVersion uint64
			givenIndex   int
			expectedVal  string
			expectedErr  error
		}

		testCases := []testCase{
			{
				givenVersion: 0,
				givenIndex:   0,
				expectedVal:  "",
				expectedErr:  ErrIndexOutOfRange,
			},
			{
				givenVersion: 1,
				givenIndex:   0,
				expectedVal:  "a",
				expectedErr:  nil,
			},
			{
				givenVersion: 2,
				givenIndex:   0,
				expectedVal:  "a",
				expectedErr:  nil,
			},
			{
				givenVersion: 3,
				givenIndex:   0,
				expectedVal:  "a",
				expectedErr:  nil,
			},
			{
				givenVersion: 4,
				givenIndex:   0,
				expectedVal:  "a",
				expectedErr:  nil,
			},
			{
				givenVersion: 5,
				givenIndex:   0,
				expectedVal:  "a",
				expectedErr:  nil,
			},
			{
				givenVersion: 6,
				givenIndex:   0,
				expectedVal:  "",
				expectedErr:  internal.ErrVersionNotFound,
			},
			{
				givenVersion: 0,
				givenIndex:   1,
				expectedVal:  "",
				expectedErr:  ErrIndexOutOfRange,
			},
			{
				givenVersion: 1,
				givenIndex:   1,
				expectedVal:  "",
				expectedErr:  ErrIndexOutOfRange,
			},
			{
				givenVersion: 2,
				givenIndex:   1,
				expectedVal:  "b",
				expectedErr:  nil,
			},
			{
				givenVersion: 3,
				givenIndex:   1,
				expectedVal:  "c",
				expectedErr:  nil,
			},
			{
				givenVersion: 4,
				givenIndex:   1,
				expectedVal:  "b",
				expectedErr:  nil,
			},
			{
				givenVersion: 5,
				givenIndex:   1,
				expectedVal:  "c",
				expectedErr:  nil,
			},
			{
				givenVersion: 6,
				givenIndex:   1,
				expectedVal:  "",
				expectedErr:  internal.ErrVersionNotFound,
			},
			{
				givenVersion: 0,
				givenIndex:   2,
				expectedVal:  "",
				expectedErr:  ErrIndexOutOfRange,
			},
			{
				givenVersion: 1,
				givenIndex:   2,
				expectedVal:  "",
				expectedErr:  ErrIndexOutOfRange,
			},
			{
				givenVersion: 2,
				givenIndex:   2,
				expectedVal:  "",
				expectedErr:  ErrIndexOutOfRange,
			},
			{
				givenVersion: 3,
				givenIndex:   2,
				expectedVal:  "",
				expectedErr:  ErrIndexOutOfRange,
			},
			{
				givenVersion: 4,
				givenIndex:   2,
				expectedVal:  "c",
				expectedErr:  nil,
			},
			{
				givenVersion: 5,
				givenIndex:   2,
				expectedVal:  "b",
				expectedErr:  nil,
			},
			{
				givenVersion: 6,
				givenIndex:   2,
				expectedVal:  "",
				expectedErr:  internal.ErrVersionNotFound,
			},
		}

		for i, c := range testCases {
			t.Logf("case %v: s.Get(%v, %v) -> %v, %v",
				i+1,
				c.givenVersion,
				c.givenIndex,
				c.expectedVal,
				c.expectedErr)
			val, err := s.Get(c.givenVersion, c.givenIndex)
			errShouldBe(t, err, c.expectedErr)
			isTrue(t, val == c.expectedVal)
		}
	})
}

func TestSlice_Set(t *testing.T) {
	t.Run("Set value to not existing index", func(t *testing.T) {
		t.Parallel()

		s, initialVersion := NewSlice[string]()
		versionShouldBe(t, initialVersion, 0)

		v, err := s.Append(0, "a")
		errIsNil(t, err)
		versionShouldBe(t, v, 1)

		v, err = s.Set(1, 1, "b")
		errShouldBe(t, err, ErrIndexOutOfRange)
		versionShouldBe(t, v, 0)
	})

	t.Run("Set value to not existing version", func(t *testing.T) {
		t.Parallel()

		s, initialVersion := NewSlice[string]()
		versionShouldBe(t, initialVersion, 0)

		v, err := s.Append(0, "a")
		errIsNil(t, err)
		versionShouldBe(t, v, 1)

		v, err = s.Set(2, 0, "b")
		errShouldBe(t, err, internal.ErrVersionNotFound)
		versionShouldBe(t, v, 0)
	})

	t.Run("Set value to not existing index and not existing version", func(t *testing.T) {
		t.Parallel()

		s, initialVersion := NewSlice[string]()
		versionShouldBe(t, initialVersion, 0)

		v, err := s.Append(0, "a")
		errIsNil(t, err)
		versionShouldBe(t, v, 1)

		v, err = s.Set(2, 1, "b")
		errShouldBe(t, err, ErrIndexOutOfRange)
		versionShouldBe(t, v, 0)
	})

	t.Run("Set value to ok version and index", func(t *testing.T) {
		t.Parallel()

		s, initialVersion := NewSlice[string]()
		versionShouldBe(t, initialVersion, 0)

		v, err := s.Append(0, "a")
		errIsNil(t, err)
		versionShouldBe(t, v, 1)

		v, err = s.Set(1, 0, "b")
		errIsNil(t, err)
		versionShouldBe(t, v, 2)

		val, err := s.Get(1, 0)
		errIsNil(t, err)
		isTrue(t, val == "a")

		val, err = s.Get(2, 0)
		errIsNil(t, err)
		isTrue(t, val == "b")
	})
}
