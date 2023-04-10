package neighbors

import (
	"encoding/json"
	conn "sae-shortest-path/connection"
	data "sae-shortest-path/data"
)

type JsonbView struct {
	Conn *conn.PostgresConn
}

func NewJsonbView() *JsonbView {
	v, _ := conn.GetInstance()
	return &JsonbView{
		Conn: v,
	}
}

func (j *JsonbView) Get(gid int) []data.Neighbor {
	query := `
		SELECT voisins
		FROM voisins_jsonb
		WHERE gid = $1
	`
	row := j.Conn.DB.QueryRow(query, gid)

	var voisins []data.Neighbor
	var strNeighbors string

	 err := row.Scan(&strNeighbors)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(strNeighbors), &voisins)
	if err != nil {
		panic(err)
	}

	return voisins
}
