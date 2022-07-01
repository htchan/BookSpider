package model

import (
	"database/sql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"time"
)

var db *sql.DB

func init() {
	var err error
	// location := os.Getenv("ASSETS_LOCATION")+"/test-data/internal_database_sqlite.db"
	// db, err = sql.Open("sqlite3", location + "?cache=shared")
	connection := os.Getenv("DB_CONN")
	db, err = sql.Open("postgres", connection)
	if err != nil {
		panic(err)
	}
	db.SetConnMaxIdleTime(30 * time.Second)
	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)
}
