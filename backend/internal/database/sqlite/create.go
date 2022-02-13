package sqlite

import (
	"errors"

	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
)

func (db *SqliteDB) CreateBookRecord(bookRecord *database.BookRecord, writerRecord *database.WriterRecord) (err error) {
	defer utils.Recover(func() {})
	if bookRecord == nil || writerRecord == nil {
		err = errors.New("nil record in creating book with writer")
		panic(err)
	}
	db.execute(BookInsertStatement(bookRecord, writerRecord.Name))
	return
}

func (db *SqliteDB) CreateWriterRecord(record *database.WriterRecord) (err error) {
	defer utils.Recover(func() {})
	if record == nil {
		err = errors.New("nil error record in creating writer")
		panic(err)
	}
	db.execute(WriterInsertStatement(record))
	return
}

func (db *SqliteDB) CreateErrorRecord(record *database.ErrorRecord) (err error) {
	defer utils.Recover(func() {})
	if record == nil {
		err = errors.New("nil error record in creating error")
		panic(err)
	}
	db.execute(ErrorInsertStatement(record))
	return
}