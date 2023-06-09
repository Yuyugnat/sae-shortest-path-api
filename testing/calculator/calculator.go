package calculator

import (
	node "sae-shortest-path/testing/node"
)

type HeuristicCalculator interface {
	Compute(gid1, gid2 *node.Node) float64
}
