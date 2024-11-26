package internal

import (
	"testing"
)

func TestFatNodeSearch(t *testing.T) {
	fatNode := FatNode{
		nodes: []*Node{
			{data: "Node 1", version: 1},
			{data: "Node 2", version: 2},
			{data: "Node 3", version: 3},
			{data: "Node 4", version: 4},
			{data: "Node 5", version: 5},
			{data: "Node 6", version: 6},
		},
	}

	node := fatNode.FindByVersion(3)
	if node == nil {
		t.Fatal("node was not found!")
	} else if node.version != 3 {
		t.Fatal("wrong node version!")
	}
}
