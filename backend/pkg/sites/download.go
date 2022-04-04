package sites

import (
	"github.com/htchan/BookSpider/internal/logging"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/pkg/flags"
	"github.com/htchan/BookSpider/pkg/books"
	"errors"
	"context"
	"sync"
	"golang.org/x/sync/semaphore"
)

func (site *Site) download() (err error) {
	ctx := context.Background()
	var wg sync.WaitGroup
	s := semaphore.NewWeighted(int64(site.config.DownloadBookThreads))
	var loadContentMutex sync.Mutex
	// query all end book
	rows := site.database.QueryBooksByStatus(database.End)
	defer rows.Close()
	// loop and construct the book record
	for rows.Next() {
		record, err := rows.ScanCurrent()
		utils.CheckError(err)
		wg.Add(1)
		s.Acquire(ctx, 1)
		loadContentMutex.Lock()
		go func(s *semaphore.Weighted, mutex *sync.Mutex, wg *sync.WaitGroup, record *database.BookRecord) {
			logging.LogBookEvent(record.String(), "download", "start", nil)
			defer s.Release(1)
			defer wg.Done()
			book := books.LoadBookByRecord(site.database, record, site.config.BookMeta)
			// call book.download with thread, wait group and semaphore
			if book.Download(site.config.Threads, mutex) {
				book.Save(site.database)
			}
			logging.LogBookEvent(book.String(), "download", "completed", nil)
		}(s, &loadContentMutex, &wg, record.(*database.BookRecord))
		loadContentMutex.Unlock()
	}
	wg.Wait()
	return nil
}

func Download(site *Site, args *flags.Flags) (err error) {
	if !args.Valid() { return errors.New("invalid arguments") }
	if args.IsBook() && *args.Site == site.Name {
		siteName, id, hash := args.GetBookInfo()
		book := books.LoadBook(site.database, siteName, id, hash, site.config.BookMeta)
		if book != nil {
			maxThreads := *args.MaxThreads
			if maxThreads <= 0 {
				maxThreads = site.config.Threads
			}
			var mutex sync.Mutex
			if book.Download(maxThreads, &mutex) {
				book.Save(site.database)
			}
			return nil
		} else {
			return errors.New("book not found")
		}
	} else if args.IsEverything() || (args.IsSite() && *args.Site == site.Name) {
		return site.download()
	}
	return nil
}