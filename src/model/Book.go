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
		if (len(html) == 0) || (helper.Search(html, book.titleRegex) == "error") {
			time.Sleep(1000)
			continue
		}
		if (book.decoder != nil) {
			html, _, _ = transform.String(book.decoder, html);
			break
		}
	}
	if (len(html) == 0) {
		fmt.Println("{\"site\":\"" + book.SiteName + "\", " +
					"\"id\":" + strconv.Itoa(book.Id) + ", " +
					"\"version\":" + strconv.Itoa(book.Version) + ", " +
					"\"retry\":" + strconv.Itoa(i) + ", " +
					"\"url\":\"" + book.baseUrl + "\", " +
					"\"result\":\"load html fail\"}");
		return false;
	} else {
		fmt.Println("{\"site\":\"" + book.SiteName + "\", " +
					"\"id\":" + strconv.Itoa(book.Id) + ", " +
					"\"version\":" + strconv.Itoa(book.Version) + ", " +
					"\"retry\":" + strconv.Itoa(i) + ", " +
					"\"url\":\"" + book.baseUrl + "\", " +
					"\"result\":\"load html success\"}");
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
					fmt.Println("{\"site\":\"" + book.SiteName + "\", " +
								"\"id\":" + strconv.Itoa(book.Id) + ", " +
								"\"version\":" + strconv.Itoa(book.Version) + ", " +
								"\"old : {title\":\"" + book.Title + "\", " + 
								"\"writer\":\"" + book.Writer + "\", " + 
								"\"type\":\"" + book.Type + "\"}, " + 
								"\"new : {title\":\"" + title + "\", " + 
								"\"writer\":\"" + writer + "\", " + 
								"\"type\":\"" + typeName + "\"}, " + 
								"\"message\":\"already download\"}");
				}
				book.Version++;
				book.EndFlag = false;
				book.DownloadFlag = false;
				book.ReadFlag = false;
			}
	}
	if (title == "error" || writer == "error" || typeName == "error" ||
		lastUpdate== "error" || lastChapter == "error") {
			fmt.Println("{\"site\":\"" + book.SiteName + "\", " + 
						"\"id\":" + strconv.Itoa(book.Id) + "," +
						"\"version\":" + strconv.Itoa(book.Version) + ", " +
						"\"title\":\"" + title + "\", \"writer\":\"" + writer + "\", " +
						"\"typeName\" : \"" + typeName + "\", " +
						"\"lastUpdate\":\"" + lastUpdate + "\", " +
						"\"lastChapter\":\"" + lastChapter + "\"}")
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
	var i int;
	for i = 0; i < 10; i++ {
		html = helper.GetWeb(book.downloadUrl);
		if (len(html) == 0) {
			continue
		}
		if (book.decoder != nil) {
			html, _, _ = transform.String(book.decoder, html);
			break;
		}
	}
	if (len(html) == 0) {
		fmt.Println("{\"site\":\"" + book.SiteName + "\", " +
					"\"id\":" + strconv.Itoa(book.Id) + ", " +
					"\"version\":" + strconv.Itoa(book.Version) + ", " +
					"\"retry\":" + strconv.Itoa(i) + ", " +
					"\"url\":\"" + book.baseUrl + "\", " +
					"\"result\":\"load html fail\"}");
		return false;
	} else {
		fmt.Println("{\"site\":\"" + book.SiteName + "\", " +
					"\"id\":" + strconv.Itoa(book.Id) + ", " +
					"\"version\":" + strconv.Itoa(book.Version) + ", " +
					"\"retry\":" + strconv.Itoa(i) + ", " +
					"\"url\":\"" + book.baseUrl + "\", " +
					"\"result\":\"load html success\"}");
	}
	urls := helper.SearchAll(html, book.chapterUrlRegex);
	titles := helper.SearchAll(html, book.chapterTitleRegex);
	// if length are difference, return error
	if (len(urls) != len(titles)) {
		fmt.Println("{\"site\":\"" + book.SiteName + "\", " +
					"\"id\":" + strconv.Itoa(book.Id) + ", " +
					"\"version\":" + strconv.Itoa(book.Version) + ", " +
					"\"chapterCount\":" + strconv.Itoa(len(urls)) + ", " +
					"\"titleCount\"" + strconv.Itoa(len(titles)) + ", " +
					"\"message\":\"title and url have different length\"}")
		return false
	}
	if len(urls) == 0 {
		fmt.Println("{\"site\":\"" + book.SiteName + "\", " +
					"\"id\":" + strconv.Itoa(book.Id) + ", " +
					"\"version\":" + strconv.Itoa(book.Version) + ", " +
					"\"title\":\"" + book.Title + "\", " + 
					"\"message\":\"no chapter found\"}");
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
	fmt.Println("{\"site\":\"" + book.SiteName + "\", " +
				"\"id\":" + strconv.Itoa(book.Id) + ", " +
				"\"version\":" + strconv.Itoa(book.Version) + ", " +
				"\"title\":\"" + book.Title + "\", " + 
				"\"message\":\"all chapter download\"}");
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
			fmt.Println("{\"site\":\"" + book.SiteName + "\", " +
						"\"id\":" + strconv.Itoa(book.Id) + ", " +
						"\"version\":" + strconv.Itoa(book.Version) + ", " +
						"\"title\":\"" + book.Title + "\", " + 
						"\"message\":\"no chapter found\"}");
		}
	}
	f.Close()
	maxErrorCount := 50
	if (int(float64(len(results)) * 0.1) < 50) {
		maxErrorCount = int(float64(len(results)) * 0.1);
	}
	if errorCount > maxErrorCount {
		fmt.Println("{\"site\":\"" + book.SiteName + "\", " +
					"\"id\":" + strconv.Itoa(book.Id) + ", " +
					"\"version\":" + strconv.Itoa(book.Version) + ", " +
					"\"title\":\"" + book.Title + "\", " + 
					"\"message\":\"download cancel due to more than " + strconv.Itoa(maxErrorCount) +
					" chapters loss\"}");
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
	for i = 0; i < 10; i++ {
		html = helper.GetWeb(url);
		if (len(html) == 0) {
			continue
		}
		if (book.decoder != nil) {
			html, _, _ = transform.String(book.decoder, html);
			break;
		}
	}
	if (len(html) == 0) {
		fmt.Println("{\"site\":\"" + book.SiteName + "\", " +
					"\"id\":" + strconv.Itoa(book.Id) + ", " +
					"\"version\":" + strconv.Itoa(book.Version) + ", " +
					"\"retry\":" + strconv.Itoa(i) + ", " +
					"\"url\":\"" + book.baseUrl + "\", " +
					"\"result\":\"load html fail\"}");
		return;
	} else {
		fmt.Println("{\"site\":\"" + book.SiteName + "\", " +
					"\"id\":" + strconv.Itoa(book.Id) + ", " +
					"\"version\":" + strconv.Itoa(book.Version) + ", " +
					"\"retry\":" + strconv.Itoa(i) + ", " +
					"\"url\":\"" + book.baseUrl + "\", " +
					"\"result\":\"load html success\"}");
	}
	// extract chapter
	content := helper.Search(html, book.chapterContentRegex)
	if content == "error" {
		fmt.Println("{\"site\":\"" + book.SiteName + "\", " +
					"\"id\":" + strconv.Itoa(book.Id) + ", " +
					"\"version\":" + strconv.Itoa(book.Version) + ", " +
					"\"url\":\"" + book.baseUrl + "\", " +
					"\"result\":\"recognize html fail\"}");
	} else {
		content = strings.ReplaceAll(content, "<br />", "");
		content = strings.ReplaceAll(content, "&nbsp;", "");
		content = strings.ReplaceAll(content, "<b>", "");
		content = strings.ReplaceAll(content, "</b>", "");
		content = strings.ReplaceAll(content, "<p>", "");
		content = strings.ReplaceAll(content, "</p>", "");
		content = strings.ReplaceAll(content, "                ", "")
		content = strings.ReplaceAll(content, "<p/>", "\n")
	}
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

func (book Book) JsonString() (string) {
	return "{"+
		"\"site\" : \""+book.SiteName+"\", "+
		"\"num\" : \""+strconv.Itoa(book.Id)+"\", "+
		"\"version\" : \""+strconv.Itoa(book.Version)+"\", "+
		"\"title\" : \""+book.Title+"\", "+
		"\"writer\" : \""+book.Writer+"\", "+
		"\"type\" : \""+book.Type+"\", "+
		"\"update\" : \""+book.LastUpdate+"\", "+
		"\"chapter\" : \""+book.LastChapter+"\", "+
		"\"end\" : "+strconv.FormatBool(book.EndFlag)+", "+
		"\"read\" : "+strconv.FormatBool(book.ReadFlag)+", "+
		"\"download\" : "+strconv.FormatBool(book.DownloadFlag)+
	"}"
}
