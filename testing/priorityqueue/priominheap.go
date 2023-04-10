package priorityqueue

import (
	"sae-shortest-path/structure"
	node "sae-shortest-path/testing/node"
)

type PrioMinHeap struct {
	heap *structure.MinHeap
}

func NewPrioMinHeap() *PrioMinHeap {
	return &PrioMinHeap{
		heap: structure.NewMinHeap(),
	}
}

func (p *PrioMinHeap) Pop() node.Node {
	return p.heap.ExtractMin()
}

func (p *PrioMinHeap) Push(key float64, val node.Node) {
	p.heap.Insert(key, val)
}

func (p *PrioMinHeap) Empty() bool {
	return p.heap.IsEmpty()
}
