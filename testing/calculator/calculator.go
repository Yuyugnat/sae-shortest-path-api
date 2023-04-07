package calculator

import (
	node "sae-shortest-path/testing/Node"
)

type HeuristicCalculator interface {
	Compute(gid1, gid2 *node.AStarNode) float64
}
