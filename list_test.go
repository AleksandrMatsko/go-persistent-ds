package main

import "testing"

func TestDoubleLinkedList_PushBack(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	_, err := list.PushFront(12, 0)
	if err != nil {
		t.Error("PushBack failed with", err)
	}
	_, err = list.PushFront(9, 1)
	if err != nil {
		t.Error("PushBack failed with", err)
	}
	_, err = list.PushFront(90, 2)
	if err != nil {
		t.Error("PushBack failed with", err)
	}
	val, err := list.PushFront(69, 3)
	if err != nil {
		t.Error("PushBack failed with", err)
	}
	if val.version != 4 {
		t.Error("Expected version 4, but got: ", val)
	}
}
