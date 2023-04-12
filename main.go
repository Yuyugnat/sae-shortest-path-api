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
	conn "sae-shortest-path/connection"
	hist "sae-shortest-path/history"
	"time"

	_ "github.com/lib/pq"
)

var (
	mux  *http.ServeMux
	port = 8080
)

func main() {
	conn.FirstConnection()
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

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(res.JSON()))

		// hist.PutHistory(userID, res)
		if userID := r.URL.Query().Get("user"); userID != "" {
			hist.PutHistory(userID, res)
		}

		solver.Debug().Print()
	})

	mux.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userID := r.URL.Query().Get("user")

		history, err := hist.GetHistory(userID)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json, _ := history.JSON()
		w.Write(json)
	})

	mux.HandleFunc("/debug-shortest-path", func(w http.ResponseWriter, r *http.Request) {
		var solver s.ISolver
		var nbGetter nb.NeighborGetter
		var hCalc hc.HeuristicCalculator
		var queue prio.PriorityQueue

		fmt.Println(r.URL.Query().Get("hcalc"))

		switch r.URL.Query().Get("hcalc") {
		case "haversine":
			hCalc = hc.NewHaversineCalculator()
		case "dijkstra":
			hCalc = hc.NewDijkstraCalculator()
		default:
			hCalc = hc.NewHaversineCalculator()
		}

		switch r.URL.Query().Get("nbget") {
		case "voisins_jsonb":
			nbGetter = nb.NewJsonbView()
		case "hash_table":
			nbGetter = nb.GetFLInstance()
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

		solver = s.NewSolver(nbGetter, hCalc, queue)

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
