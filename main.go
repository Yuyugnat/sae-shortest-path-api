package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	s "sae-shortest-path/testing"
	"strconv"

	fast "sae-shortest-path/fastest"
	hc "sae-shortest-path/testing/calculator"
	nb "sae-shortest-path/testing/neighbors"
	prio "sae-shortest-path/testing/priorityqueue"
	"time"

	_ "github.com/lib/pq"
)

var (
	mux  *http.ServeMux
	port = 8080
)

func main() {
	nb.Load()

	if len(os.Args) > 1 {
		potentialPort, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalln(err)
		}
		port = potentialPort
	}

	cors := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			h.ServeHTTP(w, r)
		})
	}

	mux = http.NewServeMux()

	mux.HandleFunc("/shortest-path", func(w http.ResponseWriter, r *http.Request) {

		depart := r.URL.Query().Get("depart")
		arrivee := r.URL.Query().Get("arrivee")

		solver := fast.NewFastest(depart, arrivee)

		now := time.Now()
		res := solver.Solve()
		fmt.Printf("Execution time:%f\n", time.Since(now).Seconds())

		switch res.ErrCode {
		case fast.NoErr:
			fmt.Printf("Found : %f\n", res.Distance)
			w.WriteHeader(http.StatusOK)
		case fast.NoDepartOrArrivee:
			w.WriteHeader(http.StatusNotFound)
		case fast.NoPath:
			w.WriteHeader(http.StatusNotFound)
		case fast.NotReady:
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		os.WriteFile("example.json", res.JSON(), 0777)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(res.JSON()))
	})

	mux.HandleFunc("/debug-shortest-path", func(w http.ResponseWriter, r *http.Request) {
		var solver s.ISolver
		var nbGetter nb.NeighborGetter
		var hCalc hc.HeuristicCalculator
		var queue prio.PriorityQueue

		switch r.URL.Query().Get("hcalc") {
		default:
			hCalc = hc.NewHaversineCalculator()
		}

		switch r.URL.Query().Get("nbget") {
		default:
			nbGetter = nb.GetFLInstance()
		}

		switch r.URL.Query().Get("prioqueue") {
		case "minheap":
			queue = prio.NewPrioMinHeap()
		case "map":
			queue = prio.NewPrioMap()
		default:
			queue = prio.NewPrioMinHeap()
		}

		switch r.URL.Query().Get("solver") {
		case "dijkstra":
			solver = s.NewDijkstra(nbGetter, hCalc, queue)
		case "astar":
			solver = s.NewAStar(nbGetter, hCalc, queue)
		default:
			// A*
			solver = s.NewAStar(nbGetter, hCalc, queue)
		}

		depart := r.URL.Query().Get("depart")
		arrivee := r.URL.Query().Get("arrivee")

		now := time.Now()
		res := solver.Solve(depart, arrivee)
		solver.Debug().TotalTime = time.Since(now).Seconds()
		solver.Debug().Distance = res.Distance

		switch res.ErrCode {
		case s.NoErr:
			w.WriteHeader(http.StatusOK)
		case s.NoDepartOrArrivee:
			w.WriteHeader(http.StatusNotFound)
		case s.NoPath:
			w.WriteHeader(http.StatusNotFound)
		case s.NotReady:
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		w.Header().Set("Content-Type", "application/json")

		resp := solver.Debug().JSON()
		w.Write(resp)
		fmt.Println(string(resp))
	})

	http.ListenAndServe(fmt.Sprintf(":%d", port), cors(mux))
}
