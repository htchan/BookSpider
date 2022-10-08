package site

import (
	"errors"
	"os"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/htchan/ApiParser"
	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/mock"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	psql "github.com/htchan/BookSpider/internal/repo/psql"
)

func Test_Explore(t *testing.T) {
	t.Parallel()

	ApiParser.SetDefault(
		ApiParser.FromDirectory(os.Getenv("ASSETS_LOCATION") + "/test-data/api_parser"))

	server := mock.MockSiteServer()

	psql.StubPsqlConn()
	db, err := psql.OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "explore/st"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from writers where id>0 and name like $1", "writer%")
		db.Exec("delete from errors where site=$1", site)
		server.Close()
	})

	st, err := NewSite(
		site,
		&config.BookConfig{URLConfig: config.URLConfig{
			Base: server.URL + "/explore/%v",
		}},
		&config.SiteConfig{
			BookKey:         "test_book",
			MaxExploreError: 2,
		},
		&config.CircuitBreakerClientConfig{MaxThreads: 1}, nil, nil)
	if err != nil {
		t.Errorf("fail to init site: %v", err)
		return
	}

	stubData(st.rp, site)

	tests := []struct {
		name          string
		st            *Site
		expectErr     bool
		expectSummary repo.Summary
		expectBooks   []model.Book
	}{
		{
			name:      "works",
			st:        st,
			expectErr: false,
			expectSummary: repo.Summary{
				BookCount: 9, WriterCount: 7, ErrorCount: 3, DownloadCount: 2,
				UniqueBookCount: 8, MaxBookID: 8, LatestSuccessID: 5,
				StatusCount: map[model.StatusCode]int{
					model.Error: 3, model.InProgress: 3, model.End: 3,
				},
			},
			expectBooks: []model.Book{
				{
					Site: site, ID: 4, Title: "title-4", Writer: model.Writer{Name: "writer-4"},
					Type: "type-4", UpdateDate: "4", UpdateChapter: "chapter-4",
					Status: model.InProgress,
				},
				{
					Site: site, ID: 5, Title: "title-5", Writer: model.Writer{Name: "writer-5"},
					Type: "type-5", UpdateDate: "5", UpdateChapter: "chapter-5",
					Status: model.InProgress,
				},
				{Site: site, ID: 6, Error: errors.New("title not found")},
				{Site: site, ID: 7, Error: errors.New("title not found")},
				{Site: site, ID: 8, Error: errors.New("title not found")},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			t.Log(test.st.Info())
			err := Explore(test.st)
			t.Log(test.st.Info())
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want error: %v", err, test.expectErr)
			}

			summary := test.st.Info()
			if !cmp.Equal(summary, test.expectSummary) {
				t.Errorf("summary diff: %v", cmp.Diff(summary, test.expectSummary))
			}

			for _, expectBook := range test.expectBooks {
				bk, err := test.st.BookFromIDHash(expectBook.ID, strconv.FormatInt(int64(expectBook.HashCode), 36))
				if err != nil {
					t.Errorf("fail to fetch book for compare: %v", err)
					return
				}
				expectBook.Writer.ID = bk.Writer.ID
				if !cmp.Equal(*bk, expectBook) {
					t.Errorf("book diff: %v", cmp.Diff(*bk, expectBook))
				}
			}
		})
	}
}
