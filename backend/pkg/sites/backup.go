package sites

import (
	"github.com/htchan/BookSpider/pkg/flags"
	"errors"
	"path/filepath"
	"time"
	"os"
)

func (site *Site) Backup(args flags.Flags) (err error) {
	if !args.Valid() || args.IsBook() {
		err = errors.New("invalid arguments")
		return
	}
	if args.IsEverything() || (args.IsSite() && *args.Site == site.Name) {
		backupDirectory := filepath.Join(site.config.BackupDirectory, time.Now().Format("2006-01-02"))
		return site.database.Backup(os.Getenv("ASSETS_LOCATION") + backupDirectory, site.Name + ".sql")
	}
	return nil
}