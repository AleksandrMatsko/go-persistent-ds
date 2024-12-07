package main

import (
	"container/list"
	"errors"
	"go-persistent-ds/internal"
)

var ErrVersionNotFound = errors.New("version not found")
var ErrIndexOutOfRange = errors.New("index out of range")

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
	FatNodeNumber uint64
	prev, next    *ListNode
}

func NewDoubleLinkedList[T any]() *DoubleLinkedList[T] {
	return &DoubleLinkedList[T]{
		versionTree: internal.NewVersionTree[listInfo](),
		storage:     make([]*internal.FatNode, 0),
	}
}

func (l *DoubleLinkedList[T]) PushFront(value T, version uint64) (uint64, error) {
	return l.push(value, version, true)
}

func (l *DoubleLinkedList[T]) PushBack(value T, version uint64) (uint64, error) {
	return l.push(value, version, false)
}

func (l *DoubleLinkedList[T]) Update(value T, listNode *ListNode, version uint64) (uint64, error) {
	_, err := l.versionTree.GetVersionInfo(version)
	if err != nil {
		return 0, err
	}

	newVersion, err := l.versionTree.Update(version)
	if err != nil {
		return 0, err
	}
	l.storage[listNode.FatNodeNumber].Update(value, newVersion)
	return newVersion, nil
}

func (l *DoubleLinkedList[T]) Len(version uint64) (int, error) {
	info, err := l.versionTree.GetVersionInfo(version)
	if err != nil {
		return 0, err
	}
	return info.listSize, nil
}

func (l *DoubleLinkedList[T]) Get(index int, version uint64) (*ListNode, error) {
	info, err := l.versionTree.GetVersionInfo(version)
	if err != nil {
		return nil, err
	}
	if index >= info.listSize {
		return nil, ErrIndexOutOfRange
	}

	node := info.root.next
	for i := 0; i < index; i++ {
		node = node.next
	}

	return node, nil
}

func (l *DoubleLinkedList[T]) ToGoList(version uint64) (*list.List, error) {
	info, err := l.versionTree.GetVersionInfo(version)
	if err != nil {
		return nil, err
	}
	seq := info.root

	var newList = list.New()
	var curNode = seq.next
	for i := 0; i < info.listSize; i++ {
		fatNode := l.storage[curNode.FatNodeNumber]
		val, _, found := fatNode.FindByVersion(version)

		if !found {
			changeHistory, err := l.versionTree.GetHistory(version)
			if err != nil {
				return nil, err
			}

			if len(changeHistory) == 1 || len(changeHistory) == 2 {
				return nil, ErrVersionNotFound
			}

			for i := uint64(len(changeHistory) - 2); i >= 1; i-- {
				val, _, found = fatNode.FindByVersion(changeHistory[i])
				if found {
					if val == nil {
						return nil, ErrVersionNotFound
					}

					newList.PushBack(val)
				}
			}
		} else {
			newList.PushBack(val)
		}

		curNode = curNode.next
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
		} else {
			newRoot = insert(&ListNode{FatNodeNumber: uint64(oldVersion.listSize)}, deepCopy(oldVersion.root.prev))
			err = l.versionTree.SetVersionInfo(newVersion, listInfo{
				listSize: oldVersion.listSize + 1,
				root:     newRoot,
			})
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
