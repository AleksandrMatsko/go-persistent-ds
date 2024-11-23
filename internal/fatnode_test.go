package internal

import (
	"strconv"
	"testing"
)

func TestFatNodeVersion(t *testing.T) {

	fn := NewFatNode("root")
	if fn.root.version != 0 {
		t.Error("Node version is incorrect. Expected 0, but was:" + strconv.Itoa(fn.root.version))
	}

	fn.root.UpdateNode("root1")
	if fn.root.version != 0 {
		t.Error("Node version is incorrect. Expected 0, but was:" + strconv.Itoa(fn.root.version))
	}
	if len(fn.root.children) != 1 {
		t.Error("No children after updating node!")
	}
	if fn.root.children[0].version != 1 {
		t.Error("Child node version is incorrect. Expected 1, but was:" + strconv.Itoa(fn.root.version))
	}
}

func TestFatNodeSearchByVersion(t *testing.T) {

	fn := NewFatNode("root")
	fn.root.UpdateNode("ch1")
	fn.root.UpdateNode("ch2")
	if len(fn.root.children) != 2 {
		t.Error("Expected 2 children after 2 root node updates!")
	}

	child := fn.FindNodeByVersion(1)
	if child.data != "ch1" {
		t.Error("Wrong child was found! Expected child with data = ch1 but got with data = " + child.data.(string))
	}

	child = fn.FindNodeByVersion(2)
	if child.data != "ch2" {
		t.Error("Wrong child was found! Expected child with data = ch2 but got with data = " + child.data.(string))
	}
}
