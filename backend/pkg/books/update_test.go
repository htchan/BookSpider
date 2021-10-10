package books

import (
	"testing"
	"golang.org/x/text/encoding/traditionalchinese"
	"io/ioutil"
	"log"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

var testBook_Update = Book{
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

func Test_extractInfo(t *testing.T) {
	//TODO: add fail testcase here
	var expect = struct {
		title, writer, typeName, lastUpdate, lastChapter string
		errorExist bool
	} { "神經病戀愛指南", "小貓一尾", "言情小說", "2016-10-09", "第89章", false }
	actualTitle, actualWriter, actualTypeName,
		actualLastUpdate, actualLastChapter, actualErr := testBook_Update.extractInfo()
	if actualTitle != expect.title || actualWriter != expect.writer ||
		actualTypeName != expect.typeName ||
		actualLastUpdate != expect.lastUpdate || actualLastChapter != expect.lastChapter ||
		(actualErr != nil) != expect.errorExist {
			t.Fatalf("book.extractInfo() returns\n\"%v\", \"%v\", \"%v\", \"%v\", \"%v\", %v\n" +
					"but not\n\"%v\", \"%v\", \"%v\", \"%v\", \"%v\", %v",
					actualTitle, actualWriter, actualTypeName,
					actualLastUpdate, actualLastChapter, actualErr,
					expect.title, expect.writer, expect.typeName,
					expect.lastUpdate, expect.lastChapter, expect.errorExist)
		}

}

func TestUpdate(t *testing.T) {
	//TODO: add fail testcase here
	var testcase = struct {
		expected bool
		updatedString string
	} { true, "ck101\t5\t0\n神經病戀愛指南\t小貓一尾\n2016-10-09\t第89章" }
	actual := testBook_Update.Update()
	if actual != testcase.expected || testBook_Update.String() != testcase.updatedString {
		t.Fatalf("Book.validateHTML() result gives\n(%v, %v),\nbut not\n(%v, %v)\n", 
			actual, testBook_Update.String(), testcase.expected, testcase.updatedString)
	}
}