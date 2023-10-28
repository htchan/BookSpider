package service

import (
	"errors"
	"os"
	"testing"

	config "github.com/htchan/BookSpider/internal/config_new"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestServiceImp_BookFileLocation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		serv   ServiceImp
		bk     *model.Book
		wantBk *model.Book
		want   string
	}{
		{
			name:   "filename with 0 hashcode",
			serv:   ServiceImp{conf: config.SiteConfig{Storage: "./test-book-file-location"}},
			bk:     &model.Book{Site: "test-book-file-location", ID: 1, HashCode: 0},
			wantBk: &model.Book{Site: "test-book-file-location", ID: 1, HashCode: 0},
			want:   "test-book-file-location/1.txt",
		},
		{
			name:   "filename with 100 hashcode",
			serv:   ServiceImp{conf: config.SiteConfig{Storage: "./test-book-file-location"}},
			bk:     &model.Book{Site: "test-book-file-location", ID: 1, HashCode: 100},
			wantBk: &model.Book{Site: "test-book-file-location", ID: 1, HashCode: 100},
			want:   "test-book-file-location/1-v2s.txt",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := test.serv.BookFileLocation(test.bk)
			assert.Equal(t, test.wantBk, test.bk)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestServiceImp_BookInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		serv   ServiceImp
		bk     *model.Book
		wantBk *model.Book
		want   string
	}{
		{
			name: "book without error",
			serv: ServiceImp{},
			bk: &model.Book{
				Site: "test-info", ID: 1, HashCode: 1,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
				Status: model.StatusInProgress, IsDownloaded: true,
			},
			wantBk: &model.Book{
				Site: "test-info", ID: 1, HashCode: 1,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
				Status: model.StatusInProgress, IsDownloaded: true,
			},
			want: `{"site":"test-info","id":1,"hash_code":"1","title":"title","writer":"writer","type":"type","update_date":"date","update_chapter":"chapter","status":"INPROGRESS","is_downloaded":true,"error":""}`,
		},
		{
			name: "book with error",
			serv: ServiceImp{},
			bk: &model.Book{
				Site: "test-info", ID: 1, HashCode: 1,
				Status: model.StatusError, Error: errors.New("test error"),
			},
			wantBk: &model.Book{
				Site: "test-info", ID: 1, HashCode: 1,
				Status: model.StatusError, Error: errors.New("test error"),
			},
			want: `{"site":"test-info","id":1,"hash_code":"1","title":"","writer":"","type":"","update_date":"","update_chapter":"","status":"ERROR","is_downloaded":false,"error":"test error"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.serv.BookInfo(test.bk)
			assert.Equal(t, test.wantBk, test.bk)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestServiceImp_BookContent(t *testing.T) {
	t.Parallel()
	os.Mkdir("./test-book-content", os.ModePerm)
	os.WriteFile("./test-book-content/1.txt", []byte("some data"), os.ModePerm)
	os.WriteFile("./test-book-content/2.txt", []byte("secret data"), 0333)

	t.Cleanup(func() {
		os.RemoveAll("./test-book-content")
	})

	tests := []struct {
		name         string
		serv         ServiceImp
		bk           *model.Book
		wantBk       *model.Book
		wantContent  string
		wantError    bool
		wantErrorStr string
	}{
		{
			name:         "happy flow",
			serv:         ServiceImp{conf: config.SiteConfig{Storage: "./test-book-content"}},
			bk:           &model.Book{ID: 1, IsDownloaded: true},
			wantBk:       &model.Book{ID: 1, IsDownloaded: true},
			wantContent:  "some data",
			wantError:    false,
			wantErrorStr: "",
		},
		{
			name:         "happy flow",
			serv:         ServiceImp{conf: config.SiteConfig{Storage: "./test-book-content"}},
			bk:           &model.Book{ID: 1, IsDownloaded: false},
			wantBk:       &model.Book{ID: 1, IsDownloaded: false},
			wantContent:  "",
			wantError:    true,
			wantErrorStr: "book is not download",
		},
		{
			name:         "book not found",
			serv:         ServiceImp{conf: config.SiteConfig{Storage: "./test-book-content"}},
			bk:           &model.Book{ID: 999, IsDownloaded: true},
			wantBk:       &model.Book{ID: 999, IsDownloaded: true},
			wantContent:  "",
			wantError:    true,
			wantErrorStr: "file not found",
		},
		{
			name:         "read book failed",
			serv:         ServiceImp{conf: config.SiteConfig{Storage: "./test-book-content"}},
			bk:           &model.Book{ID: 2, IsDownloaded: true},
			wantBk:       &model.Book{ID: 2, IsDownloaded: true},
			wantContent:  "",
			wantError:    true,
			wantErrorStr: "get book content error: open test-book-content/2.txt: permission denied",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := test.serv.BookContent(test.bk)
			assert.Equal(t, test.wantBk, test.bk)
			if test.wantError {
				assert.Error(t, err)
			}

			if err != nil {
				assert.EqualError(t, err, test.wantErrorStr)
			}
			assert.Equal(t, test.wantContent, got)
		})
	}
}

func TestServiceImp_Chapters(t *testing.T) {
	t.Parallel()

	os.Mkdir("./test-chapters", os.ModePerm)
	os.WriteFile("./test-chapters/1.txt", []byte("title\nwriter\n--------------------\n\ntitle 1\n--------------------\ncontent 1\n--------------------\n\ntitle 2\n--------------------\ncontent 2\n--------------------\n"), os.ModePerm)

	t.Cleanup(func() {
		os.RemoveAll("./test-chapters")
	})

	tests := []struct {
		name         string
		serv         ServiceImp
		bk           *model.Book
		wantBk       *model.Book
		wantChapters model.Chapters
		wantError    bool
		wantErrorStr string
	}{
		{
			name:   "happy flow",
			serv:   ServiceImp{conf: config.SiteConfig{Storage: "./test-chapters"}},
			bk:     &model.Book{ID: 1, IsDownloaded: true},
			wantBk: &model.Book{ID: 1, IsDownloaded: true},
			wantChapters: model.Chapters{
				{Index: 0, URL: "", Title: "title 1", Content: "content 1"},
				{Index: 1, URL: "", Title: "title 2", Content: "content 2"},
			},
			wantError:    false,
			wantErrorStr: "",
		},
		{
			name:         "get book content fail",
			serv:         ServiceImp{conf: config.SiteConfig{Storage: "./test-chapters"}},
			bk:           &model.Book{ID: 1, IsDownloaded: false},
			wantBk:       &model.Book{ID: 1, IsDownloaded: false},
			wantChapters: nil,
			wantError:    true,
			wantErrorStr: "load content failed: book is not download",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := test.serv.Chapters(test.bk)

			assert.Equal(t, test.wantBk, test.bk)
			if test.wantError {
				assert.Error(t, err)
			}

			if err != nil {
				assert.EqualError(t, err, test.wantErrorStr)
			}
			assert.Equal(t, test.wantChapters, got)
		})
	}
}
