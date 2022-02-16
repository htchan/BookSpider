package sites

import (
	"testing"
	"io"
	"os"
	"path/filepath"
	"time"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/pkg/configs"
	"github.com/htchan/BookSpider/pkg/flags"
)

func initBackupTest() {
	source, err := os.Open(os.Getenv("ASSETS_LOCATION") + "/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./backup.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func cleanupBackupTest() {
	os.Remove("./backup.db")
}

var backupConfig = configs.LoadConfigYaml(os.Getenv("ASSETS_LOCATION") + "/test-data/config.yml").SiteConfigs["test"]

func Test_Sites_Site_Backup(t *testing.T) {
	backupConfig.DatabaseLocation = "./query.db"
	site := NewSite("test", backupConfig)
	site.OpenDatabase()
	defer site.CloseDatabase()

	t.Run("success for specific site", func(t *testing.T) {
		flagSite := "test"
		backupLocation := filepath.Join(os.Getenv("ASSETS_LOCATION") + site.config.BackupDirectory, time.Now().Format("2006-01-02"), "test.sql")

		err := site.Backup(flags.Flags{Site: &flagSite})

		if err != nil {
			t.Fatalf("site Backup return error: %v", err)
		}
		
		b, err := os.ReadFile(backupLocation)
		utils.CheckError(err)
		reference, err := os.ReadFile(os.Getenv("ASSETS_LOCATION") + "/test-data/backup-reference.sql")
		utils.CheckError(err)

		if string(b) != string(reference){
			t.Fatalf("site backup save such content: %v", string(b))
		}
		os.Remove(backupLocation)
		os.Remove(filepath.Dir(backupLocation))
	})

	t.Run("success for all site", func(t *testing.T) {
		backupLocation := filepath.Join(os.Getenv("ASSETS_LOCATION") + site.config.BackupDirectory, time.Now().Format("2006-01-02"), "test.sql")

		err := site.Backup(flags.Flags{})

		if err != nil {
			t.Fatalf("site Backup return error: %v", err)
		}
		
		b, err := os.ReadFile(backupLocation)
		utils.CheckError(err)
		reference, err := os.ReadFile(os.Getenv("ASSETS_LOCATION") + "/test-data/backup-reference.sql")
		utils.CheckError(err)

		if string(b) != string(reference){
			t.Fatalf("site backup save such content: %v", string(b))
		}
		os.Remove(backupLocation)
		os.Remove(filepath.Dir(backupLocation))
	})

	t.Run("skip if not", func(t *testing.T) {
		flagSite := "not-test"
		backupLocation := filepath.Join(site.config.BackupDirectory, time.Now().Format("2006-01-02"), "test.sql")
		err := site.Backup(flags.Flags{Site: &flagSite})

		if err != nil {
			t.Fatalf("site Backup return error: %v", err)
		}
		
		if utils.Exists(backupLocation) {
			t.Fatalf("site backup does not skip for flag site: not-test")
		}
	})
}