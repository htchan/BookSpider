package model

import (
	// for println and format string
	"fmt"
	// string operation and encoding
	"strings"
	"strconv"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	// concurrency related
	"sync"
	"context"
	"time"
	"golang.org/x/sync/semaphore"
	// read write files
	"io/ioutil"
	// self define helper package
	"../helper"
)

type Book struct {
	SiteName string
	Id int
	Version int
	Title, Writer, Type string
	LastUpdate, LastChapter string
	EndFlag, DownloadFlag, ReadFlag bool
	decoder *encoding.Decoder
	baseUrl, downloadUrl, chapterUrl string
	titleRegex, writerRegex, typeRegex, lastUpdateRegex, lastChapterRegex string
	chapterUrlRegex, chapterTitleRegex string
	chapterContentRegex string
}

// update the book with online info
func (book *Book) Update() (bool) {
	// get online resource, try maximum 10 times if it keeps failed
	var html string
	for i := 0; i < 10; i++ {
		html = helper.GetWeb(book.baseUrl);
		if (len(html) == 0) {
			fmt.Println("retry (" + strconv.Itoa(i) + ")\t" + book.baseUrl);
			time.Sleep(1000)
			continue
		}
		if (book.decoder != nil) {
			html, _, _ = transform.String(book.decoder, html);
			break
		}
	}
	if (len(html) == 0) {
		return false;
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
		lastUpdate != book.LastUpdate || lastChapter != book.LastChapter) {
			update = true;
			if (book.DownloadFlag ||
				title != book.Title || writer != book.Writer || typeName != book.Type ||
				book.Version < 0) {
				book.Version++;
				book.EndFlag = false;
				book.DownloadFlag = false;
				book.ReadFlag = false;
			}
	}
	// sync with online info
	book.Title = title;
	book.Writer = writer;
	book.Type = typeName;
	book.LastUpdate = lastUpdate;
	book.LastChapter = lastChapter;
	return update;
}

type chapter struct {
	Url, Title, Content string
}

func (book *Book) Download(savePath string) (bool) {
	// set up semaphore and routine pool
	ctx := context.Background()
	var s = semaphore.NewWeighted(int64(300))
	var wg sync.WaitGroup
	// get basic info (all chapter url and title)
	html := helper.GetWeb(book.downloadUrl);
	urls := helper.SearchAll(html, book.chapterUrlRegex);
	titles := helper.SearchAll(html, book.chapterTitleRegex);
	// if length are difference, return error
	if (len(urls) != len(titles)) {
		fmt.Println("download error")
		return false
	}
	// use go routine to load chapter content
	ch := make(chan chapter)
	for i, url := range urls {
		wg.Add(1)
		s.Acquire(ctx, 1)
		go book.downloadChapter(url, titles[i], s, &wg, ch)
	}
	// put result of chapter into results
	results := make([]chapter, len(urls))
	for i := range results {
		wg.Add(1)
		results[i] = <-ch
	}
	wg.Wait()
	content := book.Title + "\n" + book.Writer + "\n" + strings.Repeat("-", 20) + strings.Repeat("\n", 2)
	// put chapters content into content with order
	for _, url := range urls {
		for _, chapter := range results {
			if (url == chapter.Url) {
				content = content + chapter.Title + "\n" + strings.Repeat("-", 20) + "\n"
				content = content + chapter.Content + strings.Repeat("\n", 2)
				break
			}
		}
	}
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

	return true
}
func (book *Book) downloadChapter(url, title string, s *semaphore.Weighted, wg *sync.WaitGroup, ch chan<-chapter) () {
	defer wg.Done()
	defer s.Release(1)
	// get chapter resource
	html := helper.GetWeb(url)
	// extract chapter
	content := helper.Search(html, book.chapterContentRegex)
	// put the chapter info to channel
	ch <- chapter{Url:url, Title:title, Content: content}
}

// to string function
func (book Book) String() (string) {
	return book.SiteName + "\t" + strconv.Itoa(book.Id) + "\t" + strconv.Itoa(book.Version) + "\n" +
			book.Title + "\t" + book.Writer;
}