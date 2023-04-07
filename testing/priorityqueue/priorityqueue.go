package priorityqueue

import node "sae-shortest-path/testing/Node"

type PriorityQueue interface {
	Pop() node.AStarNode
	Push(float64, node.AStarNode)
	Empty() bool
}
