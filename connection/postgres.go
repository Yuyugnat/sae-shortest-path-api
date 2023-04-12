package connection

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type PostgresConn struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBname   string `json:"dbname"`
	DB       *sql.DB
}

const (
	driver = "postgres"
)

var instance *PostgresConn

func newPostgresConn(confpath string) (*PostgresConn, error) {
	var conn *PostgresConn
	byteData, err := ioutil.ReadFile(confpath)
	if err != nil {
		fmt.Println("Error reading the configuration file", err)
		return conn, err
	}
	err = json.Unmarshal(byteData, &conn)

	if err != nil {
		fmt.Println("Error unmarshaling the configuration file", err)
		return conn, err
	}

	fmt.Println("Conn is", conn)

	err = conn.open()
	if err != nil {
		fmt.Println("Error opening the connection", err)
		return conn, err
	}
	return conn, nil
}

func GetInstance() (*PostgresConn, error) {
	if instance == nil {
		var err error
		instance, err = newPostgresConn("connection/config.json")
		if err != nil {
			return nil, err
		}
	}
	return instance, nil
}

func (pc *PostgresConn) open() error {
	infos := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		pc.Host,
		pc.Port,
		pc.User,
		pc.Password,
		pc.DBname,
	)

	var err error
	pc.DB, err = sql.Open(driver, infos)
	return err
}

func (pc *PostgresConn) close() {
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

func FirstConnection() {
	inst, err := GetInstance()
	timer := 10
	for err != nil || inst.DB.Ping() != nil {
		fmt.Println("Error getting the instance", err)
		fmt.Printf("Trying to reconnect after %d seconds\n", timer)
		time.Sleep(time.Duration(timer) * time.Second)
		inst, err = GetInstance()
	}
}
