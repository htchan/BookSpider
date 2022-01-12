package sqlite

import (
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
	"errors"
)

func (db *SqliteDB) UpdateBookRecord(record *database.BookRecord) (err error) {
	defer utils.Recover(func() {})
	if record == nil {
		err = errors.New("nil error record in updating book")
		panic(err)
	}
	tx, err := db._db.Begin()
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