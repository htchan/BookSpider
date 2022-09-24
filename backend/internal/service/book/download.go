package book

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/htchan/ApiParser"
	"github.com/htchan/BookSpider/internal/client"
	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/service/chapter"
)

func downloadURL(bk *model.Book, bkConf config.BookConfig) string {
	return fmt.Sprintf(bkConf.URLConfig.Download, bk.ID)
}

func fetchChaptersHeaderInfo(bk *model.Book, bkConf config.BookConfig, stConf config.SiteConfig, c *client.CircuitBreakerClient) (model.Chapters, error) {
	html, err := c.Get(downloadURL(bk, bkConf))
	if err != nil {
		return nil, err
	}

	responseApi := ApiParser.Parse(stConf.BookKey+".info", html)
	if len(responseApi.Items) == 0 {
		return nil, errors.New("empty chapters")
	}

	chapters := make(model.Chapters, len(responseApi.Items))
	for i, item := range responseApi.Items {
		chapters[i] = model.NewChapter(i, item["ChapterUrl"], item["ChapterTitle"])
	}

	return chapters, nil
}

func downloadChapters(
	bk *model.Book, chapters model.Chapters, bkConf config.BookConfig,
	stConf config.SiteConfig, c *client.CircuitBreakerClient,
) error {
	var wg sync.WaitGroup

	for i := range chapters {
		wg.Add(1)
		c.Acquire()
		go func(i int) {
			defer wg.Done()
			defer c.Release()
			chapter.Download(bk.ID, &chapters[i], bkConf, stConf, c)
		}(i)
	}

	wg.Wait()

	// TODO: return error if there are more than x% of chapter are failed
	return nil
}

func headerInfo(bk *model.Book) string {
	return bk.Title + "\n" + bk.Writer.Name + "\n" + model.CONTENT_SEP + "\n\n"
}

func saveContent(location string, bk *model.Book, chapters model.Chapters) error {
	file, err := os.Create(location)
	if err != nil {
		return fmt.Errorf("Save book fail: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(headerInfo(bk))
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

func Download(bk *model.Book, bkConf config.BookConfig, stConf config.SiteConfig, c *client.CircuitBreakerClient) (bool, error) {
	if bk.Status != model.End || bk.IsDownloaded {
		return false, nil
	}

	log.Printf("[%v] fetching header info", bk)
	chapters, err := fetchChaptersHeaderInfo(bk, bkConf, stConf, c)
	if err != nil {
		return false, fmt.Errorf("Download book error: %w", err)
	}

	log.Printf("[%v] download chapters; total: %v", bk, len(chapters))
	err = downloadChapters(bk, chapters, bkConf, stConf, c)
	if err != nil {
		return false, fmt.Errorf("Download book error: %w", err)
	}

	log.Printf("[%v] save content", bk)
	err = saveContent(BookFileLocation(bk, stConf), bk, chapters)
	if err != nil {
		return false, fmt.Errorf("Download book error: %w", err)
	}
	bk.IsDownloaded = true

	return true, nil
}
