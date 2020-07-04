package model

import (
	"testing"
	"golang.org/x/text/encoding/traditionalchinese"
)

func TestNew(t *testing.T) {
	big5Decoder := traditionalchinese.Big5.NewDecoder()
	site := NewSite("ck101", big5Decoder, 
				"../book-config/ck101-desktop.json", 
				"../database/ck101.db", "");
	if (site.SiteName != "ck101") {
		t.Errorf("NewSite had incorrect site name %v", site.SiteName)
	}
}

func TestBook(t *testing.T) {
	big5Decoder := traditionalchinese.Big5.NewDecoder()
	site := NewSite("ck101", big5Decoder, 
				"../book-config/ck101-desktop.json", 
				"../database/ck101.db", "");
	book := site.Book(1);
	if (book.Title != "異世流放") {
		t.Errorf("book (ck101-1) has wrong name : %v", book.Title)
	}
}