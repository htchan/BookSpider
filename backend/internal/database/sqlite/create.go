package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"errors"

	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
)

func (db *SqliteDB) CreateBookRecord(record *database.BookRecord) (err error) {
	var tx *sql.Tx
	defer utils.Recover(func() {
		if tx != nil { tx.Rollback() }
	})
	if record == nil {
		err = errors.New("nil error record in creating book")
		panic(err)
	}
	tx, err = db._db.Begin()
	utils.CheckError(err)
	_, err = tx.Exec(
		"insert into books " +
		"(" + database.BOOK_RECORD_FIELDS + ") " +
		"values " +
		"(?, ?, ?, ?, ?, ?, ?, ?, ?)",
		record.Site, record.Id, record.HashCode,
		record.Title, record.WriterId, record.Type,
		record.UpdateDate, record.UpdateChapter,
		record.Status)
	utils.CheckError(err)
	return tx.Commit()
}

func (db *SqliteDB) CreateWriterRecord(record *database.WriterRecord) (err error) {
	var tx *sql.Tx
	defer utils.Recover(func() {
		if tx != nil { tx.Rollback() }
	})
	if record == nil {
		err = errors.New("nil error record in creating writer")
		panic(err)
	}
	tx, err = db._db.Begin()
	utils.CheckError(err)
	if record.Id < 0 {
		var result sql.Result
		result, err = tx.Exec(
			"insert into writers (name) values (?)",
			record.Name)
		idInt64, _ := result.LastInsertId()
		record.Id = int(idInt64)
	} else {
		_, err = tx.Exec(
			"insert into writers (id, name) values (?, ?)",
			record.Id, record.Name)
	}
	utils.CheckError(err)
	return tx.Commit()
}

func (db *SqliteDB) CreateErrorRecord(record *database.ErrorRecord) (err error) {
	var tx *sql.Tx
	defer utils.Recover(func() {
		if tx != nil { tx.Rollback() }
	})
	if record == nil {
		err = errors.New("nil error record in creating error")
		panic(err)
	}
	tx, err = db._db.Begin()
	utils.CheckError(err)
	_, err = tx.Exec(
		"insert into errors (" + database.ERROR_RECORD_FIELDS + ") values (?, ?, ?)",
		record.Site, record.Id, record.Error.Error())
	utils.CheckError(err)
	return tx.Commit()
}