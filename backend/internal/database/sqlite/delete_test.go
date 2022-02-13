package sqlite

import (
	"testing"
	"os"
	"io"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/internal/utils"
)

func initDbDeleteTest() {
	source, err := os.Open(os.Getenv("ASSETS_LOCATION") + "/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./delete_test.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func cleanupDbDeleteTest() {
	os.Remove("./delete_test.db")
}

func TestSqlite_DB_DeleteBookRecords(t *testing.T) {
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

		err := db.DeleteBookRecords(record)
		
		if err == nil || db.statementCount != 0 {
			t.Fatalf("DB.DeleteWriterRecord adds statement to db: count: %v, statement: %v", db.statementCount, db.statements)
		}
	})
}

func TestSqlite_DB_DeleteWriterRecords(t *testing.T) {
	db := NewSqliteDB("./delete_test.db")
	defer db.Close()
	
	t.Run("fail", func(t *testing.T) {
		record := []database.WriterRecord {
			database.WriterRecord{
				Id: 1,
			},
		}

		err := db.DeleteWriterRecords(record)

		if err == nil || db.statementCount != 0 {
			t.Fatalf("DB.DeleteWriterRecord adds statement to db: count: %v, statement: %v", db.statementCount, db.statements)
		}
	})
}

func TestSqlite_DB_DeleteErrorRecord(t *testing.T) {
	db := NewSqliteDB("./delete_test.db")
	defer db.Close()
	
	t.Run("success only providing site and id", func(t *testing.T) {
		records := []database.ErrorRecord {
			database.ErrorRecord{
				Site: "test",
				Id: 1,
			},
		}

		err := db.DeleteErrorRecords(records)
		
		if err != nil || db.statementCount != 1 || db.statements[0] != ErrorDeleteStatement(&records[0]) {
			t.Fatalf("DB.DeleteWriterRecord adds statement to db: count: %v, statement: %v", db.statementCount, db.statements)
		}
	})

	t.Run("success even site id not exist", func(t *testing.T) {
		records := []database.ErrorRecord {
			database.ErrorRecord{
				Site: "not-exist-site",
				Id: 1,
			},
		}

		err := db.DeleteErrorRecords(records)
		
		if err != nil || db.statementCount != 2 || db.statements[1] != ErrorDeleteStatement(&records[0]) {
			t.Fatalf("DB.DeleteWriterRecord adds statement to db: count: %v, statement: %v", db.statementCount, db.statements)
		}
	})
}