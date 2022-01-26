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

func init() {
	source, err := os.Open(os.Getenv("ASSETS_LOCATION") + "/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./create_test.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func Test_Sqlite_DB_CreateBookRecord(t *testing.T) {
	db := NewSqliteDB("./create_test.db")
	defer db.Close()

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

		err := db.CreateBookRecord(bookRecord)
		rows := db.QueryBookBySiteIdHash("test", 10, 10)
		defer rows.Close()

		if err != nil || !rows.Next() {
			t.Fatalf("DB cannot create book record: %v", bookRecord)
		}
	})

	t.Run("fail if create with existing site id", func(t *testing.T) {
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

		err := db.CreateBookRecord(bookRecord)

		if err == nil {
			t.Fatalf("DB created already existed book record: %v", bookRecord)
		}
	})
}

func Test_Sqlite_DB_CreateWriterRecord(t *testing.T) {
	db := NewSqliteDB("./create_test.db")
	defer db.Close()

	t.Run("success with updated id", func(t *testing.T) {
		writerRecord := &database.WriterRecord{
			Id: -1,
			Name: "writer-10",
		}

		err := db.CreateWriterRecord(writerRecord)
		rows := db.QueryWriterById(writerRecord.Id)
		defer rows.Close()

		if err != nil || writerRecord.Id == -1 || !rows.Next() {
			t.Fatalf("DB cannot create writerRecord\nrecord: %v\nerror: %v",
			writerRecord, err)
		}
	})

	t.Run("success with provided id", func(t *testing.T) {
		writerRecord := &database.WriterRecord{
			Id: 1,
			Name: "writer-10",
		}

		err := db.CreateWriterRecord(writerRecord)

		if err == nil {
			t.Fatalf("DB create writerRecord with existing id\nrecord: %v\nerror: %v",
			writerRecord, err)
		}
	})

	t.Run("fail if create with existing id", func(t *testing.T) {
		writerRecord := &database.WriterRecord{
			Id: 10,
			Name: "writer-10",
		}

		err := db.CreateWriterRecord(writerRecord)

		if err == nil {
			t.Fatalf("DB created already existed book record: %v", writerRecord)
		}
	})
}

func Test_Sqlite_DB_CreateErrorRecord(t *testing.T) {
	db := NewSqliteDB("./create_test.db")
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		errorRecord := &database.ErrorRecord{
			Site: "test",
			Id: 10,
			Error: errors.New("error-10"),
		}

		err := db.CreateErrorRecord(errorRecord)
		rows := db.QueryErrorBySiteId("test", 10)
		defer rows.Close()

		if err != nil || !rows.Next() {
			t.Fatalf("DB cannot create error record\nrecord: %v\nerror: %v",
			errorRecord, err)
		}
	})

	t.Run("fail if create with existing site id", func(t *testing.T) {
		errorRecord := &database.ErrorRecord{
			Site: "test",
			Id: 10,
			Error: errors.New("error-10"),
		}

		err := db.CreateErrorRecord(errorRecord)

		if err == nil {
			t.Fatalf("DB create error Record with existing site id\nrecord: %v\nerror: %v",
			errorRecord, err)
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

		for i := 0; i < n; i++ {
			query := db.QueryErrorBySiteId("test", i + offset)
			if !query.Next() {
				t.Fatalf("created record (%v) cannot be query in database", i + offset)
			}
			query.Close()
		}
	}
}

func Test_Sqlite_DB_ConcrrentCreate(t *testing.T) {
	db := NewSqliteDB("./create_test.db")
	defer db.Close()

	runtime.GOMAXPROCS(runtime.NumCPU())

	t.Run("10 create concurrent", test_concurrent_create(db, 10, 20))
	runtime.GC()
	t.Run("100 create concurrent", test_concurrent_create(db, 100, 30))
	runtime.GC()
	t.Run("1000 create concurrent", test_concurrent_create(db, 1000, 130))
	runtime.GC()
	t.Run("10000 create concurrent", test_concurrent_create(db, 10000, 1130))
}