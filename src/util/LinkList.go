package util

import (
	"strconv"
	"sync"
)

// LinkList is a linked list
type LinkList[T any] struct {
	head      *linkListNode[T]
	tail      *linkListNode[T]
	len       int
	rwLock    *sync.RWMutex
	iteRWLock *sync.Mutex
	equal     func(a, b T) bool
}

// linkListNode is a node of linked list
type linkListNode[T any] struct {
	object T
	next   *linkListNode[T]
	prev   *linkListNode[T]
}

// NewLinkList
// @Description: create a new linked list
// @return       *LinkList[T] : the new linked list
func NewLinkList[T any](f func(a, b T) bool) *LinkList[T] {
	lock := &sync.RWMutex{}
	iteLock := &sync.Mutex{}
	Node := &linkListNode[T]{
		next: nil,
		prev: nil,
	}
	return &LinkList[T]{
		head:      Node,
		tail:      Node,
		len:       0,
		rwLock:    lock,
		iteRWLock: iteLock,
		equal:     f,
	}
}

// Append
// @Description: append an object to the linked list
// @receiver     L
// @param        data
func (L *LinkList[T]) Append(data T) {
	L.rwLock.Lock()
	L.len++
	L.tail.object = data
	Node := &linkListNode[T]{
		prev: L.tail,
		next: nil,
	}
	L.tail.next = Node
	L.tail = Node
	L.rwLock.Unlock()
}

// ToString
// @Description: convert the linked list to string
// @receiver     L      : the linked list
// @return       string : the string of the linked list
func (L *LinkList[T]) ToString() string {
	L.rwLock.RLock()
	var S string
	S += "[ "
	next := L.head
	for next.next != nil {
		S += Strval(next.object) + ", "
		next = next.next
	}
	S += "]"
	L.rwLock.RUnlock()
	return S
}

// ToStringBack
// @Description: convert the linked list to string in reverse order
// @receiver     L      : the linked list
// @return       string : the string of the linked list
func (L *LinkList[T]) ToStringBack() string {
	L.rwLock.RLock()
	var S string
	S += "[ "
	prev := L.tail.prev
	for prev != nil {
		S += Strval(prev.object) + ", "
		prev = prev.prev
	}
	S += "]"
	L.rwLock.RUnlock()
	return S
}

// Length
// @Description: get the length of the linked list
// @receiver     L   : the linked list
// @return       int : the length of the linked list
func (L *LinkList[T]) Length() int {
	L.rwLock.RLock()
	defer L.rwLock.RUnlock()
	return L.len
}

// Get
// @Description: get the object at the index
// @receiver     L   : the linked list
// @param        num : the index
// @return       T   : the object at the index
func (L *LinkList[T]) Get(num int) T {
	L.rwLock.RLock()
	defer L.rwLock.RUnlock()
	if num >= L.len {
		length := strconv.Itoa(L.len)
		panic("index out! max index: " + length)
	}
	if num >= L.len/2 {
		n := L.len - num - 1
		t := L.tail.prev
		for t.prev != nil && n != 0 {
			n--
			t = t.prev
		}
		return t.object
	} else {
		t := L.head
		n := num
		for t.next != nil && n != 0 {
			n--
			t = t.next
		}
		return t.object
	}
}

// Clear
// @Description: clear the linked list
// @receiver     L    : the linked list
// @return       bool : true if the linked list is empty
func (L *LinkList[T]) Clear() bool {
	L.rwLock.Lock()
	Node := &linkListNode[T]{
		next: nil,
		prev: nil,
	}
	L.head = Node
	L.tail = Node
	L.len = 0
	L.rwLock.Unlock()
	return true
}

// DeleteByNum
// @Description: delete the object at the index
// @receiver     L    : the linked list
// @param        num : the index
// @return       bool : true if the object is deleted
func (L *LinkList[T]) DeleteByNum(num int) bool {
	L.rwLock.Lock()
	defer L.rwLock.Unlock()
	if num >= L.len {
		length := strconv.Itoa(L.len)
		panic("index out! max index: " + length)
	}
	if num == 0 {
		L.head = L.head.next
		L.len--
		return true
	}
	if num == L.len-1 {
		L.tail = L.tail.prev
		L.len--
		return true
	}
	if num >= L.len/2 {
		num := L.len - num
		t := L.tail
		for t.prev != nil && num != 0 {
			num--
			t = t.prev
		}
		t.prev.next = t.next
		t.next.prev = t.prev
		L.len--
		return true
	} else {
		t := L.head
		for t.next != nil && num != 0 {
			num--
			t = t.next
		}
		t.prev.next = t.next
		t.next.prev = t.prev
		L.len--
		return true
	}
}

// DeleteByValue
// @Description: delete the object by the value
// @receiver     L    : the linked list
// @param        o : the value
// @return       bool : true if the object is deleted
func (L *LinkList[T]) DeleteByValue(o T) bool {
	L.rwLock.Lock()
	defer L.rwLock.Unlock()
	if L.len == 0 {
		return false
	}
	t := L.head
	if L.equal(t.object, o) {
		L.head.next.prev = nil
		L.head = L.head.next
		L.len--
		return true
	}
	for t.next.next != nil {
		if L.equal(t.object, o) {
			t.prev.next = t.next
			t.next.prev = t.prev
			L.len--
			return true
		}
		t = t.next
	}
	if L.equal(t.object, o) {
		t.prev.next = L.tail
		L.tail.prev = t.prev
		L.len--
		return true
	}
	return false
}

// DeleteAllByValue
// @Description: delete all the objects by the value
// @receiver     L    : the linked list
// @param        o : the value
// @return       bool : true if the object is deleted
func (L *LinkList[T]) DeleteAllByValue(o T) bool {
	L.rwLock.Lock()
	defer L.rwLock.Unlock()
	if L.len == 0 {
		return false
	}
L1:
	t := L.head
	if L.equal(t.object, o) {
		L.head.next.prev = nil
		L.head = L.head.next
		goto L1
	}
	for t != nil && t.next != nil && t.next.next != nil {
		if L.equal(t.object, o) {
			t.prev.next = t.next
			t.next.prev = t.prev
			L.len--
		}
		t = t.next
	}
	if L.equal(t.object, o) {
		t.prev.next = L.tail
		L.tail.prev = t.prev
		L.len--
	}
	return true
}

// InsertValue
// @Description: insert the value at the index
// @receiver     L    : the linked list
// @param        num : the index
// @param        v   : the value
// @return       bool : true if the value is inserted
func (L *LinkList[T]) InsertValue(num int, v T) bool {
	L.rwLock.Lock()
	defer L.rwLock.Unlock()
	if num > L.len {
		length := strconv.Itoa(L.len)
		panic("index out! max index: " + length)
	}
	if num == L.len {
		L.rwLock.Unlock()
		L.Append(v)
		return true
	}
	if num == 0 {
		Node := &linkListNode[T]{
			object: v,
			next:   L.head,
			prev:   nil,
		}
		L.head.prev = Node
		L.head = Node
		L.len++
		return true
	}
	t := L.head
	for num != 0 {
		num--
		t = t.next
	}
	Node := &linkListNode[T]{
		object: v,
		next:   t,
		prev:   t.prev,
	}
	t.prev.next = Node
	t.prev = Node
	L.len++
	return true
}

// UpdateValue
// @Description: update the value at the index
// @receiver     L    : the linked list
// @param        num : the index
// @param        v   : the value
// @return       bool : true if the value is updated
func (L *LinkList[T]) UpdateValue(num int, v T) bool {
	L.rwLock.Lock()
	defer L.rwLock.Unlock()
	t := L.head
	if num >= L.len {
		length := strconv.Itoa(L.len)
		panic("index out! max index: " + length)
	}
	for t.next != nil && num != 0 {
		num--
		t = t.next
	}
	t.object = v
	return true
}

// Contain
// @Description: check if the linked list contains the value
// @receiver     L    : the linked list
// @param        v : the value
// @return       bool : true if the value is contained
func (L *LinkList[T]) Contain(v T) bool {
	L.rwLock.RLock()
	defer L.rwLock.RUnlock()
	t := L.head
	for t.next != nil {
		if L.equal(t.object, v) {
			return true
		}
		t = t.next
	}
	return false
}

// IsEmpty
// @Description: check if the linked list is empty
// @receiver     L    : the linked list
// @return       bool : true if the linked list is empty
func (L *LinkList[T]) IsEmpty() bool {
	L.rwLock.RLock()
	defer L.rwLock.RUnlock()
	if L.head == L.tail {
		return true
	}
	return false
}

// Iterator
// @Description: iterator the linked list
// @receiver     L : the linked list
// @param        f : the function
func (L *LinkList[T]) Iterator(f func(index int, value T)) {
	L.iteRWLock.Lock()
	next := L.head
	index := 0
	for next.next != nil {
		f(index, next.object)
		index++
		next = next.next
	}
	L.iteRWLock.Unlock()
}
