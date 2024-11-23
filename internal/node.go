package internal

type FatNode struct {
	root *Node
}

type Node struct {
	data           interface{}
	children       []*Node
	parent         *Node
	version        int
	versionMachine *VersionMachine
}

func NewFatNode(rootData interface{}) *FatNode {
	return &FatNode{
		root: newNode(
			rootData,
			nil,
			&VersionMachine{
				version: 0,
			},
		),
	}
}

func (n *Node) AddChild(data interface{}) *Node {
	newNode := newNode(data, n, n.versionMachine)
	n.children = append(n.children, newNode)
	return newNode
}

func (n *Node) UpdateNode(newData interface{}) *Node {
	newNode := newNode(newData, n, n.versionMachine)
	n.children = append(n.children, newNode)
	return newNode
}

func (fn *FatNode) FindNodeByVersion(version int) *Node {
	return findNodeByVersion(fn.root, version)
}

func newNode(data interface{}, parentNode *Node, vm *VersionMachine) *Node {
	return &Node{
		data:           data,
		children:       make([]*Node, 0),
		parent:         parentNode,
		version:        vm.GetAndIncrementVersion(),
		versionMachine: vm,
	}
}

func findNodeByVersion(node *Node, version int) *Node {
	if node.version == version {
		return node
	}
	for _, child := range node.children {
		result := findNodeByVersion(child, version)
		if result != nil {
			return result
		}
	}
	return nil
}
