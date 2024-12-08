package main

import (
	"container/list"
	"errors"
	"go-persistent-ds/internal"
)

var ErrVersionNotFound = errors.New("version not found")
var ErrIndexOutOfRange = errors.New("index out of range")

// DoubleLinkedList is a persistent implementation of double linked list.
// While working with list you can add to the start and to the end, access elements by index and modify.
// Each change of the list creates new version. Versions can be retrieved by their number.
// Also, there is an opportunity to convert this list into Go list using ToGoList.
//
// Note that DoubleLinkedList implementation is not thread safe.
type DoubleLinkedList[T any] struct {
	versionTree *internal.VersionTree[listInfo]
	storage     []*internal.FatNode
}

type listInfo struct {
	listSize int
	root     *ListNode
}

type ListNode struct {
	FatNodeNumber uint64
	prev, next    *ListNode
}

// NewDoubleLinkedList creates new empty DoubleLinkedList.
func NewDoubleLinkedList[T any]() *DoubleLinkedList[T] {
	return &DoubleLinkedList[T]{
		versionTree: internal.NewVersionTree[listInfo](),
		storage:     make([]*internal.FatNode, 0),
	}
}

// PushFront adds new element to the end of the DoubleLinkedList. Returns list's new version.
//
// Complexity: O(n), where n - DoubleLinkedList size. (It's because of need to create deep copy of list's connections).
func (l *DoubleLinkedList[T]) PushFront(version uint64, value T) (uint64, error) {
	return l.push(value, version, true)
}

// PushBack adds new element to the start of the DoubleLinkedList. Returns list's new version.
//
// Complexity: O(n), where n - DoubleLinkedList size. (It's because of need to create deep copy of list's connections).
func (l *DoubleLinkedList[T]) PushBack(version uint64, value T) (uint64, error) {
	return l.push(value, version, false)
}

// Update updates element of specified DoubleLinkedList version by index. Returns list's new version.
//
// Complexity: O(n), where n - DoubleLinkedList size.
func (l *DoubleLinkedList[T]) Update(version uint64, index int, value T) (uint64, error) {
	info, err := l.versionTree.GetVersionInfo(version)
	if err != nil {
		return 0, err
	}
	if index >= info.listSize {
		return 0, ErrIndexOutOfRange
	}

	currentNode := info.root.next
	for i := 0; i < index; i++ {
		currentNode = currentNode.next
	}

	newVersion, err := l.versionTree.Update(version)
	if err != nil {
		return 0, err
	}
	l.storage[currentNode.FatNodeNumber].Update(value, newVersion)
	return newVersion, nil
}

// Len returns DoubleLinkedList size.
//
// Complexity: O(1).
func (l *DoubleLinkedList[T]) Len(version uint64) (int, error) {
	info, err := l.versionTree.GetVersionInfo(version)
	if err != nil {
		return 0, err
	}
	return info.listSize, nil
}

// Remove removes element from specified version of DoubleLinkedList by index and returns new list's version.
// By removal, we mean delete of connection between specified element and his "neighbours".
//
// Complexity: O(n), where n - DoubleLinkedList size.
func (l *DoubleLinkedList[T]) Remove(version uint64, index int) (uint64, error) {
	info, err := l.versionTree.GetVersionInfo(version)
	if err != nil {
		return 0, err
	}
	if index >= info.listSize {
		return 0, ErrIndexOutOfRange
	}

	v, err := l.versionTree.Update(version)
	if err != nil {
		return 0, err
	}

	node := info.root.next
	for i := 0; i < index; i++ {
		node = node.next
	}
	node.prev.next = node.next
	node.next.prev = node.prev

	return v, nil
}

// Get retrieves value from the specified DoubleLinkedList version by index.
//
// Complexity: O(n), where n - DoubleLinkedList size.
func (l *DoubleLinkedList[T]) Get(version uint64, index int) (T, error) {
	info, err := l.versionTree.GetVersionInfo(version)
	if err != nil {
		return *new(T), err
	}
	if index >= info.listSize {
		return *new(T), ErrIndexOutOfRange
	}

	node := info.root.next
	for i := 0; i < index; i++ {
		node = node.next
	}

	fatNode := l.storage[node.FatNodeNumber]
	val, _, found := fatNode.FindByVersion(version)
	if !found {
		changeHistory, err := l.versionTree.GetHistory(version)
		if err != nil {
			return *new(T), err
		}

		if len(changeHistory) == 1 || len(changeHistory) == 2 {
			return *new(T), ErrVersionNotFound
		}

		for i := uint64(len(changeHistory) - 2); i >= 1; i-- {
			val, _, found = fatNode.FindByVersion(changeHistory[i])
			if found {
				if val == nil {
					return *new(T), ErrVersionNotFound
				}

				return val.(T), nil
			}
		}
	} else {
		return val.(T), nil
	}

	return *new(T), nil
}

// ToGoList converts DoubleLinkedList into Go List.
//
// Complexity: O(Get) * n, where n - DoubleLinkedList size.
func (l *DoubleLinkedList[T]) ToGoList(version uint64) (*list.List, error) {
	info, err := l.versionTree.GetVersionInfo(version)
	if err != nil {
		return nil, err
	}

	var newList = list.New()
	for i := 0; i < info.listSize; i++ {
		val, _ := l.Get(version, i)
		newList.PushBack(val)
	}

	return newList, nil
}

func (l *DoubleLinkedList[T]) push(value T, version uint64, isFront bool) (uint64, error) {
	oldVersion, err := l.versionTree.GetVersionInfo(version)
	if err != nil {
		return 0, err
	}

	newVersion, err := l.versionTree.Update(version)
	if err != nil {
		return 0, err
	}

	var newRoot *ListNode
	if oldVersion.listSize == 0 {
		newRoot = &ListNode{
			FatNodeNumber: 0,
			prev:          nil,
			next:          nil,
		}
		// making loop.
		newRoot.prev = newRoot
		newRoot.next = newRoot

		err = l.versionTree.SetVersionInfo(newVersion, listInfo{
			listSize: 1,
			root:     newRoot,
		})
		if err != nil {
			return 0, err
		}
	} else {
		if isFront {
			newRoot = insert(&ListNode{FatNodeNumber: uint64(oldVersion.listSize)}, deepCopy(oldVersion.root))
			err = l.versionTree.SetVersionInfo(newVersion, listInfo{
				listSize: oldVersion.listSize + 1,
				root:     newRoot,
			})
			if err != nil {
				return 0, err
			}
		} else {
			newRoot = insert(&ListNode{FatNodeNumber: uint64(oldVersion.listSize)}, deepCopy(oldVersion.root.prev))
			err = l.versionTree.SetVersionInfo(newVersion, listInfo{
				listSize: oldVersion.listSize + 1,
				root:     newRoot,
			})
			if err != nil {
				return 0, err
			}
		}
	}

	l.storage = append(l.storage, internal.NewFatNode(value, newVersion))

	return newVersion, nil
}

// insert inserts newNode after prevNode.
func insert(newNode *ListNode, prevNode *ListNode) *ListNode {
	newNode.prev = prevNode
	newNode.next = prevNode.next
	newNode.prev.next = newNode

	return newNode
}

// deepCopy creates full copy of the given list.
func deepCopy(oldRoot *ListNode) *ListNode {
	if oldRoot == nil {
		return nil
	}

	// Creating copy of first node.
	newHead := &ListNode{FatNodeNumber: oldRoot.FatNodeNumber}
	currentOriginal := oldRoot.next
	currentNew := newHead

	// Copying rest until reach the start point.
	for currentOriginal != oldRoot {
		newNode := &ListNode{FatNodeNumber: currentOriginal.FatNodeNumber}

		// Connecting previous node with newNode.
		currentNew.next = newNode
		newNode.prev = currentNew

		currentNew = newNode
		currentOriginal = currentOriginal.next
	}

	// Looping ends.
	newHead.prev = currentNew
	currentNew.next = newHead

	return newHead
}
