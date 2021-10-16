package sites

import (
	"os"

	"context"
	"golang.org/x/sync/semaphore"
	"sync"

	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/pkg/books"
)

func (site Site) Validate() float64 {
	os.Mkdir("./validate-download/", 0755)
	site.OpenDatabase()

	var err error
	site.bookLoadTx, err = site.database.Begin()
	utils.CheckError(err)

	rows, err := site.bookQuery(" where download=? order by random()", true)
	utils.CheckError(err)

	ctx := context.Background()
	var s = semaphore.NewWeighted(int64(20))
	var wg sync.WaitGroup

	success, tried := 0.0, 1.0
	for success < 10 && tried < 1000 && rows.Next() {
		wg.Add(1)
		s.Acquire(ctx, 1)
		book, err := books.LoadBook(rows, site.meta, site.decoder, site.CONST_SLEEP)
		if err != nil {
			book.Log(map[string]interface{}{
				"error": "cannot load from database", "stage": "update",
			})
		}
		go site.validateThread(book, &success, &tried, s, &wg)
	}
	rows.Close()
	site.bookLoadTx.Rollback()
	wg.Wait()
	site.CloseDatabase()
	os.RemoveAll("./validate-download/")
	if tried/success > 90 {
		return -1
	}
	return tried / success
}
func (site Site) validateThread(book *books.Book, success *float64,
	tried *float64, s *semaphore.Weighted, wg *sync.WaitGroup) {
	defer wg.Done()
	book.Title = ""
	if *success < 10 && book.Update() && *success < 10 {
		*success++
	} else {
		s.Release(1)
	}
	if *success < 10 {
		*tried++
	}
}
func (site Site) ValidateDownload() float64 {
	os.Mkdir("./validate-download/", 0755)
	site.OpenDatabase()

	var err error
	site.bookLoadTx, err = site.database.Begin()
	utils.CheckError(err)

	rows, err := site.bookQuery(" where download=? order by random()", true)
	utils.CheckError(err)

	ctx := context.Background()
	var s = semaphore.NewWeighted(int64(4))
	var wg sync.WaitGroup

	success, tried := 0.0, 1.0
	for success < 2 && tried < 100 && rows.Next() {
		wg.Add(1)
		s.Acquire(ctx, 1)
		book, err := books.LoadBook(rows, site.meta, site.decoder, site.CONST_SLEEP)
		if err != nil {
			book.Log(map[string]interface{}{
				"error": "cannot load book from database", "stage": "validate-download",
			})
		}
		go site.validateDownloadThread(book, &success, &tried, s, &wg)
	}
	rows.Close()
	site.bookLoadTx.Rollback()
	wg.Wait()
	site.CloseDatabase()
	os.RemoveAll("./validate-download/")
	if tried/success > 90 {
		return -1
	}
	return tried / success
}

func (site Site) validateDownloadThread(book *books.Book, success *float64,
	tried *float64, s *semaphore.Weighted, wg *sync.WaitGroup) {
	defer wg.Done()
	// here have two same condition because `book.Download` take a long time
	// the success may change after finush download
	if *success < 2 && book.Download("./validate-download/", site.MAX_THREAD_COUNT) && *success < 2 {
		*success++
	} else {
		s.Release(1)
	}
	if *success < 2 {
		*tried++
	}
}
