package solver

type Resultat struct {
	Distance float64
}

func NewResultat(distance float64) *Resultat {
	return &Resultat{
		Distance: distance,
	}
}