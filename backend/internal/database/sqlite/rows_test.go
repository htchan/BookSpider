package sqlite

import (
	"testing"
	"os"
	"io"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/internal/utils"
)

func init() {
	source, err := os.Open(os.Getenv("ASSETS_LOCATION") + "/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./rows_test.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func Test_Sqlite_BooksRow_Scan(t *testing.T) {
	db := NewSqliteDB("./rows_test.db")
	defer db.Close()
	query := db.QueryBookBySiteIdHash("test", 1, 100)
	defer query.Close()

	t.Run("success", func(t *testing.T) {
		record, err := query.Scan()
		bookRecord := record.(*database.BookRecord)
		if err != nil || bookRecord.Site != "test" ||
			bookRecord.Id != 1 || bookRecord.HashCode != 100 ||
			bookRecord.Title != "title-1" || bookRecord.WriterId != 1 ||
			bookRecord.Type != "type-1" ||
			bookRecord.UpdateDate != "104" || bookRecord.UpdateChapter != "chapter-1" ||
			bookRecord.Status != database.InProgress {
				t.Fatalf("BookRows.Scan return bookRecord: %v, err: %v", bookRecord, err)
			}
	})

	t.Run("fail after Next return nil", func(t *testing.T) {
		record, err := query.Scan()

		if err == nil {
			t.Fatalf(
				"BookRows.Scan running out of record return record: %v, err: %v",
				record, err)
		}
	})
}

func Test_Sqlite_BooksRow_ScanCurrent(t *testing.T) {
	db := NewSqliteDB("./rows_test.db")
	defer db.Close()
	query := db.QueryBookBySiteIdHash("test", 1, 100)
	defer query.Close()

	t.Run("fail before calling Next", func(t *testing.T) {
		record, err := query.ScanCurrent()

		if err == nil {
			t.Fatalf(
				"BookRows.Scan running out of record return record: %v, err: %v",
				record, err)
		}
	})

	t.Run("success", func(t *testing.T) {
		query.Next()
		record, err := query.ScanCurrent()
		bookRecord := record.(*database.BookRecord)
		if err != nil || bookRecord.Site != "test" ||
			bookRecord.Id != 1 || bookRecord.HashCode != 100 ||
			bookRecord.Title != "title-1" || bookRecord.WriterId != 1 ||
			bookRecord.Type != "type-1" ||
			bookRecord.UpdateDate != "104" || bookRecord.UpdateChapter != "chapter-1" ||
			bookRecord.Status != database.InProgress {
				t.Fatalf("BookRows.Scan return bookRecord: %v, err: %v", bookRecord, err)
			}
	})

	t.Run("fail after Next return nil", func(t *testing.T) {
		query.Next()
		record, err := query.ScanCurrent()

		if err == nil {
			t.Fatalf(
				"BookRows.Scan running out of record return record: %v, err: %v",
				record, err)
		}
	})
}

func Test_Sqlite_BooksRow_Next(t *testing.T) {
	db := NewSqliteDB("./rows_test.db")
	defer db.Close()
	query := db.QueryBookBySiteIdHash("test", 1, 100)
	defer query.Close()

	t.Run("return true when there are record", func(t *testing.T) {
		result := query.Next()
		if result == false {
			t.Fatalf("BookRows return false when it still has record")
		}
	})

	t.Run("return false and clear _rows if there is running out of record", func(t *testing.T) {
		result := query.Next()
		bookQuery := query.(*SqliteBookRows)
		if result == true || bookQuery._rows != nil {
			t.Fatalf("BookRows return true when it has no record result: %v, rows: %v",
				result, bookQuery._rows)
		}
	})

	t.Run("return false if _rows is nil", func(t *testing.T) {
		result := query.Next()
		if result == true {
			t.Fatalf("BookRows return false when _rows is empty already")
		}
	})
}

func Test_Sqlite_BooksRow_Close(t *testing.T) {
	db := NewSqliteDB("./rows_test.db")
	defer db.Close()
	query := db.QueryBookBySiteIdHash("test", 1, 100)

	t.Run("success", func(t *testing.T) {
		err := query.Close()
		if err != nil {
			t.Fatalf("BookRows.Close fail err: %v", err)
		}
	})

	t.Run("fail if keep closing query", func(t *testing.T) {
		err := query.Close()
		if err == nil {
			t.Fatalf("BookRows.Close success in second call err: %v", err)
		}
	})

	t.Run("fail in Scan and Next after Close", func(t *testing.T) {
		result := query.Next()
		if result == true {
			t.Fatalf("BookRows.Next success after close")
		}

		_, err := query.Scan()
		if err == nil {
			t.Fatalf("BookRows.Scan success after Close")
		}
	})
}

func Test_Sqlite_WriterRows_Scan(t *testing.T) {
	db := NewSqliteDB("./rows_test.db")
	defer db.Close()
	query := db.QueryWriterById(1)
	defer query.Close()

	t.Run("success", func(t *testing.T) {
		record, err := query.Scan()
		writerRecord := record.(*database.WriterRecord)
		if err != nil || writerRecord.Id != 1 || writerRecord.Name != "writer-1" {
				t.Fatalf("WriterRows.Scan return writerRecord: %v, err: %v", writerRecord, err)
			}
	})

	t.Run("fail after Next return nil", func(t *testing.T) {
		record, err := query.Scan()

		if err == nil {
			t.Fatalf(
				"WriterRows.Scan running out of record return record: %v, err: %v",
				record, err)
		}
	})
}

func Test_Sqlite_WriterRows_ScanCurrent(t *testing.T) {
	db := NewSqliteDB("./rows_test.db")
	defer db.Close()
	query := db.QueryWriterById(1)
	defer query.Close()

	t.Run("fail before calling Next", func(t *testing.T) {
		record, err := query.ScanCurrent()

		if err == nil {
			t.Fatalf(
				"WriterRows.Scan running out of record return record: %v, err: %v",
				record, err)
		}
	})

	t.Run("success", func(t *testing.T) {
		query.Next()
		record, err := query.ScanCurrent()
		writerRecord := record.(*database.WriterRecord)
		if err != nil || writerRecord.Id != 1 || writerRecord.Name != "writer-1" {
				t.Fatalf("WriterRows.Scan return writerRecord: %v, err: %v", writerRecord, err)
			}
	})

	t.Run("fail after Next return nil", func(t *testing.T) {
		query.Next()
		record, err := query.ScanCurrent()

		if err == nil {
			t.Fatalf(
				"WriterRows.Scan running out of record return record: %v, err: %v",
				record, err)
		}
	})
}

func Test_Sqlite_WriterRows_Next(t *testing.T) {
	db := NewSqliteDB("./rows_test.db")
	defer db.Close()
	query := db.QueryWriterById(1)
	defer query.Close()

	t.Run("return true when there are record", func(t *testing.T) {
		result := query.Next()
		if result == false {
			t.Fatalf("WriterRows return false when it still has record")
		}
	})

	t.Run("return false and clear _rows if there is running out of record", func(t *testing.T) {
		result := query.Next()
		writerQuery := query.(*SqliteWriterRows)
		if result == true || writerQuery._rows != nil {
			t.Fatalf("WriterRows return true when it has no record result: %v, rows: %v",
				result, writerQuery._rows)
		}
	})

	t.Run("return false if _rows is nil", func(t *testing.T) {
		result := query.Next()
		if result == true {
			t.Fatalf("BookRows return false when _rows is empty already")
		}
	})
}

func Test_Sqlite_WriterRows_Close(t *testing.T) {
	db := NewSqliteDB("./rows_test.db")
	defer db.Close()
	query := db.QueryWriterById(1)

	t.Run("success", func(t *testing.T) {
		err := query.Close()
		if err != nil {
			t.Fatalf("WriterRows.Close fail err: %v", err)
		}
	})

	t.Run("fail if keep closing query", func(t *testing.T) {
		err := query.Close()
		if err == nil {
			t.Fatalf("WriterRows.Close success in second call err: %v", err)
		}
	})

	t.Run("fail in Scan and Next after Close", func(t *testing.T) {
		result := query.Next()
		if result == true {
			t.Fatalf("BookRows.Next success after close")
		}

		_, err := query.Scan()
		if err == nil {
			t.Fatalf("BookRows.Scan success after Close")
		}
	})
}

func Test_Sqlite_ErrorRows_Scan(t *testing.T) {
	db := NewSqliteDB("./rows_test.db")
	defer db.Close()
	query := db.QueryErrorBySiteId("test", 2)
	defer query.Close()

	t.Run("success", func(t *testing.T) {
		record, err := query.Scan()
		errorRecord := record.(*database.ErrorRecord)
		if err != nil || errorRecord.Site != "test" || errorRecord.Id != 2 ||
			errorRecord.Error.Error() != "error-2" {
				t.Fatalf("ErrorRows.Scan return errorRecord: %v, err: %v", errorRecord, err)
			}
	})

	t.Run("fail after Next return nil", func(t *testing.T) {
		record, err := query.Scan()

		if err == nil {
			t.Fatalf(
				"ErrorRows.Scan running out of record return record: %v, err: %v",
				record, err)
		}
	})
}

func Test_Sqlite_ErrorRows_ScanCurrent(t *testing.T) {
	db := NewSqliteDB("./rows_test.db")
	defer db.Close()
	query := db.QueryErrorBySiteId("test", 2)
	defer query.Close()

	t.Run("fail before calling Next", func(t *testing.T) {
		record, err := query.ScanCurrent()

		if err == nil {
			t.Fatalf(
				"ErrorRows.Scan running out of record return record: %v, err: %v",
				record, err)
		}
	})

	t.Run("success", func(t *testing.T) {
		query.Next()
		record, err := query.ScanCurrent()
		errorRecord := record.(*database.ErrorRecord)
		if err != nil || errorRecord.Site != "test" || errorRecord.Id != 2 ||
			errorRecord.Error.Error() != "error-2" {
				t.Fatalf("ErrorRows.Scan return errorRecord: %v, err: %v", errorRecord, err)
			}
	})

	t.Run("fail after Next return nil", func(t *testing.T) {
		query.Next()
		record, err := query.ScanCurrent()

		if err == nil {
			t.Fatalf(
				"ErrorRows.Scan running out of record return record: %v, err: %v",
				record, err)
		}
	})
}

func Test_Sqlite_ErrorRows_Next(t *testing.T) {
	db := NewSqliteDB("./rows_test.db")
	defer db.Close()
	query := db.QueryErrorBySiteId("test", 2)
	defer query.Close()

	t.Run("return true when there are record", func(t *testing.T) {
		result := query.Next()
		if result == false {
			t.Fatalf("ErrorRows return false when it still has record")
		}
	})

	t.Run("return false and clear _rows if there is running out of record", func(t *testing.T) {
		result := query.Next()
		errorQuery := query.(*SqliteErrorRows)
		if result == true || errorQuery._rows != nil {
			t.Fatalf("ErrorRows return true when it has no record result: %v, rows: %v",
				result, errorQuery._rows)
		}
	})

	t.Run("return false if _rows is nil", func(t *testing.T) {
		result := query.Next()
		if result == true {
			t.Fatalf("ErrorRows return false when _rows is empty already")
		}
	})
}

func Test_Sqlite_ErrorRows_Close(t *testing.T) {
	db := NewSqliteDB("./rows_test.db")
	defer db.Close()
	query := db.QueryErrorBySiteId("test", 2)

	t.Run("success", func(t *testing.T) {
		err := query.Close()
		if err != nil {
			t.Fatalf("ErrorRows.Close fail err: %v", err)
		}
	})

	t.Run("fail if keep closing query", func(t *testing.T) {
		err := query.Close()
		if err == nil {
			t.Fatalf("ErrorRows.Close success in second call err: %v", err)
		}
	})

	t.Run("fail in Scan and Next after Close", func(t *testing.T) {
		result := query.Next()
		if result == true {
			t.Fatalf("ErrorRows.Next success after close")
		}

		_, err := query.Scan()
		if err == nil {
			t.Fatalf("BookRows.Scan success after Close")
		}
	})
}

func Test_Sqlite_BookRows_interface(t *testing.T) {
	db := NewSqliteDB("./rows_test.db")
	defer db.Close()
	query := db.QueryBookBySiteIdHash("test", 1, 100)
	defer query.Close()

	t.Run("success", func(t *testing.T) {
		switch v := query.(type) {
		case *SqliteBookRows:
		case interface{}:
			t.Fatalf("query from QueryBookBySiteIdHash is not *SqliteBookRows, but %v", v)
		}
	})
}

func Test_Sqlite_WriterRows_interface(t *testing.T) {
	db := NewSqliteDB("./rows_test.db")
	defer db.Close()
	query := db.QueryWriterById(1)
	defer query.Close()

	t.Run("success", func(t *testing.T) {
		switch v := query.(type) {
		case *SqliteWriterRows:
		case interface{}:
			t.Fatalf("query from QueryWriterById is not *SqliteWriterRows, but %v", v)
		}
	})
}

func Test_Sqlite_ErrorRows_interface(t *testing.T) {
	db := NewSqliteDB("./rows_test.db")
	defer db.Close()
	query := db.QueryErrorBySiteId("test", 2)
	defer query.Close()

	t.Run("success", func(t *testing.T) {
		switch v := query.(type) {
		case *SqliteErrorRows:
		case interface{}:
			t.Fatalf("query from QueryErrorBySiteId is not *SqliteErrorRows, but %v", v)
		}
	})
}