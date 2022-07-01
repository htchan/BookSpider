package book

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/htchan/BookSpider/internal/book/model"
	"github.com/htchan/BookSpider/internal/client"
	"github.com/htchan/BookSpider/internal/config"
)

func equalBook(result, expect Book) bool {
	return cmp.Equal(result.BookModel, expect.BookModel) &&
		cmp.Equal(result.WriterModel, expect.WriterModel) &&
		cmp.Equal(result.ErrorModel.Error(), expect.ErrorModel.Error()) &&
		cmp.Equal(result.BookConfig, expect.BookConfig)
}

func diffBook(result, expect Book) string {
	return fmt.Sprintf(
		`new book book model diff: %v
		new book writer model diff: %v
		new book error model diff: %v
		new book config diff: %v`,
		cmp.Diff(result.BookModel, expect.BookModel),
		cmp.Diff(result.WriterModel, expect.WriterModel),
		cmp.Diff(result.ErrorModel.Error(), expect.ErrorModel.Error()),
		cmp.Diff(result.BookConfig, expect.BookConfig),
	)
}

func TestBook_NewBook(t *testing.T) {
	t.Parallel()

	con := config.BookConfig{SourceKey: "new_bk"}
	client := client.CircuitBreakerClient{}
	t.Run("create new book instance", func(t *testing.T) {
		t.Parallel()
		expect := Book{
			BookModel: model.BookModel{
				Site: "site", ID: 1,
				Title: "", WriterID: 0, Type: "",
				UpdateDate: "", UpdateChapter: "", Status: model.Error,
			},
			WriterModel:          model.WriterModel{ID: 0, Name: ""},
			ErrorModel:           model.ErrorModel{Site: "site", ID: 1, Err: errors.New("new book")},
			BookConfig:           &con,
			CircuitBreakerClient: &client,
		}
		result := NewBook("site", 1, &con, &client)
		expect.BookModel.HashCode = result.BookModel.HashCode
		if result.BookModel.HashCode == 0 || !equalBook(result, expect) {
			t.Error(diffBook(result, expect))
		}
	})
}

func TestBook_LoadBook(t *testing.T) {
	t.Parallel()

	site, writer := "load_bk", "load_bk_writer"
	con := config.BookConfig{}
	client := client.CircuitBreakerClient{}

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from errors where site=$1", site)
		db.Exec("delete from writers where name=$1", writer)
		runtime.GC()
	})

	t.Run("load download book", func(t *testing.T) {
		t.Parallel()
		expect := Book{
			BookModel: model.BookModel{
				Site: site, ID: 1, HashCode: 0,
				Title: "title", WriterID: 0, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
				Status: model.Download,
			},
			WriterModel: model.WriterModel{ID: 0, Name: writer},
			BookConfig:  &con,
		}
		err := expect.Save(db)
		if err != nil {
			t.Errorf("save book return: %v", err)
			return
		}
		result, err := LoadBook(db, site, 1, 0, &con, &client)
		if err != nil || !equalBook(result, expect) {
			t.Error(diffBook(result, expect))
		}
	})

	t.Run("load end book", func(t *testing.T) {
		t.Parallel()
		expect := Book{
			BookModel: model.BookModel{
				Site: site, ID: 2, HashCode: 0,
				Title: "title", WriterID: 0, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
				Status: model.End,
			},
			WriterModel: model.WriterModel{ID: 0, Name: writer},
			BookConfig:  &con,
		}
		err := expect.Save(db)
		if err != nil {
			t.Errorf("save book return: %v", err)
			return
		}
		result, err := LoadBook(db, site, 2, 0, &con, &client)
		if err != nil || !equalBook(result, expect) {
			t.Error(diffBook(result, expect))
		}
	})

	t.Run("load in_progress book", func(t *testing.T) {
		t.Parallel()
		expect := Book{
			BookModel: model.BookModel{
				Site: site, ID: 3, HashCode: 0,
				Title: "title", WriterID: 0, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
				Status: model.InProgress,
			},
			WriterModel: model.WriterModel{ID: 0, Name: writer},
			BookConfig:  &con,
		}
		err := expect.Save(db)
		if err != nil {
			t.Errorf("save book return: %v", err)
			return
		}
		result, err := LoadBook(db, site, 3, 0, &con, &client)
		if err != nil || !equalBook(result, expect) {
			t.Error(diffBook(result, expect))
		}
	})

	t.Run("load error book", func(t *testing.T) {
		t.Parallel()
		expect := Book{
			BookModel: model.BookModel{
				Site: site, ID: 4, HashCode: 0,
				Status: model.Download,
			},
			ErrorModel: model.ErrorModel{Site: site, ID: 4, Err: errors.New("data")},
			BookConfig: &con,
		}
		err := expect.Save(db)
		if err != nil {
			t.Errorf("save book return: %v", err)
			return
		}
		result, err := LoadBook(db, site, 4, 0, &con, &client)
		if err != nil || !equalBook(result, expect) {
			t.Error(diffBook(result, expect))
		}
	})
}

func TestBook_Save(t *testing.T) {
	t.Parallel()

	site, writer := "save_bk", "save_bk_writer"
	con := config.BookConfig{}

	t.Cleanup(func() {
		db.Exec("delete from books where site=$1", site)
		db.Exec("delete from errors where site=$1", site)
		db.Exec("delete from writers where name=$1", writer)
		runtime.GC()
	})

	t.Run("create in progress book", func(t *testing.T) {
		t.Parallel()
		expect := Book{
			BookModel: model.BookModel{
				Site: site, ID: 1, HashCode: 0,
				Status: model.InProgress,
			},
			WriterModel: model.WriterModel{Name: writer},
			BookConfig:  &con,
		}
		err := expect.Save(db)
		if err != nil || expect.BookModel.WriterID == 0 ||
			expect.WriterModel.ID != expect.BookModel.WriterID {
			t.Errorf("save book return: %v", err)
			t.Errorf("save book book model update: %v", expect.BookModel)
			t.Errorf("save book writer model update: %v", expect.WriterModel)
			t.Errorf("save book error model update: %v", expect.ErrorModel)
			return
		}
	})

	t.Run("create error book", func(t *testing.T) {
		t.Parallel()
		expect := Book{
			BookModel: model.BookModel{
				Site: site, ID: 2, HashCode: 0,
			},
			ErrorModel: model.ErrorModel{Site: site, ID: 2, Err: errors.New("save book")},
			BookConfig: &con,
		}
		err := expect.Save(db)
		if err != nil || expect.BookModel.WriterID != 0 ||
			expect.WriterModel.ID != 0 {
			t.Errorf("save book return: %v", err)
			t.Errorf("save book book model update: %v", expect.BookModel)
			t.Errorf("save book writer model update: %v", expect.WriterModel)
			t.Errorf("save book error model update: %v", expect.ErrorModel)
			return
		}
	})

	t.Run("update error book to inprogress", func(t *testing.T) {
		t.Parallel()
		expect := Book{
			BookModel: model.BookModel{
				Site: site, ID: 3, HashCode: 0,
			},
			ErrorModel: model.ErrorModel{Site: site, ID: 2, Err: errors.New("save book")},
			BookConfig: &con,
		}
		expect.Save(db)
		expect.ErrorModel.Err = nil
		expect.WriterModel = model.WriterModel{Name: writer}
		err := expect.Save(db)
		if err != nil || expect.BookModel.WriterID == 0 ||
			expect.WriterModel.ID != expect.BookModel.WriterID {
			t.Errorf("save book return: %v", err)
			t.Errorf("save book book model update: %v", expect.BookModel)
			t.Errorf("save book writer model update: %v", expect.WriterModel)
			t.Errorf("save book error model update: %v", expect.ErrorModel)
			return
		}
	})
}

func TestBook_MarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("render downloaded book", func(t *testing.T) {
		t.Parallel()
		book := Book{
			BookModel: model.BookModel{
				Site: "bk_json", ID: 1, HashCode: 100,
				Title: "title", WriterID: 1, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
				Status: model.Download,
			},
			WriterModel: model.WriterModel{ID: 1, Name: "writer"},
			ErrorModel: model.ErrorModel{
				Site: "bk_json", ID: 1, Err: errors.New("error"),
			},
		}
		expect := `{
			"site":"bk_json","id":1,"hash":"2s",
			"title":"title","writer":"writer","type":"type",
			"update_date":"date","update_chapter":"chapter","status":"download"
		}`
		expect = strings.ReplaceAll(strings.ReplaceAll(expect, "\n", ""), "\t", "")
		result, err := json.Marshal(book)
		if err != nil || !cmp.Equal(string(result), expect) {
			t.Errorf("render book json return: error: %v, result: %v", err, string(result))
		}
	})

	t.Run("render end book", func(t *testing.T) {
		t.Parallel()
		book := Book{
			BookModel: model.BookModel{
				Site: "bk_json", ID: 1, HashCode: 100,
				Title: "title", WriterID: 1, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
				Status: model.End,
			},
			WriterModel: model.WriterModel{ID: 1, Name: "writer"},
			ErrorModel: model.ErrorModel{
				Site: "bk_json", ID: 1, Err: errors.New("error"),
			},
		}
		expect := `{
			"site":"bk_json","id":1,"hash":"2s",
			"title":"title","writer":"writer","type":"type",
			"update_date":"date","update_chapter":"chapter","status":"end"
		}`
		expect = strings.ReplaceAll(strings.ReplaceAll(expect, "\n", ""), "\t", "")
		result, err := json.Marshal(book)
		if err != nil || !cmp.Equal(string(result), expect) {
			t.Errorf("render book json return: error: %v, result: %v", err, string(result))
		}
	})

	t.Run("render in_progress book", func(t *testing.T) {
		t.Parallel()
		book := Book{
			BookModel: model.BookModel{
				Site: "bk_json", ID: 1, HashCode: 100,
				Title: "title", WriterID: 1, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
				Status: model.InProgress,
			},
			WriterModel: model.WriterModel{ID: 1, Name: "writer"},
			ErrorModel: model.ErrorModel{
				Site: "bk_json", ID: 1, Err: errors.New("error"),
			},
		}
		expect := `{
			"site":"bk_json","id":1,"hash":"2s",
			"title":"title","writer":"writer","type":"type",
			"update_date":"date","update_chapter":"chapter","status":"in_progress"
		}`
		expect = strings.ReplaceAll(strings.ReplaceAll(expect, "\n", ""), "\t", "")
		result, err := json.Marshal(book)
		if err != nil || !cmp.Equal(string(result), expect) {
			t.Errorf("render book json return: error: %v, result: %v", err, string(result))
		}
	})

	t.Run("render error book", func(t *testing.T) {
		t.Parallel()
		book := Book{
			BookModel: model.BookModel{
				Site: "bk_json", ID: 1, HashCode: 100, Status: model.Error,
			},
			ErrorModel: model.ErrorModel{
				Site: "bk_json", ID: 1, Err: errors.New("error"),
			},
		}
		expect := `{
			"site":"bk_json","id":1,"hash":"2s",
			"title":"","writer":"","type":"",
			"update_date":"","update_chapter":"","status":"error"
		}`
		expect = strings.ReplaceAll(strings.ReplaceAll(expect, "\n", ""), "\t", "")
		result, err := json.Marshal(book)
		if err != nil || !cmp.Equal(string(result), expect) {
			t.Errorf("render book json return: error: %v, result: %v", err, string(result))
		}
	})
}

func TestBook_String(t *testing.T) {
	t.Parallel()

	t.Run("work", func(t *testing.T) {
		t.Parallel()
		book := Book{
			BookModel: model.BookModel{
				Site: "bk_str", ID: 1, HashCode: 100, Status: model.Error,
			},
		}
		if book.String() != "bk.bk_str.1.2s" {
			t.Errorf("book string return: %v", book.String())
		}
	})
}
