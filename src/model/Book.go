package model

import (
	// for println and format string
	"fmt"
	// string operation and encoding
	"strings"
	"strconv"
	"encoding/json"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	// concurrency related
	"sync"
	"context"
	"time"
	"golang.org/x/sync/semaphore"
	// read write files
	//"io/ioutil"
	"os"
	// self define helper package
	"../helper"
)
const BOOK_MAX_THREAD = 1000;

type Book struct {
	SiteName string
	Id int
	Version int
	Title, Writer, Type string
	LastUpdate, LastChapter string
	EndFlag, DownloadFlag, ReadFlag bool
	decoder *encoding.Decoder
	baseUrl, downloadUrl, chapterUrl, chapterPattern string
	titleRegex, writerRegex, typeRegex, lastUpdateRegex, lastChapterRegex string
	chapterUrlRegex, chapterTitleRegex string
	chapterContentRegex string
}

// update the book with online info
func (book *Book) Update() (bool) {
	// get online resource, try maximum 10 times if it keeps failed
	var html string
	var i int;
	for i = 0; i < 10; i++ {
		html = helper.GetWeb(book.baseUrl);
		if _, err := strconv.Atoi(html); err == nil || (len(html) == 0) || (helper.Search(html, book.titleRegex) == "error") {
			time.Sleep(time.Duration(i * i) * time.Second)
			continue
		}
		if (book.decoder != nil) {
			html, _, _ = transform.String(book.decoder, html);
			break
		}
	}
	if (len(html) == 0) {
		strByte, err := json.Marshal(map[string]interface{} {
			"site": book.SiteName,
			"id": book.Id,
			"version": book.Version,
			"retry": i,
			"url": book.baseUrl,
			"message": "load html fail",
		})
		helper.CheckError(err)
		fmt.Println(string(strByte))
		return false;
	} else if _, err := strconv.Atoi(html); err == nil {
		strByte, err := json.Marshal(map[string]interface{} {
			"site": book.SiteName,
			"id": book.Id,
			"version": book.Version,
			"retry": i,
			"url": book.baseUrl,
			"message": "load html fail - code " + html,
		})
		helper.CheckError(err)
		fmt.Println(string(strByte))
		return false
	} else {
		strByte, err := json.Marshal(map[string]interface{} {
			"site": book.SiteName,
			"id": book.Id,
			"version": book.Version,
			"retry": i,
			"url": book.baseUrl,
			"message": "load html success",
		})
		helper.CheckError(err)
		fmt.Println(string(strByte))
	}
	// extract info from source
	update := false;
	title := helper.Search(html, book.titleRegex);
	writer := helper.Search(html, book.writerRegex);
	typeName := helper.Search(html, book.typeRegex);
	lastUpdate := helper.Search(html, book.lastUpdateRegex);
	lastChapter := helper.Search(html, book.lastChapterRegex);
	// check difference
	if (title != book.Title || writer != book.Writer || typeName != book.Type || 
		lastUpdate != book.LastUpdate || lastChapter != book.LastChapter) && 
		(title != "error" && writer != "error" && typeName != "error" &&
		lastUpdate!= "error" && lastChapter != "error") {
			update = true;
			if (book.DownloadFlag ||
				title != book.Title || writer != book.Writer || typeName != book.Type ||
				book.Version < 0) {
				if (book.DownloadFlag) {
					strByte, err := json.Marshal(map[string]interface{} {
						"site": book.SiteName,
						"id": book.Id,
						"version": book.Version,
						"old": map[string]interface{} {
							"title": book.Title,
							"writer": book.Writer,
							"type": book.Type,
						},
						"new": map[string]interface{} {
							"title": title,
							"writer": writer,
							"type": typeName,
						},
						"message": "already download",
					})
					helper.CheckError(err)
					fmt.Println(string(strByte))
				}
				book.Version++;
				book.EndFlag = false;
				book.DownloadFlag = false;
				book.ReadFlag = false;
			}
	}
	if (title == "error" || writer == "error" || typeName == "error" ||
		lastUpdate== "error" || lastChapter == "error") {
			strByte, err := json.Marshal(map[string]interface{} {
				"site": book.SiteName,
				"id": book.Id,
				"version": book.Version,
				"title": title,
				"writer": writer,
				"type": typeName,
				"lastUpdate": lastUpdate,
				"lastChapter": lastChapter,
				"message": "extract html fail",
			})
			helper.CheckError(err)
			fmt.Println(string(strByte))
		}
	if (update) {
		// sync with online info
		book.Title = title;
		book.Writer = writer;
		book.Type = typeName;
		book.LastUpdate = lastUpdate;
		book.LastChapter = lastChapter;
	}
	return update;
}

type chapter struct {
	Url, Title, Content string
}

func (book *Book) Download(savePath string, MAX_THREAD int) (bool) {
	// set up semaphore and routine pool
	ctx := context.Background()
	var s = semaphore.NewWeighted(int64(MAX_THREAD))
	var wg sync.WaitGroup
	// get basic info (all chapter url and title)
	var html string;
	var i int;
	for i = 0; i < 10; i++ {
		html = helper.GetWeb(book.downloadUrl);
		if _, err := strconv.Atoi(html); err == nil || (len(html) == 0) {
			time.Sleep(time.Duration(i * i) * time.Second)
			continue
		}
		if (book.decoder != nil) {
			html, _, _ = transform.String(book.decoder, html);
			break;
		}
	}
	if (len(html) == 0) {
		strByte, err := json.Marshal(map[string]interface{} {
			"site": book.SiteName,
			"id": book.Id,
			"version": book.Version,
			"retry": i,
			"url": book.downloadUrl,
			"message": "load html fail",
		})
		helper.CheckError(err)
		fmt.Println(string(strByte))
		return false;
	} else if _, err := strconv.Atoi(html); err == nil {
		strByte, err := json.Marshal(map[string]interface{} {
			"site": book.SiteName,
			"id": book.Id,
			"version": book.Version,
			"retry": i,
			"url": book.downloadUrl,
			"message": "load html fail - code " + html,
		})
		helper.CheckError(err)
		fmt.Println(string(strByte))
		return false
	} else {
		strByte, err := json.Marshal(map[string]interface{} {
			"site": book.SiteName,
			"id": book.Id,
			"version": book.Version,
			"retry": i,
			"url": book.downloadUrl,
			"message": "load html success",
		})
		helper.CheckError(err)
		fmt.Println(string(strByte))
	}
	urls := helper.SearchAll(html, book.chapterUrlRegex);
	titles := helper.SearchAll(html, book.chapterTitleRegex);
	// if length are difference, return error
	if (len(urls) != len(titles)) {
		strByte, err := json.Marshal(map[string]interface{} {
			"site": book.SiteName,
			"id": book.Id,
			"version": book.Version,
			"chapterCount": len(urls),
			"titleCount": len(titles),
			"message": "title and url have different length",
		})
		helper.CheckError(err)
		fmt.Println(string(strByte))
		return false
	}
	if len(urls) == 0 {
		strByte, err := json.Marshal(map[string]interface{} {
			"site": book.SiteName,
			"id": book.Id,
			"version": book.Version,
			"title": book.Title,
			"message": "no chapter found",
		})
		helper.CheckError(err)
		fmt.Println(string(strByte))
		return false
	}
	// use go routine to load chapter content
	// put result of chapter into results
	ch := make(chan chapter)
	results := make([]chapter, len(urls))
	for i = range urls {
		wg.Add(1)
		s.Acquire(ctx, 1)
		if strings.HasPrefix(urls[i], "/") || strings.HasPrefix(urls[i], "http") {
			urls[i] = fmt.Sprintf(book.chapterUrl, urls[i]);
		} else {
			urls[i] = book.downloadUrl + urls[i]
		}
		go book.downloadChapter(urls[i], titles[i], s, &wg, ch)
		if ((i % 100 == 0) && (i > 0)) {
			for j := 0; j < 100; j++ {
				results[i - 100 + j] = <-ch
			}

		}
	}
	offset := i % 100;
	for j := 0; j <= offset; j++ {
		results[i - offset + j] = <-ch
	}
	wg.Wait()
	strByte, err := json.Marshal(map[string]interface{} {
		"site": book.SiteName,
		"id": book.Id,
		"version": book.Version,
		"title": book.Title,
		"message": "all chapter download",
	})
	helper.CheckError(err)
	fmt.Println(string(strByte))
	errorCount := 0
	// save the content to target path
	path := savePath + strconv.Itoa(book.Id);
	if (book.Version == 0) {
		path = path + ".txt";
	} else {
		path = path + "-v" + strconv.Itoa(book.Version) + ".txt";
	}
	f, err := os.Create(path)
	helper.CheckError(err)
	f.WriteString(book.Title + "\n" + book.Writer + "\n" + 
		strings.Repeat("-", 20) + strings.Repeat("\n", 2))
	// put chapters content into content with order
	for _, url := range urls {
		found := false
		for _, chapter := range results {
			if (url == chapter.Url) {
				if chapter.Content == "error" {
					errorCount += 1
				}
				_, err = f.WriteString(chapter.Title + "\n" + strings.Repeat("-", 20) + "\n")
				helper.CheckError(err)
				_, err = f.WriteString(chapter.Content + strings.Repeat("\n", 2))
				helper.CheckError(err)
				found = true
				break
			}
		}
		if (!found) {
			strByte, err := json.Marshal(map[string]interface{} {
				"site": book.SiteName,
				"id": book.Id,
				"version": book.Version,
				"title": book.Title,
				"message": "no chapter found",
			})
			helper.CheckError(err)
			fmt.Println(string(strByte))
		}
	}
	f.Close()
	maxErrorCount := 50
	if (int(float64(len(results)) * 0.1) < 50) {
		maxErrorCount = int(float64(len(results)) * 0.1);
	}
	if errorCount > maxErrorCount {
		strByte, err := json.Marshal(map[string]interface{} {
			"site": book.SiteName,
			"id": book.Id,
			"version": book.Version,
			"title": book.Title,
			"message": "download cancel due to more than " + strconv.Itoa(maxErrorCount) +
			" chapters loss",
		})
		helper.CheckError(err)
		fmt.Println(string(strByte))
		err = os.Remove(path)
		helper.CheckError(err)
		return false
	}
	/*
	// save the content to target path
	path := savePath + strconv.Itoa(book.Id);
	if (book.Version == 0) {
		path = path + ".txt";
	} else {
		path = path + "-v" + strconv.Itoa(book.Version) + ".txt";
	}
	bytes := []byte(content);
	err := ioutil.WriteFile(path, bytes, 0644);
	helper.CheckError(err);
	*/
	return true
}
func (book *Book) downloadChapter(url, title string, s *semaphore.Weighted, wg *sync.WaitGroup, ch chan<-chapter) () {
	defer wg.Done()
	defer s.Release(1)
	// get chapter resource
	var html string;
	var i int;
	fmt.Println("start download" + title)
	for i = 0; i < 10; i++ {
		html = helper.GetWeb(url);
		if _, err := strconv.Atoi(html); err == nil || (len(html) == 0) {
			time.Sleep(time.Duration(i * i) * time.Second)
			continue
		}
		if (book.decoder != nil) {
			html, _, _ = transform.String(book.decoder, html);
			break;
		}
	}
	if (len(html) == 0) {
		strByte, err := json.Marshal(map[string]interface{} {
			"site": book.SiteName,
			"id": book.Id,
			"version": book.Version,
			"retry": i,
			"url": url,
			"message": "load html fail",
		})
		helper.CheckError(err)
		fmt.Println(string(strByte))
		ch <- chapter{Url: url, Title: title, Content: "load html fail"}
		return;
	} else if _, err := strconv.Atoi(html); err == nil {
		strByte, err := json.Marshal(map[string]interface{} {
			"site": book.SiteName,
			"id": book.Id,
			"version": book.Version,
			"retry": i,
			"url": url,
			"message": "load html fail - code " + html,
		})
		helper.CheckError(err)
		fmt.Println(string(strByte))
		ch <- chapter{Url: url, Title: title, Content: "load html fail - code " + html}
		return
	} else {
		strByte, err := json.Marshal(map[string]interface{} {
			"site": book.SiteName,
			"id": book.Id,
			"version": book.Version,
			"retry": i,
			"url": url,
			"message": "load html success",
		})
		helper.CheckError(err)
		fmt.Println(string(strByte))
	}
	// extract chapter
	content := helper.Search(html, book.chapterContentRegex)
	if content == "error" {
		strByte, err := json.Marshal(map[string]interface{} {
			"site": book.SiteName,
			"id": book.Id,
			"version": book.Version,
			"retry": i,
			"url": url,
			"message": "recognize html fail",
		})
		helper.CheckError(err)
		fmt.Println(string(strByte))
		ch <- chapter{Url: url, Title: title, Content: "recognize html fail\n"+html}
		return
	} else {
		content = strings.ReplaceAll(content, "<br />", "");
		content = strings.ReplaceAll(content, "&nbsp;", "");
		content = strings.ReplaceAll(content, "<b>", "");
		content = strings.ReplaceAll(content, "</b>", "");
		content = strings.ReplaceAll(content, "<p>", "");
		content = strings.ReplaceAll(content, "</p>", "");
		content = strings.ReplaceAll(content, "                ", "")
		content = strings.ReplaceAll(content, "<p/>", "\n")
		strByte, err := json.Marshal(map[string]interface{} {
			"site": book.SiteName,
			"id": book.Id,
			"version": book.Version,
			"chapter": title,
			"url": url,
			"message": "download success",
		})
		helper.CheckError(err)
		fmt.Println(string(strByte))
	
		// put the chapter info to channel
		ch <- chapter{Url:url, Title:title, Content: content}
		return
	}
}

// to string function
func (book Book) String() (string) {
	return book.SiteName + "\t" + strconv.Itoa(book.Id) + "\t" + strconv.Itoa(book.Version) + "\n" +
			book.Title + "\t" + book.Writer + "\n"+
			book.LastUpdate + "\t" + book.LastChapter;
}

func (book Book) JsonString() (string) {
	resultByte, err := json.Marshal(map[string]interface{} {
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
	})
	helper.CheckError(err)
	return string(resultByte)
}
