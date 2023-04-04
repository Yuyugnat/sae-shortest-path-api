package objects

type TronconRoute struct {
	gid           int
	idRte500      string
	sens          string
	numeroRoute   string
	longueurRoute float64
}

func NewTronconRoute(gid int, idRte500 string, sens string, numeroRoute string, longueurRoute float64) *TronconRoute {
	return &TronconRoute{
		gid:           gid,
		idRte500:      idRte500,
		sens:          sens,
		numeroRoute:   numeroRoute,
		longueurRoute: longueurRoute,
	}
}
