package main

import (
	"errors"
	"go-persistent-ds/internal"
)

var (
	// ErrSliceInitialize is returned then there is a problem in creating new tree.
	ErrSliceInitialize = errors.New("failed to init Slice because version tree is damaged")
	ErrIndexOutOfRange = errors.New("array index out of range")
)

type Slice[TVal any] struct {
	versionTree     *internal.VersionTree[sliceVersionInfo]
	sliceOfFatNodes []*internal.FatNode
}

type sliceVersionInfo struct {
	size int
}

func NewSlice[TVal any]() (*Slice[TVal], uint64) {
	return NewSliceWithCapacity[TVal](0)
}

func NewSliceWithCapacity[TVal any](capacity int) (*Slice[TVal], uint64) {
	s := &Slice[TVal]{
		versionTree:     internal.NewVersionTree[sliceVersionInfo](),
		sliceOfFatNodes: make([]*internal.FatNode, 0, capacity),
	}

	var (
		initialVersion   uint64 = 0
		initialSliceSize        = 0
	)

	err := s.versionTree.SetVersionInfo(
		initialVersion,
		sliceVersionInfo{
			size: initialSliceSize,
		})
	if err != nil {
		panic(ErrSliceInitialize)
	}

	return s, 0
}

func (s *Slice[TVal]) Set(forVersion uint64, index int, val TVal) (uint64, error) {

}

func (s *Slice[TVal]) Get(version uint64, index int) (TVal, error) {
	if len(s.sliceOfFatNodes) >= index {
		return *new(TVal), ErrIndexOutOfRange
	}

	if version == 0 {
		return *new(TVal), ErrIndexOutOfRange
	}

	fatNode := s.sliceOfFatNodes[index]

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

	for i := 
}
