package testing

import "encoding/json"

type Resultat struct {
	Distance       float64 `json:"distance"`
	VilleDepart    string  `json:"villeDepart"`
	VilleArrivee   string  `json:"villeArrivee"`
	PointsReversed bool    `json:"pointsReversed"`
	Points         []Point `json:"points"`
	ErrCode            ErrCode `json:"errCode"` // 0 = no error, 1 = error
	ErrMsg         string  `json:"errMsg"`
}

func (r *Resultat) JSON() []byte {
	res, _ := json.Marshal(r)
	return res
}
