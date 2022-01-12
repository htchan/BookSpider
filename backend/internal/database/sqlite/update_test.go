package sqlite

import (
	"testing"
	"os"
	"io"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/internal/utils"
)

func init() {
	source, err := os.Open("../../../assets/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./update_test.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func Test_Sqlite_DB_UpdateBookRecord(t *testing.T) {
	db := NewSqliteDB("./rows_test.db")
	defer db.Close()
	query := db.QueryErrorBySiteId("test", 2)
	defer query.Close()

	t.Run("success", func(t *testing.T) {
		bookRecord := &database.BookRecord{
			Site: "test",
			Id: 1,
			HashCode: 100,
			Title: "title-1-new",
			WriterId: 1,
			Type: "type-1-new",
			UpdateDate: "200",
			UpdateChapter: "chapter-1-new",
			Status: database.InProgress}
		err := db.UpdateBookRecord(bookRecord)
		
		if err != nil {
			t.Fatalf("DB.UpdateBookRecord failed err: %v", err)
		}

		query := db.QueryBooksByTitle("title-1-new")
		defer query.Close()
		query.Next()
		record, err := query.Scan()
		bookRecord = record.(*database.BookRecord)
		if err != nil || bookRecord.Site != "test" ||
			bookRecord.Id != 1 || bookRecord.HashCode != 100 ||
			bookRecord.Title != "title-1-new" || bookRecord.WriterId != 1 ||
			bookRecord.Type != "type-1-new" ||
			bookRecord.UpdateDate != "200" || bookRecord.UpdateChapter != "chapter-1-new" ||
			bookRecord.Status != database.InProgress {
				t.Fatalf(
					"DB.UpdateBookRecord does not change record in database record: %v, err: %v",
					bookRecord, err)
			}
	})

	t.Run("fail if input Book Record is nil", func(t *testing.T) {
		err := db.UpdateBookRecord(nil)
		
		if err == nil {
			t.Fatalf("DB.UpdateBookRecord success with nil err: %v", err)
		}
	})

	t.Run("pass but not create anything if site id hash not exist", func(t *testing.T) {
		bookRecord := &database.BookRecord{
			Site: "test",
			Id: 1,
			HashCode: 300,
			Title: "title-1-ultra-new",
			WriterId: 1,
			Type: "type-1-ultra-new",
			UpdateDate: "300",
			UpdateChapter: "chapter-1-ultra-new",
			Status: database.InProgress}
		err := db.UpdateBookRecord(bookRecord)
		
		if err != nil {
			t.Fatalf("DB.UpdateBookRecord failed err: %v", err)
		}

		query := db.QueryBooksByTitle("title-1-ultra-new")
		defer query.Close()
		if query.Next() {
			record, err := query.Scan()
			t.Fatalf(
				"Some record is being updated to title-1-ultra-new record: %v, err: %v",
				record, err)
		}
	})
}