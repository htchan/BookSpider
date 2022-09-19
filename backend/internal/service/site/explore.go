package site

import (
	"log"
	"sync"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/htchan/BookSpider/internal/service/book"
)

func exploreExistingErrorBooks(st *Site, summary repo.Summary, initErrorCount int) int {
	// start from latest success id
	var wg sync.WaitGroup
	errorCount := initErrorCount
	for i := summary.LatestSuccessID + 1; i <= summary.MaxBookID; i++ {
		if errorCount >= st.StConf.MaxExploreError {
			break
		}

		st.Client.Acquire()
		wg.Add(1)
		go func(id int) {
			defer st.Client.Release()
			defer wg.Done()

			bk, err := st.rp.FindBookById(id)
			if err != nil {
				errorCount += 1
				return
			}

			isUpdated, err := book.Update(bk, st.BkConf, st.StConf, st.Client)
			if isUpdated {
				st.rp.SaveError(bk, bk.Error)
				st.rp.SaveWriter(&bk.Writer)
				st.rp.UpdateBook(bk)
			}
			if err != nil {
				errorCount += 1
			} else {
				errorCount = 0
			}
		}(i)
	}
	wg.Wait()
	return errorCount
}

func exploreNewBooks(st *Site, summary repo.Summary, initErrorCount int) int {
	// start from max book id
	var wg sync.WaitGroup
	errorCount := initErrorCount
	for i := summary.MaxBookID + 1; errorCount < st.StConf.MaxExploreError; i++ {
		st.Client.Acquire()
		wg.Add(1)
		go func(id int) {
			defer st.Client.Release()
			defer wg.Done()

			bk := model.NewBook(st.Name, id)

			_, err := book.Update(&bk, st.BkConf, st.StConf, st.Client)
			st.rp.SaveError(&bk, bk.Error)
			st.rp.SaveWriter(&bk.Writer)
			st.rp.CreateBook(&bk)
			if err != nil {
				errorCount += 1
			} else {
				errorCount = 0
			}
		}(i)
	}
	wg.Wait()
	return errorCount
}

func Explore(st *Site) error {
	summary := st.rp.Stats()
	errCount := exploreExistingErrorBooks(st, summary, 0)
	if errCount >= st.StConf.MaxExploreError {
		log.Println("total error count: ", errCount)
		return nil
	}

	errCount = exploreNewBooks(st, summary, errCount)
	log.Println("total error count: ", errCount)

	return nil
}
