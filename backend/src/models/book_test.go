package models

import (
	"testing"
	"golang.org/x/text/encoding/traditionalchinese"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"context"
	"golang.org/x/sync/semaphore"
	"log"
	"github.com/htchan/BookSpider/helper"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func mapEqual(map1, map2 map[string]interface{}) (bool) {
	for key:= range map1 {
		if _, exist:= map2[key]; !exist {
			return false
		}
	}
	for key:= range map2 {
		if val, exist:= map2[key]; !exist || val != map2[key] {
			return false
		}
	}
	return true
}

var testBook = Book{
	SiteName: "ck101",
	Id: 5,
	Version: -1,
	EndFlag: false,
	DownloadFlag: false,
	ReadFlag: false,
	decoder: traditionalchinese.Big5.NewDecoder(),
	baseUrl: "https://www.ck101.org/book/5.html",
	downloadUrl: "https://www.ck101.org/0/5/",
	chapterUrl: "https://www.ck101.org%v",
	chapterPattern: "/.*?\\.html",
	titleRegex: "<h1><a.*?>(.*?)</a></h1>",
	writerRegex: "作者︰<a.*?>(.*?)</a>",
	typeRegex: " &gt; (.*?) &gt; ",
	lastUpdateRegex: "最新章節\\((\\d{4}-\\d{2}-\\d{2})\\)",
	lastChapterRegex: "<a.*?id=\"newchapter\".*>(.*?)./a>",
	chapterUrlRegex: "<dd><a href=\"(.*?)\">.*?</a></dd>",
	chapterTitleRegex: "<dd><a href.*?>(.*?)</a></dd>",
	chapterContentRegex: "(?s)<div.*?yuedu_zhengwen.*?>(.*?)</div>",
}

func TestBookLog(t *testing.T) {
	testBook.Log(map[string]interface{} {"test": "test"})
}
func TestBookvalidHTML(t *testing.T) {
	var testcases = []struct {
		input string
		expected bool
	} {
		{"", false},
		{"http://www.google.com/", true},
		{"404", false},
	}
	for _, testcase := range testcases {
		actual := testBook.validHTML(testcase.input, "", 0)
		if actual != testcase.expected {
			t.Fatalf("Book.validateHTML(\"%v\", \"\",, 0) result gives \"%v\", but not \"%v\"\n",
				testcase.input, actual, testcase.expected)
		}
	}
}
func TestBookUpdate(t *testing.T) {
	var testcase = struct {
		expected bool
		updatedString string
	} { true, "ck101\t5\t0\n神經病戀愛指南\t小貓一尾\n2016-10-09\t第89章" }
	actual := testBook.Update()
	if actual != testcase.expected || testBook.String() != testcase.updatedString {
		t.Fatalf("Book.validateHTML() result gives\n(%v, %v),\nbut not\n(%v, %v)\n", 
			actual, testBook.String(), testcase.expected, testcase.updatedString)
	}
}
func TestBooksaveBook(t *testing.T) {
	var testcase = struct {
		expected int
		expectedData string
	} { 0, "神經病戀愛指南\n小貓一尾\n--------------------\n\ntitle\n--------------------\ncontent\n\n" }
	testBook.Update()
	urls :=  []string{"test"}
	chapters := []Chapter{ Chapter{"test", "title", "content"}, }
	actual := testBook.saveBook("./test_res/save-book-5.txt",urls, chapters)
	actualData, err := ioutil.ReadFile("./test_res/save-book-5-test.txt")
	helper.CheckError(err)
	os.Remove("./test_res/save-book-5.txt")
	if actual != testcase.expected || string(actualData) != testcase.expectedData {
		t.Fatalf("Book.validateHTML() result gives\n(%v, %v),\nbut not\n(%v, %v)\n", 
			actual, string(actualData), testcase.expected, testcase.expectedData)
	}
}
func TestBookDownload(t *testing.T) {
	var testcase = struct {
		expected bool
		expectedFileName string
	} { true, "./test_res/ck101-test.txt" }
	actual := testBook.Download("./test_res", 1000)
	actualData, err := ioutil.ReadFile("./test_res/5.txt")
	helper.CheckError(err)
	expectedData, err := ioutil.ReadFile(testcase.expectedFileName)
	helper.CheckError(err)
	result := len(actualData) - len(expectedData)
	if result < 0 { result = -result }
	os.Remove("./test_res/5.txt")
	if actual != testcase.expected || result > 10 {
		t.Fatalf("Book.validateHTML() result gives\n(%v),\nbut not\n(%v)\nlength diff: %v", 
			actual, testcase.expected, result)
	}
}
func TestBookoptimizeContent(t *testing.T) {
	var testcase = struct {
		input, expected string
	} { "<br />&nbsp;<b></b><p></p>                <p/>", "\n" }
	actual := testBook.optimizeContent(testcase.input)
	if actual != testcase.expected {
		t.Fatalf("Book.optimizeContent(\"%v\") result gives\n\"%v\",\nbut not\n\"%v\"",
			testcase.input, actual, testcase.expected)
	}
}
func TestBookdownloadAllChapters(t *testing.T) {
	var testcase = struct {
		input1, input2 []string
		input3 int
		expectedChapters []Chapter
	} {
		[]string { "/0/5/772.html", }, []string { "第1章", }, 100,
		[]Chapter {
			Chapter{ "https://www.ck101.org/0/5/772.html", "第1章", strings.Repeat(" ", 15060), },
		},
	}
	actualChapters := testBook.downloadAllChapters(testcase.input1, testcase.input2, testcase.input3)
	if len(actualChapters) != len(testcase.expectedChapters) {
		t.Fatalf("Book.downloadAllChapter(%v, %v, %v) result gives\nlength = %v,\nbut not\nlength = %v", 
				testcase.input1, testcase.input2, testcase.input3,
				len(actualChapters), len(testcase.expectedChapters))
	}
	for i := range actualChapters {
		result := 0
		if len(actualChapters[i].Content) > len(testcase.expectedChapters[i].Content) {
			result = len(actualChapters[i].Content) - len(testcase.expectedChapters[i].Content)
		} else {
			result = len(testcase.expectedChapters[i].Content) - len(actualChapters[i].Content)
		}
		if result > 10 {
			t.Fatalf("Book.downloadAllChapter(%v, %v, %v) result[%v] gives\nlength %v,\nbut not\nlength %v", 
				testcase.input1, testcase.input2, testcase.input3,
				i, len(actualChapters[i].Content), len(testcase.expectedChapters[i].Content))
		}
	}
}
func TestBookdownloadChapter(t *testing.T) {
	ctx := context.Background()
	var wg sync.WaitGroup
	var testcase = struct {
		input1, input2 string
		input3 *semaphore.Weighted
		input4 *sync.WaitGroup
		input5 chan Chapter
		expected Chapter
	} {
		"https://www.ck101.org/0/5/772.html", "第1章",
		semaphore.NewWeighted(int64(1)), &wg,
		make(chan Chapter, 1),
		Chapter{ "https://www.ck101.org/0/5/772.html", "第1章", strings.Repeat(" ", 15059), },
	}
	testcase.input4.Add(1)
	testcase.input3.Acquire(ctx, 1)
	testBook.downloadChapter(testcase.input1, testcase.input2, testcase.input3,
		testcase.input4, testcase.input5)
	actual := <-testcase.input5
	result := 0
	if len(actual.Content) > len(testcase.expected.Content) {
		result = len(actual.Content) - len(testcase.expected.Content)
	} else {
		result = len(testcase.expected.Content) - len(actual.Content)
	}
	if result > 10 {
		t.Fatalf("Book.downloadAllChapter(%v, %v, s, wg, ch) channel gives result\n" +
			"length %v,\nbut not\nlength %v", 
			testcase.input1, testcase.input2, len(actual.Content), len(testcase.expected.Content))
	}
}
func TestBookContent(t *testing.T) {
	var testcase = struct {
		input string
		expected string
	} { "./test_res/", "" }
	actual := testBook.Content(testcase.input)
	if actual != testcase.expected {
		t.Fatalf("Book.Content(\"%v\") channel gives result\n %v,\nbut not\n %v", 
			"./test_res/", actual, testcase.expected)
	}
	tempBookName, tempDownloadFlag := testBook.Title, testBook.DownloadFlag
	testBook.Title, testBook.DownloadFlag = "", false
	testcase.expected = "hello"
	ioutil.WriteFile("./test_res/5.txt", []byte(testcase.expected), 0644)
	testBook.Title, testBook.DownloadFlag = tempBookName, true
	actual = testBook.Content(testcase.input)
	testBook.DownloadFlag = tempDownloadFlag
	if actual != testcase.expected {
		t.Fatalf("Book.Content(\"%v\") channel gives result\n%v,\nbut not\n%v", 
			testcase.input, actual, testcase.expected)
	}
	os.Remove("./test_res/5.txt")
}
func TestBookStorageLocation(t *testing.T) {
	var testcase = struct {
		input, expected string
	} { "./test_res", "./test_res/5.txt"}
	actual := testBook.storageLocation(testcase.input)
	if actual != testcase.expected {
		t.Fatalf("Book.StorageContent(\"%v\") channel gives result\n %v,\nbut not\n %v", 
			testcase.input, actual, testcase.expected)
	}
}
func TestBookString(t *testing.T) {
	var testcase = struct {
		expected string
	} { "ck101\t5\t0\n神經病戀愛指南\t小貓一尾\n2016-10-09\t第89章"}
	actual := testBook.String()
	if actual != testcase.expected {
		t.Fatalf("Book.String() result gives\n\"%v\",\nbut not\n\"%v\"",
			actual, testcase.expected)
	}
}
func TestBookMap(t *testing.T) {
	var testcase = struct {
		expected map[string]interface{}
	} {
		map[string]interface{} {
			"site" : "ck101", "id" : 5, "version" : 0,
			"title" : "神經病戀愛指南", "writer" : "小貓一尾",
			"chapter":"第89章", "update" : "2016-10-09", "type" : "言情小說",
			"download":false, "end":false, "read":false,    
		},
	}
	actual := testBook.Map()
	if !mapEqual(actual, testcase.expected) {
		t.Fatalf("Book.Map() result gives\n\"%v\",\nbut not\n\"%v\"",
			actual, testcase.expected)
	}
}