package testing

import (
	"fmt"
	c "sae-shortest-path/connection"
	"sae-shortest-path/data"
	bug "sae-shortest-path/debugging"
	o "sae-shortest-path/objects"
	hc "sae-shortest-path/testing/calculator"
	nb "sae-shortest-path/testing/neighbors"
	"sae-shortest-path/testing/node"
	prio "sae-shortest-path/testing/priorityqueue"
)

type ErrCode int

const (
	NoErr ErrCode = iota
	NoDepartOrArrivee
	NoPath
	NotReady
)

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}


type ISolver interface {
	Solve(start, end string) *Resultat
	Debug() *bug.Debug
}

type Solver struct {
	ISolver
	debugger *bug.Debug
	nrRepo   *o.NoeudRoutierRepo
	ncRepo   *o.NoeudCommuneRepo

	NbGetter nb.NeighborGetter

	DepartGid   int
	ArriveeGid  int
	ArriveeGeom string
	Depart      string
	Arrivee     string
	LastPoint   *Point
	GScore      map[int]float64
	FScore      map[int]float64
	Reversed    bool
	hCalc       hc.HeuristicCalculator
	queue       prio.PriorityQueue
}

func (s *Solver) Instantiate(nbGetter nb.NeighborGetter) {
	repoNR := o.NewNoeudRoutierRepo()
	repoNC := o.NewNoeudCommuneRepo()
	s.Reversed = false
	s.nrRepo = repoNR
	s.ncRepo = repoNC
	s.NbGetter = nbGetter
	s.debugger = bug.NewDebug()
	s.GScore = make(map[int]float64)
	s.FScore = make(map[int]float64)
}

func (s *Solver) InitSearch(depart, arrivee string) error {
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
	return nil
}

func (s *Solver) GetPointFromGid(gid int) Point {
	query := `
		SELECT lat, lon
		FROM geom_noeud_routier_xy
		WHERE gid = $1
	`
	var lat float64
	var lon float64
	conn, err := c.GetInstance()
	if err != nil {
		fmt.Println("Error while getting the database connection : ", err)
		return Point{}
	}
	err = conn.DB.QueryRow(query, gid).Scan(&lat, &lon)
	if err != nil {
		fmt.Println("Error while querying the database (GetPointFromGid) : ", err)
		return Point{}
	}
	return Point{Lat: lat, Lon: lon}
}

func (s *Solver) Debug() *bug.Debug {
	return s.debugger
}

func NewNode(gid int, distance float64, heuristicDistance float64) *node.Node {
	return &node.Node{
		Gid:      gid,
		Distance: distance,
		Prev:     nil,
		HDistance: heuristicDistance,
	}
}

func NewSolver(nbGetter nb.NeighborGetter, hCalc hc.HeuristicCalculator, queue prio.PriorityQueue) *Solver {
	var res = &Solver{}
	res.hCalc = hCalc
	res.queue = queue
	res.NbGetter = nbGetter
	res.Instantiate(nbGetter)
	return res
}

func (s *Solver) buildResult(node *node.Node) *Resultat {
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

func (s *Solver) reversePath(path []Point) []Point {
	var res []Point
	for i := len(path) - 1; i >= 0; i-- {
		res = append(res, path[i])
	}
	return res
}

func (s *Solver) Solve(depart, arrivee string) *Resultat {
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

	startNode := NewNode(s.DepartGid, 0, dist)

	ptStart := s.GetPointFromGid(s.DepartGid)
	startNode.Lat = ptStart.Lat
	startNode.Lon = ptStart.Lon

	// s.Opened.Insert(startNode.HDistance, startNode.Node)
	s.queue.Push(startNode.HDistance, *startNode)

	s.GScore[s.DepartGid] = 0
	s.FScore[s.DepartGid] = startNode.HDistance

	for !s.queue.Empty() {
		var current node.Node

		s.Debug().GetTimeUsing("get lowest f", func() {
			current = s.queue.Pop()
		})

		if current.Gid == s.ArriveeGid {
			return s.buildResult(&current)
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

func (s *Solver) GetAdjacentNodes(current *node.Node) {
	var neighbors []data.Neighbor
	s.Debug().GetTimeUsing("get voisins by gid", func() {
		neighbors = s.NbGetter.Get(current.Gid)
	})

	s.Debug().GetTimeUsing("range on voisins", func() {
		for _, v := range neighbors {
			asn := node.Node{
				Gid:      v.Gid,
				Lat:      v.Lat,
				Lon:      v.Lon,
				Distance: v.Length,
				Prev:     nil,
				HDistance: 0,
			}
			var h float64
			s.Debug().GetTimeUsing("heuristic", func() {
				h = s.hCalc.Compute(&asn, &node.Node{
						Lat: s.LastPoint.Lat,
						Lon: s.LastPoint.Lon,
				})
			})
			asn.HDistance = h
			potentialGScore := s.GScore[current.Gid] + asn.Distance
			val, exists := s.GScore[v.Gid]
			if !exists || potentialGScore < val {
				asn.Prev = current
				s.GScore[asn.Gid] = potentialGScore
				s.FScore[asn.Gid] = potentialGScore + asn.HDistance
				// s.Opened.Insert(potentialGScore+asn.HDistance, asn.Node)
				s.queue.Push(potentialGScore+asn.HDistance, asn)
			}
		}
	})
}
