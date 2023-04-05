package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	c "sae-shortest-path/connection"
	o "sae-shortest-path/objects"
	s "sae-shortest-path/solver"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "tp1"
	password = "tp12023"
	dbname   = "sae"
)

var (
	mux *http.ServeMux
)

func main() {
	pgConn := c.NewPostgresConn(host, port, user, password, dbname)
	pgConn.Open()
	defer pgConn.Close()
	pgConn.Test()

	c.Conn = pgConn

	cors := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			h.ServeHTTP(w, r)
		})
	}

	mux = http.NewServeMux()

	mux.HandleFunc("/coucou", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("coucou"))
	})

	mux.HandleFunc("/hey", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hey " + r.URL.Query().Get("name")))
	})

	mux.HandleFunc("/shortest-path", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		communeRepo := o.NewNoeudCommuneRepo()
		noeudRoutierRepo := o.NewNoeudRoutierRepo()

		depart := r.URL.Query().Get("depart")
		arrivee := r.URL.Query().Get("arrivee")
		departIdNdRte:= communeRepo.GetIdNdRteByName(depart)
		arriveeIdNdRte := communeRepo.GetIdNdRteByName(arrivee)
		departGid := noeudRoutierRepo.GetGidByIdRte500(departIdNdRte)
		arriveeGid := noeudRoutierRepo.GetGidByIdRte500(arriveeIdNdRte)

		fmt.Println(departGid, depart, arriveeGid, arrivee)
		geomDepart := noeudRoutierRepo.GetGeomFromGid(departGid)
		geomArrivee := noeudRoutierRepo.GetGeomFromGid(arriveeGid)

		fmt.Printf("The distance between %s and %s is %f\n", depart, arrivee, noeudRoutierRepo.GetDistance(geomDepart, geomArrivee))

		// solver := s.NewDijkstra(departGid, arriveeGid)
		solver := s.NewAStar(departGid, arriveeGid)
		now := time.Now()
		distance, times := solver.Solve()
		fmt.Printf("{\n")
		fmt.Printf("   Time: %f\n", time.Since(now).Seconds())
		fmt.Printf("\n   Distance: %s - %s  => %f\n\n", depart, arrivee, distance)
		for k, v := range times {
			fmt.Printf("   Action: %s => took %f seconds in %d calls\n", k, v.Time, v.Call)
		}
		fmt.Printf("}\n")
		result := s.NewResultat(distance)
		resp, _ := json.Marshal(result)
		w.Write(resp)
	})

	http.ListenAndServe(":8080", cors(mux))
}
