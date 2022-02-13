package sites

import (
	"github.com/htchan/BookSpider/pkg/configs"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/internal/database/sqlite"
	"golang.org/x/sync/semaphore"
	"strings"
	"errors"
)

type Site struct {
	Name string
	database database.DB
	config *configs.SiteConfig
	semaphore *semaphore.Weighted
}

func NewSite(name string, config *configs.SiteConfig) (site *Site) {
	site = new(Site)
	site.Name = name
	site.database = nil
	site.config = config
	site.semaphore = semaphore.NewWeighted(int64(config.ThreadsCount))
	return
}

func (site *Site)OpenDatabase() (err error) {
	switch strings.ToUpper(site.config.DatabaseEngine) {
	case "SQLITE3":
		site.database = sqlite.NewSqliteDB(site.config.DatabaseLocation)
	default:
		err = errors.New("invalid database engine")
	}
	return
}

func (site *Site)Map() (result map[string]interface{}) {
	defer utils.Recover(func() {
		result = map[string]interface{} {
			"name": site.Name,
			"bookCount": 0,
			"uniqueBookCount": 0,
			"writerCount": 0,
			"errorCount": 0,
			"maxBookId": 0,
			"latestSuccessBookId": 0,
			"statusCount": map[database.StatusCode]int {
				database.Error: 0,
				database.InProgress: 0,
				database.End: 0,
				database.Download: 0,
			},
		}
	})
	err := site.OpenDatabase()
	utils.CheckError(err)
	defer site.database.Close()
	summary := site.database.Summary(site.Name)
	result = map[string]interface{} {
		"name": site.Name,
		"bookCount": summary.BookCount,
		"uniqueBookCount": summary.UniqueBookCount,
		"writerCount": summary.WriterCount,
		"errorCount": summary.ErrorCount,
		"maxBookId": summary.MaxBookId,
		"latestSuccessBookId": summary.LatestSuccessId,
		"statusCount": summary.StatusCount,
	}
	return
}