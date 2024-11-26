package internal

type FatNode struct {
	nodes []*Node
}

type Node struct {
	data    interface{}
	version int
}

// NewFatNode создаёт новый "толстый узел".
func NewFatNode(data interface{}, version int) *FatNode {
	return &FatNode{
		nodes: []*Node{newNode(data, version)},
	}
}

// GetLast возвращает последнюю версию объекта внутри "толстого узла".
func (fn *FatNode) GetLast() *Node {
	return fn.nodes[len(fn.nodes)-1]
}

// Update добавляет новую версию объекта внутри "толстого узла".
func (fn *FatNode) Update(data interface{}, newVersion int) {
	fn.nodes = append(fn.nodes, newNode(data, newVersion))
}

// FindByVersion находит нужную версию объекта внутри "толстого узла" бинарным поиском.
func (fn *FatNode) FindByVersion(version int) *Node {
	left, right := 0, len(fn.nodes)-1

	for left <= right {
		mid := left + (right-left)/2

		if fn.nodes[mid].version == version {
			return fn.nodes[mid]
		} else if fn.nodes[mid].version < version {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return nil
}

func newNode(data interface{}, nodeVersion int) *Node {
	return &Node{
		data:    data,
		version: nodeVersion,
	}
}
