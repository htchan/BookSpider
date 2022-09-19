package site

import (
	"fmt"
	"sync"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/service/book"
)

func Update(st *Site) error {
	bks, err := st.rp.FindBooksForUpdate()
	if err != nil {
		return fmt.Errorf("fail to fetch books: %w", err)
	}

	var wg sync.WaitGroup

	for bk := range bks {
		bk := bk
		st.Client.Acquire()
		wg.Add(1)
		go func(bk *model.Book) {
			defer st.Client.Release()
			defer wg.Done()

			originalHashCode := bk.HashCode
			originalStatus := bk.Status
			isUpdated, err := book.Update(bk, st.BkConf, st.StConf, st.Client)
			if isUpdated {
				if bk.HashCode != originalHashCode {
					st.rp.SaveError(bk, bk.Error)
					st.rp.SaveWriter(&bk.Writer)
					st.rp.CreateBook(bk)
				} else {
					if originalStatus == model.Error {
						st.rp.SaveError(bk, bk.Error)
						st.rp.SaveWriter(&bk.Writer)
					}
					st.rp.UpdateBook(bk)
				}
			}

			if err != nil {
				//TODO: log the error
			}
		}(&bk)
	}
	wg.Wait()
	return nil
}
