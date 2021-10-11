package books

import (
	"testing"
	"golang.org/x/text/encoding/traditionalchinese"
	"io/ioutil"
	"os"
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
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

func TestNewMetaInfo(t *testing.T) {
	//TODO: add fail case
	metaInfo, err := NewMetaInfo(map[string]string{
		"baseUrl" : "https://www.ck101.org/book/%v.html",
		"downloadUrl" : "https://www.ck101.org/0/%v/",
		"chapterUrl" : "https://www.ck101.org%v",
		"chapterUrlPattern" : "/.*?\\.html",
		"titleRegex" : "<h1><a.*?>(.*?)</a></h1>",
		"writerRegex" : "作者︰<a.*?>(.*?)</a>",
		"typeRegex" : " &gt; (.*?) &gt; ",
		"lastUpdateRegex" : "最新章節\\((\\d{4}-\\d{2}-\\d{2})\\)",
		"lastChapterRegex" : "<a.*?id=\"newchapter\".*>(.*?)./a>",
		"chapterUrlRegex" : "<dd><a href=\"(.*?)\">.*?</a></dd>",
		"chapterTitleRegex" : "<dd><a href.*?>(.*?)</a></dd>",
		"chapterContentRegex" : "(?s)<div.*?yuedu_zhengwen.*?>(.*?)</div>",
	})
	if metaInfo.baseUrl == "" || metaInfo.downloadUrl == "" || metaInfo.chapterUrl == "" || 
	metaInfo.chapterUrlPattern == "" || metaInfo.titleRegex == "" || metaInfo.writerRegex == "" || 
	metaInfo.typeRegex == "" || metaInfo.lastUpdateRegex == "" || metaInfo.lastChapterRegex == "" || 
	metaInfo.chapterUrlRegex == "" || metaInfo.chapterTitleRegex == "" ||
	metaInfo.chapterContentRegex == "" || err != nil {
		t.Fatalf("metaInfo contains empty field after construction %v", metaInfo)
	}
}

func TestNewBook(t *testing.T) {
	databaseLocation := "../../test/model-test-data/ck101.db"
	metaInfo, err := NewMetaInfo(map[string]string{
		"baseUrl" : "https://www.ck101.org/book/%v.html",
		"downloadUrl" : "https://www.ck101.org/0/%v/",
		"chapterUrl" : "https://www.ck101.org%v",
		"chapterUrlPattern" : "/.*?\\.html",
		"titleRegex" : "<h1><a.*?>(.*?)</a></h1>",
		"writerRegex" : "作者︰<a.*?>(.*?)</a>",
		"typeRegex" : " &gt; (.*?) &gt; ",
		"lastUpdateRegex" : "最新章節\\((\\d{4}-\\d{2}-\\d{2})\\)",
		"lastChapterRegex" : "<a.*?id=\"newchapter\".*>(.*?)./a>",
		"chapterUrlRegex" : "<dd><a href=\"(.*?)\">.*?</a></dd>",
		"chapterTitleRegex" : "<dd><a href.*?>(.*?)</a></dd>",
		"chapterContentRegex" : "(?s)<div.*?yuedu_zhengwen.*?>(.*?)</div>",
	})

	database, err := sql.Open("sqlite3", databaseLocation)
	if err != nil {
		t.Fatalf("cannot open database at %v", databaseLocation)
	}

	tx, err := database.Begin()
	if err != nil {
		t.Fatalf("cannot open transaction for database at %v", databaseLocation)
	}

	book, err := NewBook("ck101", 5, -1, *metaInfo, nil, tx)

	tx.Commit()
	err = database.Close()
	if err != nil {
		t.Fatalf("cannot close database at %v", databaseLocation)
	}
	expectData := map[string]interface{}{
		"site": "ck101", "id": 5, "version": 0,
		"title": "test", "writer": "test", "type": "test",
		"update": "2011-11-11", "chapter": "test",
		"end": true, "download": true, "read": true,
	}
	if !mapEqual(book.Map(), expectData) || (err != nil) {
		t.Fatalf("books.NewBook(\"%v\", %v, %v, %v, %v, %v) returns\n" +
			"%v\nbut not\n%v", "ck101", 5, -1, metaInfo, nil, tx, book.Map(), expectData)
	}
}

func TestLog(t *testing.T) {
	testBook.Log(map[string]interface{} {"test": "test"})
}

func Test_validHTML(t *testing.T) {
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

func TestContent(t *testing.T) {
	var testcase = struct {
		input string
		expected string
	} { "../../test/model-test-data/", "" }
	testBook.Update()
	actual := testBook.Content(testcase.input)
	if actual != testcase.expected {
		t.Fatalf("Book.Content(\"%v\") channel gives result\n %v,\nbut not\n %v", 
			"../../test/model-test-data/", actual, testcase.expected)
	}
	tempBookName, tempDownloadFlag := testBook.Title, testBook.DownloadFlag
	testBook.Title, testBook.DownloadFlag = "", false
	testcase.expected = "hello"
	ioutil.WriteFile("../../test/model-test-data/5.txt", []byte(testcase.expected), 0644)
	testBook.Title, testBook.DownloadFlag = tempBookName, true
	actual = testBook.Content(testcase.input)
	testBook.DownloadFlag = tempDownloadFlag
	if actual != testcase.expected {
		t.Fatalf("Book.Content(\"%v\") channel gives result\n%v,\nbut not\n%v", 
			testcase.input, actual, testcase.expected)
	}
	os.Remove("../../test/model-test-data/5.txt")
}

func TestStorageLocation(t *testing.T) {
	var testcase = struct {
		input, expected string
	} { "../../test/model-test-data", "../../test/model-test-data/5.txt"}
	actual := testBook.storageLocation(testcase.input)
	if actual != testcase.expected {
		t.Fatalf("Book.StorageContent(\"%v\") channel gives result\n %v,\nbut not\n %v", 
			testcase.input, actual, testcase.expected)
	}
}

func TestString(t *testing.T) {
	var testcase = struct {
		expected string
	} { "ck101\t5\t0\n神經病戀愛指南\t小貓一尾\n2016-10-09\t第89章"}
	testBook.Update()
	actual := testBook.String()
	if actual != testcase.expected {
		t.Fatalf("Book.String() result gives\n\"%v\",\nbut not\n\"%v\"",
			actual, testcase.expected)
	}
}

func TestMap(t *testing.T) {
	var testcase = struct {
		expected map[string]interface{}
	} {
		map[string]interface{} {
			"site" : "ck101", "id" : 5, "version" : 0,
			"title" : "神經病戀愛指南", "writer" : "小貓一尾",
			"chapter":"第89章", "update" : "2016-10-09", "type" : "言情小說",
			"download": false, "end": false, "read": false,    
		},
	}
	testBook.Update()
	actual := testBook.Map()
	if !mapEqual(actual, testcase.expected) {
		t.Fatalf("Book.Map() result gives\n\"%v\",\nbut not\n\"%v\"",
			actual, testcase.expected)
	}
}