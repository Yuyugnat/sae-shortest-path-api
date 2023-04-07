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

type NoeudCommuneRepo struct {
	conn *c.PostgresConn
}

func NewNoeudCommune(gid int, idRte500 string, nomCommune string, idNdRte string) *NoeudCommune {
	return &NoeudCommune{
		gid:        gid,
		idRte500:   idRte500,
		nomCommune: nomCommune,
		idNdRte:    idNdRte,
	}
}

func NewNoeudCommuneRepo() *NoeudCommuneRepo {
	conn, _ := c.GetInstance()
	return &NoeudCommuneRepo{
		conn: conn,
	}
}

func (repo *NoeudCommuneRepo) GetIdNdRteByName(name string) (string, error) {
	query := `
		SELECT id_nd_rte
		FROM noeud_commune
		WHERE nom_comm = $1
	`
	row := repo.conn.DB.QueryRow(query, name)

	var idNdRte string
	err := row.Scan(&idNdRte)
	return idNdRte, err
}

func (repo *NoeudCommuneRepo) GetHeuristicDistance(gid1, gid2 int) float64 {
	query := `
		SELECT distance
		FROM heuristic_gid
		WHERE gid_1 = $1 AND gid_2 = $2
	`

	row := repo.conn.DB.QueryRow(query, gid1, gid2)

	if err := row.Err(); err != nil {
		fmt.Println("Error while querying data : ", err)
		return 0
	}

	var distance float64
	row.Scan(&distance)
	return distance
}

func (repo *NoeudCommuneRepo) GetSuperficie(nom string) int {
	query := `
		SELECT superficie
		FROM superficie_commune
		WHERE nom_comm = $1
	`
	row := repo.conn.DB.QueryRow(query, nom)

	var superficie int
	row.Scan(&superficie)
	return superficie
}
