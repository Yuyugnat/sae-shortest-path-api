package history

import (
	"encoding/json"
	"fmt"
	c "sae-shortest-path/connection"
	fast "sae-shortest-path/fastest"
)

type History struct {
	Paths []HistoryPath `json:"paths"`
}

type HistoryPath struct {
	Depart  string `json:"depart"`
	Arrivee string `json:"arrivee"`
}

func PutHistory(userID string, path *fast.Result) error {
	query := `
		INSERT INTO path_history (user_id, date, depart, arrivee)
		VALUES ($1, NOW(), $2, $3)
	`

	conn, _ := c.GetInstance()
	_, err := conn.DB.Exec(query, userID, path.VilleDepart, path.VilleArrivee)
	if err != nil {
		fmt.Println("Error inserting the path in the history", err)
		return err
	}
	return nil
}

func GetHistory(userID string) (History, error) {
	res := History{
		Paths: make([]HistoryPath, 0),
	}
	query := `
		SELECT data
		FROM path_history
		WHERE user_id = $1
	`

	conn, _ := c.GetInstance()
	rows, err := conn.DB.Query(query, userID)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var data string
		err = rows.Scan(&data)
		if err != nil {
			return res, err
		}
		var result fast.Result
		json.Unmarshal([]byte(data), &result)
		res.Paths = append(res.Paths, HistoryPath{
			Depart:  result.VilleDepart,
			Arrivee: result.VilleArrivee,
		})
	}
	return res, nil
}

func (h *History) JSON() ([]byte, error) {
	return json.Marshal(h.Paths)
}
