package sites

import (
	"log"

	"context"
	"golang.org/x/sync/semaphore"
	"sync"

	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/pkg/books"
)

func (site *Site) Explore(maxError int, s *semaphore.Weighted) {
	// init concurrent variable
	site.OpenDatabase()
	ctx := context.Background()
	var err error
	site.bookLoadTx, err = site.database.Begin()
	utils.CheckError(err)
	site.PrepareStmt()
	if s == nil {
		s = semaphore.NewWeighted(int64(maxError))
	}
	var wg sync.WaitGroup

	maxId := site.maxBookId() + 1
	if maxId == 0 {
		maxId++
	}
	log.Println(maxId)
	// keep explore until reach max error count
	errorCount := 0
	for errorCount < maxError {
		wg.Add(1)
		s.Acquire(ctx, 1)
		book := books.NewBook(site.SiteName, maxId, site.meta, site.decoder, site.bookLoadTx)
		go site.exploreThread(book, &errorCount, s, &wg)
		maxId += 1
	}
	wg.Wait()
	utils.CheckError(site.bookLoadTx.Rollback())
	site.CloseStmt()
	site.CloseDatabase()
}

func (site *Site) exploreThread(book *books.Book, errorCount *int, s *semaphore.Weighted,
	wg *sync.WaitGroup) {
	defer wg.Done()
	defer s.Release(1)
	if book.Version >= 0 {
		//TODO: here break the loop is not the best handling
		//		the best handling is do the update and update the database
		book.Log(map[string]interface{}{
			"error": "books already in database", "stage": "explore",
		})
		*errorCount = 0
		return
	}
	updated := book.Update()
	// if updated, save in books table, else, save in error table and **reset error count**
	if updated {
		site.InsertBook(*book)
		site.DeleteError(*book)
		book.Log(map[string]interface{}{
			"title": book.Title, "writer": book.Writer, "type": book.Type,
			"lastUpdate": book.LastUpdate, "lastChapter": book.LastChapter,
			"message": "explored", "stage": "explore",
		})
		*errorCount = 0
	} else { // increase error Count
		site.InsertError(*book)
		book.Log(map[string]interface{}{
			"message": "no such book", "stage": "explore",
		})
		*errorCount++
	}
}
