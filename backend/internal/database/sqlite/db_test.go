package sqlite

import (
	"os"
	"io"
	"testing"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
)

func init() {
	source, err := os.Open(os.Getenv("ASSETS_LOCATION") +  "/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./db_test.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func Test_Sqlite_DB_Constructor(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := NewSqliteDB("./db_test.db")
		defer db.Close()

		if db._db == nil {
			t.Fatalf("DB construct failed")
		}
	})
}

func Test_Sqlite_DB_Summary(t *testing.T) {
	db := NewSqliteDB("./db_test.db")
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		actual := db.Summary("test")
		
		if actual.BookCount != 6 || actual.WriterCount != 3 ||
			actual.ErrorCount != 3 || actual.UniqueBookCount != 5 ||
			actual.MaxBookId != 5 || actual.LatestSuccessId != 3 ||
			actual.StatusCount[database.Error] != 3 ||
			actual.StatusCount[database.InProgress] != 1 ||
			actual.StatusCount[database.End] != 1 ||
			actual.StatusCount[database.Download] != 1 {
			t.Fatalf(
				"DB Summary() failed\nactual: %v",
				actual)
		}
	})
}

func Test_Sqlite_DB_Close(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := NewSqliteDB("./db_test.db")
		db.Close()
		if db._db != nil {
			t.Fatalf("DB Close() failed")
		}
	})
}

func Test_Sqlite_DB_interface(t *testing.T) {
	var db database.DB
	db = NewSqliteDB("./db_test.db")
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		switch db.(type) {
		case *SqliteDB:
		case interface{}:
			t.Fatalf("NewSqliteDB cannot use db interface")
		}
	})

	t.Run("success in query", func(t *testing.T) {
		query := db.QueryBookBySiteIdHash("test", 1, 100)
		switch query.(type) {
		case *SqliteBookRows:
		case interface{}:
			t.Fatalf("NewSqliteDB query does not return *BookRecord")
		}
	})
}