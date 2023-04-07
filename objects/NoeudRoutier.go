package objects

import (
	"fmt"
	"log"
	"math"
	c "sae-shortest-path/connection"
)

// private int $gid,
// private string $id_rte500,

type NoeudRoutier struct {
	gid      int
	idRte500 string
	voisins  []Voisin
}

type NoeudRoutierRepo struct{
	conn *c.PostgresConn
}

func NewNoeudRoutierRepo() (*NoeudRoutierRepo) {
	conn, _ := c.GetInstance()
	return &NoeudRoutierRepo{
		conn: conn,
	}
}

func NewNoeudRoutier(gid int, idRte500 string) *NoeudRoutier {
	repo := NewNoeudRoutierRepo()

	nr := &NoeudRoutier{
		gid:      gid,
		idRte500: idRte500,
		voisins:  []Voisin{},
	}
	repo.GenerateVoisins(nr)
	return nr
}

func (n *NoeudRoutier) GetGid() int {
	return n.gid
}

func (n *NoeudRoutier) GetIdRte500() string {
	return n.idRte500
}

func (n *NoeudRoutier) GetVoisins() []Voisin {
	return n.voisins
}

func (repo *NoeudRoutierRepo) GenerateVoisins(nr *NoeudRoutier) {
	query := `
		(
			SELECT noeud_voisin as noeud_routier_gid_1, noeud_routier as noeud_routier_gid_2, troncon_id as troncon_gid, longueur 
			FROM voisins_noeud
			WHERE noeud_routier = $1 or noeud_voisin = $1
		);
	`
	rows, err := repo.conn.DB.Query(query, nr.gid)
	if err != nil {
		fmt.Println("Error while querying the database (GenerateVoisins) : ", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var nrGid1 int
		var nrGid2 int
		var tronconGid int
		var longueur float64
		err = rows.Scan(&nrGid1, &nrGid2, &tronconGid, &longueur)
		if err != nil {
			fmt.Println("Error while scanning the database : ", err)
			return
		}
		if nrGid1 == nr.gid {
			nr.voisins = append(nr.voisins, Voisin{nrGid2, tronconGid, longueur})
		} else {
			nr.voisins = append(nr.voisins, Voisin{nrGid1, tronconGid, longueur})
		}
	}
}

func (repo *NoeudRoutierRepo) GetVoisins(gid int) []Voisin {
	query := `
		(
			SELECT noeud_voisin as noeud_routier_gid_1, noeud_routier as noeud_routier_gid_2, troncon_id as troncon_gid, longueur 
			FROM voisins_noeud
			WHERE noeud_routier = $1 or noeud_voisin = $1
		);
	`
	rows, err := repo.conn.DB.Query(query, gid)
	if err != nil {
		fmt.Println("Error while querying the database (GetVoisins) : ", err)
		return []Voisin{}
	}
	defer rows.Close()

	res := []Voisin{}

	for rows.Next() {
		var nrGid1 int
		var nrGid2 int
		var tronconGid int
		var longueur float64
		err = rows.Scan(&nrGid1, &nrGid2, &tronconGid, &longueur)
		if err != nil {
			fmt.Println("Error while scanning the database : ", err)
			return []Voisin{}
		}
		if nrGid1 == gid {
			res = append(res, Voisin{nrGid2, tronconGid, longueur})
		} else {
			res = append(res, Voisin{nrGid1, tronconGid, longueur})
		}
	}

	return res
}

func (repo *NoeudRoutierRepo) GetGidByIdRte500(idRte500 string) int {
	query := `
		SELECT gid
		FROM noeud_routier
		WHERE id_rte500 = $1
	`
	var gid int
	row := repo.conn.DB.QueryRow(query, idRte500)
	err := row.Scan(&gid)
	if err != nil {
		fmt.Println("Error while querying the database (GetGidByIdRte500) : ", err)
		return -1
	}
	return gid
}

func (repo *NoeudRoutierRepo) GetIdRte500ByPrimaryKey(gid int) (*NoeudRoutier, error) {
	query := `
		SELECT id_rte500
		FROM noeud_routier
		WHERE gid = $1
	`
	var idRte500 string
	row := repo.conn.DB.QueryRow(query, gid)
	err := row.Scan(&idRte500)
	if err != nil {
		fmt.Println("Error while querying the database (GetByPrimaryKey) : ", err)
		return nil, err
	}
	res := NewNoeudRoutier(gid, idRte500)
	return res, nil
}

func (repo *NoeudRoutierRepo) GetGeomFromGid(gid int) string {
	query := `
		SELECT ST_AsText(geom)
		FROM geom_noeud_routier
		WHERE gid = $1
	`

	row := repo.conn.DB.QueryRow(query, gid)

	var geom string
	err := row.Scan(&geom)
	if err != nil {
		fmt.Println("Error while querying the database (GetGeomFromGid) : ", err)
		return ""
	}
	return geom
}

func (repo *NoeudRoutierRepo) GetDistance(geom1, geom2 string) float64 {
	query := `
		SELECT ST_Distance(ST_GeomFromText($1, 4326)::geography, ST_GeomFromText($2, 4326)::geography)
	`
	row := repo.conn.DB.QueryRow(query, geom1, geom2)
	var distance float64
	err := row.Scan(&distance)
	if err != nil {
		fmt.Println("Error while querying the database : ", err)
		return -1
	}
	// fmt.Println("Distance : ", distance * 6371 * 3.14159265 / 180)
	return distance * 6371 * 3.14159265 / 180.0
}

func (repo *NoeudRoutierRepo) GetDistance2(gid1 int, gid2 int) float64 {
	query := `
		SELECT lon, lat
		FROM geom_noeud_routier_xy
		WHERE gid = $1;
	`
	
	// fmt.Println("GID1 : ", gid1)
	row := repo.conn.DB.QueryRow(query, gid1)
	var lon1 float64
	var lat1 float64

	err := row.Scan(&lon1, &lat1)
	if err != nil {
		log.Fatalln("Error while querying the database caca : ", err)
	}

	query = `
		SELECT lon, lat
		FROM geom_noeud_routier_xy
		WHERE gid = $1;
	`
	row = repo.conn.DB.QueryRow(query, gid2)

	var lon2 float64
	var lat2 float64

	err = row.Scan(&lon2, &lat2)	
	if err != nil {
		log.Fatalln("Error while querying the database caca : ", err)
	}

	phi1 := lat1 * math.Pi / 180.0
	phi2 := lat2 * math.Pi / 180.0
	deltaPhi := (lat2 - lat1) * math.Pi / 180.0
	deltaLambda := (lon2 - lon1) * math.Pi / 180.0

	a := math.Sin(deltaPhi/2.0)*math.Sin(deltaPhi/2.0) + math.Cos(phi1)*math.Cos(phi2)*math.Sin(deltaLambda/2.0)*math.Sin(deltaLambda/2.0)
	c := 2.0 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := 6371.0 * c

	// fmt.Println("Distance : ", distance)
	return distance
}
