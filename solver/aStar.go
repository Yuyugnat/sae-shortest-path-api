package solver

import (
	"database/sql"
	"fmt"
	"math"
	c "sae-shortest-path/connection"
	o "sae-shortest-path/objects"
	"time"
)

var (
	noeudRoutierRepo = o.NewNoeudRoutierRepo()
	earthRadius      = 6371.0
)

type Point struct {
	Lon float64
	Lat float64
}

type ExecutionTime struct {
	Time float64
	Call int
}

type AStar struct {
	DepartGid   int
	ArriveeGid  int
	ArriveeGeom string
	LastPoint   Point
	Opened      map[int]AStarNode
	GScore      map[int]float64
	FScore      map[int]float64
	TimesSpent  map[string]ExecutionTime
}

type AStarNode struct {
	Node
	Gid       int
	Geom      string
	Distance  float64
	HDistance float64
	Prev      *AStarNode
}

func NewASrarNode(gid int, geom string, distance float64, heuristicDistance float64) AStarNode {
	return AStarNode{
		Gid:       gid,
		Geom:      geom,
		Distance:  distance,
		HDistance: heuristicDistance,
		Prev:      nil,
	}
}

func NewAStar(departGid int, arriveeGid int) *AStar {
	return &AStar{
		DepartGid:   departGid,
		ArriveeGid:  arriveeGid,
		ArriveeGeom: noeudRoutierRepo.GetGeomFromGid(arriveeGid),
		Opened:      make(map[int]AStarNode),
		GScore:      make(map[int]float64),
		FScore:      make(map[int]float64),
		TimesSpent:  make(map[string]ExecutionTime),
	}
}

// func (s *AStar) GetMinFrontiere() int {
// 	min := -1
// 	minDistance := math.MaxFloat64
// 	for gid := range s.Frontiere {
// 		if s.Distances[gid] < minDistance {
// 			minDistance = s.Distances[gid]
// 			min = gid
// 		}
// 	}
// 	return min
// }

func (s *AStar) Solve() (float64, map[string]ExecutionTime) {
	var dist float64

	s.LastPoint = s.GetPointFromGid(s.DepartGid)

	s.GetTimeUsing("heuristic", func() {
		dist = noeudRoutierRepo.GetDistance(noeudRoutierRepo.GetGeomFromGid(s.DepartGid), s.ArriveeGeom)
	})

	s.Opened[s.DepartGid] = NewASrarNode(s.DepartGid, noeudRoutierRepo.GetGeomFromGid(s.DepartGid), 0, dist)
	s.GScore[s.DepartGid] = 0
	s.FScore[s.DepartGid] = s.Opened[s.DepartGid].HDistance

	for len(s.Opened) > 0 {
		current := s.GetLowestF()

		if current.Gid == s.ArriveeGid {
			return s.GScore[current.Gid], s.TimesSpent
		}

		delete(s.Opened, current.Gid)

		var voisins []AStarNode

		s.GetTimeUsing("get voisins", func() {
			voisins = s.GetAdjacentNodes(current)
		})

		s.GetTimeUsing("boucle sur voisins", func() {
			for _, neighbor := range voisins {
				potentialGScore := s.GScore[current.Gid] + neighbor.Distance
				if _, exists := s.GScore[neighbor.Gid]; !exists || potentialGScore < s.GScore[neighbor.Gid] {
					s.GScore[neighbor.Gid] = potentialGScore
					s.FScore[neighbor.Gid] = potentialGScore + neighbor.HDistance
					neighbor.Prev = &current
					s.Opened[neighbor.Gid] = neighbor
				}
			}
		})
	}

	return -1, s.TimesSpent
}

func (s *AStar) GetAdjacentNodes(a AStarNode) []AStarNode {
	query := `
		(
			SELECT noeud_voisin as noeud_routier_gid_1, noeud_routier as noeud_routier_gid_2, troncon_id as troncon_gid, longueur, st_astext(troncon_geom) as geom
			FROM voisins_noeud
			WHERE noeud_routier = $1 or noeud_voisin = $1
		);
	`
	var rows *sql.Rows
	var err error

	s.GetTimeUsing("query voisins_noeud", func() {
		rows, err = c.Conn.DB.Query(query, a.Gid)
	})

	if err != nil {
		fmt.Println("Error while querying the database : ", err)
		return nil
	}
	defer rows.Close()

	res := make([]AStarNode, 0)

	for rows.Next() {
		var nrGid1 int
		var nrGid2 int
		var tronconGid int
		var longueur float64
		var geom string
		err = rows.Scan(&nrGid1, &nrGid2, &tronconGid, &longueur, &geom)
		if err != nil {
			fmt.Println("Error while scanning the database : ", err)
			return nil
		}
		var id int
		if nrGid1 == a.Gid {
			id = nrGid2
		} else {
			id = nrGid1
		}

		// fmt.Println("id : ", id)

		var h float64
		s.GetTimeUsing("heuristic", func() {
			h = noeudRoutierRepo.GetDistance2(id, s.ArriveeGid)
		})

		res = append(res, AStarNode{
			Gid:       id,
			Geom:      geom,
			Distance:  longueur,
			HDistance: h,
			Prev:      nil,
		})
	}
	return res
}

func (s *AStar) GetAdjacentNodes2(a AStarNode) []AStarNode {
	query := `
		SELECT noeud_voisin, noeud_routier, longueur, x_voisin, y_voisin, x_routier, y_routier
		FROM voisins_noeud_ultra
		WHERE noeud_routier = $1 or noeud_voisin = $1
	`
	var rows *sql.Rows
	var err error

	s.GetTimeUsing("query voisin_noeud_ultra", func() {
		rows, err = c.Conn.DB.Query(query, a.Gid)
	})

	if err != nil {
		fmt.Println("Error while querying the database : ", err)
		return nil
	}
	defer rows.Close()

	res := make([]AStarNode, 0)

	for rows.Next() {
		var nrGid int
		var nvGid int
		var longueur float64
		var xv float64
		var yv float64
		var xr float64
		var yr float64
		err = rows.Scan(&nvGid, &nrGid, &longueur, &xv, &yv, &xr, &yr)
		if err != nil {
			fmt.Println("Error while scanning the database : ", err)
			return nil
		}
		var h float64
		var id int
		var p Point

		if nrGid == a.Gid {
			id = nvGid
			p = Point{Lon: yv, Lat: xv}
		} else {
			id = nrGid
			p = Point{Lon: yr, Lat: xr}
		}

		s.GetTimeUsing("heuristic", func() {
			h = s.GetDistanceUltra(p, s.LastPoint)
		})

		res = append(res, AStarNode{
			Gid:       id,
			Geom:      "",
			Distance:  longueur,
			HDistance: h,
			Prev:      nil,
		})
	}
	return res
}

func (s *AStar) GetDistanceUltra(p1, p2 Point) float64 {
	phi1 := p1.Lat * math.Pi / 180
	phi2 := p2.Lat * math.Pi / 180
	deltaPhi := (p2.Lat - p1.Lat) * math.Pi / 180
	deltaLambda := (p2.Lon - p1.Lon) * math.Pi / 180

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) + math.Cos(phi1)*math.Cos(phi2)*math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := 6371.0 * c

	// fmt.Println("Distance : ", distance)
	return distance
}

func (s *AStar) FromDegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

func (s *AStar) GetPointFromGid(gid int) Point {
	query := `
		SELECT st_x(geom) as x, st_y(geom) as y
		FROM noeud_routier
		WHERE gid = $1
	`
	var x float64
	var y float64
	err := c.Conn.DB.QueryRow(query, gid).Scan(&x, &y)
	if err != nil {
		fmt.Println("Error while querying the database : ", err)
		return Point{}
	}
	return Point{Lon: x, Lat: y}
}

func (s *AStar) GetLowestF() AStarNode {
	f := math.MaxFloat64
	var res AStarNode
	for _, node := range s.Opened {
		if s.FScore[node.Gid] < f {
			f = s.FScore[node.Gid]
			res = node
		}
	}
	return res
}

func (s *AStar) GetTimeUsing(funcname string, f func()) {
	start := time.Now()
	f()
	elapsed := time.Since(start)
	if _, exists := s.TimesSpent[funcname]; !exists {
		s.TimesSpent[funcname] = ExecutionTime{
			Call: 1,
			Time: elapsed.Seconds(),
		}
	} else {
		s.TimesSpent[funcname] = ExecutionTime{
			Call: s.TimesSpent[funcname].Call + 1,
			Time: s.TimesSpent[funcname].Time + elapsed.Seconds(),
		}
	}
}
