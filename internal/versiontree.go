package internal

import "log"

type VersionTree struct {
	tree []*versionTreeNode
}

type versionTreeNode struct {
	version  uint64
	parent   *versionTreeNode
	children []*versionTreeNode
}

func NewVersionTree(rootVersion uint64) *VersionTree {
	return &VersionTree{
		tree: []*versionTreeNode{newVersionTreeNode(rootVersion, nil)},
	}
}

func newVersionTreeNode(v uint64, parent *versionTreeNode) *versionTreeNode {
	return &versionTreeNode{
		version:  v,
		parent:   parent,
		children: []*versionTreeNode{},
	}
}

func (vt *VersionTree) Update(prevVersion uint64, newVersion uint64) bool {
	node, success := vt.findVersion(prevVersion)
	if !success {
		log.Fatal("version not found")
		return false
	}

	node.children = append(node.children, newVersionTreeNode(newVersion, node))
	return true
}

func (vt *VersionTree) findVersion(version uint64) (*versionTreeNode, bool) {
	left, right := 0, len(vt.tree)-1

	for left <= right {
		mid := left + (right-left)/2

		if vt.tree[mid].version == version {
			return vt.tree[mid], true
		} else if vt.tree[mid].version < version {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return nil, false
}
