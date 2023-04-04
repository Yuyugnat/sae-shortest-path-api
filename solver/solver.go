package solver

type Solver interface {
	Solve() (float64, *Node)
}

type Node struct {}
	