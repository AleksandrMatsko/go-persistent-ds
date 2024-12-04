package main

import (
	"errors"
	"go-persistent-ds/internal"
)

var ErrIndexIsOutOfRange = errors.New("index is out of range")

type PersistentList[T any] struct {
	versionTree *internal.VersionTree[listVersionInfo]
	storage     []*internal.FatNode
}

type listVersionInfo struct {
	listSize int
}

func NewPersistentList[T any]() (*PersistentList[T], error) {
	list := &PersistentList[T]{
		versionTree: internal.NewVersionTree[listVersionInfo](),
		storage:     make([]*internal.FatNode, 0),
	}

	var (
		initialVersion  uint64 = 0
		initialListSize        = 0
	)

	err := list.versionTree.SetVersionInfo(
		initialVersion,
		listVersionInfo{
			listSize: initialListSize,
		})
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (list *PersistentList[T]) Add(value T) (uint64, error) {
	curVersion := list.versionTree.GetCurrentVersion()
	newVersion, err := list.versionTree.Update(curVersion)
	if err != nil {
		return 0, err
	}

	oldVersionInfo, _ := list.versionTree.GetVersionInfo(newVersion)
	newVersionInfo := listVersionInfo{
		listSize: oldVersionInfo.listSize + 1,
	}
	_ = list.versionTree.SetVersionInfo(newVersion, newVersionInfo)

	list.storage = append(list.storage, internal.NewFatNode(value, newVersion))

	return newVersion, nil
}

// Get FIXME: поиск по истории мб косячный.
func (list *PersistentList[T]) Get(index int, version uint64) (T, bool, error) {
	if index >= len(list.storage) {
		return *new(T), false, ErrIndexIsOutOfRange
	}

	val, _, found := list.storage[index].FindByVersion(version)
	if found {
		if val == nil {
			return *new(T), false, nil
		}
		return *val.(*T), true, nil
	}

	changeHistory, err := list.versionTree.GetHistory(version)
	if err != nil {
		return *new(T), false, err
	}

	for i := uint64(len(changeHistory)); i >= 1; i-- {
		val, _, found = list.storage[index].FindByVersion(changeHistory[i])
		if found {
			if val == nil {
				return *new(T), false, nil
			}
			return *val.(*T), true, nil
		}
	}

	return *new(T), false, nil
}
