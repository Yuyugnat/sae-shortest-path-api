package node

type Node struct {
	Gid       int
	Lat       float64
	Lon       float64
	Distance  float64
	HDistance float64
	Prev      *Node
}