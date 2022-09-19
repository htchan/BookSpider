package site

import (
	"fmt"
	"log"
	"sync"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/service/book"
)

func addMissingRecords(st *Site) {
	var wg sync.WaitGroup
	maxBookID := st.rp.Stats().MaxBookID
	for i := 1; i < maxBookID; i++ {
		i := i
		wg.Add(1)
		st.Client.Acquire()

		go func(id int) {
			defer wg.Done()
			defer st.Client.Release()
			_, err := st.rp.FindBookById(id)
			if err != nil {
				log.Printf("[%v] book <%v> not exist in database", st.Name, i)
				bk := model.NewBook(st.Name, id)
				book.Update(&bk, st.BkConf, st.StConf, st.Client)
				bk.HashCode = 0
				st.rp.SaveError(&bk, bk.Error)
				st.rp.SaveWriter(&bk.Writer)
				st.rp.CreateBook(&bk)
			}
		}(i)
	}
	wg.Wait()
}

func Fix(st *Site) error {
	addMissingRecords(st)

	bks, err := st.rp.FindAllBooks()
	if err != nil {
		return fmt.Errorf("Fix fail: %w", err)
	}

	var wg sync.WaitGroup

	for bk := range bks {
		bk := bk
		wg.Add(1)
		go func(bk *model.Book) {
			defer wg.Done()
			isUpdated, _ := book.Fix(bk, st.StConf)
			if isUpdated {
				st.rp.UpdateBook(bk)
			}
		}(&bk)
	}

	wg.Wait()

	return nil
}
