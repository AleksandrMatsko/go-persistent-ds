package main

import (
	"testing"
)

func TestDoubleLinkedList_PushFront(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	_, err := list.PushFront(0, 12)
	errIsNil(t, err)

	_, err = list.PushFront(1, 9)
	errIsNil(t, err)

	_, err = list.PushFront(2, 90)
	errIsNil(t, err)

	version, err := list.PushFront(3, 69)
	errIsNil(t, err)

	goList, err := list.ToGoList(version)
	errIsNil(t, err)

	if goList.Back().Value != 12 {
		t.Error("Expected value of 12, got ", goList.Back().Value)
	}

	if goList.Back().Prev().Value != 9 {
		t.Error("Expected value of 9, got ", goList.Back().Value)
	}

	if goList.Front().Value != 69 {
		t.Error("Expected value of 69, got ", goList.Front().Value)
	}

	if goList.Front().Next().Value != 90 {
		t.Error("Expected value of 90, got ", goList.Front().Next().Value)
	}

	if goList.Len() != 4 {
		t.Error("Expected length of 4, got ", goList.Len())
	}
}

func TestDoubleLinkedList_PushBack(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	_, err := list.PushBack(0, 12)
	errIsNil(t, err)

	_, err = list.PushBack(1, 9)
	errIsNil(t, err)

	_, err = list.PushBack(2, 90)
	errIsNil(t, err)

	version, err := list.PushBack(3, 69)
	errIsNil(t, err)

	goList, err := list.ToGoList(version)
	errIsNil(t, err)

	if goList.Back().Value != 69 {
		t.Error("Expected value of 69, got ", goList.Back().Value)
	}

	if goList.Back().Prev().Value != 90 {
		t.Error("Expected value of 9, got ", goList.Back().Prev().Value)
	}

	if goList.Front().Value != 12 {
		t.Error("Expected value of 12, got ", goList.Front().Value)
	}

	if goList.Front().Next().Value != 9 {
		t.Error("Expected value of 9, got ", goList.Front().Next().Value)
	}

	if goList.Len() != 4 {
		t.Error("Expected length of 4, got ", goList.Len())
	}
}

func TestDoubleLinkedList_Len(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	_, err := list.PushBack(0, 12)
	errIsNil(t, err)

	_, err = list.PushBack(1, 9)
	errIsNil(t, err)

	_, err = list.PushBack(2, 90)
	errIsNil(t, err)

	version, err := list.PushBack(3, 69)
	errIsNil(t, err)

	goList, err := list.ToGoList(version)
	errIsNil(t, err)

	versionLength, err := list.Len(version)
	if versionLength != 4 {
		t.Error("Expected persistent list length of 4, got ", versionLength)
	}
	if goList.Len() != 4 {
		t.Error("Expected go list length of 4, got ", goList.Len())
	}
}

func TestDoubleLinkedList_Update(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	_, err := list.PushBack(0, 12)
	errIsNil(t, err)

	_, err = list.PushBack(1, 9)
	errIsNil(t, err)

	_, err = list.PushBack(2, 90)
	errIsNil(t, err)

	version, err := list.Update(2, 1, 99)
	if err != nil {
		return
	}

	val, err := list.Get(version, 1)
	if err != nil {
		t.Error(err)
	}
	if val != 99 {
		t.Error("Expected value of 99, got ", val)
	}
}
