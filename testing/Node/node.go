package node

type Node struct {
	Gid       int
	Lat       float64
	Lon       float64
	Distance  float64
	Prev      *Node
}

type AStarNode struct {
	Node
	HDistance float64
}