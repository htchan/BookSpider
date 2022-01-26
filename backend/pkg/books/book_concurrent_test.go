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

func init() {
	source, err := os.Open(os.Getenv("ASSETS_LOCATION") + "/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./book_concurrent_test.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func test_concurrent_create(db database.DB, config *configs.BookConfig, n, offset int) func(t *testing.T) {
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
					t.Fatalf("concurrent save book failed at %v times trial, book: %v", i, book.bookRecord)
				}
			}(i)
		}
		wg.Wait()
		t.Logf("create %v finish", n)
	}
}

func Test_Books_Book_Concurrent_Save(t *testing.T) {
	db := sqlite.NewSqliteDB("./book_concurrent_test.db")
	defer db.Close()
	t.Run("success with 10 multi thread", test_concurrent_create(db, config, 10, 10))
	runtime.GC()
	t.Run("success with 100 multi thread", test_concurrent_create(db, config, 100, 100))
	runtime.GC()
	t.Run("success with 1000 multi thread", test_concurrent_create(db, config, 1000, 1000))
	runtime.GC()
	t.Run("success with 10000 multi thread", test_concurrent_create(db, config, 10000, 10000))
}

func test_concurrent_load_book(db database.DB, config *configs.BookConfig, n, offset int) func(t *testing.T) {
	return func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				book := LoadBook(db, "test", i + offset, -1, config)
				if site, id, _ := book.GetInfo();
					site != "test" || id != i + offset {
					t.Fatalf("concurrent load book failed at %v times trial, book: %v", i, book.bookRecord)
				}
			}(i)
		}
		wg.Wait()
		t.Logf("create %v finish", n)
	}
}

func Test_Books_Book_Concurrent_LoadBook(t *testing.T) {
	db := sqlite.NewSqliteDB("./book_concurrent_test.db")
	defer db.Close()
	t.Run("success with 10 multi thread", test_concurrent_load_book(db, config, 10, 10))
	runtime.GC()
	t.Run("success with 100 multi thread", test_concurrent_load_book(db, config, 100, 100))
	runtime.GC()
	t.Run("success with 1000 multi thread", test_concurrent_load_book(db, config, 1000, 1000))
	runtime.GC()
	t.Run("success with 10000 multi thread", test_concurrent_load_book(db, config, 10000, 10000))
}

func test_concurrent_load_book_by_record(db database.DB, config *configs.BookConfig, n, offset int) func(t *testing.T) {
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
					t.Fatalf("concurrent load book failed at %v times trial, book: %v", i, book.bookRecord)
				}
			}(i)
		}
		wg.Wait()
		t.Logf("create %v finish", n)
	}
}

func Test_Books_Book_Concurrent_LoadBookByRecord(t *testing.T) {
	db := sqlite.NewSqliteDB("./book_concurrent_test.db")
	defer db.Close()
	t.Run("success with 10 multi thread", test_concurrent_load_book_by_record(db, config, 10, 10))
	runtime.GC()
	t.Run("success with 100 multi thread", test_concurrent_load_book_by_record(db, config, 100, 100))
	runtime.GC()
	t.Run("success with 1000 multi thread", test_concurrent_load_book_by_record(db, config, 1000, 1000))
	runtime.GC()
	t.Run("success with 10000 multi thread", test_concurrent_load_book_by_record(db, config, 10000, 10000))
}



func test_concurrent_update_book(db database.DB, config *configs.BookConfig, n, offset int) func(t *testing.T) {
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
					t.Fatalf("concurrent update book failed at %v times trial, book: %v", i, book.bookRecord)
				}
			}(i)
		}
		wg.Wait()
		t.Logf("create %v finish", n)
	}
}

func Test_Books_Book_Concurrent_Update(t *testing.T) {
	db := sqlite.NewSqliteDB("./book_concurrent_test.db")
	defer db.Close()
	t.Run("success with 10 multi thread", test_concurrent_update_book(db, config, 10, 10))
	runtime.GC()
	t.Run("success with 100 multi thread", test_concurrent_update_book(db, config, 100, 100))
	runtime.GC()
	t.Run("success with 1000 multi thread", test_concurrent_update_book(db, config, 1000, 1000))
	runtime.GC()
	t.Run("success with 10000 multi thread", test_concurrent_update_book(db, config, 10000, 10000))
}
