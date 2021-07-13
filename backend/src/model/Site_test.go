package model

import (
	"testing"
	// "fmt"
	// "golang.org/x/text/encoding"
	"golang.org/x/text/encoding/traditionalchinese"
	"log"
	"io/ioutil"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

// func mapEqual(map1, map2 map[string]interface{}) (bool) {
// 	for key:= range map1 {
// 		if _, exist:= map2[key]; !exist {
// 			return false
// 		}
// 	}
// 	for key:= range map2 {
// 		if val, exist:= map2[key]; !exist || val != map2[key] {
// 			return false
// 		}
// 	}
// 	return true
// }

var testSite = Site{
	SiteName: "ck101",
	database: nil,
	metaBaseUrl: "https://www.ck101.org/book/%v.html",
	metaDownloadUrl: "https://www.ck101.org/0/%v/",
	metaChapterUrl: "https://www.ck101.org%v",
	chapterPattern: "/.*?\\.html",
	decoder: traditionalchinese.Big5.NewDecoder(),
	titleRegex: "<h1><a.*?>(.*?)</a></h1>",
	writerRegex: "作者︰<a.*?>(.*?)</a>",
	typeRegex: " &gt; (.*?) &gt; ",
	lastUpdateRegex: "最新章節\\((\\d{4}-\\d{2}-\\d{2})\\)",
	lastChapterRegex: "<a.*?id=\"newchapter\".*>(.*?)./a>",
	chapterUrlRegex: "<dd><a href=\"(.*?)\">.*?</a></dd>",
	chapterTitleRegex: "<dd><a href.*?>(.*?)</a></dd>",
	chapterContentRegex: "(?s)<div.*?yuedu_zhengwen.*?>(.*?)</div>",
	databaseLocation: "./test_res/ck101.db",
	DownloadLocation: "./test_res",
	MAX_THREAD_COUNT: 100,
}

func TestSite_Book(t *testing.T) {
	var testcase = struct {
		id, version int
		expected string
	} { 5, 0, "ck101\t5\t0\ntest\ttest\n2011-11-11\ttest" }
	book := testSite.Book(testcase.id, testcase.version)
	if book.String() != testcase.expected {
		t.Fatalf("Site.Book(%v, %v).String() gives\n%v\nbut not\n%v",
			testcase.id, testcase.version, book.String(), testcase.expected)
	}
}

func TestSite_OpenDateabase(t *testing.T) {
	testSite.OpenDatabase()
	if testSite.database == nil {
		t.Fatalf("Site.OperDatabase() gives\ndatabase = %v\nbut not\ndatabase = %v",
			testSite.database, nil)
	}
}
func TestSite_CloseDateabase(t *testing.T) {
	testSite.CloseDatabase()
	if testSite.database != nil {
		t.Fatalf("Site.CloseDatabase() gives\ndatabase = %v\nbut not\ndatabase = %v",
			testSite.database, nil)
	}
}

// func TestSite_Update(t *testing.T) {
	
// }
// func TestSite_updateThread(t *testing.T) {
	
// }
// func TestSite_Explore(t *testing.T) {
	
// }
// func TestSite_epdateThread(t *testing.T) {
	
// }
// func TestSite_Download(t *testing.T) {
	
// }
// func TestSite_UpdateError(t *testing.T) {
	
// }
// func TestSite_updateErrorThread(t *testing.T) {
	
// }
// func TestSite_Info(t *testing.T) {
	
// }
// func TestSite_CheckEnd(t *testing.T) {
	
// }
// func TestSite_RandomSuggestBook(t *testing.T) {
	
// }
// func TestSite_fixStorageError(t *testing.T) {
	
// }
// func TestSite_CheckDownloadExistThread(t *testing.T) {
	
// }
// func TestSite_fixDatabaseDuplicateError(t *testing.T) {
	
// }
// func TestSite_fixDatabaseMissingError(t *testing.T) {
	
// }
// func TestSite_Fix(t *testing.T) {
	
// }
// func TestSite_checkDuplicateBook(t *testing.T) {
	
// }
// func TestSite_checkDuplicateError(t *testing.T) {
	
// }
// func TestSite_checkDuplicateCrossTable(t *testing.T) {
	
// }
// func TestSite_getMaxBookId(t *testing.T) {
	
// }
// func TestSite_getMissingBookId(t *testing.T) {
	
// }
// func TestSite_checkMissingId(t *testing.T) {
	
// }
// func TestSite_Check(t *testing.T) {
	
// }
func TestSite_Search(t *testing.T) {
	var testcase = struct {
		title, writer string
		page int
		expected []string
	} { "t", "", 0, []string { "ck101\t5\t0\ntest\ttest\n2011-11-11\ttest" } }
	actual := testSite.Search(testcase.title, testcase.writer, testcase.page)
	if len(actual) != len(testcase.expected) {
		t.Fatalf("Site.Search(\"%v\", \"%v\", %v) gives\nlength = %v\nbut not\nlength = %v",
			testcase.title, testcase.writer, testcase.page, len(actual), len(testcase.expected))
	}
	for i := range actual {
		if actual[i].String() != testcase.expected[i] {
			t.Fatalf("Site.Search(\"%v\", \"%v\", %v)[%v].String() gives\n" +
				"length = %v\nbut not\nlength = %v",
				testcase.title, testcase.writer, testcase.page, i,
				actual[i].String(), testcase.expected[i])
		}
	}
}
// func TestSite_Validate(t *testing.T) {
	
// }
// func TestSite_validateThread(t *testing.T) {
	
// }
// func TestSite_ValidateDownload(t *testing.T) {
	
// }
// func TestSite_validateDownloadThread(t *testing.T) {
	
// }
func TestSite_Map(t *testing.T) {
	var testcase = struct {
		expected map[string]interface{}
	} {
		map[string]interface{} {
			"bookCount":1,
			"bookRecordCount":1,
			"downloadCount":1,
			"downloadRecordCount":1,
			"endCount":1,
			"endRecordCount":1,
			"errorCount":0,
			"errorRecordCount":0,
			"maxThread":100,
			"maxid":5,
			"name":"ck101",
			"readCount":1,
		},
	}
	actual := testSite.Map()
	if !mapEqual(actual, testcase.expected) {
		t.Fatalf("Site.Map() gives\n%v\nbut not\n%v", actual, testcase.expected)
	}
}