package objects

type Voisin struct {
	nrGid      int
	tronconGid int
	longueur   float64
}

func (v *Voisin) GetNrGid() int {
	return v.nrGid
}

func (v *Voisin) GetTronconGid() int {
	return v.tronconGid
}

func (v *Voisin) GetLongueur() float64 {
	return v.longueur
}