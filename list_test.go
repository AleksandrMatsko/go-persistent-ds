package main

import (
	list2 "container/list"
	"errors"
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

func TestDoubleLinkedList_MixedPush(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	_, err := list.PushBack(0, 12)
	errIsNil(t, err)
	_, err = list.PushFront(1, 9)
	errIsNil(t, err)
	_, err = list.PushBack(2, 90)
	errIsNil(t, err)
	version, err := list.PushFront(3, 69)
	errIsNil(t, err)

	goList, err := list.ToGoList(version)
	errIsNil(t, err)
	if goList.Back().Value != 90 {
		t.Error("Expected value of 90, got ", goList.Back().Value)
	}
	if goList.Back().Prev().Value != 12 {
		t.Error("Expected value of 12, got ", goList.Back().Prev().Value)
	}
	if goList.Front().Value != 69 {
		t.Error("Expected value of 69, got ", goList.Front().Value)
	}
	if goList.Front().Next().Value != 9 {
		t.Error("Expected value of 9, got ", goList.Front().Next().Value)
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
	errIsNil(t, err)
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
	errIsNil(t, err)

	val, err := list.Get(version, 1)
	errIsNil(t, err)
	if val != 99 {
		t.Error("Expected value of 99, got ", val)
	}
}

func TestDoubleLinkedList_Remove(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	_, err := list.PushBack(0, 12)
	errIsNil(t, err)

	_, err = list.PushBack(1, 9)
	errIsNil(t, err)

	_, err = list.PushBack(2, 90)
	errIsNil(t, err)

	version, err := list.PushBack(3, 69)
	errIsNil(t, err)

	lengthByLenFunction, err := list.Len(version)
	errIsNil(t, err)
	lengthByVersionInfo, err := list.versionTree.GetVersionInfo(version)
	errIsNil(t, err)
	if lengthByLenFunction != lengthByVersionInfo.listSize {
		t.Error("Expected equal list length by Len() function and version info!")
	}
	if lengthByLenFunction != 4 || lengthByVersionInfo.listSize != 4 {
		t.Errorf("Expected list length of 4, got %d by Len(), and %d by version info",
			lengthByLenFunction, lengthByVersionInfo.listSize)
	}

	version, err = list.Remove(version, 1)
	errIsNil(t, err)
	lengthByLenFunction, err = list.Len(version)
	errIsNil(t, err)
	lengthByVersionInfo, err = list.versionTree.GetVersionInfo(version)
	errIsNil(t, err)
	if lengthByLenFunction != lengthByVersionInfo.listSize {
		t.Error("Expected equal list length by Len() function and version info!")
	}
	if lengthByLenFunction != 3 || lengthByVersionInfo.listSize != 3 {
		t.Errorf("Expected list length of 4, got %d by Len(), and %d by version info",
			lengthByLenFunction, lengthByVersionInfo.listSize)
	}
}

func TestDoubleLinkedList_ToGoList(t *testing.T) {
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

	list1 := list2.New()
	list1.PushBack(12)
	list1.PushBack(9)
	list1.PushBack(90)
	list1.PushBack(69)

	if goList.Len() != 4 && list1.Len() != 4 {
		t.Errorf("Expected persistent list length of 4, got %d. Expected Go list length 4, got %d",
			goList.Len(), list1.Len())
	}
	if goList.Back().Value != list1.Back().Value {
		t.Errorf("Expected equal values! Value %d from Go list not equal to value %d from persistent list.",
			list1.Back().Value, goList.Back().Value)
	}
	if goList.Front().Value != list1.Front().Value {
		t.Errorf("Expected equal values! Value %d from Go list not equal to value %d from persistent list.",
			list1.Front().Value, goList.Front().Value)
	}
	if goList.Back().Prev().Value != list1.Back().Prev().Value {
		t.Errorf("Expected equal values! Value %d from Go list not equal to value %d from persistent list.",
			list1.Back().Prev().Value, goList.Back().Prev().Value)
	}
	if goList.Front().Next().Value != list1.Front().Next().Value {
		t.Errorf("Expected equal values! Value %d from Go list not equal to value %d from persistent list.",
			list1.Front().Next().Value, goList.Front().Next().Value)
	}
}

func TestDoubleLinkedList_Get(t *testing.T) {
	list := NewDoubleLinkedList[int]()
	_, err := list.PushBack(0, 12)
	errIsNil(t, err)

	_, err = list.PushBack(1, 9)
	errIsNil(t, err)

	_, err = list.PushBack(2, 90)
	errIsNil(t, err)

	version, err := list.PushBack(3, 69)
	errIsNil(t, err)

	val, err := list.Get(version, 2)
	errIsNil(t, err)
	if val != 90 {
		t.Error("Expected value of 90, got ", val)
	}
	val, err = list.Get(version, 3)
	errIsNil(t, err)
	if val != 69 {
		t.Error("Expected value of 69, got ", val)
	}

	val, err = list.Get(2, 3)
	if !errors.Is(err, ErrListIndexOutOfRange) {
		t.Error("Expected ErrIndexOutOfRange, got nil")
	}
}
