package books

import (
	"os"
	"github.com/htchan/BookSpider/internal/utils"
	"golang.org/x/sync/semaphore"
	"sync"
	"strings"
	"strconv"
	"fmt"
	"log"
	"context"
	"errors"
)

type Chapter struct {
	Url, Title, Content string
}

func (book Book) getChapterUrlsTitles() ([]string, []string, error) {
	// get basic info (all chapter url and title)
	html, trial := utils.GetWeb(book.metaInfo.downloadUrl, 10, book.decoder)
	if !book.validHTML(html, book.metaInfo.downloadUrl, trial) {
		book.Log(map[string]interface{} {
			"title": book.Title, "error": "invalid table of contents html", "stage": "download",
		})
		return nil, nil, errors.New("invalid table of content html")
	}
	urls := utils.SearchAll(html, book.metaInfo.chapterUrlRegex)
	titles := utils.SearchAll(html, book.metaInfo.chapterTitleRegex)
	// if length are difference, return error
	if len(urls) != len(titles) {
		book.Log(map[string]interface{} {
			"chapterCount": len(urls), "titleCount": len(titles),
			"error": "title and url have different length", "stage": "download",
		})
		return nil, nil, errors.New("title and url have different length")
	} else if len(urls) == 0 {
		book.Log(map[string]interface{} {
			"title": book.Title, "error": "no chapter found", "stage": "download",
		})
		return nil, nil, errors.New("no chapter found")
	}
	return urls, titles, nil
}

func (book Book) optimizeContent(content string) string {
	removeList := [][2]string{
		[2]string{ "<br />", "" },
		[2]string{ "&nbsp;", "" },
		[2]string{ "<b>", "" },
		[2]string{ "</b>", "" },
		[2]string{ "<p>", "" },
		[2]string{ "</p>", "" },
		[2]string{ "                ", "" },
		[2]string{ "<p/>", "\n" },
	}
	for _, removeItem := range removeList {
		content = strings.ReplaceAll(content, removeItem[0], removeItem[1])
	}
	return content
}


func (book *Book) downloadChapter(url, title string, s *semaphore.Weighted, 
	wg *sync.WaitGroup, ch chan<-Chapter) {
	defer wg.Done()
	defer s.Release(1)
	// get chapter resource
	html, trial := utils.GetWeb(url, 10, book.decoder)
	if !book.validHTML(html, url, trial) {
		ch <- Chapter{Url: url, Title: title, Content: "load html fail"}
		return
	}
	// extract chapter
	chapterContent, err := utils.Search(html, book.metaInfo.chapterContentRegex)
	if err != nil {
		book.Log(map[string]interface{} {
			"retry": trial, "url": url, "error": "recognize html fail", "stage": "download",
		})
		ch <- Chapter{Url: url, Title: title, Content: "recognize html fail\n"+html}
	} else {
		chapterContent = book.optimizeContent(chapterContent)
		book.Log(map[string]interface{} {
			"chapter": title, "url": url, "message": "download success", "stage": "download",
		})
		ch <- Chapter{Url:url, Title:title, Content: chapterContent}
	}
}
func (book Book) downloadChapters(urls []string, titles []string, MAX_THREAD int) (results []Chapter) {
	ctx := context.Background()
	var s = semaphore.NewWeighted(int64(MAX_THREAD))
	var wg sync.WaitGroup
	ch := make(chan Chapter, 100)
	results = make([]Chapter, len(urls))
	var i int
	for i = range urls {
		wg.Add(1)
		s.Acquire(ctx, 1)
		//TODO: check the reason of adding the http here
		if strings.HasPrefix(urls[i], "/") || strings.HasPrefix(urls[i], "http") {
			urls[i] = fmt.Sprintf(book.metaInfo.chapterUrl, urls[i])
		} else {
			urls[i] = book.metaInfo.downloadUrl + urls[i]
		}
		go book.downloadChapter(urls[i], titles[i], s, &wg, ch)
		// read channel (pipe) when it should be full
		if i % 100 == 0 && i > 0 {
			for j := 0; j < 100; j++ { 
				results[i - 100 + j] = <-ch
				log.Println(j, results[i - 100 + j].Title)
			}
		}
	}
	// read remaining objects in channel (pipe)
	offset := i % 100
	for j := 0; j <= offset; j++ {
		results[i - offset + j] = <-ch
		log.Println(j, results[i - offset + j].Title)
	}
	wg.Wait()
	book.Log(map[string]interface{} {
		"title": book.Title, "message": "all chapter download", "stage": "download",
	})
	return
}

func (book Book) saveBook(path string, urls []string, chapters []Chapter) int {
	errorChapterCount := 0
	f, err := os.Create(path)
	utils.CheckError(err)
	f.WriteString(book.Title + "\n" + book.Writer + "\n" + strings.Repeat("-", 20) + "\n\n")
	for _, url := range urls {
		found := false
		for _, chapter := range chapters {
			if (url == chapter.Url) {
				if chapter.Content == "error" { errorChapterCount += 1 }
				_, err = f.WriteString(chapter.Title + "\n" + strings.Repeat("-", 20) + "\n" +
										chapter.Content + strings.Repeat("\n", 2))
				utils.CheckError(err)
				found = true
				break
			}
		}
		if (!found) {
			book.Log(map[string]interface{} {
				"title": book.Title, "url": url, "error": "no chapter found", "stage": "save book",
			})
		}
	}
	f.Close()
	return errorChapterCount
}

func (book *Book) Download(storagePath string, MAX_THREAD int) bool {
	urls, titles, err := book.getChapterUrlsTitles()
	if err != nil {
		return false
	}
	results := book.downloadChapters(urls, titles, MAX_THREAD)
	// save the content to target path
	bookLocation := book.storageLocation(storagePath)
	errorCount := book.saveBook(bookLocation, urls, results)
	maxErrorChapterCount := 50
	if (int(float64(len(results)) * 0.1) < 50) {
		maxErrorChapterCount = int(float64(len(results)) * 0.1)
	}
	if errorCount > maxErrorChapterCount {
		book.Log(map[string]interface{} {
			"title": book.Title,
			"error": "download cancel due to more than " + strconv.Itoa(maxErrorChapterCount) +
			" chapters loss", "stage": "download",
		})
		utils.CheckError(os.Remove(bookLocation))
		return false
	}
	return true
}
