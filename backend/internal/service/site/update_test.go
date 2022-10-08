package site

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/htchan/ApiParser"
	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/mock"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	psql "github.com/htchan/BookSpider/internal/repo/psql"
)

func Test_Update(t *testing.T) {
	t.Parallel()

	ApiParser.SetDefault(
		ApiParser.FromDirectory(os.Getenv("ASSETS_LOCATION") + "/test-data/api_parser"))

	server := mock.MockSiteServer()

	psql.StubPsqlConn()
	db, err := psql.OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "update/st"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
		server.Close()
	})

	st, err := NewSite(
		site,
		&config.BookConfig{URLConfig: config.URLConfig{
			Base: server.URL + "/no-update-book/%v",
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
				BookCount: 8, WriterCount: 5, ErrorCount: 0, DownloadCount: 3,
				UniqueBookCount: 4, MaxBookID: 4, LatestSuccessID: 4,
				StatusCount: map[model.StatusCode]int{
					model.InProgress: 5, model.End: 3,
				},
			},
			expectBooks: []model.Book{
				{
					Site: "update/st", ID: 1, HashCode: model.GenerateHash(),
					Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
					UpdateDate: "1234", UpdateChapter: "chapter", Status: model.InProgress, IsDownloaded: true,
				},
				{
					Site: "update/st", ID: 2, HashCode: model.GenerateHash(),
					Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
					UpdateDate: "1234", UpdateChapter: "chapter", Status: model.InProgress,
				},
				{
					Site: "update/st", ID: 3, HashCode: model.GenerateHash(),
					Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
					UpdateDate: "1234", UpdateChapter: "chapter", Status: model.InProgress,
				},
				{
					Site: "update/st", ID: 4, HashCode: 0,
					Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
					UpdateDate: "1234", UpdateChapter: "chapter", Status: model.InProgress,
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			t.Log(test.st.Info())
			err := Update(test.st)
			t.Log(test.st.Info())
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want error: %v", err, test.expectErr)
			}

			summary := test.st.Info()
			if !cmp.Equal(summary, test.expectSummary) {
				t.Errorf("summary diff: %v", cmp.Diff(summary, test.expectSummary))
			}

			for _, expectBook := range test.expectBooks {
				bk, err := test.st.BookFromID(expectBook.ID)
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
