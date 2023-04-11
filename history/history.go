package history

import (
	"encoding/json"
	"fmt"
	c "sae-shortest-path/connection"
	fast "sae-shortest-path/fastest"
)

type History struct {
	paths []fast.Result
}

func PutHistory(userID string, path *fast.Result) error {
	query := `
		INSERT INTO path_history (user_id, date, data)
		VALUES ($1, NOW(), $2)
	`

	conn, _ := c.GetInstance()
	_, err := conn.DB.Exec(query, userID, path)
	if err != nil {
		fmt.Println("Error inserting the path in the history", err)
		return err
	}
	return nil
}

func GetHistory(userID string) (History, error) {
	res := History{
		paths: []fast.Result{},
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
		res.paths = append(res.paths, result)
	}
	return res, nil
}

func (h *History) JSON() ([]byte, error) {
	return json.Marshal(h.paths)
}