package testing

import (
	"fmt"
	o "sae-shortest-path/objects"
	n "sae-shortest-path/testing/neighbors"
)

type Dijkstra struct {
	solver
	DepartGid   int
	ArriveeGid  int
	ArriveeGeom string
	Depart      string
	Arrivee     string
	LastPoint   *Point
	Distances   map[int]float64
	Frontier    map[int]*DijkstraNode
}

type DijkstraNode struct {
	Gid  int
	Geom string
	Lat  float64
	Lon  float64
	Prev *DijkstraNode
}

func NewDijkstraNode(gid int, geom string, distance float64) *DijkstraNode {
	return &DijkstraNode{
		Gid:  gid,
		Geom: geom,
		Prev: nil,
	}
}

func NewDijkstra(nbGetter n.NeighborGetter) *Dijkstra {
	var res = &Dijkstra{}
	res.Instantiate(nbGetter)
	return res
}

func (s *AStar) InitSearch(depart, arrivee string) error {
	communeRepo := o.NewNoeudCommuneRepo()
	noeudRoutierRepo := o.NewNoeudRoutierRepo()

	departIdNdRte, err := communeRepo.GetIdNdRteByName(depart)
	if err != nil {
		return fmt.Errorf("depart '%s' not found", depart)
	}
	arriveeIdNdRte, err := communeRepo.GetIdNdRteByName(arrivee)
	if err != nil {
		return fmt.Errorf("arrivee '%s' not found", arrivee)
	}

	departGid := noeudRoutierRepo.GetGidByIdRte500(departIdNdRte)
	arriveeGid := noeudRoutierRepo.GetGidByIdRte500(arriveeIdNdRte)
	reversed := false
	if communeRepo.GetSuperficie(depart) > communeRepo.GetSuperficie(arrivee) {
		depart, arrivee = arrivee, depart
		departGid, arriveeGid = arriveeGid, departGid
		reversed = true
	}

	s.DepartGid = departGid
	s.ArriveeGid = arriveeGid
	s.Depart = depart
	s.Arrivee = arrivee
	s.Reversed = reversed

	s.InitSpecificsVars()
	return nil
}

func (d *Dijkstra) Solve(start, end string) *Resultat {
	return nil
}
