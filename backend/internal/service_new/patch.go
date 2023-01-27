package service

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
)

func (serv *ServiceImp) checkBookStorage(bk *model.Book) bool {
	isUpdated, fileExist := false, true

	if _, err := os.Stat(serv.BookFileLocation(bk)); err != nil {
		fileExist = false
	}

	if fileExist && !bk.IsDownloaded {
		log.Printf("[%v] file exist for not downloaded book", bk)
		bk.IsDownloaded = true
		isUpdated = true
	} else if !fileExist && bk.IsDownloaded {
		log.Printf("[%v] file not exist for downloaded book", bk)
		bk.IsDownloaded = false
		isUpdated = true
	}
	return isUpdated
}

func (serv *ServiceImp) PatchDownloadStatus() error {
	bks, err := serv.rpo.FindAllBooks()
	if err != nil {
		return fmt.Errorf("Patch download status fail: %w", err)
	}

	var wg sync.WaitGroup
	log.Printf("[%s] update books is_downloaded by storage", serv.name)

	for bk := range bks {
		bk := bk
		serv.client.Acquire()
		wg.Add(1)

		go func(bk *model.Book) {
			defer wg.Done()
			defer serv.client.Release()

			isUpdated := serv.checkBookStorage(bk)
			if isUpdated {
				serv.rpo.UpdateBook(bk)
			}
		}(&bk)
	}

	wg.Wait()

	return nil
}

func (serv *ServiceImp) PatchMissingRecords() error {
	log.Printf("[%v] patch missing records", serv.name)

	var wg sync.WaitGroup
	maxBookID := serv.rpo.Stats().MaxBookID
	for i := 1; i <= maxBookID; i++ {
		i := i
		serv.client.Acquire()
		wg.Add(1)

		go func(id int) {
			defer serv.client.Release()
			defer wg.Done()
			_, err := serv.rpo.FindBookById(id)
			if errors.Is(err, repo.BookNotExist) {
				log.Printf("[%v] book <%v> not exist in database", serv.name, i)
				bk := model.NewBook(serv.name, id)
				serv.ExploreBook(&bk)
			} else if err != nil {
				log.Printf("fail to fetch: id: %v; err: %v", id, err)
			}
		}(i)
	}
	wg.Wait()

	return nil
}
