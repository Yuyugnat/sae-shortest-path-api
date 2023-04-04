package connection

import (
	"database/sql"
	"fmt"
)

const (
	driver = "postgres"
)

var (
	Conn *PostgresConn
)

type PostgresConn struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
	DB *sql.DB
}

func NewPostgresConn(host string, port int, user, password, dbname string) *PostgresConn {
	return &PostgresConn{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		dbname:   dbname,
	}
}

func (pc *PostgresConn) Open() {
	infos := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		pc.host,
		pc.port,
		pc.user,
		pc.password,
		pc.dbname,
	)

	var err error
	pc.DB, err = sql.Open(driver, infos)

	if err != nil {
		fmt.Println("Error opening the database", err)
	}
}

func (pc *PostgresConn) Close() {
	if pc.DB != nil {
		pc.DB.Close()
	}
}

func (pc *PostgresConn) Test() {
	if err := pc.DB.Ping(); err != nil {
		fmt.Println("The DB is not opened any more", err)
		return
	}
	fmt.Println("The DB is alive !")
}
