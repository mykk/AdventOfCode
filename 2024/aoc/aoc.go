package aoc

import (
	"container/heap"
)

type Heap[T any] struct {
	data []T
	less func(a, b T) bool
}

func (h *Heap[T]) Len() int {
	return len(h.data)
}

func (h *Heap[T]) Less(i, j int) bool {
	return h.less(h.data[i], h.data[j])
}

func (h *Heap[T]) Swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

func (h *Heap[T]) Push(x any) {
	h.data = append(h.data, x.(T))
}

func (h *Heap[T]) Pop() any {
	n := len(h.data)
	item := h.data[n-1]
	h.data = h.data[:n-1]
	return item
}

func (h *Heap[T]) PopItem() T {
	return heap.Pop(h).(T)
}

func (h *Heap[T]) PushItem(item T) {
	heap.Push(h, item)
}

func NewHeap[T any](less func(a, b T) bool) *Heap[T] {
	newHeap := &Heap[T]{
		data: []T{},
		less: less,
	}

	heap.Init(newHeap)
	return newHeap
}

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(value T) {
	s[value] = struct{}{}
}

func (s Set[T]) Contains(value T) bool {
	_, exists := s[value]
	return exists
}

func (s Set[T]) Clone() Set[T] {
	clone := make(Set[T])
	for k, v := range s {
		clone[k] = v
	}
	return clone
}

func NewSet[T comparable](args ...T) Set[T] {
	s := make(Set[T])
	for _, arg := range args {
		s.Add(arg)
	}
	return s
}
