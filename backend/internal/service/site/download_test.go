package site

import (
	"os"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/htchan/ApiParser"
	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/mock"
	"github.com/htchan/BookSpider/internal/model"
	psql "github.com/htchan/BookSpider/internal/repo/psql"
)

func Test_Content(t *testing.T) {
	t.Parallel()

	os.Mkdir("./storage", os.ModePerm)
	os.WriteFile("./storage/1.txt", []byte("hello"), os.ModePerm)

	t.Cleanup(func() {
		os.RemoveAll("./storage")
	})

	site := "content/st"
	st, err := NewSite(site, &config.BookConfig{}, &config.SiteConfig{Storage: "./storage"}, &config.CircuitBreakerClientConfig{}, nil, nil)
	if err != nil {
		t.Errorf("fail to init site: %v", err)
		return
	}

	tests := []struct {
		name      string
		st        *Site
		bk        *model.Book
		expect    string
		expectErr bool
	}{
		{
			name:      "return existing file of download book",
			st:        st,
			bk:        &model.Book{ID: 1, IsDownloaded: true},
			expect:    "hello",
			expectErr: false,
		},
		{
			name:      "return error for not existing book",
			st:        st,
			bk:        &model.Book{ID: 2, IsDownloaded: false},
			expect:    "",
			expectErr: true,
		},
		{
			name:      "return error for not downloaded book",
			st:        st,
			bk:        &model.Book{ID: 1},
			expect:    "",
			expectErr: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result, err := Content(test.st, test.bk)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v, want err: %v", err, test.expectErr)
			}

			if result != test.expect {
				t.Errorf("content diff: %v", cmp.Diff(result, test.expect))
			}
		})
	}
}

func Test_Download(t *testing.T) {
	t.Parallel()

	ApiParser.SetDefault(
		ApiParser.FromDirectory(os.Getenv("ASSETS_LOCATION") + "/test-data/api_parser"))

	server := mock.MockSiteServer()

	psql.StubPsqlConn()
	db, err := psql.OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "download/st"

	os.Mkdir("./download", os.ModePerm)

	t.Cleanup(func() {
		os.RemoveAll("./download")
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
		server.Close()
	})

	st, err := NewSite(
		site,
		&config.BookConfig{URLConfig: config.URLConfig{
			Download:      server.URL + "/chapter-header/valid/%v",
			ChapterPrefix: server.URL + "/chapter/valid",
		}},
		&config.SiteConfig{
			BookKey:           "test_book",
			Storage:           "./download",
			ConcurrencyConfig: config.ConcurrencyConfig{DownloadThreads: 1},
		},
		&config.CircuitBreakerClientConfig{MaxThreads: 10}, nil, nil)
	if err != nil {
		t.Errorf("fail to init site: %v", err)
		return
	}

	bksDB := stubData(st.rp, site)

	tests := []struct {
		name              string
		st                *Site
		expectBooks       []model.Book
		expectErr         bool
		expectLocation    string
		expectFileContent string
	}{
		{
			name: "works",
			st:   st,
			expectBooks: []model.Book{
				{
					Site: site, ID: 2, HashCode: 100, Title: "title 2 new", Writer: bksDB[2].Writer,
					Type: "type 2 new", UpdateDate: "date 2.1", UpdateChapter: "chapter 2 new",
					Status: model.End, IsDownloaded: true,
				},
			},
			expectErr:      false,
			expectLocation: "./download/2-v2s.txt",
			expectFileContent: "title 2 new\ndownload/st writer 2 new\n" + model.CONTENT_SEP +
				"\n\n1\n" + model.CONTENT_SEP + "\nsuccess-content-regex\n" + model.CONTENT_SEP +
				"\n2\n" + model.CONTENT_SEP + "\nsuccess-content-regex\n" + model.CONTENT_SEP +
				"\n3\n" + model.CONTENT_SEP + "\nsuccess-content-regex\n" + model.CONTENT_SEP +
				"\n4\n" + model.CONTENT_SEP + "\nsuccess-content-regex\n" + model.CONTENT_SEP + "\n",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := Download(test.st)
			if err != nil {
				t.Errorf("got error: %v, want err: %v", err, test.expectErr)
			}

			for _, expectBook := range test.expectBooks {
				bk, err := test.st.BookFromIDHash(expectBook.ID, strconv.FormatInt(int64(expectBook.HashCode), 36))
				if err != nil {
					t.Errorf("fail to fetch book for compare: %v", err)
					return
				}

				if !cmp.Equal(*bk, expectBook) {
					t.Errorf("book diff: %v", cmp.Diff(*bk, expectBook))
				}
			}

			content, err := os.ReadFile(test.expectLocation)
			if err != nil {
				t.Errorf("fail to read file: %v", test.expectLocation)
			}
			if string(content) != test.expectFileContent {
				t.Errorf("content diff: %v", cmp.Diff(string(content), test.expectFileContent))
			}
		})
	}
}
