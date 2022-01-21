package sqlite

import (
	"testing"
	"os"
	"io"
	"errors"
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
		// query.Next()
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
		record, err := query.Scan()
		if err == nil {
			t.Fatalf(
				"Some record is being updated to title-1-ultra-new record: %v, err: %v",
				record, err)
		}
	})
}

func Test_Sqlite_DB_UpdateErrorRecord(t *testing.T) {
	db := NewSqliteDB("./rows_test.db")
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		errorRecord := &database.ErrorRecord{
			Site: "test",
			Id: 2,
			Error: errors.New("error-2-new")}
		err := db.UpdateErrorRecord(errorRecord)
		
		if err != nil {
			t.Fatalf("DB.UpdateErrorRecord failed err: %v", err)
		}

		query := db.QueryErrorBySiteId("test", 2)
		defer query.Close()
		// query.Next()
		record, err := query.Scan()
		errorRecord = record.(*database.ErrorRecord)
		if err != nil || errorRecord.Site != "test" ||
			errorRecord.Id != 2 || errorRecord.Error.Error() != "error-2-new" {
				t.Fatalf(
					"DB.UpdateErrorRecord does not change record in database record: %v, err: %v",
					errorRecord, err)
			}
	})

	t.Run("fail if input Error Record is nil", func(t *testing.T) {
		err := db.UpdateErrorRecord(nil)
		
		if err == nil {
			t.Fatalf("DB.UpdateErrorRecord success with nil err: %v", err)
		}
	})

	t.Run("pass but not create anything if site id not exist", func(t *testing.T) {
		errorRecord := &database.ErrorRecord{
			Site: "test",
			Id: -1,
			Error: errors.New("error-not-exist-new")}
		err := db.UpdateErrorRecord(errorRecord)
		
		if err != nil {
			t.Fatalf("DB.UpdateErrorRecord failed err: %v", err)
		}

		query := db.QueryErrorBySiteId("test", -1)
		defer query.Close()
		record, err := query.Scan()
		if err == nil {
			t.Fatalf(
				"error of site: test id: -f is created record: %v, err: %v",
				record, err)
		}
	})
}