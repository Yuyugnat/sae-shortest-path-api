package priorityqueue

import (
	"math"
	node "sae-shortest-path/testing/node"
)

type PrioMap struct {
	Opened map[int]node.Node
	FValue map[int]float64
}

func NewPrioMap() *PrioMap {
	return &PrioMap{
		Opened: make(map[int]node.Node),
		FValue: make(map[int]float64),
	}
}

func (p *PrioMap) Pop() node.Node {
	min := math.MaxFloat64
	var minNode node.Node
	for _, n := range p.Opened {
		if p.FValue[n.Gid] < min {
			min = p.FValue[n.Gid]
			minNode = n
		}
	}
	delete(p.Opened, minNode.Gid)
	return minNode
}

func (p *PrioMap) Push(fval float64, n node.Node) {
	p.Opened[n.Gid] = n
	p.FValue[n.Gid] = fval
}

func (p *PrioMap) Empty() bool {
	return len(p.Opened) == 0
}