package objects

import (
	"fmt"
	c "sae-shortest-path/connection"
)

type NoeudCommune struct {
	gid        int
	idRte500   string
	nomCommune string
	idNdRte    string
}

type NoeudCommuneRepo struct {}

func NewNoeudCommune(gid int, idRte500 string, nomCommune string, idNdRte string) *NoeudCommune {
	return &NoeudCommune{
		gid:        gid,
		idRte500:   idRte500,
		nomCommune: nomCommune,
		idNdRte:    idNdRte,
	}
}

func NewNoeudCommuneRepo() *NoeudCommuneRepo {
	return &NoeudCommuneRepo{}
}

func (repo *NoeudCommuneRepo) GetIdNdRteByName(name string) string {
	query := `
		SELECT id_nd_rte
		FROM noeud_commune
		WHERE nom_comm = $1
	`
	row := c.Conn.DB.QueryRow(query, name)

	var idNdRte string
	row.Scan(&idNdRte)
	return idNdRte
}

func (repo *NoeudCommuneRepo) GetHeuristicDistance(gid1, gid2 int) float64 {
	query := `
		SELECT distance
		FROM heuristic_gid
		WHERE gid_1 = $1 AND gid_2 = $2
	`

	row := c.Conn.DB.QueryRow(query, gid1, gid2)

	if err := row.Err(); err != nil {
		fmt.Println("Error while querying data : ", err)
		return 0
	}

	var distance float64
	row.Scan(&distance)
	return distance
}
