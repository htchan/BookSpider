package sites

import (
	"github.com/htchan/BookSpider/internal/logging"
	"github.com/htchan/BookSpider/pkg/flags"
	"github.com/htchan/BookSpider/pkg/books"
	"errors"
	"fmt"
	"sync"
	"context"
	"golang.org/x/sync/semaphore"
)

func (site *Site) exploreOldBook(id int, count *int) error {
	book := books.LoadBook(site.database, site.Name, id, -1, site.config.BookMeta)
	if book == nil {
		return errors.New(fmt.Sprintf(
		"load book %v-%v fail", site.Name, id))
	}
	if book.GetError() == nil {
		return errors.New(fmt.Sprintf(
		"load book %v-%v return status %v", site.Name, id, book.GetStatus()))
	}
	if book.Update() {
		book.Save(site.database)
		*count = 0
	} else {
		(*count)++
	}
	return nil
}

func (site *Site) exploreNewBook(id int, count *int) error {
	book := books.NewBook(site.Name, id, -1, site.config.BookMeta)
	updateSuccess := book.Update() 
	book.Save(site.database)
	if updateSuccess {
		*count = 0
	} else {
		(*count)++
	}
	return nil
}

// mark books to end status after finish update
func (site *Site) explore() (err error) {
	summary := site.database.Summary(site.Name)
	ctx := context.Background()
	var wg sync.WaitGroup
	i := summary.LatestSuccessId + 1
	errorCount := 0
	// loop latest continuous error books in database
	for ; i <= summary.MaxBookId; i++ {
		if errorCount >= site.config.MaxExploreError {
			return nil
		}
		site.semaphore.Acquire(ctx, 1)
		wg.Add(1)
		go func (s *semaphore.Weighted, wg *sync.WaitGroup, id int, errorCount *int) {
			defer s.Release(1)
			defer wg.Done()
			err := site.exploreOldBook(id, errorCount)
			if err != nil {
				logging.Info("Book %v-%v explore failed: %v", site.Name, id, err)
			} else {
				logging.Info("Book %v-%v explore complete", site.Name, id)
			}
		} (site.semaphore, &wg, i, &errorCount)
	}
	// loop books not in database
	for ; errorCount < site.config.MaxExploreError; i++ {
		site.semaphore.Acquire(ctx, 1)
		wg.Add(1)
		go func (s *semaphore.Weighted, wg *sync.WaitGroup, id int, errorCount *int) {
			defer s.Release(1)
			defer wg.Done()
			err := site.exploreNewBook(id, errorCount)
			if err != nil {
				logging.Info("Book %v-%v explore failed: %v", site.Name, id, err)
			} else {
				logging.Info("Book %v-%v explore complete", site.Name, id)
			}
		} (site.semaphore, &wg, i, &errorCount)
	}
	wg.Wait()
	return
}

func Explore(site *Site, args *flags.Flags) (err error) {
	if !args.Valid() { return errors.New("invalid arguments") }
	if args.IsBook() && *args.Site == site.Name {
		if *args.Site != site.Name { return nil }
		siteName, id, hash := args.GetBookInfo()
		book := books.LoadBook(site.database, siteName, id, hash, site.config.BookMeta)
		if book != nil {
			book = books.NewBook(siteName, id, hash, site.config.BookMeta)
		}
		book.Update()
		book.Save(site.database)
		return nil
	} else if args.IsEverything() || (args.IsSite() && *args.Site == site.Name) {
		return site.explore()
	}
	return nil
}