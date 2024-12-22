package main

import (
	golist "container/list"
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

	listFromPersistentList, err := list.ToGoList(version)
	errIsNil(t, err)

	if listFromPersistentList.Back().Value != 12 {
		t.Error("Expected value of 12, got ", listFromPersistentList.Back().Value)
	}

	if listFromPersistentList.Back().Prev().Value != 9 {
		t.Error("Expected value of 9, got ", listFromPersistentList.Back().Value)
	}

	if listFromPersistentList.Front().Value != 69 {
		t.Error("Expected value of 69, got ", listFromPersistentList.Front().Value)
	}

	if listFromPersistentList.Front().Next().Value != 90 {
		t.Error("Expected value of 90, got ", listFromPersistentList.Front().Next().Value)
	}

	if listFromPersistentList.Len() != 4 {
		t.Error("Expected length of 4, got ", listFromPersistentList.Len())
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

	listFromPersistentList, err := list.ToGoList(version)
	errIsNil(t, err)

	if listFromPersistentList.Back().Value != 69 {
		t.Error("Expected value of 69, got ", listFromPersistentList.Back().Value)
	}

	if listFromPersistentList.Back().Prev().Value != 90 {
		t.Error("Expected value of 9, got ", listFromPersistentList.Back().Prev().Value)
	}

	if listFromPersistentList.Front().Value != 12 {
		t.Error("Expected value of 12, got ", listFromPersistentList.Front().Value)
	}

	if listFromPersistentList.Front().Next().Value != 9 {
		t.Error("Expected value of 9, got ", listFromPersistentList.Front().Next().Value)
	}

	if listFromPersistentList.Len() != 4 {
		t.Error("Expected length of 4, got ", listFromPersistentList.Len())
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

	listFromPersistentList, err := list.ToGoList(version)
	errIsNil(t, err)
	if listFromPersistentList.Back().Value != 90 {
		t.Error("Expected value of 90, got ", listFromPersistentList.Back().Value)
	}
	if listFromPersistentList.Back().Prev().Value != 12 {
		t.Error("Expected value of 12, got ", listFromPersistentList.Back().Prev().Value)
	}
	if listFromPersistentList.Front().Value != 69 {
		t.Error("Expected value of 69, got ", listFromPersistentList.Front().Value)
	}
	if listFromPersistentList.Front().Next().Value != 9 {
		t.Error("Expected value of 9, got ", listFromPersistentList.Front().Next().Value)
	}
}

func TestDoubleLinkedList_ManyPushesInOneVersion(t *testing.T) {
	list := NewDoubleLinkedList[string]()
	_, err := list.PushBack(0, "anna")
	errIsNil(t, err)
	_, err = list.PushBack(1, "oleg")
	errIsNil(t, err)
	_, err = list.PushBack(1, "natalia")
	errIsNil(t, err)
	_, err = list.PushBack(1, "alexander")
	errIsNil(t, err)

	val, err := list.Get(2, 1)
	if val != "oleg" {
		t.Error("Expected oleg in version 2 at position 1, got ", val)
	}
	val, err = list.Get(3, 1)
	if val != "natalia" {
		t.Error("Expected natalia in version 3 at position 1, got ", val)
	}
	val, err = list.Get(4, 1)
	if val != "alexander" {
		t.Error("Expected alexander in version 3 at position 1, got ", val)
	}
}

func TestDoubleLinkedList_BranchingCase(t *testing.T) {
	list := NewDoubleLinkedList[string]()
	_, err := list.PushBack(0, "anna")
	errIsNil(t, err)
	_, err = list.PushBack(1, "oleg")
	errIsNil(t, err)
	_, err = list.PushBack(2, "natalia")
	errIsNil(t, err)
	_, err = list.PushBack(3, "alexander")
	errIsNil(t, err)

	version, err := list.PushFront(2, "ilya")
	errIsNil(t, err)
	finVersion, err := list.PushBack(version, "filip")
	errIsNil(t, err)

	goList, err := list.ToGoList(finVersion)
	errIsNil(t, err)
	expectedList := golist.New()
	expectedList.PushBack("anna")
	expectedList.PushBack("oleg")
	expectedList.PushFront("ilya")
	expectedList.PushBack("filip")

	compareLists(goList, expectedList, t)
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

	listFromPersistentList, err := list.ToGoList(version)
	errIsNil(t, err)

	versionLength, err := list.Len(version)
	errIsNil(t, err)
	if versionLength != 4 {
		t.Error("Expected persistent list length of 4, got ", versionLength)
	}
	if listFromPersistentList.Len() != 4 {
		t.Error("Expected go list length of 4, got ", listFromPersistentList.Len())
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

	listFromPersistentList, err := list.ToGoList(version)
	errIsNil(t, err)

	expectedGoList := golist.New()
	expectedGoList.PushBack(12)
	expectedGoList.PushBack(9)
	expectedGoList.PushBack(90)
	expectedGoList.PushBack(69)

	compareLists(listFromPersistentList, expectedGoList, t)
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

func compareLists(actualList *golist.List, expectedList *golist.List, t *testing.T) {
	if actualList.Len() != expectedList.Len() {
		t.Errorf("Lengths differ: actualList length = %d, expectedList length = %d", actualList.Len(), expectedList.Len())
		return
	}

	actualElem := actualList.Front()
	expectedElem := expectedList.Front()

	for actualElem != nil && expectedElem != nil {
		if actualElem.Value != expectedElem.Value {
			t.Errorf("Lists differ at element: actual = %v, expected = %v", actualElem.Value, expectedElem.Value)
			return
		}
		actualElem = actualElem.Next()
		expectedElem = expectedElem.Next()
	}
}
