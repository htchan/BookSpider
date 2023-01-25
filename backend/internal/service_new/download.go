package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/htchan/BookSpider/internal/model"
	"golang.org/x/sync/semaphore"
)

func (serv *ServiceImp) downloadURL(bk *model.Book) string {
	return fmt.Sprintf(serv.conf.URL.Download, bk.ID)
}

func (serv *ServiceImp) chapterURL(bk *model.Book, chapter *model.Chapter) string {
	if strings.HasPrefix(chapter.URL, "http") {
		return chapter.URL
	} else if strings.HasPrefix(chapter.URL, "/") {
		return serv.conf.URL.ChapterPrefix + chapter.URL
	}

	downloadURL := serv.downloadURL(bk)
	if !strings.HasSuffix(downloadURL, "/") {
		downloadURL = downloadURL + "/"
	}

	return downloadURL + chapter.URL
}

func (serv *ServiceImp) downloadChapter(bk *model.Book, ch *model.Chapter) error {
	html, err := serv.client.Get(serv.chapterURL(bk, ch))
	if err != nil {
		ch.Error = fmt.Errorf("fetch chapter html fail: %w", err)
		return ch.Error
	}

	parsedChapterFields, err := serv.parser.ParseChapter(html)
	if err != nil {
		ch.Error = fmt.Errorf("parse chapter html fail: %w", err)
		return ch.Error
	}

	parsedChapterFields.Populate(ch)

	ch.OptimizeContent()

	return ch.Error
}

func (serv *ServiceImp) downloadChapterList(bk *model.Book) (model.Chapters, error) {
	html, err := serv.client.Get(serv.downloadURL(bk))
	if err != nil {
		return nil, fmt.Errorf("fetch chapter list html fail: %w", err)
	}

	parsedChapterList, err := serv.parser.ParseChapterList(html)
	if err != nil {
		return nil, fmt.Errorf("parse chapter list html fail: %w", err)
	}
	var chapters model.Chapters
	parsedChapterList.Populate(&chapters)

	var wg sync.WaitGroup
	for i := range chapters {
		wg.Add(1)
		serv.client.Acquire()
		go func(i int) {
			defer wg.Done()
			defer serv.client.Release()
			err := serv.downloadChapter(bk, &chapters[i])
			if err != nil {
				log.Printf("download Chapter fail. URL: %v, Title: %v, error: %v", chapters[i].URL, chapters[i].Title, err)
			}
		}(i)
	}

	wg.Wait()

	// TODO: return error if there are more than x% of chapter are failed
	return chapters, nil
}

func (serv *ServiceImp) saveContent(location string, bk *model.Book, chapters model.Chapters) error {
	file, err := os.Create(location)
	if err != nil {
		return fmt.Errorf("Save book fail: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(bk.HeaderInfo())
	if err != nil {
		return fmt.Errorf("Save book fail: %w", err)
	}

	for _, chapter := range chapters {
		_, err := file.WriteString(chapter.ContentString())
		if err != nil {
			return err
		}
	}

	return nil
}

func (serv *ServiceImp) DownloadBook(bk *model.Book) error {
	if bk.Status != model.End {
		return fmt.Errorf("book status not ready for download. status: %v", bk.Status)
	} else if bk.IsDownloaded {
		return errors.New("book was downloaded")
	}

	log.Printf("[%v] download chapters", bk)
	chapters, err := serv.downloadChapterList(bk)
	if err != nil {
		return fmt.Errorf("Download chapters fail: %w", err)
	}

	log.Printf("[%v] save content", bk)
	err = serv.saveContent(serv.BookFileLocation(bk), bk, chapters)
	if err != nil {
		return fmt.Errorf("save content fail: %w", err)
	}
	bk.IsDownloaded = true
	err = serv.rpo.UpdateBook(bk)
	if err != nil {
		return fmt.Errorf("update book fail: %w", err)
	}

	return nil
}

func (serv *ServiceImp) Download() error {
	ctx := context.Background()
	s := semaphore.NewWeighted(int64(serv.conf.MaxDownloadConcurrency))
	var wg sync.WaitGroup

	bkChan, err := serv.rpo.FindBooksForDownload()
	if err != nil {
		return fmt.Errorf("fail to fetch books: %w", err)
	}

	for bk := range bkChan {
		bk := bk
		serv.client.Acquire()
		s.Acquire(ctx, 1)
		wg.Add(1)

		go func(bk *model.Book) {
			defer wg.Done()
			defer s.Release(1)
			defer serv.client.Release()

			err := serv.DownloadBook(bk)
			if err != nil {
				log.Printf("[%v] download failed: %v", bk, err)
			}
		}(&bk)
	}

	wg.Wait()

	return nil
}
