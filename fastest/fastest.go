package fastest

import (
	"fmt"
	"math"
	c "sae-shortest-path/connection"
	"sae-shortest-path/data"
	bug "sae-shortest-path/debugging"
	o "sae-shortest-path/objects"
	// nb "sae-shortest-path/solver/neighbors"
)

var degToRad    = math.Pi / 180.0

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type AStar struct {
	debugger *bug.Debug
	nrRepo   *o.NoeudRoutierRepo
	NbGetter *data.NeighborTable

	DepartGid   int
	ArriveeGid  int
	Depart      string
	Arrivee     string
	Reversed    bool
	LastPoint   *Point
	Opened      MinHeap
	GScore      map[int]float64
}

type AStarNode struct {
	Gid       int
	Lat       float64
	Lon       float64
	Length  float64
	Prev      *AStarNode
	HDistance float64
}

// func NewASrarNode(gid int, distance float64, heuristicDistance float64) *AStarNode {
// 	return &AStarNode{
// 		Gid:       gid,
// 		Distance:  distance,
// 		Prev:      nil,
// 		HDistance: heuristicDistance,
// 	}
// }

func NewFastest(depart, arrivee string) *AStar {
	communeRepo := o.NewNoeudCommuneRepo()
	noeudRoutierRepo := o.NewNoeudRoutierRepo()

	departIdNdRte, _ := communeRepo.GetIdNdRteByName(depart)
	arriveeIdNdRte, _ := communeRepo.GetIdNdRteByName(arrivee)

	departGid := noeudRoutierRepo.GetGidByIdRte500(departIdNdRte)
	arriveeGid := noeudRoutierRepo.GetGidByIdRte500(arriveeIdNdRte)
	reversed := false
	if communeRepo.GetSuperficie(depart) > communeRepo.GetSuperficie(arrivee) {
		depart, arrivee = arrivee, depart
		departGid, arriveeGid = arriveeGid, departGid
		reversed = true
	}
	var res = &AStar{
		nrRepo:     noeudRoutierRepo,
		// ncRepo:     o.NewNoeudCommuneRepo(),
		Depart:     depart,
		Arrivee:    arrivee,
		DepartGid:  departGid,
		ArriveeGid: arriveeGid,
		Reversed:   reversed,
		debugger:   bug.NewDebug(),
		NbGetter:   data.GetInstance(),
		Opened:     MinHeap{},
		GScore:     make(map[int]float64),
	}

	return res
}

func (s *AStar) Solve() *Result {
	if !s.NbGetter.Ready() {
		return &Result{
			ErrCode: NotReady,
			ErrMsg:  "Data server not ready",
		}
	}

	pt := s.GetPointFromGid(s.ArriveeGid)
	s.LastPoint = &pt

	pt = s.GetPointFromGid(s.DepartGid)
	dist := haversine(&pt, s.LastPoint)

	startNode := &AStarNode{
		Gid:       s.DepartGid,
		Length:  0,
		Prev:      nil,
		HDistance: dist,
	}

	ptStart := s.GetPointFromGid(s.DepartGid)
	startNode.Lat = ptStart.Lat
	startNode.Lon = ptStart.Lon

	s.Opened.Insert(startNode.HDistance, *startNode)

	for !s.Opened.IsEmpty() {
		current := s.Opened.ExtractMin()
		if current.Gid == s.ArriveeGid {
			return s.buildResult(&current)
		}
		s.GetAdjacentNodes(&current)
	}

	return &Result{
		ErrCode: 1,
		ErrMsg:  "No path found",
	}
}

func (s *AStar) GetAdjacentNodes(current *AStarNode) {
	neighbors := data.GetInstance().Table[current.Gid]

	for _, v := range neighbors {
		h := haversine(&Point{
			Lat: v.Lat,
			Lon: v.Lon,
		}, s.LastPoint)
		// })
		asn := AStarNode{
			Gid:       v.Gid,
			Lat:       v.Lat,
			Lon:       v.Lon,
			Length:  v.Length,
			Prev:      nil,
			HDistance: h,
		}

		potentialGScore := s.GScore[current.Gid] + asn.Length
		val, exists := s.GScore[v.Gid]
		if !exists || potentialGScore < val {
			asn.Prev = current
			s.GScore[asn.Gid] = potentialGScore
			s.Opened.Insert(potentialGScore+asn.HDistance, asn)
		}
	}
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

func haversine(p1, p2 *Point) float64 {
	sinDeltaPhi := math.Sin(((p2.Lat - p1.Lat) * degToRad) / 2)
	sinDeltaLambda := math.Sin(((p2.Lon - p1.Lon) * degToRad) / 2)
	a := sinDeltaPhi*sinDeltaPhi + math.Cos(p1.Lat * degToRad)*math.Cos(p2.Lat * degToRad)*sinDeltaLambda*sinDeltaLambda
	c := 2.0 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := 6371.0 * c

	return distance
}
