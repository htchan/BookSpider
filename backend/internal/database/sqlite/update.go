package sqlite

import (
	// "database/sql"
	// _ "github.com/mattn/go-sqlite3"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
	"errors"
	"fmt"
	"time"
	"strconv"
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

func (db *SqliteDB) UpdateBookRecordsStatusByChapter() error {
	matchingChapterCriteria := []string{"后记", "後記", "新书", "新書", "结局", "結局", "感言",
	"尾声", "尾聲", "终章", "終章", "外传", "外傳", "完本", "结束", "結束", "完結",
	"完结", "终结", "終結", "番外", "结尾", "結尾", "全书完", "全書完", "全本完"}
	sqlStmt := fmt.Sprintf("update books set status=%v where (", database.End)
	for _, criteria := range matchingChapterCriteria {
		sqlStmt += "update_chapter like '%" + criteria + "%' or "
	}
	sqlStmt += fmt.Sprintf(
		"update_date < '" + strconv.Itoa(time.Now().Year()-1) + "') " +
		"and status != %v and status != %v and status != %v",
		database.Error, database.End, database.Download)
	db.execute(sqlStmt)
	return nil
}
