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
var BOOK_MAX_THREAD int = 1000;

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
	for i := 0; i < 10; i++ {
		html = helper.GetWeb(book.baseUrl);
		if (len(html) == 0) || (helper.Search(html, book.titleRegex) == "error") {
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
		lastUpdate != book.LastUpdate || lastChapter != book.LastChapter) && 
		(title != "error" && writer != "error" && typeName != "error" &&
		lastUpdate!= "error" && lastChapter != "error") {
			update = true;
			if (book.DownloadFlag ||
				title != book.Title || writer != book.Writer || typeName != book.Type ||
				book.Version < 0) {
				if (book.DownloadFlag) {
					fmt.Println(book.SiteName + "\t" + strconv.Itoa(book.Id) + "\t is already downloaded")
				}
				book.Version++;
				book.EndFlag = false;
				book.DownloadFlag = false;
				book.ReadFlag = false;
			}
	}
	if (title == "error" || writer == "error" || typeName == "error" ||
		lastUpdate== "error" || lastChapter == "error") {
			fmt.Println(title+"\t"+writer+"\t"+typeName)
			fmt.Println(lastUpdate+"\t"+lastChapter)
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

func (book *Book) Download(savePath string) (bool) {
	// set up semaphore and routine pool
	ctx := context.Background()
	var s = semaphore.NewWeighted(int64(BOOK_MAX_THREAD))
	var wg sync.WaitGroup
	// get basic info (all chapter url and title)
	var html string;
	for i := 0; i < 10; i++ {
		html = helper.GetWeb(book.downloadUrl);
		if (len(html) == 0) {
			fmt.Println("retry download info (" + strconv.Itoa(i) + ")\t" + book.downloadUrl);
			continue
		}
		if (book.decoder != nil) {
			html, _, _ = transform.String(book.decoder, html);
			break;
		}
	}

	urls := helper.SearchAll(html, book.chapterUrlRegex);
	titles := helper.SearchAll(html, book.chapterTitleRegex);
	// if length are difference, return error
	if (len(urls) != len(titles)) {
		fmt.Println("download error")
		return false
	}
	// use go routine to load chapter content
	// put result of chapter into results
	ch := make(chan chapter)
	results := make([]chapter, len(urls))
	var i int
	for i = range urls {
		wg.Add(1)
		s.Acquire(ctx, 1)
		urls[i] = fmt.Sprintf(book.chapterUrl, urls[i]);
		go book.downloadChapter(urls[i], titles[i], s, &wg, ch)
		if ((i % 100 == 0) && (i > 0)) {
			for j := 0; j < 100; j++ {
				results[i - 100 + j] = <-ch
			}

		}
	}
	fmt.Println(i)
	offset := i % 100;
	for j := -1; j < offset; j++ {
		results[i - offset + j] = <-ch
	}
	wg.Wait()
	fmt.Println("finish")
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
	var html string;
	for i := 0; i < 10; i++ {
		html = helper.GetWeb(url);
		if (len(html) == 0) {
			fmt.Println("retry (" + strconv.Itoa(i) + ")\t" + url + "\t" + title);
			continue
		}
		if (book.decoder != nil) {
			html, _, _ = transform.String(book.decoder, html);
			break;
		}
	}
	fmt.Println(url + "\t" + title);
	// extract chapter
	content := helper.Search(html, book.chapterContentRegex)
	content = strings.Replace(content, "<br />", "", -1);
	content = strings.Replace(content, "&nbsp;", "", -1);
	content = strings.Replace(content, "<b>", "", -1);
	content = strings.Replace(content, "</b>", "", -1);
	content = strings.Replace(content, "<p>", "", -1);
	content = strings.Replace(content, "</p>", "", -1);
	// put the chapter info to channel
	ch <- chapter{Url:url, Title:title, Content: content}
	return
}

// to string function
func (book Book) String() (string) {
	return book.SiteName + "\t" + strconv.Itoa(book.Id) + "\t" + strconv.Itoa(book.Version) + "\n" +
			book.Title + "\t" + book.Writer + "\n"+
			book.LastUpdate + "\t" + book.LastChapter;
}
