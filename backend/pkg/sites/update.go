package sites

import (
	"context"
	"golang.org/x/sync/semaphore"
	"sync"
	"time"

	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/pkg/books"
)

func (site *Site) Update(s *semaphore.Weighted) {
	// init concurrent variable
	site.OpenDatabase()
	ctx := context.Background()
	site.bookLoadTx, _ = site.database.Begin()
	if s == nil {
		s = semaphore.NewWeighted(int64(site.MAX_THREAD_COUNT))
	}
	site.semaphore = semaphore.NewWeighted(int64(site.MAX_THREAD_COUNT))
	var wg sync.WaitGroup
	site.PrepareStmt()
	// update all normal books
	rows, err := site.bookQuery(" group by site, num order by date desc")
	utils.CheckError(err)
	for i := 0; rows.Next(); i++ {
		wg.Add(1)
		s.Acquire(ctx, 1)
		site.semaphore.Acquire(ctx, 1)
		book, err := books.LoadBook(rows, site.meta, site.decoder, site.CONST_SLEEP)
		if err != nil {
			book.Log(map[string]interface{}{
				"message": "cannot load from database", "stage": "update",
			})
			continue
		}
		go site.updateThread(book, s, &wg)
		if i % 100 == 0 {
			time.Sleep(time.Duration(100) * time.Millisecond)
		}
	}
	rows.Close()
	wg.Wait()
	utils.CheckError(site.bookLoadTx.Rollback())
	site.CloseStmt()
	site.CloseDatabase()
}
func (site *Site) updateThread(book *books.Book, s *semaphore.Weighted, wg *sync.WaitGroup) {
	defer wg.Done()
	defer s.Release(1)
	defer site.semaphore.Release(1)
	checkVersion := book.Version
	updated := book.Update()
	if updated {
		if book.Version != checkVersion {
			site.InsertBook(*book)
			book.Log(map[string]interface{}{
				"title": book.Title, "message": "new version updated", "stage": "update",
			})
		} else {
			site.UpdateBook(*book)
			book.Log(map[string]interface{}{
				"message": "regular update", "stage": "update",
			})
		}
	} else {
		book.Log(map[string]interface{}{
			"message": "not updated", "stage": "update",
		})
	}
}

func (site *Site) UpdateError(s *semaphore.Weighted) {
	// init concurrent variable
	site.OpenDatabase()
	var err error
	ctx := context.Background()
	site.bookLoadTx, err = site.database.Begin()
	if s == nil {
		s = semaphore.NewWeighted(int64(site.MAX_THREAD_COUNT))
	}
	site.semaphore = semaphore.NewWeighted(int64(site.MAX_THREAD_COUNT))
	var wg sync.WaitGroup
	var siteName string
	var id int
	site.PrepareStmt()
	// try update all error books
	rows, err := site.database.Query("SELECT site, num FROM error order by num desc")
	utils.CheckError(err)
	for i := 0; rows.Next(); i++ {
		wg.Add(1)
		s.Acquire(ctx, 1)
		site.semaphore.Acquire(ctx, 1)
		rows.Scan(&siteName, &id)
		book := books.NewBook(site.SiteName, id, site.meta, site.decoder, site.bookLoadTx)
		//TODO: try to change this to updateThread (after finish testcase, it should be same)
		go site.updateErrorThread(book, s, &wg)
		if i % 100 == 0 {
			time.Sleep(time.Duration(100) * time.Millisecond)
		}
	}
	rows.Close()
	wg.Wait()
	utils.CheckError(site.bookLoadTx.Rollback())
	site.CloseStmt()
	site.CloseDatabase()
}
func (site *Site) updateErrorThread(book *books.Book, s *semaphore.Weighted, wg *sync.WaitGroup) {
	defer wg.Done()
	defer s.Release(1)
	defer site.semaphore.Release(1)
	updated := book.Update()
	if updated {
		// if update successfully
		site.InsertBook(*book)
		site.DeleteError(*book)
		book.Log(map[string]interface{}{
			"title": book.Title, "message": "error updated", "stage": "update",
		})
	} else {
		// tell others nothing updated
		book.Log(map[string]interface{}{
			"message": "error not updated", "stage": "update",
		})
	}
}
