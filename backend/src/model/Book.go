package model

import (
	// for println and format string
	"fmt"
	// string operation and encoding
	"strings"
	"strconv"
	"encoding/json"
	"golang.org/x/text/encoding"
	// concurrency related
	"sync"
	"context"
	"golang.org/x/sync/semaphore"
	// read write files
	//"io/ioutil"
	"os"
	// self define helper package
	"github.com/htchan/BookSpider/helper"
	"io/ioutil"
)
const BOOK_MAX_THREAD = 1000

type Book struct {
	SiteName string
	Id, Version int
	Title, Writer, Type, LastUpdate, LastChapter string
	EndFlag, DownloadFlag, ReadFlag bool
	decoder *encoding.Decoder
	baseUrl, downloadUrl, chapterUrl, chapterPattern string
	titleRegex, writerRegex, typeRegex, lastUpdateRegex, lastChapterRegex string
	chapterUrlRegex, chapterTitleRegex, chapterContentRegex string
}

func (book Book) Log(info map[string]interface{}) {
	info["site"], info["id"], info["version"] = book.SiteName, book.Id, book.Version
	strByte, err := json.Marshal(info)
	helper.CheckError(err)
	fmt.Println(string(strByte))
}

func (book *Book) checkHTML(html string, url string, trial int) bool {
	if (len(html) == 0) {
		book.Log(map[string]interface{} {
			"retry": trial, "url": url, "message": "load html fail - zero length",
		})
		return false
	} else if _, err := strconv.Atoi(html); err == nil {
		book.Log(map[string]interface{} {
			"retry": trial, "url": url, "message": "load html fail - code " + html,
		})
		return false
	} else {
		book.Log(map[string]interface{} {
			"retry": trial, "url": url, "message": "load html success",
		})
	}
	return true
}

// update the book with online info
func (book *Book) Update() bool {
	// get online resource, try maximum 10 times if it keeps failed
	html, trial := helper.GetWeb(book.baseUrl, 10, book.decoder)
	if helper.Search(html, book.titleRegex) == "error" || !book.checkHTML(html, book.baseUrl, trial) {
		return false
	}
	// extract info from source
	update := false
	title := helper.Search(html, book.titleRegex)
	writer := helper.Search(html, book.writerRegex)
	typeName := helper.Search(html, book.typeRegex)
	lastUpdate := helper.Search(html, book.lastUpdateRegex)
	lastChapter := helper.Search(html, book.lastChapterRegex)
	if (title == "error" || writer == "error" || typeName == "error" ||
		lastUpdate== "error" || lastChapter == "error") {
			book.Log(map[string]interface{} {
				"title": title, "writer": writer, "type": typeName,
				"lastUpdate": lastUpdate, "lastChapter": lastChapter,
				"message": "extract html fail",
			})
			return false
		}
	// check difference
	if lastUpdate != book.LastUpdate || lastChapter != book.LastChapter { update = true }
	if (title != book.Title || writer != book.Writer || typeName != book.Type) {
		update = true
		if (book.DownloadFlag) {
			book.Log(map[string]interface{} {
				"old": map[string]interface{} {
					"title": book.Title, "writer": book.Writer, "type": book.Type,
				},
				"new": map[string]interface{} {
					"title": title, "writer": writer, "type": typeName,
				},
				"message": "already download",
			})
		}
		book.Version++
		book.EndFlag, book.DownloadFlag, book.ReadFlag = false, false, false
	}
	if (update) {
		// sync with online info
		book.Title, book.Writer, book.Type = title, writer, typeName
		book.LastUpdate, book.LastChapter = lastUpdate, lastChapter
	}
	return update
}

type Chapter struct {
	Url, Title, Content string
}

func (book Book) saveBook(path string, urls []string, results []Chapter) (errorCount int) {
	f, err := os.Create(path)
	helper.CheckError(err)
	f.WriteString(book.Title + "\n" + book.Writer + "\n" + 
		strings.Repeat("-", 20) + strings.Repeat("\n", 2))
	for _, url := range urls {
		found := false
		for _, chapter := range results {
			if (url == chapter.Url) {
				if chapter.Content == "error" { errorCount += 1 }
				_, err = f.WriteString(chapter.Title + "\n" + strings.Repeat("-", 20) + "\n" +
										chapter.Content + strings.Repeat("\n", 2))
				helper.CheckError(err)
				found = true
				break
			}
		}
		if (!found) {
			book.Log(map[string]interface{} {
				"title": book.Title, "message": "no chapter found",
			})
		}
	}
	f.Close()
	return
}

func (book *Book) Download(savePath string, MAX_THREAD int) bool {
	// get basic info (all chapter url and title)
	html, trial := helper.GetWeb(book.downloadUrl, 10, book.decoder)
	if !book.checkHTML(html, book.downloadUrl, trial) { return false }
	urls := helper.SearchAll(html, book.chapterUrlRegex)
	titles := helper.SearchAll(html, book.chapterTitleRegex)
	// if length are difference, return error
	if len(urls) != len(titles) {
		book.Log(map[string]interface{} {
			"chapterCount": len(urls), "titleCount": len(titles),
			"message": "title and url have different length",
		})
		return false
	}
	if len(urls) == 0 {
		book.Log(map[string]interface{} {
			"title": book.Title, "message": "no chapter found",
		})
		return false
	}
	// use go routine to load chapter content
	// put result of chapter into results
	results := book.downloadAllChapters(urls, titles, MAX_THREAD)
	// save the content to target path
	path := savePath + strconv.Itoa(book.Id)
	if (book.Version != 0) { path = path + "-v" + strconv.Itoa(book.Version); }
	path = path + ".txt"
	errorCount := book.saveBook(path, urls, results)
	maxErrorCount := 50
	if (int(float64(len(results)) * 0.1) < 50) {
		maxErrorCount = int(float64(len(results)) * 0.1)
	}
	if errorCount > maxErrorCount {
		book.Log(map[string]interface{} {
			"title": book.Title,
			"message": "download cancel due to more than " + strconv.Itoa(maxErrorCount) +
			" chapters loss",
		})
		helper.CheckError(os.Remove(path))
		return false
	}
	return true
}
func (book Book) optimizeContent(content string) string {
	content = strings.ReplaceAll(content, "<br />", "")
	content = strings.ReplaceAll(content, "&nbsp;", "")
	content = strings.ReplaceAll(content, "<b>", "")
	content = strings.ReplaceAll(content, "</b>", "")
	content = strings.ReplaceAll(content, "<p>", "")
	content = strings.ReplaceAll(content, "</p>", "")
	content = strings.ReplaceAll(content, "                ", "")
	content = strings.ReplaceAll(content, "<p/>", "\n")
	return content
}
func (book Book) downloadAllChapters(urls []string, titles []string, MAX_THREAD int) (results []Chapter) {
	ctx := context.Background()
	var s = semaphore.NewWeighted(int64(MAX_THREAD))
	var wg sync.WaitGroup
	ch := make(chan Chapter)
	results = make([]Chapter, len(urls))
	var i int
	for i = range urls {
		wg.Add(1)
		s.Acquire(ctx, 1)
		if strings.HasPrefix(urls[i], "/") || strings.HasPrefix(urls[i], "http") {
			urls[i] = fmt.Sprintf(book.chapterUrl, urls[i])
		} else {
			urls[i] = book.downloadUrl + urls[i]
		}
		go book.downloadChapter(urls[i], titles[i], s, &wg, ch)
		if ((i % 100 == 0) && (i > 0)) {
			for j := 0; j < 100; j++ { results[i - 100 + j] = <-ch }
		}
	}
	offset := i % 100
	for j := 0; j <= offset; j++ { results[i - offset + j] = <-ch }
	wg.Wait()
	book.Log(map[string]interface{} {
		"title": book.Title, "message": "all chapter download",
	})
	return
}
func (book *Book) downloadChapter(url, title string, s *semaphore.Weighted, 
	wg *sync.WaitGroup, ch chan<-Chapter) {
	defer wg.Done()
	defer s.Release(1)
	// get chapter resource
	fmt.Println("start download" + title)
	html, trial := helper.GetWeb(url, 10, book.decoder)
	if !book.checkHTML(html, url, trial) {
		ch <- Chapter{Url: url, Title: title, Content: "load html fail"}
		return
	}
	// extract chapter
	content := helper.Search(html, book.chapterContentRegex)
	if content == "error" {
		book.Log(map[string]interface{} {
			"retry": trial, "url": url, "message": "recognize html fail",
		})
		ch <- Chapter{Url: url, Title: title, Content: "recognize html fail\n"+html}
	} else {
		content = book.optimizeContent(content)
		book.Log(map[string]interface{} {
			"chapter": title, "url": url, "message": "download success",
		})
		// put the chapter info to channel
		ch <- Chapter{Url:url, Title:title, Content: content}
	}
}

func (book Book) Content(location string) string {
	if book.Title == "" || !book.DownloadFlag { return "" }
	path := location + "/" + strconv.Itoa(book.Id)
	if book.Version > 0 { path += "-v" + strconv.Itoa(book.Version) }
	path += ".txt"
	content, err := ioutil.ReadFile(path)
	helper.CheckError(err)
	return string(content)
}

// to string function
func (book Book) String() string {
	return book.SiteName + "\t" + strconv.Itoa(book.Id) + "\t" + strconv.Itoa(book.Version) + "\n" + 
			book.Title + "\t" + book.Writer + "\n"+ book.LastUpdate + "\t" + book.LastChapter
}

func(book Book) Map() map[string]interface{} {
	return map[string]interface{} {
		"site": book.SiteName,
		"id": book.Id,
		"version": book.Version,
		"title": book.Title,
		"writer": book.Writer,
		"type": book.Type,
		"update": book.LastUpdate,
		"chapter": book.LastChapter,
		"end": book.EndFlag,
		"read": book.ReadFlag,
		"download": book.DownloadFlag,
	}
}
