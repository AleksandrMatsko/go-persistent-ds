package internal

import "slices"

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
func (vt *VersionTree[T]) Update(prevVersion uint64) uint64 {
	node := vt.findVersion(prevVersion)
	newNode := newVersionTreeNode(vt.versionMachine.GetAndIncrementVersion(), node)
	node.children = append(node.children, newNode)
	vt.tree = append(vt.tree, newNode)

	return vt.versionMachine.GetVersion()
}

// GetHistory returns change history for specified object's version.
func (vt *VersionTree[T]) GetHistory(version uint64) []uint64 {
	node := vt.findVersion(version)

	var history []uint64
	for node != nil {
		history = append(history, node.version)
		node = node.parent
	}
	slices.Reverse(history)

	return history
}

// GetVersionInfo returns info for specified version.
func (vt *VersionTree[T]) GetVersionInfo(version uint64) T {
	node := vt.findVersion(version)
	return node.versionInfo
}

func (vt VersionTree[T]) SetVersionInfo(version uint64, info T) {
	vt.findVersion(version).versionInfo = info
}

func newVersionTreeNode[T any](v uint64, parent *versionTreeNode[T]) *versionTreeNode[T] {
	return &versionTreeNode[T]{
		version:  v,
		parent:   parent,
		children: []*versionTreeNode[T]{},
	}
}

func (vt *VersionTree[T]) findVersion(version uint64) *versionTreeNode[T] {
	return vt.tree[version]
}
