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
	destination, err := os.Create("./delete_test.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func Test_Sqlite_DB_DeleteBookRecord(t *testing.T) {
	db := NewSqliteDB("./delete_test.db")
	defer db.Close()
	
	t.Run("fail", func(t *testing.T) {
		record := []database.BookRecord {
			database.BookRecord{
				Site: "test",
				Id: 1,
				HashCode: 100,
			},
		}

		err := db.DeleteBookRecord(record)
		rows := db.QueryBookBySiteIdHash("test", 1, 100)
		defer rows.Close()

		if err == nil || !rows.Next() {
			t.Fatalf("DB success in deleting book record")
		}
	})
}

func Test_Sqlite_DB_DeleteWriterRecord(t *testing.T) {
	db := NewSqliteDB("./delete_test.db")
	defer db.Close()
	
	t.Run("fail", func(t *testing.T) {
		record := []database.WriterRecord {
			database.WriterRecord{
				Id: 1,
			},
		}

		err := db.DeleteWriterRecord(record)
		rows := db.QueryWriterById(1)
		defer rows.Close()

		if err == nil || !rows.Next() {
			t.Fatalf("DB success in deleting book record")
		}
	})
}

func Test_Sqlite_DB_DeleteErrorRecord(t *testing.T) {
	db := NewSqliteDB("./delete_test.db")
	defer db.Close()
	
	t.Run("success only providing site and id", func(t *testing.T) {
		record := []database.ErrorRecord {
			database.ErrorRecord{
				Site: "test",
				Id: 1,
			},
		}

		err := db.DeleteErrorRecord(record)
		rows := db.QueryErrorBySiteId("test", 1)
		defer rows.Close()

		if err != nil || rows.Next() {
			t.Fatalf("DB fail in deleting error record err: %v", err)
		}
	})

	t.Run("success even site id not exist", func(t *testing.T) {
		record := []database.ErrorRecord {
			database.ErrorRecord{
				Site: "not-exist-site",
				Id: 1,
			},
		}

		err := db.DeleteErrorRecord(record)

		if err != nil {
			t.Fatalf("DB fail in deleting not exist error record err: %v", err)
		}
	})
}