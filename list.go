package go_persistent_ds

import (
	"container/list"
	"errors"

	"github.com/AleksandrMatsko/go-persistent-ds/internal"
)

var ErrListIndexOutOfRange = errors.New("index out of range")

// DoubleLinkedList is a persistent implementation of double linked list.
// While working with list you can add to the start and to the end, access elements by index and modify.
// Each change of the list creates new version. Versions can be retrieved by their number.
// Also, there is an opportunity to convert this list into Go list using ToGoList.
//
// DoubleLinkedList can perform total of 2^65-1 modifications, and will panic on attempt to modify it for 2^65 time.
// If you need to continue editing DoubleLinkedList, the good idea is to use ToGoList method to dump list for special version.
//
// Note that DoubleLinkedList implementation is not thread safe.
type DoubleLinkedList[T any] struct {
	versionTree *internal.VersionTree[listInfo]
	storage     []*internal.FatNode
}

type listInfo struct {
	listSize int

	head *infoNode
	tail *infoNode
}

type infoNode struct {
	prev *internal.FatNode
	next *internal.FatNode

	value *internal.FatNode
}

// NewDoubleLinkedList creates new empty DoubleLinkedList.
func NewDoubleLinkedList[T any]() (*DoubleLinkedList[T], uint64) {
	newList := &DoubleLinkedList[T]{
		versionTree: internal.NewVersionTree[listInfo](),
		storage:     make([]*internal.FatNode, 0),
	}
	infoNode := &infoNode{
		prev: nil,
		next: nil,

		value: nil,
	}
	err := newList.versionTree.SetVersionInfo(0, listInfo{
		listSize: 0,

		head: infoNode,
		tail: infoNode,
	})
	if err != nil {
		panic(err)
	}
	return newList, 0
}

// NewDoubleLinkedListWithAnyValues creates DoubleLinkedList that can store values of any type.
func NewDoubleLinkedListWithAnyValues() (*DoubleLinkedList[any], uint64) {
	return NewDoubleLinkedList[any]()
}

// PushFront adds new element to the head of the DoubleLinkedList. Returns list's new version.
// Note: head->[1][2][3]<-tail.
//
// Complexity: O(n).
func (l *DoubleLinkedList[T]) PushFront(version uint64, value T) (uint64, error) {
	return l.push(value, version, true)
}

// PushBack adds new element to the tail of the DoubleLinkedList. Returns list's new version.
// Note: head->[1][2][3]<-tail.
//
// Complexity: O(1).
func (l *DoubleLinkedList[T]) PushBack(version uint64, value T) (uint64, error) {
	return l.push(value, version, false)
}

// Update updates element of specified DoubleLinkedList version by index. Returns list's new version.
//
// Complexity: O(n * log(m)), where n - DoubleLinkedList size and m - is number of changes in FatNode.
func (l *DoubleLinkedList[T]) Update(version uint64, index int, value T) (uint64, error) {
	info, err := l.versionTree.GetVersionInfo(version)
	if err != nil {
		return 0, err
	}
	if index >= info.listSize {
		return 0, ErrListIndexOutOfRange
	}

	changeHistory, err := l.versionTree.GetHistory(version)
	if err != nil {
		return 0, err
	}
	head := info.head
	iterInfo := head
	for i := 0; i < index; i++ {
		node := l.findNodeByChangeHistory(iterInfo.next, changeHistory, version)
		iterInfo = node.(*infoNode)
	}

	newVersion, err := l.versionTree.Update(version)
	if err != nil {
		return 0, err
	}

	iterInfo.value.Update(value, newVersion)

	err = l.versionTree.SetVersionInfo(newVersion, listInfo{
		listSize: info.listSize,
		head:     info.head,
		tail:     info.tail,
	})
	if err != nil {
		return 0, err
	}

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
// Complexity: O(n * log(m)), where n - DoubleLinkedList size and m - is number of changes in FatNode.
func (l *DoubleLinkedList[T]) Remove(version uint64, index int) (uint64, error) {
	info, err := l.versionTree.GetVersionInfo(version)
	if err != nil {
		return 0, err
	}
	if index >= info.listSize {
		return 0, ErrListIndexOutOfRange
	}

	changeHistory, err := l.versionTree.GetHistory(version)
	if err != nil {
		return 0, err
	}
	head := info.head
	iterInfo := head
	for i := 0; i < index; i++ {
		node := l.findNodeByChangeHistory(iterInfo.next, changeHistory, version)
		iterInfo = node.(*infoNode)
	}

	newVersion, err := l.versionTree.Update(version)
	if err != nil {
		return 0, err
	}

	previousNode := l.findNodeByChangeHistory(iterInfo.prev, changeHistory, version)
	nextNode := l.findNodeByChangeHistory(iterInfo.next, changeHistory, version)
	prevInfo := previousNode.(*infoNode)
	nextInfo := nextNode.(*infoNode)

	prevInfo.next.Update(nextInfo, newVersion)
	nextInfo.prev.Update(prevInfo, newVersion)

	iterInfo.next.Update(nil, newVersion)
	iterInfo.prev.Update(nil, newVersion)

	err = l.versionTree.SetVersionInfo(newVersion, listInfo{
		listSize: info.listSize - 1,
		head:     info.head,
		tail:     info.tail,
	})
	if err != nil {
		return 0, err
	}

	return newVersion, nil
}

// Get retrieves value from the specified DoubleLinkedList version by index.
//
// Complexity: O(n * log(m)), where n - DoubleLinkedList size and m - is number of changes in FatNode.
func (l *DoubleLinkedList[T]) Get(version uint64, index int) (T, error) {
	info, err := l.versionTree.GetVersionInfo(version)
	if err != nil {
		return *new(T), err
	}
	if index >= info.listSize {
		return *new(T), ErrListIndexOutOfRange
	}

	changeHistory, err := l.versionTree.GetHistory(version)
	if err != nil {
		return *new(T), err
	}
	head := info.head
	iterInfo := head
	for i := 0; i < index; i++ {
		node := l.findNodeByChangeHistory(iterInfo.next, changeHistory, version)
		iterInfo = node.(*infoNode)
	}

	val := l.findNodeByChangeHistory(iterInfo.value, changeHistory, version)
	return val.(T), nil
}

// ToGoList converts DoubleLinkedList into Go List.
//
// Complexity: O(Get) * n, where n - DoubleLinkedList size.
func (l *DoubleLinkedList[T]) ToGoList(version uint64) (*list.List, error) {
	info, err := l.versionTree.GetVersionInfo(version)
	if err != nil {
		return nil, err
	}

	newList := list.New()
	for i := 0; i < info.listSize; i++ {
		val, _ := l.Get(version, i)
		newList.PushBack(val)
	}

	return newList, nil
}

func (l *DoubleLinkedList[T]) push(value T, version uint64, isFront bool) (uint64, error) {
	oldVersionInfo, err := l.versionTree.GetVersionInfo(version)
	if err != nil {
		return 0, err
	}

	newVersion, err := l.versionTree.Update(version)
	if err != nil {
		return 0, err
	}
	newFatNode := internal.NewFatNode(value, newVersion)
	l.storage = append(l.storage, newFatNode)

	prevHead := oldVersionInfo.head
	prevTail := oldVersionInfo.tail

	if oldVersionInfo.listSize == 0 {
		newInfo := &infoNode{
			next: nil,
			prev: nil,

			value: newFatNode,
		}
		err = l.versionTree.SetVersionInfo(newVersion, listInfo{
			listSize: oldVersionInfo.listSize + 1,
			head:     newInfo,
			tail:     newInfo,
		})
	} else {
		if isFront {
			newHeadInfo := &infoNode{
				next:  internal.NewFatNode(prevHead, newVersion),
				prev:  nil,
				value: newFatNode,
			}
			if prevHead.prev == nil {
				prevHead.prev = internal.NewFatNode(newHeadInfo, newVersion)
			} else {
				prevHead.prev.Update(newHeadInfo, newVersion)
			}
			newListInfo := listInfo{
				listSize: oldVersionInfo.listSize + 1,
				head:     newHeadInfo,
				tail:     prevTail,
			}
			err = l.versionTree.SetVersionInfo(newVersion, newListInfo)
		} else {
			newTailInfo := &infoNode{
				next:  nil,
				prev:  internal.NewFatNode(prevTail, newVersion),
				value: newFatNode,
			}
			if prevTail.next == nil {
				prevTail.next = internal.NewFatNode(newTailInfo, newVersion)
			} else {
				prevTail.next.Update(newTailInfo, newVersion)
			}
			newListInfo := listInfo{
				listSize: oldVersionInfo.listSize + 1,
				head:     prevHead,
				tail:     newTailInfo,
			}
			err = l.versionTree.SetVersionInfo(newVersion, newListInfo)
		}
	}
	if err != nil {
		return 0, err
	}

	return newVersion, nil
}

func (l *DoubleLinkedList[T]) findNodeByChangeHistory(fn *internal.FatNode, changeHistory []uint64, version uint64) interface{} {
	val, _, found := fn.FindByVersion(version)
	if found {
		return val
	}

	if len(changeHistory) == 1 || len(changeHistory) == 2 {
		return nil
	}

	for i := len(changeHistory) - 2; i >= 1; i-- {
		val, _, found = fn.FindByVersion(changeHistory[i])
		if found {
			return val
		}
	}

	return nil
}
