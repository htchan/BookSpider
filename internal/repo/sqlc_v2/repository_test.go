package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/stretchr/testify/assert"
)

func stubData(t testing.TB, r repo.RepositoryV2, site string) []model.Book {
	t.Helper()

	bks := []model.Book{
		{
			Site: site, ID: 1, HashCode: 0,
			Title: "title 1", Writer: model.Writer{Name: site + " writer 1"}, Type: "type 1",
			UpdateDate: "date 1", UpdateChapter: "chapter 1",
			Status: model.StatusEnd, IsDownloaded: true, Error: nil,
		},
		{
			Site: site, ID: 2, HashCode: 0,
			Title: "title 2", Writer: model.Writer{Name: site + " writer 2"}, Type: "type 2",
			UpdateDate: "date 2", UpdateChapter: "chapter 2",
			Status: model.StatusEnd, IsDownloaded: true, Error: nil,
		},
		{
			Site: site, ID: 2, HashCode: 100,
			Title: "title 2 new", Writer: model.Writer{Name: site + " writer 2 new"}, Type: "type 2 new",
			UpdateDate: "date 2.1", UpdateChapter: "chapter 2 new",
			Status: model.StatusEnd, IsDownloaded: false, Error: nil,
		},
		{
			Site: site, ID: 3, HashCode: 0,
			Title: "title 3", Writer: model.Writer{Name: site + " writer 3"}, Type: "type 3",
			UpdateDate: "date 3", UpdateChapter: "chapter 3",
			Status: model.StatusInProgress, IsDownloaded: false, Error: nil,
		},
		{
			Site: site, ID: 4, HashCode: 0,
			Title: "", Writer: model.Writer{Name: ""}, Type: "",
			UpdateDate: "", UpdateChapter: "",
			Status: model.StatusError, IsDownloaded: false, Error: errors.New("error"),
		},
	}

	for i := range bks {
		var err error
		for _ = range 5 {
			err = r.SaveWriter(t.Context(), &bks[i].Writer)
			if err == nil {
				break
			}
		}
		if !assert.NoError(t, err, "Failed to save writer %d: %v", i, bks[i].Writer) {
			t.FailNow()
		}

		err = r.CreateBook(t.Context(), &bks[i])
		if !assert.NoError(t, err, "Failed to create book %d: %v", i, bks[i]) {
			t.FailNow()
		}

		err = r.SaveError(t.Context(), &bks[i], bks[i].Error)
		if !assert.NoError(t, err, "Failed to save error for book %d: %v", i, bks[i]) {
			t.FailNow()
		}
	}

	return bks
}

func Test_NewRepo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		site   string
		db     *sql.DB
		expect *SqlcRepo
	}{
		{
			name:   "works",
			site:   "test",
			db:     nil,
			expect: &SqlcRepo{db: nil},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := NewRepo(test.db)
			if test.db != result.db {
				t.Errorf("got: %v, want: %v", result, test.expect)
			}
		})
	}
}

func TestSqlcRepo_CreateBook(t *testing.T) {
	t.Parallel()
	StubPsqlConn()
	db := testDB

	site := "bk/create"
	db.Exec("ALTER SEQUENCE writers_id_seq RESTART WITH 1;")

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)

	})

	tests := []struct {
		name       string
		r          *SqlcRepo
		bk         model.Book
		expectBook model.Book
		expectErr  bool
	}{
		{
			name:       "create book with new id to hash code 0",
			r:          NewRepo(db),
			bk:         model.Book{Site: site, ID: 1, HashCode: 100, Writer: model.Writer{ID: 10}},
			expectBook: model.Book{Site: site, ID: 1, HashCode: 0, Writer: model.Writer{ID: 10}},
			expectErr:  false,
		},
		{
			name:       "create book with existing id with input hash code",
			r:          NewRepo(db),
			bk:         model.Book{Site: site, ID: 1, HashCode: 100, Writer: model.Writer{ID: 10}},
			expectBook: model.Book{Site: site, ID: 1, HashCode: 100, Writer: model.Writer{ID: 10}},
			expectErr:  false,
		},
		{
			name:       "fail to create book with existing id and hash",
			r:          NewRepo(db),
			bk:         model.Book{Site: site, ID: 1, HashCode: 100, Writer: model.Writer{ID: 10}},
			expectBook: model.Book{Site: site, ID: 1, HashCode: 100, Writer: model.Writer{ID: 10}},
			expectErr:  true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			err := test.r.CreateBook(context.Background(), &test.bk)
			t.Run("result", func(t *testing.T) {
				if (err != nil) != test.expectErr {
					t.Errorf("got error: %v, expect err: %v", err, test.expectErr)
				}

				assert.Equal(t, test.expectBook, test.bk)
			})

			t.Run("book in db", func(t *testing.T) {
				bk, err := test.r.FindBookByIdHash(context.Background(), site, test.bk.ID, test.bk.HashCode)
				if err != nil {
					t.Fatalf("query got error: %v", err)
				}

				bk.Writer.Name = "" // Ignore writer name for comparison

				assert.Equal(t, test.expectBook, *bk)
			})
		})
	}
}

func TestSqlcRepo_UpdateBook(t *testing.T) {
	t.Parallel()
	StubPsqlConn()
	db := testDB
	site := "bk/update"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)

	})

	bksDB := stubData(t, NewRepo(db), site)

	tests := []struct {
		name              string
		r                 repo.RepositoryV2
		inputBook         *model.Book
		expectErr         bool
		expectQueryResult *model.Book
		expectQueryErr    bool
	}{
		{
			name: "update not existing book",
			r:    NewRepo(db),
			inputBook: &model.Book{
				Site: site, ID: -1, HashCode: 0, Title: "hello",
			},
			expectErr:         true,
			expectQueryResult: nil,
			expectQueryErr:    true,
		},
		{
			name: "update error book to in progress without changing writer id and error",
			r:    NewRepo(db),
			inputBook: &model.Book{
				Site: bksDB[4].Site, ID: bksDB[4].ID, HashCode: bksDB[4].HashCode,
				Title: "t", Writer: model.Writer{Name: bksDB[0].Writer.Name}, Type: "t",
				UpdateDate: "d", UpdateChapter: "c",
				Status: model.StatusInProgress, IsDownloaded: false, Error: nil,
			},
			expectErr: false,
			expectQueryResult: &model.Book{
				Site: bksDB[4].Site, ID: bksDB[4].ID, HashCode: bksDB[4].HashCode,
				Title: "t", Writer: model.Writer{ID: 0, Name: ""}, Type: "t",
				UpdateDate: "d", UpdateChapter: "c",
				Status: model.StatusInProgress, IsDownloaded: false, Error: errors.New("error"),
			},
			expectQueryErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			err := test.r.UpdateBook(context.Background(), test.inputBook)
			t.Run("result", func(t *testing.T) {
				if (err != nil) != test.expectErr {
					t.Errorf("got error: %v; want error: %v", err, test.expectErr)
				}
			})

			t.Run("book in db", func(t *testing.T) {
				bk, err := test.r.FindBookByIdHash(context.Background(), site, test.inputBook.ID, test.inputBook.HashCode)
				if (err != nil) != test.expectQueryErr {
					t.Errorf("query got error: %v; want error: %v", err, test.expectQueryErr)
				}
				assert.Equal(t, test.expectQueryResult, bk)
			})
		})
	}
}

func TestSqlcRepo_FindBookByID(t *testing.T) {
	t.Parallel()
	StubPsqlConn()
	db := testDB
	site := "bk_id/find"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)

	})

	bksDB := stubData(t, NewRepo(db), site)

	tests := []struct {
		name         string
		r            repo.RepositoryV2
		id           int
		expectResult *model.Book
		expectHash   int
		expectErr    bool
	}{
		{
			name:         "find not existing book",
			r:            NewRepo(db),
			id:           0,
			expectResult: nil,
			expectHash:   0,
			expectErr:    true,
		},
		{
			name:         "find book with largest id",
			r:            NewRepo(db),
			id:           2,
			expectResult: &bksDB[2],
			expectHash:   100,
			expectErr:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			result, err := test.r.FindBookById(context.Background(), site, test.id)
			t.Run("result", func(t *testing.T) {
				if (err != nil) != test.expectErr {
					t.Errorf("got error: %v; want error: %v", err, test.expectErr)
				}

				assert.Equal(t, test.expectResult, result)
			})
		})
	}
}

func TestSqlcRepo_FindBookByIDHash(t *testing.T) {
	t.Parallel()
	StubPsqlConn()
	db := testDB
	site := "bk_id_hash/find"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)

	})

	bksDB := stubData(t, NewRepo(db), site)

	tests := []struct {
		name         string
		r            repo.RepositoryV2
		id           int
		hashcode     int
		expectResult *model.Book
		expectErr    bool
	}{
		{
			name:         "find not existing book",
			r:            NewRepo(db),
			id:           0,
			hashcode:     0,
			expectResult: nil,
			expectErr:    true,
		},
		{
			name:         "find book with correct id hash",
			r:            NewRepo(db),
			id:           2,
			hashcode:     0,
			expectResult: &bksDB[1],
			expectErr:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			result, err := test.r.FindBookByIdHash(context.Background(), site, test.id, test.hashcode)
			t.Run("result", func(t *testing.T) {
				if (err != nil) != test.expectErr {
					t.Errorf("got error: %v; want error: %v", err, test.expectErr)
				}

				assert.Equal(t, test.expectResult, result)
			})
		})
	}
}

func TestSqlcRepo_FindBookByStatus(t *testing.T) {
	t.Parallel()
	//Check if this will really be used
}

func TestSqlcRepo_FindAllBooks(t *testing.T) {
	t.Parallel()
	StubPsqlConn()
	db := testDB
	site := "all_bk/find"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)

	})

	bksDB := stubData(t, NewRepo(db), site)

	tests := []struct {
		name         string
		r            repo.RepositoryV2
		site         string
		expectResult []model.Book
		expectErr    bool
	}{
		{
			name:         "works",
			r:            NewRepo(db),
			site:         site,
			expectResult: bksDB,
			expectErr:    false,
		},
		{
			name:         "works for empty",
			r:            NewRepo(db),
			site:         "empty",
			expectResult: nil,
			expectErr:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			result, err := test.r.FindAllBooks(context.Background(), test.site)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want error: %v", err, test.expectErr)
			}

			var bks []model.Book
			for bk := range result {
				bks = append(bks, bk)
			}

			assert.Equal(t, test.expectResult, bks)
		})
	}
}

func TestSqlcRepo_FindBooksForUpdate(t *testing.T) {
	t.Parallel()
	StubPsqlConn()
	db := testDB
	site := "update_bk/find"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)

	})

	bksDB := stubData(t, NewRepo(db), site)

	tests := []struct {
		name         string
		r            repo.RepositoryV2
		expectResult []model.Book
		expectErr    bool
	}{
		{
			name:         "works",
			r:            NewRepo(db),
			expectResult: []model.Book{bksDB[4], bksDB[3], bksDB[2], bksDB[0]},
			expectErr:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			result, err := test.r.FindBooksForUpdate(context.Background(), site)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want error: %v", err, test.expectErr)
			}

			var bks []model.Book
			for bk := range result {
				bks = append(bks, bk)
			}

			assert.Equal(t, test.expectResult, bks)
		})
	}
}

func TestSqlcRepo_FindBooksForDownload(t *testing.T) {
	t.Parallel()
	StubPsqlConn()
	db := testDB
	site := "down_bk/find"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
	})

	bksDB := stubData(t, NewRepo(db), site)

	tests := []struct {
		name         string
		r            repo.RepositoryV2
		expectResult []model.Book
		expectErr    bool
	}{
		{
			name:         "works",
			r:            NewRepo(db),
			expectResult: []model.Book{bksDB[2]},
			expectErr:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			result, err := test.r.FindBooksForDownload(context.Background(), site)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want error: %v", err, test.expectErr)
			}

			var bks []model.Book
			for bk := range result {
				bks = append(bks, bk)
			}

			assert.Equal(t, test.expectResult, bks)
		})
	}
}

func TestSqlcRepo_FindBooksByTitleWriter(t *testing.T) {
	StubPsqlConn()
	db := testDB
	site := "bk_tit_wrt/find"
	siteV2 := "bk_tit_wrt/v2"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
	})
	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", siteV2)
		db.Exec("delete from writers where id>0 and name like $1", siteV2+"%")
		db.Exec("delete from errors where site=$1", siteV2)
	})

	bksDB := stubData(t, NewRepo(db), site)
	bksDBV2 := stubData(t, NewRepo(db), siteV2)

	tests := []struct {
		name         string
		r            repo.RepositoryV2
		title        string
		writer       string
		limit        int
		offset       int
		expectResult []model.Book
		expectErr    bool
	}{
		{
			name:   "works",
			r:      NewRepo(db),
			title:  "title",
			writer: "writer",
			limit:  10,
			offset: 0,
			expectResult: []model.Book{
				bksDBV2[3], bksDB[3],
				bksDBV2[2], bksDB[2],
				bksDBV2[1], bksDB[1],
				bksDBV2[0], bksDB[0],
			},
			expectErr: false,
		},
		{
			name:         "works with limit",
			r:            NewRepo(db),
			title:        "title",
			writer:       "writer",
			limit:        1,
			offset:       0,
			expectResult: []model.Book{bksDBV2[3]},
			expectErr:    false,
		},
		{
			name:         "works with offset",
			r:            NewRepo(db),
			title:        "title",
			writer:       "writer",
			limit:        1,
			offset:       1,
			expectResult: []model.Book{bksDB[3]},
			expectErr:    false,
		},
		{
			name:         "return all books match either title or writer",
			r:            NewRepo(db),
			title:        "title 1",
			writer:       "writer 3",
			limit:        5,
			offset:       0,
			expectResult: []model.Book{bksDBV2[3], bksDB[3], bksDBV2[0], bksDB[0]},
			expectErr:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			result, err := test.r.FindBooksByTitleWriter(context.Background(), test.title, test.writer, test.limit, test.offset)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want err: %v", err, test.expectErr)
			}
			assert.Equal(t, test.expectResult, result)
		})
	}
}

func TestSqlcRepo_FindBooksByRandom(t *testing.T) {
	t.Parallel()

	StubPsqlConn()
	db := testDB
	site := "rand_bk/find"
	siteV2 := "rand_bk/find_v2"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
	})
	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", siteV2)
		db.Exec("delete from writers where id>0 and name like $1", siteV2+"%")
		db.Exec("delete from errors where site=$1", siteV2)
	})

	stubData(t, NewRepo(db), site)
	stubData(t, NewRepo(db), siteV2)

	tests := []struct {
		name         string
		r            repo.RepositoryV2
		limit        int
		expectLength int
		expectErr    bool
	}{
		{
			name:         "works",
			r:            NewRepo(db),
			limit:        10,
			expectLength: 4,
			expectErr:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			result, err := test.r.FindBooksByRandom(context.Background(), test.limit)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want err: %v", err, test.expectErr)
			}
			if len(result) != test.expectLength {
				t.Errorf("query got:  %v\nwant length: %v", result, test.expectLength)
			}
			for _, bk := range result {
				if bk.Site != site && bk.Site != siteV2 {
					t.Errorf("query got book with unexpected site: %v", bk.Site)
				}
			}
		})
	}
}

func TestSqlcRepo_UpdateBooksStatus(t *testing.T) {
	t.Parallel()

	StubPsqlConn()
	db := testDB
	site := "stat_bk/update"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)

	})

	r := NewRepo(db)

	bksDB := stubData(t, NewRepo(db), site)

	tests := []struct {
		name       string
		r          repo.RepositoryV2
		bkID       int
		bkHash     int
		expectErr  bool
		expectBook *model.Book
	}{
		{
			name:      "works for updating in progress to end",
			r:         NewRepo(db),
			bkID:      3,
			bkHash:    0,
			expectErr: false,
			expectBook: &model.Book{
				Site: site, ID: 3, HashCode: 0,
				Title: "title 3", Writer: bksDB[3].Writer,
				Type: "type 3", UpdateDate: "date 3", UpdateChapter: "end " + model.ChapterEndKeywords[0] + " end",
				Status: model.StatusEnd, IsDownloaded: false,
			},
		},
		{
			name:      "works for updating download to false",
			r:         NewRepo(db),
			bkID:      1,
			bkHash:    0,
			expectErr: false,
			expectBook: &model.Book{
				Site: site, ID: 1, HashCode: 0,
				Title: "title 1", Writer: bksDB[0].Writer,
				Type: "type 1", UpdateDate: "date 1", UpdateChapter: "end " + model.ChapterEndKeywords[0] + " end",
				Status: model.StatusEnd, IsDownloaded: false,
			},
		},
		{
			name:      "works for not updating in progress to end",
			r:         NewRepo(db),
			bkID:      2,
			bkHash:    100,
			expectErr: false,
			expectBook: &model.Book{
				Site: site, ID: 2, HashCode: 100,
				Title: "title 2 new", Writer: bksDB[2].Writer, Type: "type 2 new",
				UpdateDate: "date 2.1", UpdateChapter: bksDB[2].UpdateChapter,
				Status: model.StatusInProgress, IsDownloaded: false,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			bksDB[2].Status = model.StatusInProgress
			r.UpdateBook(context.Background(), &bksDB[2])
			bksDB[3].UpdateChapter = fmt.Sprintf("end %v end", model.ChapterEndKeywords[0])
			r.UpdateBook(context.Background(), &bksDB[3])
			bksDB[0].UpdateChapter = fmt.Sprintf("end %v end", model.ChapterEndKeywords[0])
			bksDB[0].Status = model.StatusInProgress
			r.UpdateBook(context.Background(), &bksDB[0])

			err := test.r.UpdateBooksStatus(context.Background())
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want error: %v", err, test.expectErr)
			}

			bk, err := test.r.FindBookByIdHash(context.Background(), site, test.bkID, test.bkHash)
			if err != nil {
				t.Errorf("book fail to fetch: id: %v; hash: %v; err: %v", test.bkID, test.bkHash, err)
				return
			}
			assert.Equal(t, test.expectBook, bk)
		})
	}
}

func BenchmarkSqlcRepo_UpdateBooksStatus(b *testing.B) {
	StubPsqlConn()
	db := testDB
	site := "bm/bk_st/update"

	b.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)

	})

	r := NewRepo(db)

	stubData(b, NewRepo(db), site)

	for n := 0; n < b.N; n++ {
		r.UpdateBooksStatus(context.Background())
	}
}

func TestSqlcRepo_FindAllBookIDs(t *testing.T) {
	t.Parallel()

	StubPsqlConn()
	db := testDB
	site := "bk/find_all_ids"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)

	})

	r := NewRepo(db)

	stubData(t, NewRepo(db), site)

	tests := []struct {
		name      string
		r         repo.RepositoryV2
		wantError error
		want      []int
	}{
		{
			name:      "happy flow",
			r:         r,
			wantError: nil,
			want:      []int{1, 2, 3, 4},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := test.r.FindAllBookIDs(context.Background(), site)
			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, test.wantError, err)
		})
	}
}

func TestSqlcRepo_SaveWriter(t *testing.T) {
	t.Parallel()
	StubPsqlConn()
	db := testDB
	site := "writer/save"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)

	})

	bksDB := stubData(t, NewRepo(db), site)

	tests := []struct {
		name      string
		r         repo.RepositoryV2
		writer    *model.Writer
		expectErr bool
	}{
		{
			name:      "save existing writer",
			r:         NewRepo(db),
			writer:    &model.Writer{ID: 0, Name: bksDB[0].Writer.Name},
			expectErr: false,
		},
		{
			name:      "save new writer",
			r:         NewRepo(db),
			writer:    &model.Writer{ID: 0, Name: site + " new writer"},
			expectErr: false,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			err := test.r.SaveWriter(context.Background(), test.writer)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v, expect err: %v", err, test.expectErr)
			}

			if test.writer.ID <= 0 {
				t.Errorf("got writer:  %v", test.writer)
			}
		})
	}
}

func TestSqlcRepo_SaveError(t *testing.T) {
	t.Parallel()
	StubPsqlConn()
	db := testDB
	site := "error/save"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)

	})

	stubData(t, NewRepo(db), site)

	tests := []struct {
		name         string
		r            repo.RepositoryV2
		bk           *model.Book
		e            error
		expectErrStr string
		expectErr    bool
	}{
		{
			name:         "create error for existing book",
			r:            NewRepo(db),
			bk:           &model.Book{Site: site, ID: 1},
			e:            errors.New("create error"),
			expectErrStr: "create error",
			expectErr:    false,
		},
		{
			name:         "update error for existing book",
			r:            NewRepo(db),
			bk:           &model.Book{Site: site, ID: 1},
			e:            errors.New("update error"),
			expectErrStr: "update error",
			expectErr:    false,
		},
		{
			name:         "delete error for existing book",
			r:            NewRepo(db),
			bk:           &model.Book{Site: site, ID: 1},
			e:            nil,
			expectErrStr: "",
			expectErr:    false,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			err := test.r.SaveError(context.Background(), test.bk, test.e)
			t.Run("result", func(t *testing.T) {
				if (err != nil) != test.expectErr {
					t.Errorf("got error: %v, expect err: %v", err, test.expectErr)
				}
			})

			t.Run("error in db", func(t *testing.T) {
				bk, err := test.r.FindBookByIdHash(context.Background(), site, test.bk.ID, test.bk.HashCode)
				if err != nil {
					t.Errorf("query got error: %v; want error: %v", err, false)
				}
				if !((bk.Error == nil && test.expectErrStr == "") ||
					(bk.Error != nil && bk.Error.Error() == test.expectErrStr)) {
					t.Errorf("query got:  %v\nwant: %v", bk.Error, test.expectErrStr)
					t.Error(cmp.Diff(bk.Error.Error(), test.expectErrStr))
				}
			})
		})
	}
}

// func TestSqlcRepo_Backup(t *testing.T) {
// 	t.Parallel()
// 	StubPsqlConn()
// 	db := testDB
// 	site := "backup"

// 	t.Cleanup(func() {
// 		db.Exec("delete from books where site=$1", site)
// 		db.Exec("delete from writers where id>0 and name like $1", site+"%")
// 		db.Exec("delete from errors where site=$1", site)

// 	})

// 	stubData(t, NewRepo(db), site)

// 	tests := []struct {
// 		name      string
// 		r         repo.RepositoryV2
// 		path      string
// 		expectErr bool
// 	}{
// 		{
// 			name:      "works",
// 			r:         NewRepo(db),
// 			path:      "/" + site,
// 			expectErr: false,
// 		},
// 	}

// 	for _, test := range tests {
// 		test := test
// 		t.Run(test.name, func(t *testing.T) {
// 			err := test.r.Backup(context.Background(), site, test.path)
// 			if (err != nil) != test.expectErr {
// 				t.Errorf("got err: %v; want err: %v", err, test.expectErr)
// 			}
// 		})
// 	}
// }

func TestSqlcRepo_Stats(t *testing.T) {
	t.Parallel()

	StubPsqlConn()
	db := testDB
	site := "stats"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)

	})

	stubData(t, NewRepo(db), site)

	tests := []struct {
		name   string
		r      repo.RepositoryV2
		expect repo.Summary
	}{
		{
			name: "works",
			r:    NewRepo(db),
			expect: repo.Summary{
				BookCount: 5, WriterCount: 5, ErrorCount: 1, DownloadCount: 2,
				UniqueBookCount: 4, MaxBookID: 4,
				LatestSuccessID: 3, StatusCount: map[model.StatusCode]int{
					model.StatusError:      1,
					model.StatusInProgress: 1,
					model.StatusEnd:        3,
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := test.r.Stats(context.Background(), site)
			assert.Equal(t, test.expect, result)
		})
	}
}
