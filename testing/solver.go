package testing

import (
	"fmt"
	c "sae-shortest-path/connection"
	bug "sae-shortest-path/debugging"
	o "sae-shortest-path/objects"
	nb "sae-shortest-path/testing/neighbors"
)

type ErrCode int

const (
	NoErr ErrCode = iota
	NoDepartOrArrivee
	NoPath
	NotReady
)

type ISolver interface {
	Solve(start, end string) *Resultat
	Debug() *bug.Debug
}

type solver struct {
	ISolver
	debugger *bug.Debug
	reversed bool
	nrRepo   *o.NoeudRoutierRepo
	ncRepo   *o.NoeudCommuneRepo

	NbGetter nb.NeighborGetter

	DepartGid   int
	ArriveeGid  int
	ArriveeGeom string
	Depart      string
	Arrivee     string
	Reversed    bool
	LastPoint   *Point
}

func (s *solver) Instantiate(nbGetter nb.NeighborGetter) {
	repoNR := o.NewNoeudRoutierRepo()
	repoNC := o.NewNoeudCommuneRepo()
	s.reversed = false
	s.nrRepo = repoNR
	s.ncRepo = repoNC
	s.NbGetter = nbGetter
	s.debugger = bug.NewDebug()
}

func (s *solver) InitSearch(depart, arrivee string) error {
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

func (s *solver) GetPointFromGid(gid int) Point {
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

func (s *solver) Debug() *bug.Debug {
	return s.debugger
}
