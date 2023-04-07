package fastest

import "encoding/json"

type ErrCode int

const (
	NoErr ErrCode = iota
	NoDepartOrArrivee
	NoPath
	NotReady
)

type Result struct {
	Distance       float64 `json:"distance"`
	VilleDepart    string  `json:"villeDepart"`
	VilleArrivee   string  `json:"villeArrivee"`
	PointsReversed bool    `json:"pointsReversed"`
	Points         []Point `json:"points"`
	ErrCode        ErrCode `json:"errCode"` // 0 = no error, other = error
	ErrMsg         string  `json:"errMsg"`
}

func (r *Result) JSON() []byte {
	res, _ := json.Marshal(r)
	return res
}

func reversePath(path []Point) []Point {
	var res []Point
	for i := len(path) - 1; i >= 0; i-- {
		res = append(res, path[i])
	}
	return res
}

func (s *AStar) buildResult(node *AStarNode) *Result {
	var path []Point
	// s.Debug().GetTimeUsing("buildResult", func() {
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
	// })
	if s.Reversed {
		path = reversePath(path)
	}
	return &Result{
		Distance:       s.GScore[node.Gid],
		VilleDepart:    s.Depart,
		VilleArrivee:   s.Arrivee,
		PointsReversed: s.Reversed,
		Points:         path,
		ErrCode:        0,
		ErrMsg:         "",
	}
}