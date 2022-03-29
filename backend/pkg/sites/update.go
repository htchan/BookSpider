package sites

import (
	"github.com/htchan/BookSpider/pkg/flags"
	"github.com/htchan/BookSpider/pkg/books"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/internal/logging"
	"github.com/htchan/BookSpider/internal/utils"
	"errors"
	"context"
	"sync"
	"fmt"
	"golang.org/x/sync/semaphore"
)

func (site *Site) updateBook(record *database.BookRecord) error {
	book := books.LoadBookByRecord(site.database, record, site.config.BookMeta)
	if book == nil {
		return errors.New(fmt.Sprintf(
			"[update] failed to load books %v-%v", site.Name, record.Id))
	}
	if book.Update() {
		book.Save(site.database)
		logging.LogSiteEvent(site.Name, "updated", "changed", record.String())
	} else {
		logging.LogSiteEvent(site.Name, "updated", "not-changed", record.String())
	}
	return nil
}

// update / create books
func (site *Site) update(errorFocus bool) (err error) {
	// loop all books in site db
	ctx := context.Background()
	var rows database.Rows
	if errorFocus {
		rows = site.database.QueryBooksByStatus(database.Error)
	} else {
		rows = site.database.QueryBooksOrderByUpdateDate()
	}
	var wg sync.WaitGroup
	for rows.Next() {
		record, err := rows.ScanCurrent()
		if err != nil {
			logging.LogSiteEvent(site.Name, "updated", "db-fail", err)
			return err
		}
		site.semaphore.Acquire(ctx, 1)
		wg.Add(1)
		go func(s *semaphore.Weighted, wg *sync.WaitGroup, record *database.BookRecord) {
			defer s.Release(1)
			defer wg.Done()
			err := site.updateBook(record)
			if err != nil {
			logging.LogSiteEvent(site.Name, "updated", "update-failed", record.String() + err.Error())
			}
		} (site.semaphore, &wg, record.(*database.BookRecord))
		if site.config.UseRequestInterval { utils.RequestInterval() }
	}
	wg.Wait()
	rows.Close()
	return
}

func Update(site *Site, args *flags.Flags) (err error) {
	if !args.Valid() { return errors.New("invalid arguments") }
	if args.IsBook() && *args.Site == site.Name {
		if *args.Site != site.Name { return nil }
		siteName, id, hash := args.GetBookInfo()
		book := books.LoadBook(site.database, siteName, id, hash, site.config.BookMeta)
		if book == nil {
			return errors.New("book not found")
		}
		if book.Update() {
			book.Save(site.database)
		}
		return nil
	} else if args.IsEverything() || (args.IsSite() && *args.Site == site.Name) {
		return site.update(false)
	}
	return nil
}

func UpdateError(site *Site, args *flags.Flags) (err error) {
	if !args.Valid() { return errors.New("invalid arguments") }
	if args.IsBook() && *args.Site == site.Name {
		return errors.New("use Update instead of UpdateError")
	} else if args.IsEverything() || (args.IsSite() && *args.Site == site.Name) {
		return site.update(true)
	}
	return nil
}