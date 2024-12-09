package main

import (
	"errors"

	"go-persistent-ds/internal"
)

var (
	// ErrSliceInitialize is returned then there is a problem in creating new slice because of tree problems.
	ErrSliceInitialize = errors.New("failed to init Slice because version tree is damaged")
	// ErrIndexOutOfRange is returned then index is greater than len array for the version.
	ErrIndexOutOfRange = errors.New("array index out of range")
)

// Slice is a persistent implementation of go slice.
// While working with slice you can access and/or modify each previous version.
// Note that modifying version creates new one.
//
// Slice can perform total of 2^65-1 modifications, and will panic on attempt to modify it for 2^65 time.
// If you need to continue editing Slice, the good idea is to use ToGoSlice method to dump Slice for special version.
//
// Note that Slice is not thread safe.
type Slice[TVal any] struct {
	versionTree     *internal.VersionTree[sliceVersionInfo]
	sliceOfFatNodes []*internal.FatNode
}

type sliceVersionInfo struct {
	size       int
	startIndex int
}

// NewSlice creates empty Slice.
func NewSlice[TVal any]() (*Slice[TVal], uint64) {
	return NewSliceWithCapacity[TVal](0)
}

// NewSliceWithCapacity creates empty Slice with given capacity.
func NewSliceWithCapacity[TVal any](capacity int) (*Slice[TVal], uint64) {
	s := &Slice[TVal]{
		versionTree:     internal.NewVersionTree[sliceVersionInfo](),
		sliceOfFatNodes: make([]*internal.FatNode, 0, capacity),
	}

	var (
		initialVersion    uint64 = 0
		initialSliceSize         = 0
		initialStartIndex        = 0
	)

	err := s.versionTree.SetVersionInfo(
		initialVersion,
		sliceVersionInfo{
			size:       initialSliceSize,
			startIndex: initialStartIndex,
		})
	if err != nil {
		panic(ErrSliceInitialize)
	}

	return s, 0
}

// NewSliceWithAnyValues creates a Slice, that can store values of any type.
func NewSliceWithAnyValues() (*Slice[any], uint64) {
	return NewSliceWithCapacity[any](0)
}

// Set value for given index and version in Slice.
//
// Complexity: O(1).
func (s *Slice[TVal]) Set(forVersion uint64, index int, val TVal) (uint64, error) {
	if index < 0 {
		return 0, ErrIndexOutOfRange
	}

	oldVersionInfo, err := s.versionTree.GetVersionInfo(forVersion)
	if err != nil {
		return 0, err
	}

	actualIndex := oldVersionInfo.startIndex + index

	if len(s.sliceOfFatNodes) <= actualIndex {
		return 0, ErrIndexOutOfRange
	}

	if forVersion == 0 {
		return 0, ErrIndexOutOfRange
	}

	fatNode := s.sliceOfFatNodes[actualIndex]

	newVersion, err := s.versionTree.Update(forVersion)
	if err != nil {
		return 0, err
	}

	newVersionInfo := sliceVersionInfo{
		size:       oldVersionInfo.size,
		startIndex: oldVersionInfo.startIndex,
	}

	fatNode.Update(val, newVersion)
	_ = s.versionTree.SetVersionInfo(newVersion, newVersionInfo)

	return newVersion, nil
}

// Get returns a pair of value and error for provided version and index.
// If error is nil then value for such index and version exists.
//
// Complexity: O(log(m) * k) there:
//   - m - amount of modifications for value by the index from slice creation.
//   - k - amount of modifications visible from current branch.
func (s *Slice[TVal]) Get(version uint64, index int) (TVal, error) {
	if index < 0 {
		return *new(TVal), ErrIndexOutOfRange
	}

	info, err := s.versionTree.GetVersionInfo(version)
	if err != nil {
		return *new(TVal), err
	}

	actualIndex := info.startIndex + index

	if len(s.sliceOfFatNodes) <= actualIndex {
		return *new(TVal), ErrIndexOutOfRange
	}

	if version == 0 {
		return *new(TVal), ErrIndexOutOfRange
	}

	fatNode := s.sliceOfFatNodes[actualIndex]

	val, _, found := fatNode.FindByVersion(version)
	if found {
		return val.(TVal), nil
	}

	changeHistory, err := s.versionTree.GetHistory(version)
	if err != nil {
		return *new(TVal), err
	}

	if len(changeHistory) == 1 || len(changeHistory) == 2 {
		return *new(TVal), ErrIndexOutOfRange
	}

	for i := len(changeHistory) - 2; i >= 1; i-- {
		val, _, found = fatNode.FindByVersion(changeHistory[i])
		if found {
			return val.(TVal), nil
		}
	}

	return *new(TVal), ErrIndexOutOfRange
}

// Len returns the len of Slice.
//
// Complexity: O(1).
func (s *Slice[TVal]) Len(version uint64) (int, error) {
	info, err := s.versionTree.GetVersionInfo(version)
	if err != nil {
		return 0, err
	}

	return info.size, nil
}

// Append adds the value to the end of Slice of given version.
//
// Complexity: O(1).
func (s *Slice[TVal]) Append(version uint64, val TVal) (uint64, error) {
	oldVersionInfo, err := s.versionTree.GetVersionInfo(version)
	if err != nil {
		return 0, err
	}

	newVersion, err := s.versionTree.Update(version)
	if err != nil {
		return 0, err
	}

	newVersionInfo := sliceVersionInfo{
		size:       oldVersionInfo.size + 1,
		startIndex: oldVersionInfo.startIndex,
	}

	actualIndex := oldVersionInfo.startIndex + oldVersionInfo.size

	if actualIndex >= len(s.sliceOfFatNodes) {
		newFatNode := internal.NewFatNode(val, newVersion)
		s.sliceOfFatNodes = append(s.sliceOfFatNodes, newFatNode)
	} else {
		existedFatNode := s.sliceOfFatNodes[actualIndex]
		existedFatNode.Update(val, newVersion)
	}

	_ = s.versionTree.SetVersionInfo(newVersion, newVersionInfo)

	return newVersion, nil
}

// ToGoSlice converts persistent Slice for specified version into go slice.
//
// Complexity: O(Get) * n, there:
//   - n - size of Slice for version.
func (s *Slice[TVal]) ToGoSlice(forVersion uint64) ([]TVal, error) {
	size, err := s.Len(forVersion)
	if err != nil {
		return nil, err
	}

	resSlice := make([]TVal, 0, size)

	for i := 0; i < size; i++ {
		val, err := s.Get(forVersion, i)
		if err != nil {
			return nil, err
		}

		resSlice = append(resSlice, val)
	}

	return resSlice, nil
}

// Range takes the range of Slice for given version from startIndex (inclusive) to
// endIndex (not inclusive).
//
// Complexity: O(1).
func (s *Slice[TVal]) Range(forVersion uint64, startIndex, endIndex int) (uint64, error) {
	if startIndex < 0 || endIndex < 0 {
		return 0, ErrIndexOutOfRange
	}

	if startIndex > endIndex {
		return 0, ErrIndexOutOfRange
	}

	oldVersionInfo, err := s.versionTree.GetVersionInfo(forVersion)
	if err != nil {
		return 0, err
	}

	if startIndex >= oldVersionInfo.size {
		return 0, ErrIndexOutOfRange
	}

	if endIndex > oldVersionInfo.size {
		return 0, ErrIndexOutOfRange
	}

	newVersion, _ := s.versionTree.Update(forVersion)
	newVersionInfo := sliceVersionInfo{
		size:       endIndex - startIndex,
		startIndex: oldVersionInfo.startIndex + startIndex,
	}

	_ = s.versionTree.SetVersionInfo(newVersion, newVersionInfo)

	return newVersion, nil
}
