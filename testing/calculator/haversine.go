package calculator

import (
	"math"
	node "sae-shortest-path/testing/Node"
)

const (
	degToRad = math.Pi / 180.0
)

type HaversineCalculator struct{}

func NewHaversineCalculator() *HaversineCalculator {
	return &HaversineCalculator{}
}

func (h *HaversineCalculator) Compute(gid1, gid2 *node.AStarNode) float64 {
	return haversine(gid1, gid2)
}

func haversine(gid1, gid2 *node.AStarNode) float64 {
	phi1 := gid1.Lat * degToRad
	phi2 := gid2.Lat * degToRad
	deltaPhi := (gid2.Lat - gid1.Lat) * degToRad
	deltaLambda := (gid2.Lon - gid1.Lon) * degToRad

	sinDeltaPhi := math.Sin(deltaPhi / 2)
	sinDeltaLambda := math.Sin(deltaLambda / 2)

	a := sinDeltaPhi*sinDeltaPhi + math.Cos(phi1)*math.Cos(phi2)*sinDeltaLambda*sinDeltaLambda
	c := 2.0 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := 6371.0 * c

	// fmt.Println("Distance : ", distance)
	return distance
}
