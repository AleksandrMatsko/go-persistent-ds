package internal

import "log"

// VersionTree is a struct to store object change history.
type VersionTree struct {
	tree []*versionTreeNode
}

type versionTreeNode struct {
	version  uint64
	parent   *versionTreeNode
	children []*versionTreeNode
}

// NewVersionTree creates new object change history tree.
func NewVersionTree(rootVersion uint64) *VersionTree {
	if rootVersion != 1 {
		return nil
	}

	return &VersionTree{
		tree: []*versionTreeNode{newVersionTreeNode(rootVersion, nil)},
	}
}

// Update creates new version for specified version.
func (vt *VersionTree) Update(prevVersion uint64, newVersion uint64) bool {
	if prevVersion >= newVersion {
		return false
	}
	len1 := uint64(len(vt.tree)) + 1
	if newVersion != len1 {
		return false
	}

	node, success := vt.findVersion(prevVersion)
	if !success {
		log.Fatal("version not found")
		return false
	}

	newNode := newVersionTreeNode(newVersion, node)
	node.children = append(node.children, newNode)
	vt.tree = append(vt.tree, newNode)

	return true
}

// GetHistory returns change history for specified object's version.
func (vt *VersionTree) GetHistory(version uint64) []uint64 {
	node, success := vt.findVersion(version)
	if !success {
		log.Fatal("version not found")
		return nil
	}

	var history []uint64
	for node != nil {
		history = append(history, node.version)
		node = node.parent
	}
	reverse(history)

	return history
}

func newVersionTreeNode(v uint64, parent *versionTreeNode) *versionTreeNode {
	return &versionTreeNode{
		version:  v,
		parent:   parent,
		children: []*versionTreeNode{},
	}
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

func reverse(s []uint64) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
