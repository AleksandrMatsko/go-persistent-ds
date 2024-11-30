package internal

import (
	"testing"
)

func TestFatNodeSearch(t *testing.T) {
	fatNode := FatNode{
		nodes: []*node{
			{data: "Node 1", version: 1},
			{data: "Node 2", version: 2},
			{data: "Node 3", version: 5},
			{data: "Node 4", version: 9},
			{data: "Node 5", version: 10},
			{data: "Node 6", version: 13},
		},
	}

	data, version, success := fatNode.FindByVersion(3)
	if data != nil && version != 0 && success != false {
		t.Fatalf("Expected nil, -1, false, but got: %s, %d, %v", data, version, success)
	}

	data, version, success = fatNode.FindByVersion(9)
	if data != "Node 4" && version != 9 && success != true {
		t.Fatal("Expected 9, got: ", version)
	}
}
