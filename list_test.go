package main

import (
	"testing"
)

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
	_, err = list.PushFront(69, 3)
	if err != nil {
		t.Error("PushBack failed with", err)
	}

	newList, err := list.ToGoList(2)
	if err != nil {
		t.Error("ToGoList failed with", err)
	}
	println(newList)
}
