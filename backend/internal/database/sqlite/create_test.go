package sqlite

import (
	"os"
	"io"
	"testing"
	"errors"
	"sync"
	"runtime"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
)

func initDbCreateTest() {
	source, err := os.Open(os.Getenv("ASSETS_LOCATION") + "/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./create_test.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func cleanupDbCreateTest() {
	os.Remove("./create_test.db")
}

func TestSqlite_DB_CreateBookRecord(t *testing.T) {
	db := NewSqliteDB("./create_test.db")
	defer db.Close()
	writerRecord := &database.WriterRecord{
		Id: 1,
		Name: "writer-1",
	}

	t.Run("success", func(t *testing.T) {
		bookRecord := &database.BookRecord{
			Site: "test",
			Id: 10,
			HashCode: 10,
			Title: "title-10",
			WriterId: 1,
			Type: "type-10",
			UpdateDate: "10",
			UpdateChapter: "chapter-10",
			Status: database.InProgress,
		}

		err := db.CreateBookRecord(bookRecord, writerRecord)

		if err != nil {
			t.Fatalf("DB cannot create book record - err: %v", err)
		}
		if db.statementCount != 1 || db.statements[0] != BookInsertStatement(bookRecord, writerRecord.Name) {
			t.Fatalf("DB.CreateBookRecord does not add record to statement - count: %v, statements: %v", db.statementCount, db.statements)
		}
	})

	t.Run("success with negative writer id", func(t *testing.T) {
		bookRecord := &database.BookRecord{
			Site: "test",
			Id: 10,
			HashCode: 10,
			Title: "title-10",
			WriterId: -1,
			Type: "type-10",
			UpdateDate: "10",
			UpdateChapter: "chapter-10",
			Status: database.InProgress,
		}

		err := db.CreateBookRecord(bookRecord, writerRecord)

		if err != nil {
			t.Fatalf("DB cannot create book record - err: %v", err)
		}
		if db.statementCount != 2 || db.statements[1] != BookInsertStatement(bookRecord, writerRecord.Name) {
			t.Fatalf("DB.CreateBookRecord does not add record to statement - count: %v, statements: %v", db.statementCount, db.statements)
		}
	})

	t.Run("success even creating with existing site id", func(t *testing.T) {
		bookRecord := &database.BookRecord{
			Site: "test",
			Id: 10,
			HashCode: 10,
			Title: "title-10",
			WriterId: 1,
			Type: "type-10",
			UpdateDate: "10",
			UpdateChapter: "chapter-10",
			Status: database.InProgress,
		}

		err := db.CreateBookRecord(bookRecord, writerRecord)

		if err != nil {
			t.Fatalf("DB cannot create book record - err: %v", err)
		}
		if db.statementCount != 3 || db.statements[2] != BookInsertStatement(bookRecord, writerRecord.Name) {
			t.Fatalf("DB.CreateBookRecord does not add record to statement - count: %v, statements: %v", db.statementCount, db.statements)
		}
	})
}

func TestSqlite_DB_CreateWriterRecord(t *testing.T) {
	db := NewSqliteDB("./create_test.db")
	defer db.Close()

	t.Run("success with updated id", func(t *testing.T) {
		writerRecord := &database.WriterRecord{
			Id: -1,
			Name: "writer-10",
		}

		err := db.CreateWriterRecord(writerRecord)
		
		if err != nil {
			t.Fatalf("DB cannot create writer record - err: %v", err)
		}
		if db.statementCount != 1 || db.statements[0] != WriterInsertStatement(writerRecord) {
			t.Fatalf("DB.CreateWriterRecord does not add record to statement - count: %v, statements: %v", db.statementCount, db.statements)
		}
	})

	t.Run("success with provided id", func(t *testing.T) {
		writerRecord := &database.WriterRecord{
			Id: 1,
			Name: "writer-10",
		}

		err := db.CreateWriterRecord(writerRecord)

		if err != nil {
			t.Fatalf("DB cannot create writer record - err: %v", err)
		}
		if db.statementCount != 2 || db.statements[1] != WriterInsertStatement(writerRecord) {
			t.Fatalf("DB.CreateWriterRecord does not add record to statement - count: %v, statements: %v", db.statementCount, db.statements)
		}
	})

	t.Run("success even creating with existing id", func(t *testing.T) {
		writerRecord := &database.WriterRecord{
			Id: 10,
			Name: "writer-10",
		}

		err := db.CreateWriterRecord(writerRecord)

		if err != nil {
			t.Fatalf("DB cannot create writer record - err: %v", err)
		}
		if db.statementCount != 3 || db.statements[2] != WriterInsertStatement(writerRecord) {
			t.Fatalf("DB.CreateWriterRecord does not add record to statement - count: %v, statements: %v", db.statementCount, db.statements)
		}
	})
}

func TestSqlite_DB_CreateErrorRecord(t *testing.T) {
	db := NewSqliteDB("./create_test.db")
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		errorRecord := &database.ErrorRecord{
			Site: "test",
			Id: 10,
			Error: errors.New("error-10"),
		}

		err := db.CreateErrorRecord(errorRecord)

		if err != nil {
			t.Fatalf("DB cannot create error record - err: %v", err)
		}
		if db.statementCount != 1 || db.statements[0] != ErrorInsertStatement(errorRecord) {
			t.Fatalf("DB.CreateErrorRecord does not add record to statement - count: %v, statements: %v", db.statementCount, db.statements)
		}
	})

	t.Run("fail if create with existing site id", func(t *testing.T) {
		errorRecord := &database.ErrorRecord{
			Site: "test",
			Id: 10,
			Error: errors.New("error-10"),
		}

		err := db.CreateErrorRecord(errorRecord)

		if err != nil {
			t.Fatalf("DB cannot create error record - err: %v", err)
		}
		if db.statementCount != 2 || db.statements[0] != ErrorInsertStatement(errorRecord) {
			t.Fatalf("DB.CreateErrorRecord does not add record to statement - count: %v, statements: %v", db.statementCount, db.statements)
		}
	})
}

func test_concurrent_create(db database.DB, n int, offset int) func(t *testing.T) {
	return func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				record := &database.ErrorRecord{
					Site: "test",
					Id: i + offset,
					Error: errors.New("test")}
				err := db.CreateErrorRecord(record)
				if err != nil {
					t.Fatalf("CreateErrorRecord(\"test\", %v, \"\") fail err: %v", i + offset, err)
				}
			}(i)
		}
		wg.Wait()

		if db.(*SqliteDB).statementCount != n {
			t.Fatalf("statement count not match: statement count: %v, n: %v", db.(*SqliteDB).statementCount, n)
		}

		for i := 0; i < n; i++ {
			if db.(*SqliteDB).statements[i][:6] != "insert" {
				t.Fatalf("statement missing at %v: %v", i, db.(*SqliteDB).statements[i])
			}
		}
	}
}

func TestSqlite_DB_ConcrrentCreate(t *testing.T) {
	db := NewSqliteDB("./create_test.db")
	defer db.Close()

	runtime.GOMAXPROCS(runtime.NumCPU())

	t.Run("10 create concurrent", test_concurrent_create(db, 10, 20))
	db.Commit()
	runtime.GC()
	t.Run("100 create concurrent", test_concurrent_create(db, 100, 30))
	db.Commit()
	runtime.GC()
	t.Run("1000 create concurrent", test_concurrent_create(db, 1000, 130))
	db.Commit()
	runtime.GC()
	t.Run("10000 create concurrent", test_concurrent_create(db, 10000, 1130))
}