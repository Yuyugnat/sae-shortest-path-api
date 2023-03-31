package objects

import "database/sql"

type NoeudCommune struct {
	gid        int
	idRte500   string
	nomCommune string
	idNdRte    string
	db 	   *sql.DB
}

func NewNoeudCommune(gid int, idRte500 string, nomCommune string, idNdRte string, db *sql.DB) *NoeudCommune {
	return &NoeudCommune{
		gid:        gid,
		idRte500:   idRte500,
		nomCommune: nomCommune,
		idNdRte:    idNdRte,
		db:         db,
	}
}