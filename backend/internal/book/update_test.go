package book

import (
	"testing"
	"time"
	"github.com/htchan/BookSpider/internal/book/model"
	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/mock"
	"github.com/htchan/BookSpider/internal/client"
)

func TestBook_baseURL(t *testing.T) {
	t.Parallel()
	con := config.BookConfig{
		URL: config.URLConfig{ Base: "http://some.book.site/base/%v" },
	}

	book := Book{
		BookModel: model.BookModel{ID:1},
		BookConfig: &con,
	}
	result := book.baseURL()
	if result != "http://some.book.site/base/1" {
		t.Errorf("base url return: %v", result)
	}
}

func TestBook_fetchInfo(t *testing.T) {
	t.Parallel()

	server := mock.MockBookUpdateServer()

	t.Cleanup(func () {
		server.Close()
	})

	con := config.BookConfig{
		SourceKey: "test_book",
	}
	client := client.CircuitBreakerClient{
		CircuitBreakerConfig: config.CircuitBreakerConfig{
			Retry503: 2,
			RetryErr: 2,
			IntervalSleep: 1,
		},
	}

	client.Init(0)

	t.Run("success", func (t *testing.T) {
		tempCon := con
		t.Parallel()

		tempCon.URL.Base = server.URL + "/no_updated/%v"
		book := Book{BookConfig: &tempCon, CircuitBreakerClient: &client}
		title, writer, typeStr, date, chapterStr, err := book.fetchInfo()
		if err != nil || title != "title" || writer != "writer" ||
		typeStr != "type" || date != "date" || chapterStr != "chapter" {
			t.Errorf(
				"fetch info get result: %v, %v, %v, %v, %v, %v",
				title, writer, typeStr, date, chapterStr, err,
			)
		}
	})

	t.Run("zero length", func (t *testing.T) {
		tempCon := con
		t.Parallel()

		tempCon.URL.Base = server.URL + "/zero_length/%v"
		book := Book{BookConfig: &tempCon, CircuitBreakerClient: &client}
		startFetch := time.Now()
		title, writer, typeStr, date, chapterStr, err := book.fetchInfo()
		if err == nil || err.Error() != "zero length" ||
		time.Now().Before(startFetch.Add(3 * time.Second)) {
			t.Errorf(
				"fetch info get result: %v, %v, %v, %v, %v, %v",
				title, writer, typeStr, date, chapterStr, err,
			)
		}
	})

	t.Run("400", func (t *testing.T) {
		tempCon := con
		t.Parallel()

		tempCon.URL.Base = server.URL + "/400/%v"
		book := Book{BookConfig: &tempCon, CircuitBreakerClient: &client}
		startFetch := time.Now()
		title, writer, typeStr, date, chapterStr, err := book.fetchInfo()
		if err == nil || err.Error() != "code 400" ||
		time.Now().Before(startFetch.Add(3 * time.Second)) {
			t.Errorf(
				"fetch info get result: %v, %v, %v, %v, %v, %v",
				title, writer, typeStr, date, chapterStr, err,
			)
		}
	})

	t.Run("503", func (t *testing.T) {
		tempCon := con
		t.Parallel()

		tempCon.URL.Base = server.URL + "/503/%v"
		book := Book{BookConfig: &tempCon, CircuitBreakerClient: &client}
		startFetch := time.Now()
		title, writer, typeStr, date, chapterStr, err := book.fetchInfo()
		if err == nil || err.Error() != "code 503" ||
		time.Now().Before(startFetch.Add(3 * time.Second)) {
			t.Errorf(
				"fetch info get result: %v, %v, %v, %v, %v, %v",
				title, writer, typeStr, date, chapterStr, err,
			)
		}
	})

	t.Run("missing date", func (t *testing.T) {
		tempCon := con
		t.Parallel()

		tempCon.URL.Base = server.URL + "/missing_date/%v"
		book := Book{BookConfig: &tempCon, CircuitBreakerClient: &client}
		title, writer, typeStr, date, chapterStr, err := book.fetchInfo()
		if err == nil || err.Error() != "date not found" {
			t.Errorf(
				"fetch info get result: %v, %v, %v, %v, %v, %v",
				title, writer, typeStr, date, chapterStr, err,
			)
		}
	})
	
}

func TestBook_isNewBook(t *testing.T) {
	t.Parallel()
	book := Book{
		BookModel: model.BookModel{Title: "title", Type: "type"},
		WriterModel: model.WriterModel{Name: "writer"},
	}

	t.Run("title changed", func (t *testing.T) {
		t.Parallel()
		result := book.isNewBook("title_2", book.Name, book.Type)
		if !result {
			t.Errorf("is new book return: %v", result)
		}
	})

	t.Run("writer changed", func (t *testing.T) {
		t.Parallel()
		result := book.isNewBook(book.Title, "writer_2", book.Type)
		if !result {
			t.Errorf("is new book return: %v", result)
		}
	})

	t.Run("type string changed", func (t *testing.T) {
		t.Parallel()
		result := book.isNewBook(book.Title, book.Name, "type_2")
		if !result {
			t.Errorf("is new book return: %v", result)
		}
	})

	t.Run("nothing changed", func (t *testing.T) {
		t.Parallel()
		result := book.isNewBook(book.Title, book.Name, book.Type)
		if result {
			t.Errorf("is new book return: %v", result)
		}
	})
}

func TestBook_isUpdated(t *testing.T) {
	t.Parallel()
	book := Book{
		BookModel: model.BookModel{
			Title: "title", Type: "type",
			UpdateDate: "date", UpdateChapter: "chapter",
		},
		WriterModel: model.WriterModel{Name: "writer"},
	}

	t.Run("title changed", func (t *testing.T) {
		t.Parallel()
		result := book.isUpdated("t", book.Name, book.Type, book.UpdateDate, book.UpdateChapter)
		if !result {
			t.Errorf("is new book return: %v", result)
		}
	})

	t.Run("writer changed", func (t *testing.T) {
		t.Parallel()
		result := book.isUpdated(book.Title, "w", book.Type, book.UpdateDate, book.UpdateChapter)
		if !result {
			t.Errorf("is new book return: %v", result)
		}
	})

	t.Run("type string changed", func (t *testing.T) {
		t.Parallel()
		result := book.isUpdated(book.Title, book.Name, "t", book.UpdateDate, book.UpdateChapter)
		if !result {
			t.Errorf("is new book return: %v", result)
		}
	})

	t.Run("udpate date changed", func (t *testing.T) {
		t.Parallel()
		result := book.isUpdated(book.Title, book.Name, book.Type, "d", book.UpdateChapter)
		if !result {
			t.Errorf("is new book return: %v", result)
		}
	})

	t.Run("update chapter changed", func (t *testing.T) {
		t.Parallel()
		result := book.isUpdated(book.Title, book.Name, book.Type, book.UpdateDate, "c")
		if !result {
			t.Errorf("is new book return: %v", result)
		}
	})

	t.Run("nothing changed", func (t *testing.T) {
		t.Parallel()
		result := book.isUpdated(book.Title, book.Name, book.Type, book.UpdateDate, book.UpdateChapter)
		if result {
			t.Errorf("is new book return: %v", result)
		}
	})
}

func TestBook_Update(t *testing.T) {
	t.Parallel()

	server := mock.MockBookUpdateServer()

	t.Cleanup(func () {
		server.Close()
	})

	con := config.BookConfig{
		SourceKey: "test_book",
	}
	client := client.CircuitBreakerClient{}
	client.Init(0)

	t.Run("fetch not updated book", func (t *testing.T) {
		tempCon := con
		t.Parallel()

		tempCon.URL.Base = server.URL + "/updated/%v"
		book := Book{BookConfig: &tempCon, CircuitBreakerClient: &client}
		result, err := book.Update()
		if !result || err != nil ||
		book.Title != "title_new" || book.Name != "writer_new" ||
		book.Type != "type_new" || book.UpdateDate != "date_new" ||
		book.UpdateChapter != "chapter_new" || book.Status != model.InProgress {
			t.Errorf("book model: %v", book.BookModel)
			t.Errorf("writer model: %v", book.WriterModel)
			t.Errorf("error model: %v", book.ErrorModel)
			t.Errorf("update return result: %v, err: %v", result, err)
		}
	})

	t.Run("fetch updated book", func (t *testing.T) {
		tempCon := con
		t.Parallel()

		tempCon.URL.Base = server.URL + "/no_updated/%v"
		book := Book{
			BookModel: model.BookModel{
				Title: "title", Type: "type", UpdateDate: "date",
				UpdateChapter: "chapter", Status: model.End,
			},
			WriterModel: model.WriterModel{Name: "writer"},
			BookConfig: &tempCon,
			CircuitBreakerClient: &client,
		}
		result, err := book.Update()
		if result || err != nil || book.Title != "title" || book.Name != "writer" ||
		book.Type != "type" || book.UpdateDate != "date" ||
		book.UpdateChapter != "chapter" || book.Status != model.End {
			t.Errorf("book model: %v", book.BookModel)
			t.Errorf("writer model: %v", book.WriterModel)
			t.Errorf("error model: %v", book.ErrorModel)
			t.Errorf("update return result: %v, err: %v", result, err)
		}
	})

	t.Run("fetch error", func (t *testing.T) {
		tempCon := con
		t.Parallel()

		tempCon.URL.Base = server.URL + "/zero_length/%v"
		book := Book{BookConfig: &tempCon, CircuitBreakerClient: &client}
		result, err := book.Update()
		if err == nil || err.Error() != "zero length" || book.Title != "" || book.Name != "" ||
		book.Type != "" || book.UpdateDate != "" || book.UpdateChapter != "" {
			t.Errorf("book model: %v", book.BookModel)
			t.Errorf("writer model: %v", book.WriterModel)
			t.Errorf("error model: %v", book.ErrorModel)
			t.Errorf("update return result: %v, err: %v", result, err)
		}
	})
}