package testing

import (
	"math"
	"sae-shortest-path/data"
	node "sae-shortest-path/testing/Node"
	hc "sae-shortest-path/testing/calculator"
	n "sae-shortest-path/testing/neighbors"
	prio "sae-shortest-path/testing/priorityqueue"
)

var (
	earthRadius = 6371.0
	degToRad    = math.Pi / 180.0
)

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// type VoisinTuple struct {
// 	Gid      int
// 	Voisins []Voisin
// }

type AStar struct {
	solver
	GScore   map[int]float64
	FScore   map[int]float64
	Reversed bool
	hCalc    hc.HeuristicCalculator
	queue    prio.PriorityQueue
}

func NewASrarNode(gid int, distance float64, heuristicDistance float64) *node.AStarNode {
	return &node.AStarNode{
		Node: node.Node{
			Gid:      gid,
			Distance: distance,
			Prev:     nil,
		},
		HDistance: heuristicDistance,
	}
}

func NewAStar(nbGetter n.NeighborGetter, hCalc hc.HeuristicCalculator, queue prio.PriorityQueue) *AStar {
	var res = &AStar{}
	res.hCalc = hCalc
	res.queue = queue
	res.Instantiate(nbGetter)
	return res
}

func (s *AStar) InitSpecificsVars() {
	s.GScore = make(map[int]float64)
	s.FScore = make(map[int]float64)
}

func (s *AStar) buildResult(node *node.Node) *Resultat {
	var path []Point
	s.Debug().GetTimeUsing("buildResult", func() {
		curr := node
		for curr.Prev != nil {
			path = append(path, Point{
				Lon: curr.Lon,
				Lat: curr.Lat,
			})
			curr = curr.Prev
		}
		path = append(path, Point{
			Lon: curr.Lon,
			Lat: curr.Lat,
		})
	})
	if s.Reversed {
		s.Debug().GetTimeUsing("reversePath", func() {
			path = s.reversePath(path)
		})
	}
	return &Resultat{
		Distance:       s.GScore[node.Gid],
		VilleDepart:    s.Depart,
		VilleArrivee:   s.Arrivee,
		PointsReversed: s.Reversed,
		Points:         path,
		ErrCode:        0,
		ErrMsg:         "",
	}
}

func (s *AStar) reversePath(path []Point) []Point {
	var res []Point
	for i := len(path) - 1; i >= 0; i-- {
		res = append(res, path[i])
	}
	return res
}

func (s *AStar) Solve(depart, arrivee string) *Resultat {
	// if !s.NbGetter.Ready() {
	// 	return &Resultat{
	// 		ErrCode: NotReady,
	// 		ErrMsg:  "Data server not ready",
	// 	}
	// }
	err := s.InitSearch(depart, arrivee)
	if err != nil {
		return &Resultat{
			ErrCode: 1,
			ErrMsg:  err.Error(),
		}
	}

	var dist float64

	ptEnd := s.GetPointFromGid(s.ArriveeGid)
	s.LastPoint = &ptEnd

	dist = s.nrRepo.GetDistance2(s.DepartGid, s.ArriveeGid)

	startNode := NewASrarNode(s.DepartGid, 0, dist)

	ptStart := s.GetPointFromGid(s.DepartGid)
	startNode.Lat = ptStart.Lat
	startNode.Lon = ptStart.Lon

	// s.Opened.Insert(startNode.HDistance, startNode.Node)
	s.queue.Push(startNode.HDistance, *startNode)

	s.GScore[s.DepartGid] = 0
	s.FScore[s.DepartGid] = startNode.HDistance

	for !s.queue.Empty() {
		var current node.AStarNode

		s.Debug().GetTimeUsing("get lowest f", func() {
			current = s.queue.Pop()
		})

		if current.Gid == s.ArriveeGid {
			return s.buildResult(&current.Node)
		}

		s.Debug().GetTimeUsing("get voisins", func() {
			s.GetAdjacentNodes(&current)
		})
	}

	return &Resultat{
		ErrCode: 1,
		ErrMsg:  "No path found",
	}
}

func (s *AStar) GetAdjacentNodes(current *node.AStarNode) {
	var neighbors []data.Neighbor
	s.Debug().GetTimeUsing("get voisins in map", func() {
		neighbors = s.NbGetter.Get(current.Gid)
	})

	s.Debug().GetTimeUsing("range on voisins", func() {
		for _, v := range neighbors {
			asn := node.AStarNode{
				Node: node.Node{
					Gid:      v.Gid,
					Lat:      v.Lat,
					Lon:      v.Lon,
					Distance: v.Length,
					Prev:     nil,
				},
				HDistance: 0,
			}
			var h float64
			s.Debug().GetTimeUsing("heuristic", func() {
				h = s.hCalc.Compute(&asn, &node.AStarNode{
					Node: node.Node{
						Lat: s.LastPoint.Lat,
						Lon: s.LastPoint.Lon,
					},
				})
			})
			asn.HDistance = h
			potentialGScore := s.GScore[current.Gid] + asn.Distance
			val, exists := s.GScore[v.Gid]
			if !exists || potentialGScore < val {
				asn.Prev = &current.Node
				s.GScore[asn.Gid] = potentialGScore
				s.FScore[asn.Gid] = potentialGScore + asn.HDistance
				// s.Opened.Insert(potentialGScore+asn.HDistance, asn.Node)
				s.queue.Push(potentialGScore+asn.HDistance, asn)
			}
		}
	})
}
