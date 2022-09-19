package site

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	psql "github.com/htchan/BookSpider/internal/repo/psql"
)

func Test_Check(t *testing.T) {
	t.Parallel()

	psql.StubPsqlConn()
	db, err := psql.OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "check/st"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
		db.Close()
	})

	st, err := NewSite(site, config.BookConfig{}, config.SiteConfig{BackupDirectory: "/backup"}, config.CircuitBreakerClientConfig{}, nil, nil)
	if err != nil {
		t.Errorf("fail to init site: %v", err)
		return
	}

	bksDB := stubData(st.rp, site)

	tests := []struct {
		name       string
		st         *Site
		bkID       int
		bkHash     int
		expectErr  bool
		expectBook model.Book
	}{
		{
			name:      "works for updating in progress to end",
			st:        st,
			bkID:      3,
			bkHash:    0,
			expectErr: false,
			expectBook: model.Book{
				Site: site, ID: 3, HashCode: 0,
				Title: "title 3", Writer: bksDB[3].Writer,
				Type: "type 3", UpdateDate: "date 3", UpdateChapter: "end " + repo.ChapterEndKeywords[0] + " end",
				Status: model.End,
			},
		},
		{
			name:      "works for not updating in progress to end",
			st:        st,
			bkID:      2,
			bkHash:    100,
			expectErr: false,
			expectBook: model.Book{
				Site: site, ID: 2, HashCode: 100,
				Title: "title 2 new", Writer: bksDB[2].Writer, Type: "type 2 new",
				UpdateDate: "date 2.1", UpdateChapter: bksDB[2].UpdateChapter,
				Status: model.InProgress,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			bksDB[2].Status = model.InProgress
			test.st.rp.UpdateBook(&bksDB[2])
			bksDB[3].UpdateChapter = fmt.Sprintf("end %v end", repo.ChapterEndKeywords[0])
			test.st.rp.UpdateBook(&bksDB[3])

			err := Check(test.st)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want error: %v", err, test.expectErr)
			}

			bk, err := test.st.BookFromIDHash(test.bkID, strconv.FormatInt(int64(test.bkHash), 36))
			if err != nil {
				t.Errorf("book fail to fetch: id: %v; hash: %v; err: %v", test.bkID, test.bkHash, err)
				return
			}
			if !cmp.Equal(*bk, test.expectBook) {
				t.Errorf("book diff: %v", cmp.Diff(bk, test.expectBook))
			}
		})
	}
}
