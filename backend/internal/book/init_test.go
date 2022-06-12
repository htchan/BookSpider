package book

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
	"os"
	"time"
	"github.com/htchan/ApiParser"
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
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(50)

	ApiParser.SetDefault(
		ApiParser.FromDirectory(os.Getenv("ASSETS_LOCATION") + "/test-data/api_parser"))
}
