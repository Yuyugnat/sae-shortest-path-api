package testing

import (
	"fmt"
	c "sae-shortest-path/connection"
	"sae-shortest-path/data"
	node "sae-shortest-path/testing/Node"
	prio "sae-shortest-path/testing/priorityqueue"
	hc "sae-shortest-path/testing/calculator"
	n "sae-shortest-path/testing/neighbors"
)

// type VoisinTuple struct {
// 	Gid      int
// 	Voisins []Voisin
// }

type Dijkstra struct {
	solver
	Dist   map[int]float64
	Reversed bool
	hCalc    hc.HeuristicCalculator
	queue    prio.PriorityQueue
}

func NewDijkstraNode(gid int, distance float64, heuristicDistance float64) *node.DijkstraNode {
	return &node.DijkstraNode{
		Node: node.Node{
			Gid:      gid,
			Distance: distance,
			Prev:     nil,
		},
	}
}

func NewDijkstra(nbGetter n.NeighborGetter, hCalc hc.HeuristicCalculator, queue prio.PriorityQueue) *Dijkstra {
	var res = &Dijkstra{}
	res.hCalc = hCalc
	res.queue = queue
	res.Instantiate(nbGetter)
	return res
}

func (d *Dijkstra) InitSpecificsVars() {
	d.Dist = make(map[int]float64)
}

func (d *Dijkstra) buildResult(node *node.Node) *Resultat {
	var path []Point
	d.Debug().GetTimeUsing("buildResult", func() {
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
	if d.Reversed {
		d.Debug().GetTimeUsing("reversePath", func() {
			path = d.reversePath(path)
		})
	}
	return &Resultat{
		Distance:       d.Dist[node.Gid],
		VilleDepart:    d.Depart,
		VilleArrivee:   d.Arrivee,
		PointsReversed: d.Reversed,
		Points:         path,
		ErrCode:        0,
		ErrMsg:         "",
	}
}

func (s *Dijkstra) reversePath(path []Point) []Point {
	var res []Point
	for i := len(path) - 1; i >= 0; i-- {
		res = append(res, path[i])
	}
	return res
}

func (d *Dijkstra) Solve(depart, arrivee string) *Resultat {
	// if !s.NbGetter.Ready() {
	// 	return &Resultat{
	// 		ErrCode: NotReady,
	// 		ErrMsg:  "Data server not ready",
	// 	}
	// }
	err := d.InitSearch(depart, arrivee)
	if err != nil {
		return &Resultat{
			ErrCode: 1,
			ErrMsg:  err.Error(),
		}
	}

	var dist float64

	ptEnd := d.GetPointFromGid(d.ArriveeGid)
	d.LastPoint = &ptEnd

	dist = d.nrRepo.GetDistance2(d.DepartGid, d.ArriveeGid)

	startNode := NewASrarNode(d.DepartGid, 0, dist)

	ptStart := d.GetPointFromGid(d.DepartGid)
	startNode.Lat = ptStart.Lat
	startNode.Lon = ptStart.Lon

	// s.Opened.Insert(startNode.HDistance, startNode.Node)
	d.queue.Push(startNode.HDistance, *startNode)

	d.Dist[d.DepartGid] = 0

	for !d.queue.Empty() {
		var current node.AStarNode

		d.Debug().GetTimeUsing("get lowest f", func() {
			current = d.queue.Pop()
		})

		if current.Gid == d.ArriveeGid {
			return d.buildResult(&current.Node)
		}

		d.Debug().GetTimeUsing("get voisins", func() {
			d.GetAdjacentNodes(&current)
		})
	}

	return &Resultat{
		ErrCode: 1,
		ErrMsg:  "No path found",
	}
}

func (d *Dijkstra) GetAdjacentNodes(current *node.AStarNode) {
	var neighbors []data.Neighbor
	d.Debug().GetTimeUsing("get voisins in map", func() {
		neighbors = d.NbGetter.Get(current.Gid)
	})

	d.Debug().GetTimeUsing("range on voisins", func() {
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
			d.Debug().GetTimeUsing("heuristic", func() {
				h = d.hCalc.Compute(&asn, current)
			})
			asn.HDistance = h
			potentialGScore := d.Dist[current.Gid] + asn.Distance
			val, exists := d.Dist[v.Gid]
			if !exists || potentialGScore < val {
				asn.Prev = &current.Node
				d.Dist[asn.Gid] = potentialGScore
				d.queue.Push(potentialGScore+asn.HDistance, asn)
			}
		}
	})
}

func (s *AStar) GetPointFromGid(gid int) Point {
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
