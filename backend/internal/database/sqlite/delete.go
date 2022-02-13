package sqlite

import (
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
	"errors"
)

func (db *SqliteDB) DeleteBookRecords(records []database.BookRecord) error {
	return errors.New("It is not expected to have any books being deleted")
}

func (db *SqliteDB) DeleteWriterRecords(records []database.WriterRecord) error {
	return errors.New("It is not expected to have any writers being deleted")
}

func (db *SqliteDB) DeleteErrorRecords(records []database.ErrorRecord) (err error) {
	defer utils.Recover(func() {})
	if len(records) == 0 {
		err = errors.New("nil error record in creating error")
		panic(err)
	}
	for _, record := range records {
		db.execute(ErrorDeleteStatement(&record))
	}
	return
}