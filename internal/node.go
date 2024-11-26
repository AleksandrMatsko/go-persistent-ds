package internal

type FatNode struct {
	root    *Node
	history []int
}

type Node struct {
	data     interface{}
	children []*Node
	parent   *Node
	version  int
}

func NewFatNode(rootData interface{}, version int) *FatNode {
	return &FatNode{
		root: newNode(
			rootData,
			nil,
			version,
		),
		history: []int{version},
	}
}

func (fn *FatNode) GetLast() {
	// TODO
}

func (fn *FatNode) Update(data interface{}, prevVersion int, newVersion int) {
	node := fn.FindNodeByVersion(prevVersion)
	if node == nil {
		panic("Node not found!")
	}
	node.updateNode(data, newVersion)
}

func (fn *FatNode) FindNodeByVersion(version int) *Node {
	return findNodeByVersion(fn.root, version)
}

func newNode(data interface{}, parentNode *Node, nodeVersion int) *Node {
	return &Node{
		data:     data,
		children: make([]*Node, 0),
		parent:   parentNode,
		version:  nodeVersion,
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

func (n *Node) addChild(data interface{}, version int, parent *Node) *Node {
	newNode := newNode(data, parent, version)
	n.children = append(n.children, newNode)
	return newNode
}

func (n *Node) updateNode(newData interface{}, version int) *Node {
	newNode := newNode(newData, n, version)
	n.children = append(n.children, newNode)
	return newNode
}
