package objects

import "database/sql"

type TronconRoute struct {
	gid           int
	idRte500      string
	sens          string
	numeroRoute   string
	longueurRoute float64
	db            *sql.DB
}

func NewTronconRoute(gid int, idRte500 string, sens string, numeroRoute string, longueurRoute float64, db *sql.DB) *TronconRoute {
	return &TronconRoute{
		gid:           gid,
		idRte500:      idRte500,
		sens:          sens,
		numeroRoute:   numeroRoute,
		longueurRoute: longueurRoute,
		db:            db,
	}
}
