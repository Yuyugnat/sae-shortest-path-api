package structure

import node "sae-shortest-path/testing/Node"

type Heapnode struct {
	Key float64
	Val node.AStarNode
}

type MinHeap struct {
	Heap []Heapnode
}

func NewMinHeap() *MinHeap {
	return &MinHeap{}
}

func (h *MinHeap) Insert(key float64, val node.AStarNode) {
	h.Heap = append(h.Heap, Heapnode{key, val})
	h.bubbleUp(len(h.Heap) - 1)
}

func (h *MinHeap) bubbleUp(index int) {
	if index == 0 {
		return
	}
	parent := (index - 1) / 2
	if h.Heap[parent].Key > h.Heap[index].Key {
		h.Heap[parent], h.Heap[index] = h.Heap[index], h.Heap[parent]
		h.bubbleUp(parent)
	}
}

func (h *MinHeap) ExtractMin() node.AStarNode {
	if len(h.Heap) == 0 {
		return *new(node.AStarNode)
	}
	min := h.Heap[0]
	h.Heap[0] = h.Heap[len(h.Heap)-1]
	h.Heap = h.Heap[:len(h.Heap)-1]
	h.bubbleDown(0)
	return min.Val
}

func (h *MinHeap) bubbleDown(index int) {
	left := index*2 + 1
	right := index*2 + 2
	if left >= len(h.Heap) {
		return
	}
	min := left
	if right < len(h.Heap) && h.Heap[right].Key < h.Heap[left].Key {
		min = right
	}
	if h.Heap[index].Key > h.Heap[min].Key {
		h.Heap[index], h.Heap[min] = h.Heap[min], h.Heap[index]
		h.bubbleDown(min)
	}
}

func (h *MinHeap) Len() int {
	return len(h.Heap)
}

func (h *MinHeap) IsEmpty() bool {
	return len(h.Heap) == 0
}
