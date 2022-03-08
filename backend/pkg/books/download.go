package books

import (
	"context"
	"errors"
	"fmt"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/internal/logging"
	"golang.org/x/sync/semaphore"
	// "log"
	"os"
	// "strconv"
	"strings"
	"sync"
)

func (book Book) getEmptyChapters() (chapters []Chapter, err error) {
	// get basic info (all chapter url and title)
	html, _ := utils.GetWeb(book.config.DownloadUrl, 10, book.config.Decoder, book.config.CONST_SLEEP)
	if err = book.validHTML(html); err != nil {
		err = errors.New(fmt.Sprintf("invalid table of content html: %v", html))
		return
	}
	urls := utils.SearchAll(html, book.config.ChapterUrlRegex)
	titles := utils.SearchAll(html, book.config.ChapterTitleRegex)
	// if length are difference, return error
	if len(urls) != len(titles) {
		err = errors.New("title and url have different length")
		return
	} else if len(urls) == 0 {
		err = errors.New("no chapter found")
		return
	}
	chapters = make([]Chapter, len(urls))
	for i := 0; i < len(urls); i++ {
		chapters[i] = NewChapter(i, urls[i], titles[i], &book.config)
	}
	return
}

func (book *Book) downloadChapter(i int, url, title string, s *semaphore.Weighted,
	wg *sync.WaitGroup, ch chan<- Chapter) {
	defer wg.Done()
	defer s.Release(1)
	// get chapter resource
	html, _ := utils.GetWeb(url, 10, book.config.Decoder, book.config.CONST_SLEEP)
	chapter := Chapter{Url: url, Title: title}
	// chapter.generateIndex()
	chapter.Index = i
	if err := book.validHTML(html); err != nil {
		chapter.Content = "load html fail"
		ch <- chapter
		return
	}
	// extract chapter
	content, err := utils.Search(html, book.config.ChapterContentRegex)
	if err != nil {
		chapter.Content = "recognize html fail\n" + html
		ch <- chapter
	} else {
		chapter.Content = content
		chapter.optimizeContent()
		ch <- chapter
	}
}
func (book Book) downloadChapters(chapters []Chapter, MAX_THREAD int) []Chapter {
	ctx := context.Background()
	var s = semaphore.NewWeighted(int64(MAX_THREAD))
	var wg sync.WaitGroup

	for i, _ := range chapters {
		wg.Add(1)
		s.Acquire(ctx, 1)
		go func(wg *sync.WaitGroup, s *semaphore.Weighted, i int) {
			defer wg.Done()
			defer s.Release(1)
			chapters[i].Download(&book.config, book.validHTML)
		}(&wg, s, i)
	}
	wg.Wait()
	return chapters
}

func (book Book) saveContent(chapters []Chapter) int {
	errorChapterCount := 0
	f, err := os.Create(book.getContentLocation())
	utils.CheckError(err)
	f.WriteString(book.GetTitle() + "\n" + book.GetWriter() + "\n" + strings.Repeat("-", 20) + "\n\n")
	sortChapters(chapters)
	for _, chapter := range chapters {
		_, err = f.WriteString(chapter.Title + "\n" + strings.Repeat("-", 20) + "\n" +
		chapter.Content + strings.Repeat("\n", 2))
		utils.CheckError(err)
	}
	f.Close()
	return errorChapterCount
}

func (book *Book) Download(MAX_THREAD int, loadContentMutex *sync.Mutex) bool {
	chapters, err := book.getEmptyChapters()
	if err != nil {
		logging.Info("Book %v-%v-%v Download fail: %v", book.bookRecord.Site, book.bookRecord.Id, book.bookRecord.HashCode, err)
		return false
	}
	loadContentMutex.Lock()
	defer loadContentMutex.Unlock()
	results := book.downloadChapters(chapters, MAX_THREAD)
	// save the content to target path
	errorCount := book.saveContent(results)
	maxErrorChapterCount := 50
	if int(float64(len(results))*0.1) < 50 {
		maxErrorChapterCount = int(float64(len(results)) * 0.1)
	}
	if errorCount > maxErrorChapterCount {
		logging.Info("Book %v-%v-%v Download fail: too much chapters return error", book.bookRecord.Site, book.bookRecord.Id, book.bookRecord.HashCode)
		utils.CheckError(os.Remove(book.getContentLocation()))
		return false
	}
	book.SetStatus(database.Download)
	return true
}
