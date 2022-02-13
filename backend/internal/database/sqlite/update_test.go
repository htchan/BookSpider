package sqlite

import (
	"testing"
	"os"
	"io"
	"errors"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/internal/utils"
)

func initDbUpdateTest() {
	source, err := os.Open(os.Getenv("ASSETS_LOCATION") + "/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./update_test.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func cleanupDbUpdateTest() {
	os.Remove("./update_test.db")
}

func Test_Sqlite_DB_UpdateBookRecord(t *testing.T) {
	db := NewSqliteDB("./rows_test.db")
	defer db.Close()

	writerRecord := &database.WriterRecord{
		Id: 1,
		Name: "writer-1",
	}

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
		err := db.UpdateBookRecord(bookRecord, writerRecord)
		
		if err != nil {
			t.Fatalf("DB.UpdateBookRecord failed err: %v", err)
		}
		if db.statementCount != 1 ||
			db.statements[0] != BookUpdateStatement(bookRecord, "") {
				t.Fatalf("DB.UpdateBookRecord does not add record to statement - count: %v, statements: %v", db.statementCount, db.statements)
		}
	})

	t.Run("fail if input Book Record is nil", func(t *testing.T) {
		err := db.UpdateBookRecord(nil, nil)
		
		if err == nil || db.statementCount != 1 {
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
		err := db.UpdateBookRecord(bookRecord, writerRecord)
		
		if err != nil {
			t.Fatalf("DB.UpdateBookRecord failed err: %v", err)
		}
		if db.statementCount != 2 ||
			db.statements[1] != BookUpdateStatement(bookRecord, "") {
				t.Fatalf("DB.UpdateBookRecord does not add record to statement - count: %v, statements: %v", db.statementCount, db.statements)
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
		if db.statementCount != 1 ||
			db.statements[0] != ErrorUpdateStatement(errorRecord) {
				t.Fatalf("DB.UpdateBookRecord does not add record to statement - count: %v, statements: %v", db.statementCount, db.statements)
		}
	})

	t.Run("fail if input Error Record is nil", func(t *testing.T) {
		err := db.UpdateErrorRecord(nil)
		
		if err == nil || db.statementCount != 1 {
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
		if db.statementCount != 2 ||
			db.statements[1] != ErrorUpdateStatement(errorRecord) {
				t.Fatalf("DB.UpdateErrorRecord does not add record to statement - count: %v, statements: %v", db.statementCount, db.statements)
		}
	})
}