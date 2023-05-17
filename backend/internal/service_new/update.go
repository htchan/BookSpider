package service

import (
	"fmt"
	"sync"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/parse"
	"github.com/rs/zerolog/log"
)

func (serv *ServiceImp) baseURL(bk *model.Book) string {
	return fmt.Sprintf(serv.conf.URL.Base, bk.ID)
}
func (serv *ServiceImp) UpdateBook(bk *model.Book) error {
	html, err := serv.client.Get(serv.baseURL(bk))
	if err != nil {
		return fmt.Errorf("fetch html fail: %w", err)
	}

	parsedBookFields, err := serv.parser.ParseBook(html)
	if err != nil {
		return fmt.Errorf("parse html fail: %w", err)
	}

	if parse.IsNewBook(parsedBookFields, bk) {
		bk.HashCode = model.GenerateHash()
		bk.Status = model.InProgress
		bk.Error = nil
		parsedBookFields.Populate(bk)
		err := serv.rpo.SaveWriter(&bk.Writer)
		if err != nil {
			return err
		}
		err = serv.rpo.CreateBook(bk)
		if err != nil {
			return err
		}

		log.Info().Str("book", bk.String()).Str("title", bk.Title).Msg("new book found")
	} else if parse.IsUpdatedBook(parsedBookFields, bk) {
		// TODO: log updated
		bk.Status = model.InProgress
		bk.Error = nil
		parsedBookFields.Populate(bk)
		log.Info().Str("book", bk.String()).Str("title", bk.Title).Msg("updated book found")
		err := serv.rpo.SaveWriter(&bk.Writer)
		if err != nil {
			return err
		}
		err = serv.rpo.UpdateBook(bk)
		if err != nil {
			return err
		}
	} else {
		log.Debug().Str("book", bk.String()).Msg("book not updated")
	}

	return nil
}

func (serv *ServiceImp) Update() error {
	var wg sync.WaitGroup

	bkChan, err := serv.rpo.FindBooksForUpdate()
	if err != nil {
		return fmt.Errorf("fail to load books from DB: %w", err)
	}

	for bk := range bkChan {
		bk := bk
		serv.client.Acquire()
		wg.Add(1)

		go func(bk *model.Book) {
			defer wg.Done()
			defer serv.client.Release()

			err := serv.UpdateBook(bk)
			if err != nil {
				log.Error().Err(err).Str("book", bk.String()).Msg("update book failed")
			}
		}(&bk)
	}

	wg.Wait()

	return nil
}
