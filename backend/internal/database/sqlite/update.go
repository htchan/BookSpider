package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
	"errors"
)

func (db *SqliteDB) UpdateBookRecord(record *database.BookRecord) (err error) {
	var tx *sql.Tx
	defer utils.Recover(func() {
		if tx != nil { tx.Rollback() }
	})
	if record == nil {
		err = errors.New("nil error record in updating book")
		panic(err)
	}
	tx, err = db._db.Begin()
	utils.CheckError(err)
	_, err = tx.Exec(
		"update books set " +
		"title=?, writer_id=?, type=?, update_date=?, update_chapter=?, status=?" +
		"where site=? and id=? and hash_code=?",
		record.Title, record.WriterId, record.Type,
		record.UpdateDate, record.UpdateChapter,
		record.Status,
		record.Site, record.Id, record.HashCode)
	utils.CheckError(err)
	return tx.Commit()
}

func (db *SqliteDB) UpdateErrorRecord(record *database.ErrorRecord) (err error) {
	var tx *sql.Tx
	defer utils.Recover(func() {
		if tx != nil { tx.Rollback() }
	})
	if record == nil {
		err = errors.New("nil error record in updating error")
		panic(err)
	}
	tx, err = db._db.Begin()
	utils.CheckError(err)
	_, err = tx.Exec(
		"update errors set " +
		"data=? where site=? and id=?",
		record.Error.Error(),
		record.Site, record.Id)
	utils.CheckError(err)
	return tx.Commit()
}