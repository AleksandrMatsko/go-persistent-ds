package internal

import "slices"

// VersionTree is a struct to store object change history.
type VersionTree struct {
	tree           []*versionTreeNode
	versionMachine *VersionMachine
}

type versionTreeNode struct {
	version  uint64
	parent   *versionTreeNode
	children []*versionTreeNode
}

// NewVersionTree creates new object change history tree.
func NewVersionTree() *VersionTree {
	vm := &VersionMachine{
		version: 0,
	}

	return &VersionTree{
		tree:           []*versionTreeNode{newVersionTreeNode(vm.GetAndIncrementVersion(), nil)},
		versionMachine: vm,
	}
}

// Update creates new version for specified version.
func (vt *VersionTree) Update(prevVersion uint64) uint64 {
	node := vt.findVersion(prevVersion)
	newNode := newVersionTreeNode(vt.versionMachine.GetAndIncrementVersion(), node)
	node.children = append(node.children, newNode)
	vt.tree = append(vt.tree, newNode)

	return vt.versionMachine.GetVersion()
}

// GetHistory returns change history for specified object's version.
func (vt *VersionTree) GetHistory(version uint64) []uint64 {
	node := vt.findVersion(version)

	var history []uint64
	for node != nil {
		history = append(history, node.version)
		node = node.parent
	}
	slices.Reverse(history)

	return history
}

func newVersionTreeNode(v uint64, parent *versionTreeNode) *versionTreeNode {
	return &versionTreeNode{
		version:  v,
		parent:   parent,
		children: []*versionTreeNode{},
	}
}

func (vt *VersionTree) findVersion(version uint64) *versionTreeNode {
	return vt.tree[version]
}
