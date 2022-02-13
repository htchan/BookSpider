package sqlite

import (
	// "database/sql"
	// _ "github.com/mattn/go-sqlite3"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
	"errors"
	"fmt"
)

func (db *SqliteDB) UpdateBookRecord(bookRecord *database.BookRecord, writerRecord *database.WriterRecord) (err error) {
	defer utils.Recover(func() {})
	if bookRecord == nil || writerRecord == nil {
		err = errors.New("nil book record in updating book")
		panic(err)
	}
	db.execute(BookUpdateStatement(bookRecord, writerRecord.Name))
	return
}

func (db *SqliteDB) UpdateErrorRecord(record *database.ErrorRecord) (err error) {
	defer utils.Recover(func() {})
	if record == nil {
		err = errors.New("nil error record in updating error")
		panic(err)
	}
	db.execute(fmt.Sprintf(
		"update errors set " +
		"data=\"%v\" where site=\"%v\" and id=%v",
		record.Error.Error(), record.Site, record.Id))
	return
}