package model

import (
	"testing"
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/traditionalchinese"
)

func TestNewSite(t *testing.T) {
	var testCases = []struct {
		name string
		decoder *encoding.Decoder
		config, database, download string
		expectedName string
	} {
		{"ck101", traditionalchinese.Big5.NewDecoder(), 
		"../test-resource/config/ck101-desktop.json", 
		"../test-resource/database/ck101.db", 
		"../test-resource/ck101", "ck101"}}
	success := true
	for _, testCase := range testCases {
		site := NewSite(testCase.name, testCase.decoder, 
			testCase.config, testCase.database, testCase.download)
		if site.SiteName != testCase.expectedName {
			t.Errorf("name of site <%v> is not \"%v\"", site.SiteName, testCase.expectedName)
			success = false
		}
	}
	if success {
		t.Logf("model.NewSite test pass")
	}
}

func TestBook(t *testing.T) {
	site := NewSite("ck101", traditionalchinese.Big5.NewDecoder(), 
	"../test-resource/config/ck101-desktop.json", 
	"../test-resource/database/ck101.db", "../test-resource/ck101")
	var testCases = []struct {
		id int
		expectedVersion int
		expectedTitle, expectedWriter, expectedTypeName string
		expectedLastUpdate, expectedLastChapter string
	} {
		{1, 0, "異世流放", "易人北", "歷史軍事", "2016-09-24", "第659章 番外十"}}
	success := true
	for _, testCase := range testCases {
		book := site.Book(1)
		if book.Version != testCase.expectedVersion || book.Title != testCase.expectedTitle ||
			book.Writer != testCase.expectedWriter || book.Type != testCase.expectedTypeName ||
			book.LastUpdate != testCase.expectedLastUpdate ||
			book.LastChapter != testCase.expectedLastChapter {
			t.Errorf("model.NewSite(%v...).Book(%v) gives (%v, %v, %v, %v, %v, %v), but not (%v, %v, %v, %v, %v, %v)",
				site.SiteName, testCase.id,
				book.Version, book.Title, book.Writer, book.Type, book.LastUpdate, book.LastChapter,
				testCase.expectedVersion, testCase.expectedTitle, testCase.expectedWriter,
				testCase.expectedTypeName, testCase.expectedLastUpdate, testCase.expectedLastChapter)
			success = false
		}
	}
	if success {
		t.Logf("model.Site.Book test pass")
	}
}

func TestBookContent(t *testing.T) {
	site := NewSite("ck101", traditionalchinese.Big5.NewDecoder(), 
	"../test-resource/config/ck101-desktop.json", 
	"../test-resource/database/ck101.db", "../test-resource/ck101")
	var testCases = []struct {
		book Book
		expectedContent string
	} {
		{site.Book(1), "test book result"}}
	success := true
	for _, testCase := range testCases {
		actual := site.BookContent(testCase.book)
		if actual != testCase.expectedContent {
			t.Errorf("model.NewSite(%v...).SiteContent(Book(%v)) gives \"%v\", but not \"%v\"",
				site.SiteName, testCase.book.Id, actual, testCase.expectedContent)
			success = false
		}
	}
	if success {
		t.Logf("model.Site.BookContent test pass")
	}
}

func TestSiteInfo(t *testing.T) {
	site := NewSite("ck101", traditionalchinese.Big5.NewDecoder(), 
	"../test-resource/config/ck101-desktop.json", 
	"../test-resource/database/ck101.db", "../test-resource/ck101")
	fmt.Println("please check the book info ...")
	site.Info()
	fmt.Println()
}