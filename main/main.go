package main

import (
	"net/http"
	c "sae-shortest-path/connection"

	_ "github.com/lib/pq"
)

const (
	host     = "192.168.123.202"
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

	mux = http.NewServeMux()
	mux.HandleFunc("/coucou", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("coucou"))
	})

	mux.HandleFunc("/hey", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hey " + r.URL.Query().Get("name")))
	})
}
