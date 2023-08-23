package service

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/rs/zerolog/log"
)

func (serv *ServiceImp) checkBookStorage(bk *model.Book) bool {
	isUpdated, fileExist := false, true

	if _, err := os.Stat(serv.BookFileLocation(bk)); err != nil {
		fileExist = false
	}

	if fileExist && !bk.IsDownloaded {
		log.Info().Str("book", bk.String()).Msg("file exist for not downloaded book")
		bk.IsDownloaded = true
		isUpdated = true
	} else if !fileExist && bk.IsDownloaded {
		log.Info().Str("book", bk.String()).Msg("file not exist for downloaded book")
		bk.IsDownloaded = false
		isUpdated = true
	}
	return isUpdated
}

func (serv *ServiceImp) PatchDownloadStatus() error {
	bks, err := serv.rpo.FindAllBooks()
	if err != nil {
		return fmt.Errorf("patch download status fail: %w", err)
	}

	var wg sync.WaitGroup
	log.Info().Str("site", serv.name).Msg("update books is_downloaded by storage")

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
	log.Info().Str("site", serv.name).Msg("patch missing records")

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
				log.Error().Err(err).Str("site", serv.name).Int("id", id).Msg("book not exist in database")
				bk := model.NewBook(serv.name, id)
				serv.ExploreBook(&bk)
			} else if err != nil {
				log.Error().Err(err).Str("site", serv.name).Int("id", id).Msg("fetch book failed")
			}
		}(i)
	}
	wg.Wait()

	return nil
}
