package sites

import (
	"runtime"

	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/pkg/books"
)

func (site *Site) Download() {
	site.OpenDatabase()
	var err error
	site.bookLoadTx, err = site.database.Begin()
	utils.CheckError(err)
	rows, err := site.bookQuery(" where end=? and download=?", true, false)

	site.PrepareStmt()
	for rows.Next() {
		book, err := books.LoadBook(rows, site.meta, site.decoder)
		if err != nil {
			book.Log(map[string]interface{}{
				"error": "cannot load books from database", "stage": "downlaod",
			})
			continue
		}
		if book.DownloadFlag {
			book.Log(map[string]interface{}{
				"error": "book already been download", "stage": "download",
			})
			continue
		}
		book.Log(map[string]interface{}{
			"title": book.Title, "message": "start download", "stage": "download",
		})

		check := book.Download(site.DownloadLocation, site.MAX_THREAD_COUNT)
		if !check {
			book.Log(map[string]interface{}{
				"title": book.Title, "message": "download failure", "stage": "download",
			})
		} else {
			site.UpdateBook(*book)
		}
		runtime.GC()
	}
	utils.CheckError(rows.Close())
	utils.CheckError(site.bookLoadTx.Rollback())
	site.CloseStmt()
	site.CloseDatabase()
}
