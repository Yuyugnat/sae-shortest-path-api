package neighbors

import (
	data "sae-shortest-path/data"
	bug "sae-shortest-path/debugging"
)

type FullyLoaded struct {
	data.NeighborTable
}

func Load() {
	GetFLInstance()
}

func GetFLInstance() *FullyLoaded {
	table := data.GetInstance()
	return &FullyLoaded{
		NeighborTable: *table,
	}
}

func (f *FullyLoaded) Debug() *bug.Debug {
	return f.NeighborTable.Debug()
}

func (f *FullyLoaded) Get(gid int) []data.Neighbor {
	return f.NeighborTable.Get(gid)
}

func (f *FullyLoaded) Ready() bool {
	return f.NeighborTable.Ready()
}
