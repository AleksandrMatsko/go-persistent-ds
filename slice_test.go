package go_persistent_ds

import (
	"slices"
	"testing"

	"github.com/AleksandrMatsko/go-persistent-ds/internal"
)

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
		errShouldBe(t, err, internal.ErrVersionNotFound)
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

func TestSlice_Len(t *testing.T) {
	t.Run("With empty slice", func(t *testing.T) {
		t.Parallel()

		s, initialVersion := NewSlice[string]()
		versionShouldBe(t, initialVersion, 0)

		size, err := s.Len(0)
		errIsNil(t, err)
		isTrue(t, size == 0)
	})

	t.Run("With not existing version", func(t *testing.T) {
		t.Parallel()

		s, initialVersion := NewSlice[string]()
		versionShouldBe(t, initialVersion, 0)

		v, err := s.Append(0, "a")
		errIsNil(t, err)
		versionShouldBe(t, v, 1)

		size, err := s.Len(2)
		errShouldBe(t, err, internal.ErrVersionNotFound)
		isTrue(t, size == 0)
	})

	t.Run("With Append and Set", func(t *testing.T) {
		t.Parallel()

		s, initialVersion := NewSlice[string]()
		versionShouldBe(t, initialVersion, 0)

		v, err := s.Append(0, "a")
		errIsNil(t, err)
		versionShouldBe(t, v, 1)

		v, err = s.Set(1, 0, "b")
		errIsNil(t, err)
		versionShouldBe(t, v, 2)

		size, err := s.Len(1)
		errIsNil(t, err)
		isTrue(t, size == 1)

		size, err = s.Len(2)
		errIsNil(t, err)
		isTrue(t, size == 1)
	})

	t.Run("With branched slice", func(t *testing.T) {
		t.Parallel()

		s := getBranchedSlice(t)

		type testCase struct {
			givenVersion uint64
			expectedLen  int
			expectedErr  error
		}

		testCases := []testCase{
			{
				givenVersion: 0,
				expectedLen:  0,
				expectedErr:  nil,
			},
			{
				givenVersion: 1,
				expectedLen:  1,
				expectedErr:  nil,
			},
			{
				givenVersion: 2,
				expectedLen:  2,
				expectedErr:  nil,
			},
			{
				givenVersion: 3,
				expectedLen:  2,
				expectedErr:  nil,
			},
			{
				givenVersion: 4,
				expectedLen:  3,
				expectedErr:  nil,
			},
			{
				givenVersion: 5,
				expectedLen:  3,
				expectedErr:  nil,
			},
		}

		for i, c := range testCases {
			t.Logf("case %v: s.Len(%v) -> %v, %v",
				i+1,
				c.givenVersion,
				c.expectedLen,
				c.expectedErr)
			gotLen, gotErr := s.Len(c.givenVersion)
			errShouldBe(t, gotErr, c.expectedErr)
			isTrue(t, gotLen == c.expectedLen)
		}
	})
}

func TestSlice_ToGoSlice(t *testing.T) {
	t.Run("With empty slice", func(t *testing.T) {
		t.Parallel()

		s, initialVersion := NewSlice[string]()
		versionShouldBe(t, initialVersion, 0)

		slice, err := s.ToGoSlice(0)
		errIsNil(t, err)
		isTrue(t, slices.Equal(slice, []string{}))

		slice, err = s.ToGoSlice(1)
		errShouldBe(t, err, internal.ErrVersionNotFound)
		isTrue(t, slice == nil)
	})

	t.Run("With branched slice", func(t *testing.T) {
		t.Parallel()

		s := getBranchedSlice(t)

		type testCase struct {
			givenVersion  uint64
			expectedSlice []string
			expectedErr   error
		}

		testCases := []testCase{
			{
				givenVersion:  0,
				expectedSlice: []string{},
				expectedErr:   nil,
			},
			{
				givenVersion:  1,
				expectedSlice: []string{"a"},
				expectedErr:   nil,
			},
			{
				givenVersion:  2,
				expectedSlice: []string{"a", "b"},
				expectedErr:   nil,
			},
			{
				givenVersion:  3,
				expectedSlice: []string{"a", "c"},
				expectedErr:   nil,
			},
			{
				givenVersion:  4,
				expectedSlice: []string{"a", "b", "c"},
				expectedErr:   nil,
			},
			{
				givenVersion:  5,
				expectedSlice: []string{"a", "c", "b"},
				expectedErr:   nil,
			},
		}

		for i, c := range testCases {
			t.Logf("case %v: s.Len(%v) -> %s, %v",
				i+1,
				c.givenVersion,
				c.expectedSlice,
				c.expectedErr)
			gotSlice, gotErr := s.ToGoSlice(c.givenVersion)
			errShouldBe(t, gotErr, c.expectedErr)
			isTrue(t, slices.Equal(gotSlice, c.expectedSlice))
		}
	})
}

func TestSlice_Range(t *testing.T) {
	t.Run("With bad indexes", func(t *testing.T) {
		t.Parallel()

		t.Run("negative start index", func(t *testing.T) {
			t.Parallel()

			s, initialVersion := NewSlice[string]()
			versionShouldBe(t, initialVersion, 0)

			v, err := s.Range(0, -1, 0)
			errShouldBe(t, err, ErrIndexOutOfRange)
			versionShouldBe(t, v, 0)
		})

		t.Run("negative end index", func(t *testing.T) {
			t.Parallel()

			s, initialVersion := NewSlice[string]()
			versionShouldBe(t, initialVersion, 0)

			v, err := s.Range(0, 0, -1)
			errShouldBe(t, err, ErrIndexOutOfRange)
			versionShouldBe(t, v, 0)
		})

		t.Run("start index greater than end index", func(t *testing.T) {
			t.Parallel()

			s, initialVersion := NewSlice[string]()
			versionShouldBe(t, initialVersion, 0)

			v, err := s.Range(0, 2, 1)
			errShouldBe(t, err, ErrIndexOutOfRange)
			versionShouldBe(t, v, 0)
		})
	})

	t.Run("With empty slice", func(t *testing.T) {
		t.Parallel()

		s, initialVersion := NewSlice[string]()
		versionShouldBe(t, initialVersion, 0)

		v, err := s.Range(0, 0, 0)
		errShouldBe(t, err, ErrIndexOutOfRange)
		versionShouldBe(t, v, 0)
	})

	t.Run("With branched slice", func(t *testing.T) {
		t.Parallel()

		s := getBranchedSlice(t)

		v, err := s.Range(4, 1, 2)
		errIsNil(t, err)
		versionShouldBe(t, v, 6)

		slice, err := s.ToGoSlice(6)
		errIsNil(t, err)
		isTrue(t, slices.Equal(slice, []string{"b"}))

		size, err := s.Len(6)
		errIsNil(t, err)
		isTrue(t, size == 1)

		v, err = s.Append(6, "d")
		errIsNil(t, err)
		versionShouldBe(t, v, 7)

		slice, err = s.ToGoSlice(7)
		errIsNil(t, err)
		isTrue(t, slices.Equal(slice, []string{"b", "d"}))

		size, err = s.Len(7)
		errIsNil(t, err)
		isTrue(t, size == 2)

		v, err = s.Append(7, "e")
		errIsNil(t, err)
		versionShouldBe(t, v, 8)

		slice, err = s.ToGoSlice(8)
		errIsNil(t, err)
		isTrue(t, slices.Equal(slice, []string{"b", "d", "e"}))

		size, err = s.Len(8)
		errIsNil(t, err)
		isTrue(t, size == 3)

		v, err = s.Set(8, 2, "f")
		errIsNil(t, err)
		versionShouldBe(t, v, 9)

		slice, err = s.ToGoSlice(9)
		errIsNil(t, err)
		isTrue(t, slices.Equal(slice, []string{"b", "d", "f"}))

		size, err = s.Len(9)
		errIsNil(t, err)
		isTrue(t, size == 3)
	})
}

func TestSliceWithAnyTypes(t *testing.T) {
	t.Run("Append and Get values ok", func(t *testing.T) {
		t.Parallel()

		s, v := NewSliceWithAnyValues()
		versionShouldBe(t, v, 0)

		v, err := s.Append(0, "a")
		errIsNil(t, err)
		versionShouldBe(t, v, 1)

		v, err = s.Append(1, 1)
		errIsNil(t, err)
		versionShouldBe(t, v, 2)

		v, err = s.Append(2, map[string]string{})
		errIsNil(t, err)
		versionShouldBe(t, v, 3)

		v, err = s.Append(3, []uint64{})
		errIsNil(t, err)
		versionShouldBe(t, v, 4)

		val, err := s.Get(4, 0)
		errIsNil(t, err)
		isTrue(t, val == "a")

		val, err = s.Get(4, 1)
		errIsNil(t, err)
		isTrue(t, val == 1)

		val, err = s.Get(4, 2)
		errIsNil(t, err)
		isTrue(t, func() bool {
			castedVal, ok := val.(map[string]string)
			if !ok {
				return false
			}

			return len(castedVal) == 0
		}())

		val, err = s.Get(4, 3)
		errIsNil(t, err)
		isTrue(t, func() bool {
			castedVal, ok := val.([]uint64)
			if !ok {
				return false
			}

			return len(castedVal) == 0
		}())
	})
}
