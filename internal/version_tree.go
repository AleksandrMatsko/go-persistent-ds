package internal

import (
	"errors"
	"slices"
)

// VersionTree is a struct to store object change history.
type VersionTree[T any] struct {
	tree           []*versionTreeNode[T]
	versionMachine *VersionMachine
}

type versionTreeNode[T any] struct {
	version     uint64
	versionInfo T
	parent      *versionTreeNode[T]
	children    []*versionTreeNode[T]
}

// ErrVersionNotFound will be returned if searched version was not found in VersionTree.
var ErrVersionNotFound = errors.New("version not found")

// NewVersionTree creates new object change history tree.
func NewVersionTree[T any]() *VersionTree[T] {
	vm := &VersionMachine{
		version: 0,
	}

	return &VersionTree[T]{
		tree:           []*versionTreeNode[T]{newVersionTreeNode[T](vm.GetAndIncrementVersion(), nil)},
		versionMachine: vm,
	}
}

// Update creates new version for specified version.
func (vt *VersionTree[T]) Update(prevVersion uint64) (uint64, error) {
	node, success := vt.findVersion(prevVersion)
	if !success {
		return 0, ErrVersionNotFound
	}
	newNode := newVersionTreeNode(vt.versionMachine.GetAndIncrementVersion(), node)
	node.children = append(node.children, newNode)
	vt.tree = append(vt.tree, newNode)

	return vt.versionMachine.GetVersion(), nil
}

// GetHistory returns change history for specified object's version.
func (vt *VersionTree[T]) GetHistory(version uint64) ([]uint64, error) {
	node, success := vt.findVersion(version)
	if !success {
		return nil, ErrVersionNotFound
	}

	var history []uint64
	for node != nil {
		history = append(history, node.version)
		node = node.parent
	}
	slices.Reverse(history)

	return history, nil
}

// GetVersionInfo returns info for specified version.
func (vt *VersionTree[T]) GetVersionInfo(version uint64) (*T, error) {
	node, success := vt.findVersion(version)
	if !success {
		return nil, ErrVersionNotFound
	}

	return &node.versionInfo, nil
}

func (vt *VersionTree[T]) SetVersionInfo(version uint64, info T) error {
	node, success := vt.findVersion(version)
	if !success {
		return ErrVersionNotFound
	}
	node.versionInfo = info
	return nil
}

// GetCurrentVersion returns current(latest) version in VersionTree.
func (vt *VersionTree[T]) GetCurrentVersion() uint64 {
	return vt.versionMachine.GetVersion()
}

func newVersionTreeNode[T any](v uint64, parent *versionTreeNode[T]) *versionTreeNode[T] {
	return &versionTreeNode[T]{
		version:  v,
		parent:   parent,
		children: []*versionTreeNode[T]{},
	}
}

func (vt *VersionTree[T]) findVersion(version uint64) (*versionTreeNode[T], bool) {
	if version >= uint64(len(vt.tree)) {
		return nil, false
	}
	return vt.tree[version], true
}
