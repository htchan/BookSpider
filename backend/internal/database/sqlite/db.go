package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"sync"

	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
)

type SqliteDB struct {
	_db *sql.DB
	statements map[int]string
	statementCount int
	lock sync.Mutex
	maxStatements int
}

const (
	SQLITE_MAX_IDLE_CONN = 10
	SQLITE_MAX_OPEN_CONN = 10000
)

func NewSqliteDB(location string, maxStatements int) (db *SqliteDB) {
	var err error
	db = new(SqliteDB)
	db._db, err = sql.Open("sqlite3", location + "?cache=shared")
	utils.CheckError(err)
	db._db.SetMaxIdleConns(SQLITE_MAX_IDLE_CONN)
	db._db.SetMaxOpenConns(SQLITE_MAX_OPEN_CONN)
	db.statements = make(map[int]string)
	db.statementCount = 0
	db.maxStatements = maxStatements
	return
}

func (db *SqliteDB) execute(statement string) {
	db.lock.Lock()
	defer db.lock.Unlock()
	db.statements[db.statementCount] = statement
	db.statementCount++
	if db.statementCount >= db.maxStatements {
		db.commit()
		db.statementCount = 0
	}
}

func contains(array []string, target string) bool {
	if target == "" {
		return true
	}
	for _, item := range array {
		if item == "" {
			break
		}
		if item == target {
			return true
		}
	}
	return false
}

func (db *SqliteDB) commit() (err error) {
	var tx *sql.Tx
	defer utils.Recover(func() {
		if tx != nil { tx.Rollback() }
		// write all statement to file (eg. error log) if it fails
		// file format: yyyy-mm-ddTHH:MM:SS-commit.sql
	})

	tx, err = db._db.Begin()
	utils.CheckError(err)
	executedStatements := make([]string, db.statementCount)
	executedCount := 0
	for i := 0; i < db.statementCount; i++ {
		// it will skip repeated statement
		if contains(executedStatements, db.statements[i]) { continue }
		executedStatements[executedCount] = db.statements[i]
		executedCount++
		_, err = tx.Exec(db.statements[i])
		utils.CheckError(err)
	}
	db.statements = make(map[int]string)
	db.statementCount = 0
	return tx.Commit()
}

func (db *SqliteDB) Commit() (err error) {
	db.lock.Lock()
	defer db.lock.Unlock()
	return db.commit()
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
