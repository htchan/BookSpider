package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/rs/zerolog/log"
)

func (serv *ServiceImp) ExploreBook(bk *model.Book) error {
	if bk.Status != model.StatusError {
		return errors.New("book status is not error")
	}

	isNew := bk.Error == nil
	if isNew {
		serv.rpo.CreateBook(bk)
	}

	err := serv.UpdateBook(bk)
	if err != nil {
		log.Error().Err(err).Str("book", bk.String()).Msg("explore book fail")
		bk.Error = err
		saveErr := serv.rpo.SaveError(bk, bk.Error)
		if saveErr != nil {
			return fmt.Errorf("save error fail: %w", saveErr)
		}

		return fmt.Errorf("explore book fail: %w", err)
	}

	err = serv.rpo.SaveError(bk, bk.Error)

	return err
}

func (serv *ServiceImp) exploreExisting(summary repo.Summary, errorCount *int) {
	var wg sync.WaitGroup

	for i := summary.LatestSuccessID + 1; i <= summary.MaxBookID && *errorCount < serv.conf.MaxExploreError; i++ {
		i := i

		serv.sema.Acquire(serv.ctx, 1)
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			defer serv.sema.Release(1)

			bk, err := serv.rpo.FindBookById(id)
			if err != nil {
				*errorCount += 1
				return
			}

			err = serv.ExploreBook(bk)
			if err != nil {
				*errorCount += 1
			} else {
				*errorCount = 0
			}
		}(i)
	}

	wg.Wait()
}

func (serv *ServiceImp) exploreNew(summary repo.Summary, errorCount *int) {
	var wg sync.WaitGroup

	for i := summary.MaxBookID + 1; *errorCount < serv.conf.MaxExploreError; i++ {
		i := i

		serv.sema.Acquire(serv.ctx, 1)
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			defer serv.sema.Release(1)

			bk := model.NewBook(serv.name, i)
			err := serv.ExploreBook(&bk)
			if err != nil {
				*errorCount += 1
			} else {
				*errorCount = 0
			}
		}(i)
	}

	wg.Wait()
}

func (serv *ServiceImp) Explore() error {
	summary := serv.rpo.Stats()
	errorCount := 0

	// explore error books in db
	serv.exploreExisting(summary, &errorCount)

	// explore new books not in db
	serv.exploreNew(summary, &errorCount)

	return nil
}
