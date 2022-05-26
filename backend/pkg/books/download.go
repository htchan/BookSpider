package books

import (
	"context"
	"errors"
	"fmt"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/internal/logging"
	"github.com/htchan/ApiParser"
	"golang.org/x/sync/semaphore"
	// "log"
	"os"
	// "strconv"
	"strings"
	"sync"
)

func (book Book) getEmptyChapters() (chapters []Chapter, err error) {
	// get basic info (all chapter url and title)
	html, trial := utils.GetWeb(book.config.DownloadUrl, 10, book.config.Decoder, book.config.ConstSleep)
	if trial > 0 {
		// logging.LogBookEvent(book.String(), "source_download_content", "trial", trial)
	}
	if err = book.validHTML(html); err != nil {
		err = errors.New(fmt.Sprintf("invalid table of content html: %v", html))
		return
	}

	responseApi := ApiParser.Parse(book.config.SourceKey + ".info", html)
	urls, titles := make([]string, len(responseApi.Items)), make([]string, len(responseApi.Items))
	for i, item := range responseApi.Items {
		urls[i] = item["ChapterUrl"]
		titles[i] = item["ChapterTitle"]
	}
	// // if length are difference, return error
	// if len(urls) != len(titles) {
	// 	err = errors.New("title and url have different length")
	// 	return
	// } else 
	if len(urls) == 0 {
		err = errors.New("no chapter found")
		return
	}
	chapters = make([]Chapter, len(urls))
	for i := 0; i < len(urls); i++ {
		chapters[i] = NewChapter(i, urls[i], titles[i], &book.config)
	}
	return
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
		if book.config.UseRequestInterval { utils.RequestInterval() }
	}
	wg.Wait()
	return chapters
}

func (book Book) saveContent(storageDirectory string, chapters []Chapter) int {
	errorChapterCount := 0
	f, err := os.Create(book.getContentLocation(storageDirectory))
	utils.CheckError(err)
	f.WriteString(book.GetTitle() + "\n" + book.GetWriter() + "\n" + strings.Repeat("-", 20) + "\n\n")
	sortChapters(chapters)
	for _, chapter := range chapters {
		_, err = f.WriteString(chapter.Title + "\n" + strings.Repeat("-", 20) + "\n" +
		chapter.Content + strings.Repeat("\n", 2))
		utils.CheckError(err)
		if strings.Contains(chapter.Content, "fail") {
			errorChapterCount += 1
		}
	}
	f.Close()
	return errorChapterCount
}

func (book *Book) Download(storageDirectory string, MAX_THREAD int, loadContentMutex *sync.Mutex) bool {
	chapters, err := book.getEmptyChapters()
	if err != nil {
		logging.LogBookEvent(book.String(), "download", "fail", err)
		return false
	}
	loadContentMutex.Lock()
	defer loadContentMutex.Unlock()
	results := book.downloadChapters(chapters, MAX_THREAD)
	// save the content to target path
	errorCount := book.saveContent(storageDirectory, results)
	maxErrorChapterCount := 50
	if int(float64(len(results))*0.1) < 50 {
		maxErrorChapterCount = int(float64(len(results)) * 0.1)
	}
	if errorCount > maxErrorChapterCount {
		logging.LogBookEvent(book.String(), "download", "fail", "too much chapters return error")
		utils.CheckError(os.Remove(book.getContentLocation(storageDirectory)))
		return false
	}
	book.SetStatus(database.Download)
	return true
}
