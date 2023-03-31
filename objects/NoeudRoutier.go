package objects

import (
	"database/sql"
	"fmt"
)

// private int $gid,
// private string $id_rte500,

type NoeudRoutier struct {
	gid      int
	idRte500 string
	voisins  []Voisin
	db       *sql.DB
}

type Voisin struct {
	nrGid      int
	tronconGid int
	longueur   float64
}

func NewNoeudRoutier(gid int, idRte500 string, db *sql.DB) *NoeudRoutier {
	nr := &NoeudRoutier{
		gid:      gid,
		idRte500: idRte500,
		db:       db,
		voisins: []Voisin{},
	}
	nr.GenerateVoisins()
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

// public function getVoisins(int $noeudRoutierGid): array
//
//	{
//	    $requeteSQL = <<<SQL
//	        (
//	            (SELECT noeud_voisin as noeud_routier_gid, troncon_id as troncon_gid, longueur
//	            from voisins_noeud
//	            WHERE noeud_routier = :gidTag
//	        ) UNION (SELECT noeud_routier as noeud_routier_gid, troncon_id as troncon_gid, longueur
//	         from voisins_noeud
//	         WHERE noeud_voisin = :gidTag)
//	        );
//	    SQL;
//	    $pdoStatement = ConnexionBaseDeDonnees::getPdo()->prepare($requeteSQL);
//	    $pdoStatement->execute(array(
//	        "gidTag" => $noeudRoutierGid
//	    ));
//	    return $pdoStatement->fetchAll(PDO::FETCH_ASSOC);
//	}
func (n *NoeudRoutier) GenerateVoisins() {
	query := `
	(
		SELECT noeud_voisin as noeud_routier_gid, troncon_id as troncon_gid, longueur 
		FROM voisins_noeud
		WHERE noeud_routier = ?
	) UNION (
		SELECT noeud_routier as noeud_routier_gid, troncon_id as troncon_gid, longueur 
		FROM voisins_noeud
		WHERE noeud_voisin = ?
	)
	`
	rows, err := n.db.Query(query, n.gid, n.gid)
	if err != nil {
		fmt.Println("Error while querying the database : ", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var nrGid int
		var tronconGid int
		var longueur float64
		err = rows.Scan(&nrGid, &tronconGid, &longueur)
		if err != nil {
			fmt.Println("Error while scanning the database : ", err)
			return
		}
		n.voisins = append(n.voisins, Voisin{nrGid, tronconGid, longueur})
	}
}
