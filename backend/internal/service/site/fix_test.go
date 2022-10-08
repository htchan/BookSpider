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
	psql "github.com/htchan/BookSpider/internal/repo/psql"
)

func Test_Fix(t *testing.T) {
	t.Parallel()

	ApiParser.SetDefault(
		ApiParser.FromDirectory(os.Getenv("ASSETS_LOCATION") + "/test-data/api_parser"))

	server := mock.MockSiteServer()

	psql.StubPsqlConn()
	db, err := psql.OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "fix/st"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
		server.Close()
	})

	st, err := NewSite(
		site,
		&config.BookConfig{URLConfig: config.URLConfig{Base: server.URL + "/error/%v"}},
		&config.SiteConfig{BookKey: "test_book"},
		&config.CircuitBreakerClientConfig{MaxThreads: 2}, nil, nil)
	if err != nil {
		t.Errorf("fail to init site: %v", err)
		return
	}
	st.rp.CreateBook(&model.Book{Site: site, ID: 2})
	st.rp.CreateBook(&model.Book{Site: site, ID: 4})

	tests := []struct {
		name        string
		st          *Site
		expectErr   bool
		expectBooks []model.Book
	}{
		{
			name:      "works",
			st:        st,
			expectErr: false,
			expectBooks: []model.Book{
				{Site: site, ID: 1, Error: errors.New("title not found")},
				{Site: site, ID: 2},
				{Site: site, ID: 3, Error: errors.New("title not found")},
				{Site: site, ID: 4},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := Fix(test.st)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want error: %v", err, test.expectErr)
			}

			for _, expectBook := range test.expectBooks {
				bk, err := test.st.BookFromIDHash(expectBook.ID, strconv.FormatInt(int64(expectBook.HashCode), 36))
				if err != nil {
					t.Errorf("fail to fetch book for compare: id: %v; err: %v", expectBook.ID, err)
					return
				}
				if !cmp.Equal(*bk, expectBook) {
					t.Errorf("book diff: %v", cmp.Diff(*bk, expectBook))
				}
			}
		})
	}
}
