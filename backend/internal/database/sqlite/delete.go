package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
	"fmt"
	"strings"
	"errors"
)

func (db *SqliteDB) DeleteBookRecord(records []database.BookRecord) error {
	return errors.New("It is not expected to have any books being deleted")
}

func (db *SqliteDB) DeleteWriterRecord(records []database.WriterRecord) error {
	return errors.New("It is not expected to have any writers being deleted")
}

func (db *SqliteDB) DeleteErrorRecord(records []database.ErrorRecord) (err error) {
	var tx *sql.Tx
	defer utils.Recover(func() {
		if tx != nil { tx.Rollback() }
	})
	if len(records) == 0 {
		err = errors.New("nil error record in creating error")
		panic(err)
	}
	results := make([]interface{}, len(records))
	for i, record := range records {
		results[i] = fmt.Sprintf("%v-%v", record.Site, record.Id)
	}
	tx, err = db._db.Begin()
	utils.CheckError(err)
	_, err = tx.Exec(
		"delete from errors " +
		"where site||'-'||id in (?" + strings.Repeat(", ?", len(results) - 1) + ")",
		results...
	)
	utils.CheckError(err)
	return tx.Commit()
}