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
	root     *list.List
}

// NewDoubleLinkedList creates new empty DoubleLinkedList.
func NewDoubleLinkedList[T any]() *DoubleLinkedList[T] {
	newList := &DoubleLinkedList[T]{
		versionTree: internal.NewVersionTree[listInfo](),
		storage:     make([]*internal.FatNode, 0),
	}
	err := newList.versionTree.SetVersionInfo(0, listInfo{
		listSize: 0,
		root:     list.New(),
	})
	if err != nil {
		return nil
	}
	return newList
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

	currentNode := info.root.Front()
	for i := 0; i < index; i++ {
		currentNode = currentNode.Next()
	}

	newVersion, err := l.versionTree.Update(version)
	if err != nil {
		return 0, err
	}
	err = l.versionTree.SetVersionInfo(newVersion, listInfo{
		listSize: info.listSize,
		root:     info.root,
	})
	if err != nil {
		return 0, err
	}
	l.storage[currentNode.Value.(int)].Update(value, newVersion)
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

	node := info.root.Front()
	for i := 0; i < index; i++ {
		node = node.Next()
	}
	info.root.Remove(node)

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

	node := info.root.Front()
	for i := 0; i < index; i++ {
		node = node.Next()
	}

	fatNode := l.storage[node.Value.(int)]
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

	newSeq := deepCopy(oldVersion.root)
	if isFront {
		newSeq.PushFront(oldVersion.listSize)
	} else {
		newSeq.PushBack(oldVersion.listSize)
	}

	err = l.versionTree.SetVersionInfo(newVersion, listInfo{
		listSize: oldVersion.listSize + 1,
		root:     newSeq,
	})
	if err != nil {
		return 0, err
	}

	l.storage = append(l.storage, internal.NewFatNode(value, newVersion))

	return newVersion, nil
}

func deepCopy(original *list.List) *list.List {
	newList := list.New()
	for e := original.Front(); e != nil; e = e.Next() {
		if e.Value != nil {
			newList.PushBack(e.Value)
		}
	}

	return newList
}
