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
	Date   string `json:"date"`
	Depart  string `json:"depart"`
	Arrivee string `json:"arrivee"`
	Distance float64 `json:"distance"`
}

func PutHistory(userID string, path *fast.Result) error {
	query := `
		INSERT INTO path_history (user_id, date, depart, arrivee, distance)
		VALUES ($1, NOW(), $2, $3, $4)
	`

	conn, _ := c.GetInstance()
	_, err := conn.DB.Exec(query, userID, path.VilleDepart, path.VilleArrivee, path.Distance)
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
		SELECT depart, arrivee, date, distance
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
		var depart, arrivee, date string
		var distance float64
		err = rows.Scan(&depart, &arrivee, &date, &distance)
		if err != nil {
			return res, err
		}
		res.Paths = append(res.Paths, HistoryPath{
			Depart: depart,
			Arrivee: arrivee,
			Date: date,
			Distance: distance,
		})
	}
	return res, nil
}

func (h *History) JSON() ([]byte, error) {
	return json.Marshal(h.Paths)
}
