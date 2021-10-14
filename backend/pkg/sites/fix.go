package sites

import (
	"log"
	"strconv"

	"context"
	"golang.org/x/sync/semaphore"
	"sync"

	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/pkg/books"
)

func (site *Site) fixStroageError(s *semaphore.Weighted) {
	// init var for concurrency
	ctx := context.Background()
	if s == nil {
		s = semaphore.NewWeighted(int64(site.MAX_THREAD_COUNT))
	}
	var wg sync.WaitGroup
	var err error
	site.bookLoadTx, err = site.database.Begin()
	utils.CheckError(err)
	site.PrepareStmt()
	rows, err := site.bookQuery("")
	// loop all book
	for rows.Next() {
		wg.Add(1)
		s.Acquire(ctx, 1)
		book, err := books.LoadBook(rows, site.meta, site.decoder)
		if err != nil {
			book.Log(map[string]interface{}{
				"error": "cannot load book from database", "stage": "fix",
			})
			continue
		}
		go site.CheckDownloadExistThread(book, s, &wg)
	}
	wg.Wait()
	// commit changes to database
	utils.CheckError(rows.Close())
	utils.CheckError(site.bookLoadTx.Rollback())
	site.CloseStmt()
	site.bookLoadTx = nil
}

func (site *Site) CheckDownloadExistThread(book *books.Book, s *semaphore.Weighted,
	wg *sync.WaitGroup) {
	defer wg.Done()
	defer s.Release(1)
	bookLocation := book.StorageLocation(site.DownloadLocation)
	// check book file exist
	exist := utils.Exists(bookLocation)
	if exist && !book.DownloadFlag {
		// if book mark as not download, but it exist, mark as download
		book.EndFlag = true
		book.DownloadFlag = true
		site.UpdateBook(*book)
		book.Log(map[string]interface{}{
			"message": "mark to download", "stage": "fix",
		})
	} else if !exist && book.DownloadFlag {
		// if book mark as download, but not exist, mark as not download
		book.DownloadFlag = false
		site.UpdateBook(*book)
		book.Log(map[string]interface{}{
			"message": "mark to not download yet", "stage": "fix",
		})
	}
}

//TODO: break it into delete duplicated book, delete duplicated error and delete cross table duplicated
func (site *Site) fixDuplicateRecordError() {
	// init variable
	var bookId, bookVersion, bookCount, errorCount int
	rows, tx := site.query("select num, version from books group by num, version having count(*) > 1")
	for rows.Next() {
		rows.Scan(&bookId, &bookVersion)
		log.Println(site.SiteName, "("+strconv.Itoa(bookId)+", "+strconv.Itoa(bookVersion)+")")
		bookCount += 1
	}
	log.Println(site.SiteName, "duplicate book count : "+strconv.Itoa(bookCount))
	// delete duplicate record in book
	deleteStmt, err := tx.Prepare("delete from books where rowid not in " +
		"(select min(rowid) from books group by num, version)")
	utils.CheckError(err)
	_, err = deleteStmt.Exec()
	utils.CheckError(err)
	utils.CheckError(deleteStmt.Close())
	closeQuery(rows, tx)
	// check any duplicate record in error table and show them
	rows, tx = site.query("select num from error group by num having count(*) > 1")
	for rows.Next() {
		rows.Scan(&bookId)
		log.Println(site.SiteName, "("+strconv.Itoa(bookId)+")")
	}
	log.Println(site.SiteName, "duplicate error count : "+strconv.Itoa(errorCount))
	// delete duplicate record
	deleteStmt, err = tx.Prepare("delete from error where rowid not in " +
		"(select min(rowid) from books group by site, num)")
	utils.CheckError(err)
	_, err = deleteStmt.Exec()
	utils.CheckError(err)
	utils.CheckError(deleteStmt.Close())
	// check if any record in book table duplicate in error table
	log.Println(site.SiteName, "duplicate cross - - - - - - - - - -")
	deleteStmt, err = tx.Prepare("delete from error where num in (select distinct num from books)")
	utils.CheckError(err)
	tx.Stmt(deleteStmt).Exec()
	utils.CheckError(deleteStmt.Close())
	closeQuery(rows, tx)
}

func (site *Site) fixMissingRecordError(s *semaphore.Weighted) {
	// init variable
	missingIds := site.missingIds()
	// insert missing record by thread
	var err error
	site.bookLoadTx, err = site.database.Begin()
	utils.CheckError(err)
	if s == nil {
		s = semaphore.NewWeighted(int64(site.MAX_THREAD_COUNT))
	}
	ctx := context.Background()
	var wg sync.WaitGroup
	var errorCount int
	site.PrepareStmt()
	log.Println(site.SiteName, "missing count : "+strconv.Itoa(len(missingIds)))
	for _, bookId := range missingIds {
		log.Println(site.SiteName, bookId)
		wg.Add(1)
		s.Acquire(ctx, 1)
		book := books.NewBook(site.SiteName, bookId, site.meta, site.decoder, site.bookLoadTx)
		go site.exploreThread(book, &errorCount, s, &wg)
	}
	wg.Wait()
	utils.CheckError(site.bookLoadTx.Rollback())
	site.CloseStmt()
	// print missing record count
	log.Println(site.SiteName, "finish add missing count", len(missingIds))
}

func (site *Site) Fix(s *semaphore.Weighted) {
	site.OpenDatabase()
	log.Println(site.SiteName, "Add Missing Record")
	site.fixMissingRecordError(s)
	log.Println(site.SiteName, "Fix duplicate record")
	site.fixDuplicateRecordError()
	log.Println(site.SiteName, "Fix storage error")
	site.fixStroageError(s)
	log.Println()
	site.CloseDatabase()
}
