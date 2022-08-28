package util

import (
	"reflect"
	"strconv"
	"sync"
)

type LinkList[T any] struct {
	head      *linkListNode[T]
	tail      *linkListNode[T]
	len       int
	rwLock    *sync.RWMutex
	iteRWLock *sync.Mutex
}

type linkListNode[T any] struct {
	object T
	next   *linkListNode[T]
	prev   *linkListNode[T]
}

func NewLinkList[T any]() *LinkList[T] {
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
	}
}

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

func (L LinkList[T]) ToString() string {
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

func (L LinkList[T]) ToStringBack() string {
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

func (L LinkList[T]) Length() int {
	L.rwLock.RLock()
	defer L.rwLock.RUnlock()
	return L.len
}

func (L LinkList[T]) Get(num int) T {
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

func (L *LinkList[T]) DeleteByValue(o T) bool {
	L.rwLock.Lock()
	defer L.rwLock.Unlock()
	if L.len == 0 {
		return false
	}
	t := L.head
	if reflect.DeepEqual(t.object, o) {
		L.head.next.prev = nil
		L.head = L.head.next
		L.len--
		return true
	}
	for t.next.next != nil {
		if reflect.DeepEqual(t.object, o) {
			t.prev.next = t.next
			t.next.prev = t.prev
			L.len--
			return true
		}
		t = t.next
	}
	if reflect.DeepEqual(t.object, o) {
		t.prev.next = L.tail
		L.tail.prev = t.prev
		L.len--
		return true
	}
	return false
}

func (L *LinkList[T]) DeleteAllByValue(o T) bool {
	L.rwLock.Lock()
	defer L.rwLock.Unlock()
	if L.len == 0 {
		return false
	}
L1:
	t := L.head
	if reflect.DeepEqual(t.object, o) {
		L.head.next.prev = nil
		L.head = L.head.next
		goto L1
	}
	for t != nil && t.next != nil && t.next.next != nil {
		if reflect.DeepEqual(t.object, o) {
			t.prev.next = t.next
			t.next.prev = t.prev
			L.len--
		}
		t = t.next
	}
	if reflect.DeepEqual(t.object, o) {
		t.prev.next = L.tail
		L.tail.prev = t.prev
		L.len--
	}
	return true
}

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

func (L LinkList[T]) Contain(v T) bool {
	L.rwLock.RLock()
	defer L.rwLock.RUnlock()
	t := L.head
	for t.next != nil {
		if reflect.DeepEqual(t.object, v) {
			return true
		}
		t = t.next
	}
	return false
}

func (L LinkList[T]) IsEmpty() bool {
	L.rwLock.RLock()
	defer L.rwLock.RUnlock()
	if L.head == L.tail {
		return true
	}
	return false
}

func (L LinkList[T]) Iterator(f func(index int, value T)) {
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
