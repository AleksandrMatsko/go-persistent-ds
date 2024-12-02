package main

import "go-persistent-ds/internal"

type Map[TKey comparable, TVal any] struct {
	versionTree *internal.VersionTree[mapVersionInfo]
	nodes       map[TKey]internal.FatNode
}

type mapVersionInfo struct{}

func NewMap[TKey comparable, TVal any]() *Map[TKey, TVal] {
	return NewMapWithCapacity[TKey, TVal](0)
}

func NewMapWithCapacity[TKey comparable, TVal any](capacity int) *Map[TKey, TVal] {
	return &Map[TKey, TVal]{
		versionTree: internal.NewVersionTree[mapVersionInfo](),
		nodes:       make(map[TKey]internal.FatNode, capacity),
	}
}

func (m *Map[TKey, TVal]) Get(version uint64, key TKey) (TVal, bool) {
	fatNode, exists := m.nodes[key]
	if !exists {
		return *new(TVal), false
	}

	// on version = 0 Map is empty
	if version == 0 {
		return *new(TVal), false
	}

	val, _, found := fatNode.FindByVersion(version)
	if found {
		// found value exactly for the version
		return val.(TVal), true
	}

	changeHistory, err := m.versionTree.GetHistory(version)
	if err != nil {
		return *new(TVal), false
	}

	// zero version is always inside change history, and for the last version
	if len(changeHistory) == 1 || len(changeHistory) == 2 {
		return *new(TVal), false
	}

	for i := uint64(len(changeHistory) - 2); i >= 1; i-- {
		val, _, found = fatNode.FindByVersion(version)
		if found {
			return val.(TVal), true
		}
	}

	return *new(TVal), false
}
