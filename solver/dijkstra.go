package solver

import (
	"math"

	o "sae-shortest-path/objects"
)

type Dijkstra struct {
	DepartGid  int
	ArriveeGid int
	Distances  map[int]float64
	Frontiere  map[int]bool
}

func NewDijkstra(departGid int, arriveeGid int) *Dijkstra {
	return &Dijkstra{
		DepartGid:  departGid,
		ArriveeGid: arriveeGid,
		Distances:  make(map[int]float64),
		Frontiere:  make(map[int]bool),
	}
}

func (s *Dijkstra) GetMinFrontiere() int {
	min := -1
	minDistance := math.MaxFloat64
	for gid := range s.Frontiere {
		if s.Distances[gid] < minDistance {
			minDistance = s.Distances[gid]
			min = gid
		}
	}
	return min
}

func (s *Dijkstra) Solve() (float64, *Node) {
	s.Distances[s.DepartGid] = 0
	s.Frontiere[s.DepartGid] = true

	for len(s.Frontiere) > 0 {
		currentGid := s.GetMinFrontiere()
		if currentGid == s.ArriveeGid {
			return s.Distances[currentGid], nil
		}

		delete(s.Frontiere, currentGid)

		currentNode := o.NewNoeudRoutierRepo().GetByPrimaryKey(currentGid)
		voisins := currentNode.GetVoisins()
		for _, voisin := range voisins {
			potentialDistance := s.Distances[currentGid] + voisin.GetLongueur()

			if _, exist := s.Distances[voisin.GetNrGid()]; !exist || potentialDistance < s.Distances[voisin.GetNrGid()] {
				s.Distances[voisin.GetNrGid()] = potentialDistance
				s.Frontiere[voisin.GetNrGid()] = true
			}
		}
	}
	return -1, nil
}
