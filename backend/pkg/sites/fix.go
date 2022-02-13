package sites

import (
	"github.com/htchan/BookSpider/pkg/flags"
	"github.com/htchan/BookSpider/pkg/books"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
	"errors"
	"sync"
	"golang.org/x/sync/semaphore"
	"context"
	"strings"
	"strconv"
	"os"
	// "path/filepath"
	// "time"
	// "fmt"
)

func (site *Site) addMissingRecords() (err error) {
	// fix missing record
	summary := site.database.Summary(site.Name)
	var wg sync.WaitGroup
	ctx := context.Background()
	for i := 1; i <= summary.MaxBookId; i++ {
		wg.Add(1)
		site.semaphore.Acquire(ctx, 1)
		go func(s *semaphore.Weighted, wg *sync.WaitGroup, i int) {
			defer s.Release(1)
			defer wg.Done()
			book := books.LoadBook(site.database, site.Name, i, -1, site.config.BookMeta)
			if book == nil {
				book = books.NewBook(site.Name, i, -1, site.config.BookMeta)
				book.Update()
				book.Save(site.database)
			}
		}(site.semaphore, &wg, i)
	}
	wg.Wait()
	return nil
}

func (site *Site) printPotentialDuplicatedRecords() (err error) {
	// todo: print the site num version of two books
	// 	     if site num are equal and version smaller than specific number
	return nil
}

func (site *Site) updateBooksByStorage() (err error) {
	// fix storage error (update database status <download> and <end> according to storage)
	defer utils.Recover(func() {})
	var wg sync.WaitGroup
	ctx := context.Background()
	// loop all download book to ensure they have txt download else turn them to end
	rows := site.database.QueryBooksByStatus(database.Download)
	defer rows.Close()
	for rows.Next() {
		wg.Add(1)
		site.semaphore.Acquire(ctx, 1)
		record, err := rows.ScanCurrent()
		utils.CheckError(err)
		go func(s *semaphore.Weighted, wg *sync.WaitGroup, record *database.BookRecord) {
			defer s.Release(1)
			defer wg.Done()
			book := books.LoadBookByRecord(site.database, record, site.config.BookMeta)
			// ensure the storage exist
			if !book.HasContent() {
				book.SetStatus(database.End)
				book.Save(site.database)
			}
		}(site.semaphore, &wg, record.(*database.BookRecord))
	}
	wg.Wait()
	// loop all storage to ensure all storage can find a mapped download record
	storageList, err := os.ReadDir(os.Getenv("ASSETS_LOCATION") + site.config.StorageDirectory)
	utils.CheckError(err)
	for _, file := range storageList {
		if file.IsDir() { continue }
		wg.Add(1)
		site.semaphore.Acquire(ctx, 1)
		go func(s *semaphore.Weighted, wg *sync.WaitGroup, filename string) {
			defer s.Release(1)
			defer wg.Done()
			if filename[len(filename) - 3:] != "txt" { return }
			info := strings.Split(strings.ReplaceAll(filename, ".txt", ""), "-v")
			i, _ := strconv.Atoi(info[0])
			hashCode := 0
			if len(info) == 2 {
				hashCode, err = strconv.Atoi(info[1])
				if err != nil { return }
			}
			book := books.LoadBook(site.database, site.Name, i, hashCode, site.config.BookMeta)
			if book!= nil && book.GetStatus() != database.Download {
				book.SetStatus(database.Download)
				book.Save(site.database)
			}
		}(site.semaphore, &wg, file.Name())
	}
	wg.Wait()
	return nil
}

func (site *Site) fix() (err error) {
	defer utils.Recover(func() {})
	err = site.addMissingRecords()
	utils.CheckError(err)
	err = site.printPotentialDuplicatedRecords()
	utils.CheckError(err)
	return site.updateBooksByStorage()
}

func (site *Site) Fix(args *flags.Flags) (err error) {
	if !args.Valid() || args.IsBook() {
		err = errors.New("invalid arguments")
		return
	}
	if args.IsEverything() || (args.IsSite() && *args.Site == site.Name) {
		return site.fix()
	}
	return nil
}