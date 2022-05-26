package books

import (
	"os"
	"io"
	"testing"
	"sync"
	"runtime"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/pkg/configs"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/internal/database/sqlite"
)

func initConcurrentTest() {
	source, err := os.Open(os.Getenv("ASSETS_LOCATION") + "/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./book_concurrent_test.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func cleanupConcurrentTest() {
	os.Remove("./book_concurrent_test.db")
}

func test_concurrent_create(db database.DB, config *configs.SourceConfig, n, offset int) func(t *testing.T) {
	return func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				book := NewBook("test", i + offset, -1, config)
				book.SetStatus(database.InProgress)
				result := book.Save(db)
				if !result {
					t.Errorf("concurrent save book failed at %v times trial, book: %v", i, book.bookRecord)
				}
			}(i)
		}
		wg.Wait()
		db.Commit()
		for i := 0; i < n; i++ {
			book := LoadBook(db, "test", i + offset, -1, config)
			if book == nil {
				t.Errorf("book.Save does not create test-%v", i + offset)
			}
			if book.bookRecord.Site != "test" || book.bookRecord.Id != i + offset ||
				book.GetStatus() != database.InProgress || book.GetTitle() != "" ||
				book.bookRecord.WriterId != 0 {
					t.Errorf("book.Save for test-%v save wrong book data: %v", i + offset, book)
				}
		}
	}
}

func test_concurrent_load_book(db database.DB, config *configs.SourceConfig, n, offset int) func(t *testing.T) {
	return func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				book := LoadBook(db, "test", i + offset, -1, config)
				if book == nil {
					t.Errorf("fail to load test-%v from db", i + offset)
				}
				if site, id, _ := book.GetInfo();
					site != "test" || id != i + offset {
					t.Errorf("concurrent load book failed at %v times trial, book: %v", i, book.bookRecord)
				}
			}(i)
		}
		wg.Wait()
	}
}

func test_concurrent_load_book_by_record(db database.DB, config *configs.SourceConfig, n, offset int) func(t *testing.T) {
	return func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				rows := db.QueryBookBySiteIdHash("test", i + offset, -1)
				record, err := rows.Scan()
				utils.CheckError(err)
				rows.Close()
				book := LoadBookByRecord(db, record.(*database.BookRecord), config)
				if site, id, _ := book.GetInfo();
					site != "test" || id != i + offset {
					t.Errorf("concurrent load book failed at %v times trial, book: %v", i, book.bookRecord)
				}
			}(i)
		}
		wg.Wait()
		t.Logf("create %v finish", n)
	}
}

func test_concurrent_update_book(db database.DB, config *configs.SourceConfig, n, offset int) func(t *testing.T) {
	return func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				book := LoadBook(db, "test", i + offset, -1, config)
				book.SetTitle(book.GetTitle() + "-new")
				result := book.Save(db)
				if !result {
					t.Errorf("concurrent update book failed at %v times trial, book: %v", i, book.bookRecord)
				}
			}(i)
		}
		wg.Wait()
		db.Commit()
		for i := 0; i < n; i++ {
			rows := db.QueryBookBySiteIdHash("test", i + offset, -1)
			record, err := rows.Scan()
			if err != nil {
				t.Errorf("book.Save does not create test-%v, err: %v", i + offset, err)
			}
			rows.Close()
			actualRecord := record.(*database.BookRecord)
			if actualRecord.Site != "test" || actualRecord.Id != i + offset ||
				actualRecord.Status != database.InProgress || actualRecord.Title != "-new" ||
				actualRecord.WriterId != 0 {
					t.Errorf("book.Save for test-%v save wrong book data: %v", i + offset, actualRecord)
				}
		}
	}
}

func TestBooks_Book_Concurrent(t *testing.T) {
	db := sqlite.NewSqliteDB("./book_concurrent_test.db", 100)
	defer db.Close()

	t.Run("func Save", func(t *testing.T) {
		t.Run("success with 10 threads", test_concurrent_create(db, config, 10, 10))
		runtime.GC()
		t.Run("success with 100 threads", test_concurrent_create(db, config, 100, 100))
		runtime.GC()
		t.Run("success with 1000 threads", test_concurrent_create(db, config, 1000, 1000))
		runtime.GC()
		t.Run("success with 10000 threads", test_concurrent_create(db, config, 10000, 10000))
	})

	t.Run("func LoadBook", func(t *testing.T) {
		t.Run("success with 10 threads", test_concurrent_load_book(db, config, 10, 10))
		runtime.GC()
		t.Run("success with 100 threads", test_concurrent_load_book(db, config, 100, 100))
		runtime.GC()
		t.Run("success with 1000 threads", test_concurrent_load_book(db, config, 1000, 1000))
		runtime.GC()
		t.Run("success with 10000 threads", test_concurrent_load_book(db, config, 10000, 10000))
	})

	t.Run("func LoadBookByRecord", func(t *testing.T) {
		t.Run("success with 10 threads", test_concurrent_load_book_by_record(db, config, 10, 10))
		runtime.GC()
		t.Run("success with 100 threads", test_concurrent_load_book_by_record(db, config, 100, 100))
		runtime.GC()
		t.Run("success with 1000 threads", test_concurrent_load_book_by_record(db, config, 1000, 1000))
		runtime.GC()
		t.Run("success with 10000 threads", test_concurrent_load_book_by_record(db, config, 10000, 10000))
	})

	t.Run("func Update", func(t *testing.T) {
		t.Run("success with 10 threads", test_concurrent_update_book(db, config, 10, 10))
		runtime.GC()
		t.Run("success with 100 threads", test_concurrent_update_book(db, config, 100, 100))
		runtime.GC()
		t.Run("success with 1000 threads", test_concurrent_update_book(db, config, 1000, 1000))
		runtime.GC()
		t.Run("success with 10000 threads", test_concurrent_update_book(db, config, 10000, 10000))
	})
}