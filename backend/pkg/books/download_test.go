package books

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
	"github.com/htchan/BookSpider/internal/utils"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

var testBook_Save = Book{
	SiteName: "ck101",
	Id: 5,
	Version: -1,
	EndFlag: false,
	DownloadFlag: false,
	ReadFlag: false,
	decoder: traditionalchinese.Big5.NewDecoder(),
	metaInfo: MetaInfo{
		baseUrl: "https://www.ck101.org/book/5.html",
		downloadUrl: "https://www.ck101.org/0/5/",
		chapterUrl: "https://www.ck101.org%v",
		chapterUrlPattern: "/.*?\\.html",
		titleRegex: "<h1><a.*?>(.*?)</a></h1>",
		writerRegex: "作者︰<a.*?>(.*?)</a>",
		typeRegex: " &gt; (.*?) &gt; ",
		lastUpdateRegex: "最新章節\\((\\d{4}-\\d{2}-\\d{2})\\)",
		lastChapterRegex: "<a.*?id=\"newchapter\".*>(.*?)./a>",
		chapterUrlRegex: "<dd><a href=\"(.*?)\">.*?</a></dd>",
		chapterTitleRegex: "<dd><a href.*?>(.*?)</a></dd>",
		chapterContentRegex: "(?s)<div.*?yuedu_zhengwen.*?>(.*?)</div>",
	},
}

func Test_saveBook(t *testing.T) {
	var testcase = struct {
		expected int
		expectedData string
	} { 0, "神經病戀愛指南\n小貓一尾\n--------------------\n\ntitle\n--------------------\ncontent\n\n" }
	testBook_Save.Update()
	urls :=  []string{"test"}
	chapters := []Chapter{ Chapter{"test", "title", "content"}, }
	actual := testBook_Save.saveBook("../../test/model-test-data/save-book-5.txt", urls, chapters)
	actualData, err := ioutil.ReadFile("../../test/model-test-data/save-book-5-test.txt")
	utils.CheckError(err)
	os.Remove("../../test/model-test-data/save-book-5.txt")
	if actual != testcase.expected || string(actualData) != testcase.expectedData {
		t.Fatalf("Book.saveBook() result gives\n(%v, %v),\nbut not\n(%v, %v)\n", 
			actual, string(actualData), testcase.expected, testcase.expectedData)
	}
}

func TestDownload(t *testing.T) {
	var testcase = struct {
		expected bool
		expectedFileName string
	} { true, "../../test/model-test-data/ck101-test.txt" }
	actual := testBook_Save.Download("../../test/model-test-data", 1000)
	actualData, err := ioutil.ReadFile("../../test/model-test-data/5.txt")
	utils.CheckError(err)
	expectedData, err := ioutil.ReadFile(testcase.expectedFileName)
	utils.CheckError(err)
	result := len(actualData) - len(expectedData)
	if result < 0 { result = -result }
	os.Remove("../../test/model-test-data/5.txt")
	if actual != testcase.expected || result > 10 {
		t.Fatalf("Book.validateHTML() result gives\n(%v),\nbut not\n(%v)\nlength diff: %v", 
			actual, testcase.expected, result)
	}
}

func Test_optimizeContent(t *testing.T) {
	var testcase = struct {
		input, expected string
	} { "<br />&nbsp;<b></b><p></p>                <p/>", "\n" }
	actual := testBook_Save.optimizeContent(testcase.input)
	if actual != testcase.expected {
		t.Fatalf("Book.optimizeContent(\"%v\") result gives\n\"%v\",\nbut not\n\"%v\"",
			testcase.input, actual, testcase.expected)
	}
}

func Test_downloadChapters(t *testing.T) {
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
	actualChapters := testBook_Save.downloadChapters(testcase.input1, testcase.input2, testcase.input3)
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

func Test_downloadChapter(t *testing.T) {
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
	testBook_Save.downloadChapter(testcase.input1, testcase.input2, testcase.input3,
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