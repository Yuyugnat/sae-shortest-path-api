package neighbors

import data "sae-shortest-path/data"

type NeighborGetter interface {
	Get(gid int) []data.Neighbor
}
