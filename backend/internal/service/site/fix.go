package site

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/htchan/BookSpider/internal/service/book"
)

func PatchMissingRecords(st *Site) error {
	log.Printf("[%v] patch missing records", st.Name)

	var wg sync.WaitGroup
	maxBookID := st.rp.Stats().MaxBookID
	for i := 1; i < maxBookID; i++ {
		i := i
		st.Client.Acquire()
		wg.Add(1)

		go func(id int) {
			defer st.Client.Release()
			defer wg.Done()
			_, err := st.rp.FindBookById(id)
			if errors.Is(err, repo.BookNotExist) {
				log.Printf("[%v] book <%v> not exist in database", st.Name, i)
				bk := model.NewBook(st.Name, id)
				book.Update(&bk, st.BkConf, st.StConf, st.Client)
				bk.HashCode = 0
				st.rp.SaveError(&bk, bk.Error)
				st.rp.SaveWriter(&bk.Writer)
				st.rp.CreateBook(&bk)
			} else if err != nil {
				log.Printf("fail to fetch: id: %v; err: %v", id, err)
			}
		}(i)
	}
	wg.Wait()

	return nil
}

func PatchDownloadStatus(st *Site) error {
	bks, err := st.rp.FindAllBooks()
	if err != nil {
		return fmt.Errorf("Patch download status fail: %w", err)
	}

	var wg sync.WaitGroup
	log.Printf("[%s] update books is_downloaded by storage", st.Name)

	for bk := range bks {
		bk := bk
		st.Client.Acquire()
		wg.Add(1)
		go func(bk *model.Book) {
			defer st.Client.Release()
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

func Fix(st *Site) error {
	err := PatchMissingRecords(st)
	if err != nil {
		return fmt.Errorf("Fix failed: %w", err)
	}

	err = PatchDownloadStatus(st)
	if err != nil {
		return fmt.Errorf("Fix failed: %w", err)
	}
	return nil

}
