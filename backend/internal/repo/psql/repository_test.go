package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/stretchr/testify/assert"
)

func stubData(r repo.Repository, site string) []model.Book {
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

func Test_NewRepo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		site   string
		db     *sql.DB
		expect *PsqlRepo
	}{
		{
			name:   "works",
			site:   "test",
			db:     nil,
			expect: &PsqlRepo{site: "test", db: nil},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := NewRepo(test.site, test.db)
			if *result != *test.expect {
				t.Errorf("got: %v, want: %v", result, test.expect)
			}
		})
	}
}

func TestPsqlRepo_CreateBook(t *testing.T) {
	t.Parallel()
	StubPsqlConn()
	db, err := OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "bk/create"
	db.Exec("ALTER SEQUENCE writers_id_seq RESTART WITH 1;")

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Close()
	})

	tests := []struct {
		name       string
		r          *PsqlRepo
		bk         model.Book
		expectBook model.Book
		expectErr  bool
	}{
		{
			name:       "create book with new id to hash code 0",
			r:          NewRepo(site, db),
			bk:         model.Book{Site: site, ID: 1, HashCode: 100, Writer: model.Writer{ID: 10}},
			expectBook: model.Book{Site: site, ID: 1, HashCode: 0, Writer: model.Writer{ID: 10}},
			expectErr:  false,
		},
		{
			name:       "create book with existing id with input hash code",
			r:          NewRepo(site, db),
			bk:         model.Book{Site: site, ID: 1, HashCode: 100, Writer: model.Writer{ID: 10}},
			expectBook: model.Book{Site: site, ID: 1, HashCode: 100, Writer: model.Writer{ID: 10}},
			expectErr:  false,
		},
		{
			name:       "fail to create book with existing id and hash",
			r:          NewRepo(site, db),
			bk:         model.Book{Site: site, ID: 1, HashCode: 100, Writer: model.Writer{ID: 10}},
			expectBook: model.Book{Site: site, ID: 1, HashCode: 100, Writer: model.Writer{ID: 10}},
			expectErr:  true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			err := test.r.CreateBook(&test.bk)
			t.Run("result", func(t *testing.T) {
				if (err != nil) != test.expectErr {
					t.Errorf("got error: %v, expect err: %v", err, test.expectErr)
				}

				assert.Equal(t, test.expectBook, test.bk)
			})

			t.Run("book in db", func(t *testing.T) {
				bk, err := test.r.FindBookByIdHash(test.bk.ID, test.bk.HashCode)
				if err != nil {
					t.Fatalf("query got error: %v", err)
				}

				assert.Equal(t, test.expectBook, *bk)
			})
		})
	}
}

func TestPsqlRepo_UpdateBook(t *testing.T) {
	t.Parallel()
	StubPsqlConn()
	db, err := OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "bk/update"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
		db.Close()
	})

	bksDB := stubData(NewRepo(site, db), site)

	tests := []struct {
		name              string
		r                 repo.Repository
		inputBook         *model.Book
		expectErr         bool
		expectQueryResult *model.Book
		expectQueryErr    bool
	}{
		{
			name: "update not existing book",
			r:    NewRepo(site, db),
			inputBook: &model.Book{
				Site: site, ID: -1, HashCode: 0, Title: "hello",
			},
			expectErr:         false,
			expectQueryResult: nil,
			expectQueryErr:    true,
		},
		{
			name: "update error book to in progress without changing writer id and error",
			r:    NewRepo(site, db),
			inputBook: &model.Book{
				Site: bksDB[4].Site, ID: bksDB[4].ID, HashCode: bksDB[4].HashCode,
				Title: "t", Writer: model.Writer{Name: bksDB[0].Writer.Name}, Type: "t",
				UpdateDate: "d", UpdateChapter: "c",
				Status: model.InProgress, IsDownloaded: false, Error: nil,
			},
			expectErr: false,
			expectQueryResult: &model.Book{
				Site: bksDB[4].Site, ID: bksDB[4].ID, HashCode: bksDB[4].HashCode,
				Title: "t", Writer: model.Writer{ID: 0, Name: ""}, Type: "t",
				UpdateDate: "d", UpdateChapter: "c",
				Status: model.InProgress, IsDownloaded: false, Error: errors.New("error"),
			},
			expectQueryErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			err := test.r.UpdateBook(test.inputBook)
			t.Run("result", func(t *testing.T) {
				if (err != nil) != test.expectErr {
					t.Errorf("got error: %v; want error: %v", err, test.expectErr)
				}
			})

			t.Run("book in db", func(t *testing.T) {
				bk, err := test.r.FindBookByIdHash(test.inputBook.ID, test.inputBook.HashCode)
				if (err != nil) != test.expectQueryErr {
					t.Errorf("query got error: %v; want error: %v", err, test.expectQueryErr)
				}
				assert.Equal(t, test.expectQueryResult, bk)
			})
		})
	}
}

func TestPsqlRepo_FindBookByID(t *testing.T) {
	t.Parallel()
	StubPsqlConn()
	db, err := OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "bk_id/find"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
		db.Close()
	})

	bksDB := stubData(NewRepo(site, db), site)

	tests := []struct {
		name         string
		r            repo.Repository
		id           int
		expectResult *model.Book
		expectHash   int
		expectErr    bool
	}{
		{
			name:         "find not existing book",
			r:            NewRepo(site, db),
			id:           0,
			expectResult: nil,
			expectHash:   0,
			expectErr:    true,
		},
		{
			name:         "find book with largest id",
			r:            NewRepo(site, db),
			id:           2,
			expectResult: &bksDB[2],
			expectHash:   100,
			expectErr:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			result, err := test.r.FindBookById(test.id)
			t.Run("result", func(t *testing.T) {
				if (err != nil) != test.expectErr {
					t.Errorf("got error: %v; want error: %v", err, test.expectErr)
				}

				assert.Equal(t, test.expectResult, result)
			})
		})
	}
}

func TestPsqlRepo_FindBookByIDHash(t *testing.T) {
	t.Parallel()
	StubPsqlConn()
	db, err := OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "bk_id_hash/find"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
		db.Close()
	})

	bksDB := stubData(NewRepo(site, db), site)

	tests := []struct {
		name         string
		r            repo.Repository
		id           int
		hashcode     int
		expectResult *model.Book
		expectErr    bool
	}{
		{
			name:         "find not existing book",
			r:            NewRepo(site, db),
			id:           0,
			hashcode:     0,
			expectResult: nil,
			expectErr:    true,
		},
		{
			name:         "find book with correct id hash",
			r:            NewRepo(site, db),
			id:           2,
			hashcode:     0,
			expectResult: &bksDB[1],
			expectErr:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			result, err := test.r.FindBookByIdHash(test.id, test.hashcode)
			t.Run("result", func(t *testing.T) {
				if (err != nil) != test.expectErr {
					t.Errorf("got error: %v; want error: %v", err, test.expectErr)
				}

				assert.Equal(t, test.expectResult, result)
			})
		})
	}
}

func TestPsqlRepo_FindBookByStatus(t *testing.T) {
	t.Parallel()
	//Check if this will really be used
}

func TestPsqlRepo_FindAllBooks(t *testing.T) {
	t.Parallel()
	StubPsqlConn()
	db, err := OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "all_bk/find"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
		db.Close()
	})

	bksDB := stubData(NewRepo(site, db), site)

	tests := []struct {
		name         string
		r            repo.Repository
		expectResult []model.Book
		expectErr    bool
	}{
		{
			name:         "works",
			r:            NewRepo(site, db),
			expectResult: bksDB,
			expectErr:    false,
		},
		{
			name:         "works for empty",
			r:            NewRepo("empty", db),
			expectResult: nil,
			expectErr:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			result, err := test.r.FindAllBooks()
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

func TestPsqlRepo_FindBooksForUpdate(t *testing.T) {
	t.Parallel()
	StubPsqlConn()
	db, err := OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "update_bk/find"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
		db.Close()
	})

	bksDB := stubData(NewRepo(site, db), site)

	tests := []struct {
		name         string
		r            repo.Repository
		expectResult []model.Book
		expectErr    bool
	}{
		{
			name:         "works",
			r:            NewRepo(site, db),
			expectResult: []model.Book{bksDB[4], bksDB[3], bksDB[2], bksDB[0]},
			expectErr:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			result, err := test.r.FindBooksForUpdate()
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

func TestPsqlRepo_FindBooksForDownload(t *testing.T) {
	t.Parallel()
	StubPsqlConn()
	db, err := OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "down_bk/find"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
		db.Close()
	})

	bksDB := stubData(NewRepo(site, db), site)

	tests := []struct {
		name         string
		r            repo.Repository
		expectResult []model.Book
		expectErr    bool
	}{
		{
			name:         "works",
			r:            NewRepo(site, db),
			expectResult: []model.Book{bksDB[2]},
			expectErr:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			result, err := test.r.FindBooksForDownload()
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

func TestPsqlRepo_FindBooksByTitleWriter(t *testing.T) {
	t.Parallel()
	//TODO: fill in testcase
	StubPsqlConn()
	db, err := OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "bk_tit_wrt/find"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
		db.Close()
	})

	bksDB := stubData(NewRepo(site, db), site)

	tests := []struct {
		name         string
		r            repo.Repository
		title        string
		writer       string
		limit        int
		offset       int
		expectResult []model.Book
		expectErr    bool
	}{
		{
			name:         "works",
			r:            NewRepo(site, db),
			title:        "title",
			writer:       "writer",
			limit:        10,
			offset:       0,
			expectResult: []model.Book{bksDB[3], bksDB[2], bksDB[1], bksDB[0]},
			expectErr:    false,
		},
		{
			name:         "works with limit",
			r:            NewRepo(site, db),
			title:        "title",
			writer:       "writer",
			limit:        1,
			offset:       0,
			expectResult: []model.Book{bksDB[3]},
			expectErr:    false,
		},
		{
			name:         "works with offset",
			r:            NewRepo(site, db),
			title:        "title",
			writer:       "writer",
			limit:        1,
			offset:       1,
			expectResult: []model.Book{bksDB[2]},
			expectErr:    false,
		},
		{
			name:         "return all books match either title or writer",
			r:            NewRepo(site, db),
			title:        "title 1",
			writer:       "writer 3",
			limit:        5,
			offset:       0,
			expectResult: []model.Book{bksDB[3], bksDB[0]},
			expectErr:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			result, err := test.r.FindBooksByTitleWriter(test.title, test.writer, test.limit, test.offset)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want err: %v", err, test.expectErr)
			}
			assert.Equal(t, test.expectResult, result)
		})
	}
}

func TestPsqlRepo_FindBooksByRandom(t *testing.T) {
	t.Parallel()

	StubPsqlConn()
	db, err := OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "random_bk/find"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
		db.Close()
	})

	stubData(NewRepo(site, db), site)

	tests := []struct {
		name         string
		r            repo.Repository
		limit        int
		expectLength int
		expectErr    bool
	}{
		{
			name:         "works",
			r:            NewRepo(site, db),
			limit:        10,
			expectLength: 2,
			expectErr:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			result, err := test.r.FindBooksByRandom(test.limit)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want err: %v", err, test.expectErr)
			}
			if len(result) != test.expectLength {
				t.Errorf("query got:  %v\nwant length: %v", result, test.expectLength)
			}
		})
	}
}

func TestPsqlRepo_UpdateBooksStatus(t *testing.T) {
	t.Parallel()

	StubPsqlConn()
	db, err := OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "stat_bk/update"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
		db.Close()
	})

	r := NewRepo(site, db)

	bksDB := stubData(NewRepo(site, db), site)

	tests := []struct {
		name       string
		r          repo.Repository
		bkID       int
		bkHash     int
		expectErr  bool
		expectBook *model.Book
	}{
		{
			name:      "works for updating in progress to end",
			r:         NewRepo(site, db),
			bkID:      3,
			bkHash:    0,
			expectErr: false,
			expectBook: &model.Book{
				Site: site, ID: 3, HashCode: 0,
				Title: "title 3", Writer: bksDB[3].Writer,
				Type: "type 3", UpdateDate: "date 3", UpdateChapter: "end " + model.ChapterEndKeywords[0] + " end",
				Status: model.End, IsDownloaded: false,
			},
		},
		{
			name:      "works for updating download to false",
			r:         NewRepo(site, db),
			bkID:      1,
			bkHash:    0,
			expectErr: false,
			expectBook: &model.Book{
				Site: site, ID: 1, HashCode: 0,
				Title: "title 1", Writer: bksDB[0].Writer,
				Type: "type 1", UpdateDate: "date 1", UpdateChapter: "end " + model.ChapterEndKeywords[0] + " end",
				Status: model.End, IsDownloaded: false,
			},
		},
		{
			name:      "works for not updating in progress to end",
			r:         NewRepo(site, db),
			bkID:      2,
			bkHash:    100,
			expectErr: false,
			expectBook: &model.Book{
				Site: site, ID: 2, HashCode: 100,
				Title: "title 2 new", Writer: bksDB[2].Writer, Type: "type 2 new",
				UpdateDate: "date 2.1", UpdateChapter: bksDB[2].UpdateChapter,
				Status: model.InProgress, IsDownloaded: false,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			bksDB[2].Status = model.InProgress
			r.UpdateBook(&bksDB[2])
			bksDB[3].UpdateChapter = fmt.Sprintf("end %v end", model.ChapterEndKeywords[0])
			r.UpdateBook(&bksDB[3])
			bksDB[0].UpdateChapter = fmt.Sprintf("end %v end", model.ChapterEndKeywords[0])
			bksDB[0].Status = model.InProgress
			r.UpdateBook(&bksDB[0])

			err := test.r.UpdateBooksStatus()
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want error: %v", err, test.expectErr)
			}

			bk, err := test.r.FindBookByIdHash(test.bkID, test.bkHash)
			if err != nil {
				t.Errorf("book fail to fetch: id: %v; hash: %v; err: %v", test.bkID, test.bkHash, err)
				return
			}
			assert.Equal(t, test.expectBook, bk)
		})
	}
}

func BenchmarkPsqlRepo_UpdateBooksStatus(b *testing.B) {
	StubPsqlConn()
	db, err := OpenDatabase("")
	if err != nil {
		b.Fatalf("error in open database: %v", err)
	}
	site := "bm/bk_st/update"

	b.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
		db.Close()
	})

	r := NewRepo(site, db)

	stubData(NewRepo(site, db), site)

	for n := 0; n < b.N; n++ {
		r.UpdateBooksStatus()
	}
}

func TestPsqlRepo_SaveWriter(t *testing.T) {
	t.Parallel()
	StubPsqlConn()
	db, err := OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "writer/save"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
		db.Close()
	})

	bksDB := stubData(NewRepo(site, db), site)

	tests := []struct {
		name      string
		r         repo.Repository
		writer    *model.Writer
		expectErr bool
	}{
		{
			name:      "save existing writer",
			r:         NewRepo(site, db),
			writer:    &model.Writer{ID: 0, Name: bksDB[0].Writer.Name},
			expectErr: false,
		},
		{
			name:      "save new writer",
			r:         NewRepo(site, db),
			writer:    &model.Writer{ID: 0, Name: site + " new writer"},
			expectErr: false,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			err := test.r.SaveWriter(test.writer)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v, expect err: %v", err, test.expectErr)
			}

			if test.writer.ID <= 0 {
				t.Errorf("got writer:  %v", test.writer)
			}
		})
	}
}

func TestPsqlRepo_SaveError(t *testing.T) {
	t.Parallel()
	StubPsqlConn()
	db, err := OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "error/save"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
		db.Close()
	})

	stubData(NewRepo(site, db), site)

	tests := []struct {
		name         string
		r            repo.Repository
		bk           *model.Book
		e            error
		expectErrStr string
		expectErr    bool
	}{
		{
			name:         "create error for existing book",
			r:            NewRepo(site, db),
			bk:           &model.Book{Site: site, ID: 1},
			e:            errors.New("create error"),
			expectErrStr: "create error",
			expectErr:    false,
		},
		{
			name:         "update error for existing book",
			r:            NewRepo(site, db),
			bk:           &model.Book{Site: site, ID: 1},
			e:            errors.New("update error"),
			expectErrStr: "update error",
			expectErr:    false,
		},
		{
			name:         "delete error for existing book",
			r:            NewRepo(site, db),
			bk:           &model.Book{Site: site, ID: 1},
			e:            nil,
			expectErrStr: "",
			expectErr:    false,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			err := test.r.SaveError(test.bk, test.e)
			t.Run("result", func(t *testing.T) {
				if (err != nil) != test.expectErr {
					t.Errorf("got error: %v, expect err: %v", err, test.expectErr)
				}
			})

			t.Run("error in db", func(t *testing.T) {
				bk, err := test.r.FindBookByIdHash(test.bk.ID, test.bk.HashCode)
				if err != nil {
					t.Errorf("query got error: %v; want error: %v", err, false)
				}
				if !((bk.Error == nil && test.expectErrStr == "") ||
					(bk.Error != nil && bk.Error.Error() == test.expectErrStr)) {
					t.Errorf("query got:  %v\nwant: %v", bk.Error, test.expectErrStr)
					t.Errorf(cmp.Diff(bk.Error.Error(), test.expectErrStr))
				}
			})
		})
	}
}

func TestPsqlRepo_Backup(t *testing.T) {
	t.Parallel()
	StubPsqlConn()
	db, err := OpenDatabase("")
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

	stubData(NewRepo(site, db), site)

	tests := []struct {
		name      string
		r         repo.Repository
		path      string
		expectErr bool
	}{
		{
			name:      "works",
			r:         NewRepo(site, db),
			path:      "/" + site,
			expectErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			err := test.r.Backup(test.path)
			if (err != nil) != test.expectErr {
				t.Errorf("got err: %v; want err: %v", err, test.expectErr)
			}
		})
	}
}

func TestPsqlRepo_Stats(t *testing.T) {
	t.Parallel()

	StubPsqlConn()
	db, err := OpenDatabase("")
	if err != nil {
		t.Fatalf("error in open database: %v", err)
	}
	site := "stats"

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from writers where id>0 and name like $1", site+"%")
		db.Exec("delete from errors where site=$1", site)
		db.Close()
	})

	stubData(NewRepo(site, db), site)

	tests := []struct {
		name   string
		r      repo.Repository
		expect repo.Summary
	}{
		{
			name: "works",
			r:    NewRepo(site, db),
			expect: repo.Summary{
				BookCount: 5, WriterCount: 5, ErrorCount: 1, DownloadCount: 2,
				UniqueBookCount: 4, MaxBookID: 4,
				LatestSuccessID: 3, StatusCount: map[model.StatusCode]int{
					model.Error:      1,
					model.InProgress: 1,
					model.End:        3,
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result := test.r.Stats()
			assert.Equal(t, test.expect, result)
		})
	}
}
