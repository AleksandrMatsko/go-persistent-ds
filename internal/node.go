package internal

// FatNode is a structure that stores values by versions.
type FatNode struct {
	nodes []*node
}

type node struct {
	data    interface{}
	version int
}

// NewFatNode creates new FatNode.
func NewFatNode(data interface{}, version int) *FatNode {
	return &FatNode{
		nodes: []*node{newNode(data, version)},
	}
}

// GetLast returns latest version of object inside FatNode.
func (fn *FatNode) GetLast() interface{} {
	return fn.nodes[len(fn.nodes)-1].data
}

// Update adds new object version into FatNode.
func (fn *FatNode) Update(data interface{}, newVersion int) {
	fn.nodes = append(fn.nodes, newNode(data, newVersion))
}

// FindByVersion finds needed version of object inside FatNode using binary search.
// If the version is not found, then the pair (nil, -1) is returned.
func (fn *FatNode) FindByVersion(version int) (interface{}, int) {
	left, right := 0, len(fn.nodes)-1

	for left <= right {
		mid := left + (right-left)/2

		if fn.nodes[mid].version == version {
			return fn.nodes[mid].data, fn.nodes[mid].version
		} else if fn.nodes[mid].version < version {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return nil, -1
}

// newNode creates new node inside FatNode.
func newNode(data interface{}, nodeVersion int) *node {
	return &node{
		data:    data,
		version: nodeVersion,
	}
}
