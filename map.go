package main

import (
	"errors"
	"go-persistent-ds/internal"
)

var (
	// ErrMapInitialize is returned then there is a problem in creating new tree.
	ErrMapInitialize = errors.New("failed to init Map because version tree is damaged")
)

// Map is a persistent implementation of go map.
type Map[TKey comparable, TVal any] struct {
	versionTree *internal.VersionTree[mapVersionInfo]
	nodes       map[TKey]*internal.FatNode
}

type mapVersionInfo struct {
	size int
}

type mapOperationKind string

const (
	mapSet    mapOperationKind = "MapSet"
	mapDelete mapOperationKind = "MapDelete"
)

// NewMap creates empty Map.
func NewMap[TKey comparable, TVal any]() *Map[TKey, TVal] {
	return NewMapWithCapacity[TKey, TVal](0)
}

// NewMapWithCapacity creates empty Map with given capacity.
func NewMapWithCapacity[TKey comparable, TVal any](capacity int) *Map[TKey, TVal] {
	m := &Map[TKey, TVal]{
		versionTree: internal.NewVersionTree[mapVersionInfo](),
		nodes:       make(map[TKey]*internal.FatNode, capacity),
	}

	var (
		initialVersion uint64 = 0
		initialMapSize        = 0
	)

	err := m.versionTree.SetVersionInfo(
		initialVersion,
		mapVersionInfo{
			size: initialMapSize,
		})
	if err != nil {
		panic(ErrMapInitialize)
	}

	return m
}

// Get returns a pair of value and bool for provided version and key.
// Bool tells if the value for such key and version exists.
//
// Complexity: O(log(m) * log(n)) there:
//   - n - amount of different keys in map from creation;
//   - m - amount of modifications for current key from map creation.
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
		if val == nil {
			return *new(TVal), false
		}

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
			if val == nil {
				return *new(TVal), false
			}

			return val.(TVal), true
		}
	}

	return *new(TVal), false
}

// Set value for given key and version in Map.
//
// Complexity: same as for Get.
func (m *Map[TKey, TVal]) Set(forVersion uint64, key TKey, val TVal) (uint64, error) {
	newVersion, err := m.versionTree.Update(forVersion)
	if err != nil {
		return 0, err
	}

	oldVersionInfo, _ := m.versionTree.GetVersionInfo(forVersion)
	newVersionInfo := mapVersionInfo{
		size: oldVersionInfo.size,
	}

	fatNode, exists := m.nodes[key]
	if !exists {
		// adding new key
		newFatNode := internal.NewFatNode(val, newVersion)
		m.nodes[key] = newFatNode

		newVersionInfo.size += 1

		_ = m.versionTree.SetVersionInfo(
			newVersion,
			newVersionInfo)

		return newVersion, nil
	}

	fatNode.Update(val, newVersion)

	_, ok := m.Get(forVersion, key)
	if !ok {
		// the key exists in other branch of versions but not in the current branch
		newVersionInfo.size += 1
	}

	_ = m.versionTree.SetVersionInfo(
		newVersion,
		newVersionInfo)

	return newVersion, nil
}

// Len returns the len of Map.
//
// Complexity: O(1).
func (m *Map[TKey, TVal]) Len(forVersion uint64) (int, error) {
	info, err := m.versionTree.GetVersionInfo(forVersion)
	if err != nil {
		return 0, err
	}

	return info.size, nil
}

// Delete the value from Map for given key for given version.
//
// Complexity: same as for Get.
func (m *Map[TKey, TVal]) Delete(forVersion uint64, key TKey) (uint64, TVal, bool) {
	existedFatNode, keyExists := m.nodes[key]
	if !keyExists {
		// no key to delete
		return 0, *new(TVal), false
	}

	val, valExists := m.Get(forVersion, key)
	if !valExists {
		// key exists but no value visible for this version
		return 0, *new(TVal), false
	}

	newVersion, err := m.versionTree.Update(forVersion)
	if err != nil {
		return 0, *new(TVal), false
	}

	oldVersionInfo, _ := m.versionTree.GetVersionInfo(forVersion)
	newVersionInfo := mapVersionInfo{
		size: oldVersionInfo.size - 1,
	}

	existedFatNode.Update(nil, newVersion)
	m.nodes[key] = existedFatNode

	_ = m.versionTree.SetVersionInfo(newVersion, newVersionInfo)

	return newVersion, val, true
}

// ToGoMap converts persistent Map for specified version into go map.
//
// Complexity: O(Get) * n, there:
//   - n - same as in Get.
func (m *Map[TKey, TVal]) ToGoMap(version uint64) (map[TKey]TVal, error) {
	versionInfo, err := m.versionTree.GetVersionInfo(version)
	if err != nil {
		return nil, err
	}

	resMap := make(map[TKey]TVal, versionInfo.size)
	for k := range m.nodes {
		val, exists := m.Get(version, k)
		if exists {
			resMap[k] = val
		}
	}

	return resMap, nil
}