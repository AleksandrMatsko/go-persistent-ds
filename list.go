package main

import (
	"errors"
	"go-persistent-ds/internal"
)

var ErrVersionNotFound = errors.New("version not found")

// DoubleLinkedList is a persistent implementation of double linked list.
type DoubleLinkedList[T any] struct {
	versionTree *internal.VersionTree[listInfo]
	storage     []*internal.FatNode
}

type listInfo struct {
	listSize int
	root     *ListNode
}

type ListNode struct {
	version    uint64
	prev, next *ListNode
}

func NewDoubleLinkedList[T any]() *DoubleLinkedList[T] {
	return &DoubleLinkedList[T]{
		versionTree: internal.NewVersionTree[listInfo](),
		storage:     make([]*internal.FatNode, 0),
	}
}

func (list *DoubleLinkedList[T]) PushFront(value T, version uint64) (*ListNode, error) {
	return list.push(value, version, true)
}

func (list *DoubleLinkedList[T]) PushBack(value T, version uint64) (*ListNode, error) {
	return list.push(value, version, false)
}

func (list *DoubleLinkedList[T]) Back(version uint64) (*ListNode, error) {
	info, err := list.versionTree.GetVersionInfo(version)
	if err != nil {
		return nil, err
	}
	return info.root.prev, nil
}

func (list *DoubleLinkedList[T]) Front(version uint64) (*ListNode, error) {
	info, err := list.versionTree.GetVersionInfo(version)
	if err != nil {
		return nil, err
	}
	return info.root.next, nil
}

func (list *DoubleLinkedList[T]) Get(node *ListNode) (T, error) {
	data, _, found := list.storage[node.version].FindByVersion(node.version)
	if !found {
		return *new(T), ErrVersionNotFound
	}
	return data.(T), nil
}

func (list *DoubleLinkedList[T]) push(value T, version uint64, isFront bool) (*ListNode, error) {
	oldVersion, err := list.versionTree.GetVersionInfo(version)
	if err != nil {
		return nil, err
	}

	newVersion, err := list.versionTree.Update(version)
	if err != nil {
		return nil, err
	}

	var newRoot *ListNode
	if oldVersion.listSize == 0 {
		newRoot = &ListNode{
			version: newVersion,
			prev:    nil,
			next:    nil,
		}
		// making loop
		newRoot.prev = newRoot
		newRoot.next = newRoot

		err = list.versionTree.SetVersionInfo(newVersion, listInfo{
			listSize: 1,
			root:     newRoot,
		})
		if err != nil {
			return nil, err
		}
	} else {
		if isFront {
			newRoot = insert(&ListNode{version: newVersion}, deepCopy(oldVersion.root))
			err = list.versionTree.SetVersionInfo(newVersion, listInfo{
				listSize: oldVersion.listSize + 1,
				root:     newRoot,
			})
		} else {
			newRoot = insert(&ListNode{version: newVersion}, deepCopy(oldVersion.root.prev))
			err = list.versionTree.SetVersionInfo(newVersion, listInfo{
				listSize: oldVersion.listSize + 1,
				root:     newRoot,
			})
		}
	}

	list.storage = append(list.storage, internal.NewFatNode(value, newVersion))

	return newRoot, nil
}

// insert inserts newNode after prevNode
func insert(newNode *ListNode, prevNode *ListNode) *ListNode {
	newNode.prev = prevNode
	newNode.next = prevNode.next
	newNode.prev.next = newNode

	return newNode
}

// deepCopy creates full copy of the given list
func deepCopy(oldRoot *ListNode) *ListNode {
	if oldRoot == nil {
		return nil
	}

	// Creating copy of first node
	newHead := &ListNode{version: oldRoot.version}
	currentOriginal := oldRoot.next
	currentNew := newHead

	// Copying rest until reach the start point
	for currentOriginal != oldRoot {
		newNode := &ListNode{version: currentOriginal.version}

		// Connecting previous node with newNode
		currentNew.next = newNode
		newNode.prev = currentNew

		currentNew = newNode
		currentOriginal = currentOriginal.next
	}

	// Looping ends
	newHead.prev = currentNew
	currentNew.next = newHead

	return newHead
}
