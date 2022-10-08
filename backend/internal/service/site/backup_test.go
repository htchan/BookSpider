package site

import (
	"errors"
	"testing"

	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	psql "github.com/htchan/BookSpider/internal/repo/psql"
)

func stubData(r repo.Repostory, site string) []model.Book {
	bks := []model.Book{
		{
			Site: site, ID: 1, HashCode: 0,
			Title: "title 1", Writer: model.Writer{Name: site + " writer 1"}, Type: "type 1",
			UpdateDate: "date 1", UpdateChapter: "chapter 1",
			Status: model.End, IsDownloaded: true, Error: nil,
		},
		{
			Site: site, ID: 2, HashCode: 0,
			Title: "title 2", Writer: model.Writer{Name: site + " writer 2"}, Type: "type 2",
			UpdateDate: "date 2", UpdateChapter: "chapter 2",
			Status: model.End, IsDownloaded: true, Error: nil,
		},
		{
			Site: site, ID: 2, HashCode: 100,
			Title: "title 2 new", Writer: model.Writer{Name: site + " writer 2 new"}, Type: "type 2 new",
			UpdateDate: "date 2.1", UpdateChapter: "chapter 2 new",
			Status: model.End, IsDownloaded: false, Error: nil,
		},
		{
			Site: site, ID: 3, HashCode: 0,
			Title: "title 3", Writer: model.Writer{Name: site + " writer 3"}, Type: "type 3",
			UpdateDate: "date 3", UpdateChapter: "chapter 3",
			Status: model.InProgress, IsDownloaded: false, Error: nil,
		},
		{
			Site: site, ID: 4, HashCode: 0,
			Title: "", Writer: model.Writer{Name: ""}, Type: "",
			UpdateDate: "", UpdateChapter: "",
			Status: model.Error, IsDownloaded: false, Error: errors.New("error"),
		},
	}

	for i := range bks {
		err := r.SaveWriter(&bks[i].Writer)
		if err != nil {
			// panic(err)
		}
		err = r.CreateBook(&bks[i])
		if err != nil {
			// panic(err)
		}
		err = r.SaveError(&bks[i], bks[i].Error)
		if err != nil {
			// panic(err)
		}
	}
	return bks
}

func TestSite_Backup(t *testing.T) {
	// This is copied from internal/repo/psql/database_test
	t.Parallel()

	psql.StubPsqlConn()
	db, err := psql.OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "backup"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
		db.Close()
	})

	stubData(psql.NewRepo(site, db), site)

	st, err := NewSite(site, &config.BookConfig{}, &config.SiteConfig{BackupDirectory: "/backup"}, &config.CircuitBreakerClientConfig{}, nil, nil)
	if err != nil {
		t.Errorf("fail to init site: %v", err)
		return
	}

	tests := []struct {
		name      string
		st        *Site
		path      string
		expectErr bool
	}{
		{
			name:      "works",
			st:        st,
			expectErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			err := Backup(test.st)
			if (err != nil) != test.expectErr {
				t.Errorf("got err: %v; want err: %v", err, test.expectErr)
			}
		})
	}
}
