package calculator

import node "sae-shortest-path/testing/node"

type dijkstraHeuristic struct{}

func NewDijkstraCalculator() *dijkstraHeuristic {
	return &dijkstraHeuristic{}
}

func (d *dijkstraHeuristic) Compute(gid1, gid2 *node.Node) float64 {
	return 0
}
