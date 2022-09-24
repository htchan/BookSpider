package site

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/model"
	psql "github.com/htchan/BookSpider/internal/repo/psql"
)

func TestSite_BookFromID(t *testing.T) {
	t.Parallel()

	psql.StubPsqlConn()
	db, err := psql.OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "find_id/st"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
	})

	st, err := NewSite(
		site,
		config.BookConfig{},
		config.SiteConfig{},
		config.CircuitBreakerClientConfig{}, nil, nil)
	if err != nil {
		t.Errorf("fail to init site: %v", err)
		return
	}

	bksDB := stubData(st.rp, site)

	tests := []struct {
		name      string
		st        *Site
		id        int
		expect    *model.Book
		expectErr bool
	}{
		{
			name: "return book with hash 0",
			st:   st,
			id:   1,
			expect: &model.Book{
				Site: site, ID: 1,
				Title: "title 1", Writer: bksDB[0].Writer, Type: "type 1",
				UpdateDate: "date 1", UpdateChapter: "chapter 1", Status: model.End, IsDownloaded: true,
			},
			expectErr: false,
		},
		{
			name: "return book with largest hash",
			st:   st,
			id:   2,
			expect: &model.Book{
				Site: site, ID: 2, HashCode: 100,
				Title: "title 2 new", Writer: bksDB[2].Writer, Type: "type 2 new",
				UpdateDate: "date 2.1", UpdateChapter: "chapter 2 new", Status: model.End,
			},
			expectErr: false,
		},
		{
			name:      "return error for books not exist",
			st:        st,
			id:        999,
			expect:    nil,
			expectErr: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			bk, err := test.st.BookFromID(test.id)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want error: %v", err, test.expectErr)
			}

			if !cmp.Equal(bk, test.expect) {
				t.Errorf("book diff: %v", cmp.Diff(bk, test.expect))
			}
		})
	}
}

func TestSite_BookFromIDHash(t *testing.T) {
	t.Parallel()

	psql.StubPsqlConn()
	db, err := psql.OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "find_id_hash/st"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
	})

	st, err := NewSite(
		site,
		config.BookConfig{},
		config.SiteConfig{},
		config.CircuitBreakerClientConfig{}, nil, nil)
	if err != nil {
		t.Errorf("fail to init site: %v", err)
		return
	}

	bksDB := stubData(st.rp, site)

	tests := []struct {
		name      string
		st        *Site
		id        int
		hash      string
		expect    *model.Book
		expectErr bool
	}{
		{
			name: "return book with correct hash",
			st:   st,
			id:   2,
			hash: "2s",
			expect: &model.Book{
				Site: site, ID: 2, HashCode: 100,
				Title: "title 2 new", Writer: bksDB[2].Writer, Type: "type 2 new",
				UpdateDate: "date 2.1", UpdateChapter: "chapter 2 new", Status: model.End,
			},
			expectErr: false,
		},
		{
			name:      "return error for books not exist",
			st:        st,
			id:        999,
			hash:      "0",
			expect:    nil,
			expectErr: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			bk, err := test.st.BookFromIDHash(test.id, test.hash)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want error: %v", err, test.expectErr)
			}

			if !cmp.Equal(bk, test.expect) {
				t.Errorf("book diff: %v", cmp.Diff(bk, test.expect))
			}
		})
	}
}

func TestSite_QueryBooks(t *testing.T) {
	t.Parallel()

	psql.StubPsqlConn()
	db, err := psql.OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "find_book/st"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
	})

	st, err := NewSite(
		site,
		config.BookConfig{},
		config.SiteConfig{},
		config.CircuitBreakerClientConfig{}, nil, nil)
	if err != nil {
		t.Errorf("fail to init site: %v", err)
		return
	}

	bksDB := stubData(st.rp, site)

	tests := []struct {
		name          string
		st            *Site
		title, writer string
		limit, offset int
		expect        []model.Book
		expectErr     bool
	}{
		{
			name:      "works for title and writer",
			st:        st,
			title:     "title 1",
			writer:    "writer 2",
			limit:     5,
			offset:    0,
			expect:    bksDB[0:3],
			expectErr: false,
		},
		{
			name:      "works for offset",
			st:        st,
			title:     "title 1",
			writer:    "writer 2",
			limit:     1,
			offset:    2,
			expect:    bksDB[2:3],
			expectErr: false,
		},
		{
			name:      "works for not exist",
			st:        st,
			title:     "writer 1",
			writer:    "title 2",
			limit:     10,
			offset:    0,
			expect:    []model.Book{},
			expectErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			bks, err := test.st.QueryBooks(test.title, test.writer, test.limit, test.offset)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want error: %v", err, test.expectErr)
			}

			if !cmp.Equal(bks, test.expect) {
				t.Errorf("book diff: %v", cmp.Diff(bks, test.expect))
			}
		})
	}
}

func TestSite_RandomBooks(t *testing.T) {

	t.Parallel()

	psql.StubPsqlConn()
	db, err := psql.OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "random_book/st"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
	})

	st, err := NewSite(
		site,
		config.BookConfig{},
		config.SiteConfig{},
		config.CircuitBreakerClientConfig{}, nil, nil)
	if err != nil {
		t.Errorf("fail to init site: %v", err)
		return
	}

	stubData(st.rp, site)

	tests := []struct {
		name      string
		st        *Site
		limit     int
		expectLen int
		expectErr bool
	}{
		{
			name:      "works without offset",
			st:        st,
			limit:     10,
			expectLen: 2,
			expectErr: false,
		},
		{
			name:      "works with offset",
			st:        st,
			limit:     10,
			expectLen: 1,
			expectErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			bks, err := test.st.RandomBooks(test.limit)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want error: %v", err, test.expectErr)
			}

			if len(bks) != test.expectLen {
				t.Errorf("book len diff: books: %v; expect len: %v", bks, test.expectLen)
			}
		})
	}
}
