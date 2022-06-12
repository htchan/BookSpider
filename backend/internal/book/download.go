package book

import (
	"fmt"
	"sync"
	"os"
	"errors"
	"strings"
	"strconv"
	"path/filepath"
	"github.com/htchan/BookSpider/internal/book/model"
	"github.com/htchan/ApiParser"
)

var CONTENT_SEP = strings.Repeat("-", 20)

func (book Book) downloadURL() string {
	return fmt.Sprintf(book.BookConfig.URL.Download, book.BookModel.ID)
}

func (book Book) fetchChapters() []Chapter {
	chapters, err := book.generateEmptyChapters()
	if err != nil {
		//TODO: log err
		return nil
	}
	var wg sync.WaitGroup
	for i := range chapters {
		wg.Add(1)
		book.CircuitBreakerClient.Acquire()
		go func(i int) {
			defer wg.Done()
			defer book.CircuitBreakerClient.Release()
			chapters[i].Fetch()
		}(i,)
	}
	wg.Wait()
	return chapters
}

func (book Book) generateEmptyChapters() ([]Chapter, error) {
	html, err := book.Get(book.downloadURL())
	if err != nil {
		return nil, err
	}
	responseApi := ApiParser.Parse(book.BookConfig.SourceKey + ".info", html)
	if len(responseApi.Items) == 0 {
		return nil, errors.New("empty chapters")
	}
	chapters := make([]Chapter, len(responseApi.Items))
	for i, item := range responseApi.Items {
		chapters[i] = NewChapter(i, item["ChapterUrl"], item["ChapterTitle"], &book)
	}
	return chapters, nil
}

func (book Book) content() string {
	return book.Title + "\n" + book.Name + "\n" + CONTENT_SEP + "\n\n"
}

func (book Book) saveChapters(chapters []Chapter) error {
	file, err := os.Create(book.location())
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(book.content())
	if err != nil {
		return err
	}
	errorChapterCount := 0
	for _, chapter := range chapters {
		_, err = file.WriteString(chapter.content())
		if err != nil {
			return err
		}
		if strings.Contains(chapter.Content, "failed") {
			errorChapterCount += 1
		}
	}
	if errorChapterCount > book.BookConfig.MaxChaptersError {
		return errors.New("too many fail chapter")
	}
	return nil
}

func (book *Book) Download() error {
	var wg sync.WaitGroup
	chapters, err := book.generateEmptyChapters()
	if err != nil {
		return err
	}
	for i := range chapters {
		wg.Add(1)
		book.CircuitBreakerClient.Acquire()
		go func(i int) {
			defer wg.Done()
			defer book.CircuitBreakerClient.Release()
			chapters[i].Fetch()
		}(i)
	}
	wg.Wait()
	sortChapters(chapters)
	err = book.saveChapters(chapters)
	if err != nil {
		return err
	}
	book.Status = model.End
	return nil
}

func (book Book) location() string {
	filename := fmt.Sprintf("%v.txt", book.BookModel.ID)
	if book.BookModel.HashCode > 0 {
		filename = fmt.Sprintf(
			"%v-%v.txt",
			book.BookModel.ID, strconv.FormatInt(int64(book.BookModel.HashCode), 36),
		)
	}
	return filepath.Join(book.BookConfig.Storage, filename)
}

func (book Book) Content() string {
	if _, err := os.Stat(book.location()); err != nil {
		return ""
	}
	content, err := os.ReadFile(book.location())
	if err != nil {
		//TODO: log error
		return ""
	}
	return string(content)
}

func (book Book) LoadChapters() []Chapter {
	return []Chapter{}
}