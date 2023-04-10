package priorityqueue

import node "sae-shortest-path/testing/node"

type PriorityQueue interface {
	Pop() node.Node
	Push(float64, node.Node)
	Empty() bool
}
