package solver

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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

type Voisin struct {
	Gid      int     `json:"gid"`
	Longueur float64 `json:"length"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
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

	s.LastPoint = s.GetPointFromGid(s.ArriveeGid)

	s.GetTimeUsing("heuristic", func() {
		dist = noeudRoutierRepo.GetDistance(noeudRoutierRepo.GetGeomFromGid(s.DepartGid), s.ArriveeGeom)
	})

	s.Opened[s.DepartGid] = NewASrarNode(s.DepartGid, noeudRoutierRepo.GetGeomFromGid(s.DepartGid), 0, dist)
	s.GScore[s.DepartGid] = 0
	s.FScore[s.DepartGid] = s.Opened[s.DepartGid].HDistance

	for len(s.Opened) > 0 {
		var current AStarNode

		s.GetTimeUsing("get lowest f", func() {
			current = s.GetLowestF()
		})

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
		SELECT noeud_routier, noeud_voisin, longueur, st_astext(troncon_geom), nr_lat, nr_lon, nv_lat, nv_lon
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

	var nrGid int
	var nvGid int
	var longueur float64
	var geom string
	var rLat float64
	var rLon float64
	var vLat float64
	var vLon float64

	for rows.Next() {

		err = rows.Scan(&nrGid, &nvGid, &longueur, &geom, &rLat, &rLon, &vLat, &vLon)
		if err != nil {
			fmt.Println("Error while scanning the database : ", err)
		}

		var h float64
		var id int
		var p Point

		if nrGid == a.Gid {
			id = nvGid
			p = Point{Lat: vLat, Lon: vLon}
		} else {
			id = nrGid
			p = Point{Lat: rLat, Lon: rLon}
		}

		_ = p

		s.GetTimeUsing("heuristic", func() {
			h = s.GetDistanceUltra(p, s.LastPoint)
			// h = noeudRoutierRepo.GetDistance(geom, s.ArriveeGeom)
			// h = noeudRoutierRepo.GetDistance2(id, s.ArriveeGid)

			// fmt.Println("h : ", h)
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

func (s *AStar) GetAdjacentNodes(a AStarNode) []AStarNode {
	var voisins []Voisin
	var js string

	query := `
		SELECT voisins
		FROM voisins_jsonb
		WHERE gid = $1 
	`

	var row *sql.Row
	s.GetTimeUsing("query voisin_jsonb", func() {
		row = c.Conn.DB.QueryRow(query, a.Gid)
	})

	var err error
	s.GetTimeUsing("scan voisin_jsonb", func() {
		err = row.Scan(&js)
	})

	if err != nil {
		log.Fatalln(err)
	}

	s.GetTimeUsing("unmarshal voisin_jsonb", func() {
		err = json.Unmarshal([]byte(js), &voisins)
	})

	if err != nil {
		log.Fatalln(err)
	}

	res := make([]AStarNode, 0)

	for _, v := range voisins {
		h := s.GetDistanceUltra(Point{
			Lat: v.Lat,
			Lon: v.Lon,
		}, s.LastPoint)
		res = append(res, AStarNode{
			Gid:       v.Gid,
			Geom:      "",
			Distance:  v.Longueur,
			HDistance: h,
			Prev:      nil,
		})
	}

	return res
}

func (s *AStar) GetDistanceUltra(p1, p2 Point) float64 {
	phi1 := p1.Lat * math.Pi / 180.0
	phi2 := p2.Lat * math.Pi / 180.0
	deltaPhi := (p2.Lat - p1.Lat) * math.Pi / 180.0
	deltaLambda := (p2.Lon - p1.Lon) * math.Pi / 180.0

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) + math.Cos(phi1)*math.Cos(phi2)*math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2.0 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := 6371.0 * c

	// fmt.Println("Distance : ", distance)
	return distance
}

func (s *AStar) GetPointFromGid(gid int) Point {
	query := `
		SELECT lat, lon
		FROM geom_noeud_routier_xy
		WHERE gid = $1
	`
	var lat float64
	var lon float64
	err := c.Conn.DB.QueryRow(query, gid).Scan(&lat, &lon)
	if err != nil {
		fmt.Println("Error while querying the database : ", err)
		return Point{}
	}
	return Point{Lat: lat, Lon: lon}
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
