package go_persistent_ds

import (
	golist "container/list"
	"errors"
	"testing"
)

func TestDoubleLinkedList_PushFront(t *testing.T) {
	list, _ := NewDoubleLinkedList[int]()
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
	list, _ := NewDoubleLinkedList[int]()
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
	list, _ := NewDoubleLinkedList[int]()
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
	list, _ := NewDoubleLinkedList[string]()
	_, err := list.PushBack(0, "anna")
	errIsNil(t, err)
	_, err = list.PushBack(1, "oleg")
	errIsNil(t, err)
	_, err = list.PushBack(1, "natalia")
	errIsNil(t, err)
	_, err = list.PushBack(1, "alexander")
	errIsNil(t, err)

	val, err := list.Get(2, 1)
	errIsNil(t, err)
	if val != "oleg" {
		t.Error("Expected oleg in version 2 at position 1, got ", val)
	}
	val, err = list.Get(3, 1)
	errIsNil(t, err)
	if val != "natalia" {
		t.Error("Expected natalia in version 3 at position 1, got ", val)
	}
	val, err = list.Get(4, 1)
	errIsNil(t, err)
	if val != "alexander" {
		t.Error("Expected alexander in version 3 at position 1, got ", val)
	}
}

func TestDoubleLinkedList_BranchingCase(t *testing.T) {
	list, _ := NewDoubleLinkedList[string]()
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
	list, _ := NewDoubleLinkedList[int]()
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
	list, _ := NewDoubleLinkedList[int]()
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

func TestDoubleLinkedList_BranchedUpdate(t *testing.T) {
	list, _ := NewDoubleLinkedList[string]()
	_, err := list.PushBack(0, "anna")
	errIsNil(t, err)
	_, err = list.PushBack(1, "oleg")
	errIsNil(t, err)
	_, err = list.PushBack(2, "natalia")
	errIsNil(t, err)
	checkVersion, err := list.PushBack(3, "alexander")
	errIsNil(t, err)
	version, err := list.PushFront(2, "ilya")
	errIsNil(t, err)
	updVersion, err := list.PushBack(version, "filip")
	errIsNil(t, err)
	updVersion, err = list.Update(updVersion, 0, "ilya2")
	errIsNil(t, err)

	actualList1, err := list.ToGoList(updVersion)
	errIsNil(t, err)
	actualList2, err := list.ToGoList(checkVersion)
	errIsNil(t, err)

	expectedList1 := golist.New()
	expectedList1.PushBack("anna")
	expectedList1.PushBack("oleg")
	expectedList1.PushFront("ilya2")
	expectedList1.PushBack("filip")
	compareLists(actualList1, expectedList1, t)

	expectedList2 := golist.New()
	expectedList2.PushBack("anna")
	expectedList2.PushBack("oleg")
	expectedList2.PushBack("natalia")
	expectedList2.PushBack("alexander")
	compareLists(actualList2, expectedList2, t)
}

func TestDoubleLinkedList_Remove(t *testing.T) {
	list, _ := NewDoubleLinkedList[int]()
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

func TestDoubleLinkedList_BranchedRemove(t *testing.T) {
	list, _ := NewDoubleLinkedList[string]()
	_, err := list.PushBack(0, "anna")
	errIsNil(t, err)
	_, err = list.PushBack(1, "oleg")
	errIsNil(t, err)
	_, err = list.PushBack(2, "vsevolod")
	errIsNil(t, err)
	branchVersion, err := list.PushBack(1, "natalia")
	errIsNil(t, err)
	branchVersion, err = list.PushBack(branchVersion, "alexander")
	errIsNil(t, err)

	removeVersion, err := list.Remove(branchVersion, 1)
	errIsNil(t, err)
	prevLen, err := list.Len(branchVersion)
	errIsNil(t, err)
	actualLen, err := list.Len(removeVersion)
	errIsNil(t, err)

	if prevLen == actualLen {
		t.Error("Expected different list sizes!")
	}
	if prevLen != 3 {
		t.Error("Expected before remove list length of 3, got ", prevLen)
	}
	if actualLen != 2 {
		t.Error("Expected after remove list length of 2, got ", actualLen)
	}
}

func TestDoubleLinkedList_ToGoList(t *testing.T) {
	list, _ := NewDoubleLinkedList[int]()
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
	list, _ := NewDoubleLinkedList[int]()
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

	_, err = list.Get(2, 3)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !errors.Is(err, ErrListIndexOutOfRange) {
		t.Error("Expected ErrIndexOutOfRange, got nil")
	}
}

func TestDoubleLinkedList_BranchedGet(t *testing.T) {
	list, _ := NewDoubleLinkedList[string]()
	_, err := list.PushBack(0, "anna")
	errIsNil(t, err)
	_, err = list.PushBack(1, "oleg")
	errIsNil(t, err)
	_, err = list.PushBack(2, "natalia")
	errIsNil(t, err)
	branch1Version, err := list.PushBack(3, "alexander")
	errIsNil(t, err)
	version, err := list.PushFront(2, "ilya")
	errIsNil(t, err)
	branch2Version, err := list.PushBack(version, "filip")
	errIsNil(t, err)

	val, err := list.Get(branch1Version, 0)
	errIsNil(t, err)
	if val != "anna" {
		t.Error("Expected value of anna, got ", val)
	}

	val, err = list.Get(branch2Version, 0)
	errIsNil(t, err)
	if val != "ilya" {
		t.Error("Expected value of ilya, got ", val)
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

func TestNewDoubleLinkedListWithAnyValues(t *testing.T) {
	t.Run("Push and Get values ok", func(t *testing.T) {
		t.Parallel()

		l, v := NewDoubleLinkedListWithAnyValues()
		versionShouldBe(t, v, 0)

		v, err := l.PushBack(0, "a")
		errIsNil(t, err)
		versionShouldBe(t, v, 1)

		v, err = l.PushBack(1, 1)
		errIsNil(t, err)
		versionShouldBe(t, v, 2)

		v, err = l.PushBack(2, map[string]string{})
		errIsNil(t, err)
		versionShouldBe(t, v, 3)

		v, err = l.PushBack(3, []uint64{})
		errIsNil(t, err)
		versionShouldBe(t, v, 4)

		val, err := l.Get(4, 0)
		errIsNil(t, err)
		isTrue(t, val == "a")

		val, err = l.Get(4, 1)
		errIsNil(t, err)
		isTrue(t, val == 1)

		val, err = l.Get(4, 2)
		errIsNil(t, err)
		isTrue(t, func() bool {
			castedVal, ok := val.(map[string]string)
			if !ok {
				return false
			}

			return len(castedVal) == 0
		}())

		val, err = l.Get(4, 3)
		errIsNil(t, err)
		isTrue(t, func() bool {
			castedVal, ok := val.([]uint64)
			if !ok {
				return false
			}

			return len(castedVal) == 0
		}())
	})
}
