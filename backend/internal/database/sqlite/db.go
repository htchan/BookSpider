package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
)

type SqliteDB struct {
	_db *sql.DB
	queryTx *sql.Tx
}

const (
	SQLITE_MAX_IDLE_CONN = 10
	SQLITE_MAX_OPEN_CONN = 1
)

func NewSqliteDB(location string) (db *SqliteDB) {
	var err error
	db = new(SqliteDB)
	db._db, err = sql.Open("sqlite3", location + "?cache=shared")
	utils.CheckError(err)
	db._db.SetMaxIdleConns(SQLITE_MAX_IDLE_CONN)
	db._db.SetMaxOpenConns(SQLITE_MAX_OPEN_CONN)
	return
}

func (db *SqliteDB) Summary(site string) (record database.SummaryRecord) {
	defer utils.Recover(func() {})
	// select count(*), count(distinct (site||id)), max(id) from books where site=?
	rows, err := db._db.Query("select count(*), count(distinct (site||id)), max(id) from books where site=?", site)
	utils.CheckError(err)
	if rows.Next() {
		rows.Scan(&record.BookCount, &record.UniqueBookCount, &record.MaxBookId)
	}
	utils.CheckError(rows.Close())
	// select count(*) from writers
	rows, err = db._db.Query("select count(*) from writers")
	utils.CheckError(err)
	if rows.Next() {
		rows.Scan(&record.WriterCount)
		record.WriterCount -= 1
	}
	utils.CheckError(rows.Close())
	// select count(*) from errors where site=?
	rows, err = db._db.Query("select count(*) from errors where site=?", site)
	utils.CheckError(err)
	if rows.Next() {
		rows.Scan(&record.ErrorCount)
	}
	utils.CheckError(rows.Close())
	// select max(id) from books where status != Error
	rows, err = db._db.Query("select max(id) from books where status != ?", database.Error)
	utils.CheckError(err)
	if rows.Next() {
		rows.Scan(&record.LatestSuccessId)
	}
	utils.CheckError(rows.Close())
	// select status, count(*) from books groun by status
	rows, err = db._db.Query("select status, count(*) from books group by status")
	utils.CheckError(err)
	var statusKey database.StatusCode
	var statusValue int
	record.StatusCount = make(map[database.StatusCode]int)
	for rows.Next() {
		rows.Scan(&statusKey, &statusValue)
		record.StatusCount[statusKey] = statusValue
	}
	utils.CheckError(rows.Close())
	return
}

func (db *SqliteDB) Close() (err error) {
	err = db._db.Close()
	db._db = nil
	return
}
