package sqlite

import (
	"os"
	"io"
	"testing"
	"github.com/htchan/BookSpider/internal/utils"
)

func initDbBackupTest() {
	source, err := os.Open(os.Getenv("ASSETS_LOCATION") + "/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./backup_test.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func cleanupDbBackupTest() {
	os.Remove("./backup_test.db")
	os.Remove(os.Getenv("ASSETS_LOCATION") + "/test-data/" + "db-backup.sql")
}

func TestSqlite_DB_Backup(t *testing.T) {
	db := NewSqliteDB("./backup_test.db")
	defer db.Close()

	t.Run("success", func(t *testing.T) {
		db.Backup(os.Getenv("ASSETS_LOCATION") + "/test-data/", "db-backup.sql")
	})

	b, err := os.ReadFile(os.Getenv("ASSETS_LOCATION") + "/test-data/db-backup.sql")
	utils.CheckError(err)
	reference, err := os.ReadFile(os.Getenv("ASSETS_LOCATION") + "/test-data/backup-reference.sql")
	utils.CheckError(err)

	if string(b) != string(reference){
		t.Fatalf("db backup save such content: %v", string(b))
	}
}