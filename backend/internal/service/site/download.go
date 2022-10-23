package site

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/service/book"
	"golang.org/x/sync/semaphore"
)

func Download(st *Site) error {
	ctx := context.Background()
	s := semaphore.NewWeighted(int64(st.StConf.ConcurrencyConfig.DownloadThreads))
	var wg sync.WaitGroup

	bks, err := st.rp.FindBooksForDownload()
	if err != nil {
		return fmt.Errorf("fail to fetch books: %w", err)
	}

	for bk := range bks {
		bk := bk
		st.Client.Acquire()
		s.Acquire(ctx, 1)
		wg.Add(1)

		go func(bk *model.Book) {
			defer st.Client.Release()
			defer s.Release(1)
			defer wg.Done()

			isUpdated, err := book.Download(bk, st.BkConf, st.StConf, st.Client)
			if isUpdated {
				st.rp.UpdateBook(bk)
			}

			if err != nil {
				log.Printf("[%v] download failed: %v", bk, err)
			}
		}(&bk)
	}
	wg.Wait()
	return nil
}
