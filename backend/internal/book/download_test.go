package book

import (
	"fmt"
	"github.com/htchan/BookSpider/internal/book/model"
	"github.com/htchan/BookSpider/internal/client"
	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/mock"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestBook_downloadURL(t *testing.T) {
	t.Parallel()
	con := config.BookConfig{
		URLConfig: config.URLConfig{Download: "http://some.book.site/download/%v"},
	}

	book := Book{
		BookModel:  model.BookModel{ID: 1},
		BookConfig: &con,
	}
	result := book.downloadURL()
	if result != "http://some.book.site/download/1" {
		t.Errorf("download url return %v", result)
	}
}

func TestBook_fetchChapters(t *testing.T) {
	t.Parallel()

	server := mock.MockBookDownloadServer()
	t.Cleanup(func() {
		server.Close()
	})

	client := client.CircuitBreakerClient{}
	client.Init(10)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		con := config.BookConfig{
			URLConfig: config.URLConfig{
				Download:      server.URL + "/list_success_chapters/%v",
				ChapterPrefix: server.URL + "/chapter/success",
			},
			SourceKey: "test_book",
		}

		book := Book{BookConfig: &con, CircuitBreakerClient: &client}
		chapters := book.fetchChapters()

		if len(chapters) != 3 {
			t.Errorf("fetch chapter return chapters: %v", chapters)
		}

		for i, chapter := range chapters {
			if chapter.Content != fmt.Sprintf("success/%v", i) {
				t.Errorf("chapter at position %v content is %v", i, chapter.Content)
			}
		}
	})

	t.Run("failed at download page", func(t *testing.T) {
		t.Parallel()

		con := config.BookConfig{
			URLConfig: config.URLConfig{Download: server.URL + "/400/%v"},
			SourceKey: "test_book",
		}

		book := Book{BookConfig: &con, CircuitBreakerClient: &client}
		chapters := book.fetchChapters()

		if len(chapters) != 0 {
			t.Errorf("fetch chapter return chapters: %v", chapters)
		}
	})
}

func TestBook_generateEmptyChapters(t *testing.T) {
	t.Parallel()
	server := mock.MockBookDownloadServer()
	t.Cleanup(func() {
		server.Close()
	})

	client := client.CircuitBreakerClient{}
	client.Init(0)

	t.Run("success to fetch download page", func(t *testing.T) {
		t.Parallel()
		con := config.BookConfig{
			URLConfig: config.URLConfig{Download: server.URL + "/list_chapters/%v"},
			SourceKey: "test_book",
		}
		book := Book{
			BookConfig:           &con,
			CircuitBreakerClient: &client,
		}

		chapters, err := book.generateEmptyChapters()
		if err != nil || len(chapters) != 3 {
			t.Errorf("generateEmptyChapters returns chapter: %v, err: %v", chapters, err)
		}
		for i, chapter := range chapters {
			if chapter.Index != i || chapter.Title != strconv.Itoa(i) || chapter.URL != strconv.Itoa(i) {
				t.Errorf(
					"chapter at position %v has index: %v, title: %v, url: %v",
					i, chapter.Index, chapter.Title, chapter.URL,
				)
			}
		}
	})

	t.Run("fail to fetch download page", func(t *testing.T) {
		t.Parallel()
		con := config.BookConfig{
			URLConfig: config.URLConfig{Download: server.URL + "/400/%v"},
			SourceKey: "test_book",
		}
		book := Book{
			BookConfig:           &con,
			CircuitBreakerClient: &client,
		}

		chapters, err := book.generateEmptyChapters()
		if err == nil || err.Error() != "code 400" || len(chapters) != 0 {
			t.Errorf("generateEmptyChapters returns chapter: %v, err: %v", chapters, err)
		}
	})

	t.Run("fail to parse download page", func(t *testing.T) {
		t.Parallel()
		con := config.BookConfig{
			URLConfig: config.URLConfig{Download: server.URL + "/unknown/%v"},
			SourceKey: "test_book",
		}
		book := Book{
			BookConfig:           &con,
			CircuitBreakerClient: &client,
		}

		chapters, err := book.generateEmptyChapters()
		if err == nil || err.Error() != "empty chapters" || len(chapters) != 0 {
			t.Errorf("generateEmptyChapters returns chapter: %v, err: %v", chapters, err)
		}
	})
}

func TestBook_content(t *testing.T) {
	t.Parallel()

	book := Book{
		BookModel:   model.BookModel{Title: "title"},
		WriterModel: model.WriterModel{Name: "writer"},
	}
	if book.content() != "title\nwriter\n--------------------\n\n" {
		t.Errorf("content return %v", book.content())
	}
}

func TestBook_saveChapters(t *testing.T) {
	t.Parallel()

	con := config.BookConfig{Storage: "./"}

	t.Cleanup(func() {
		os.Remove("5.txt")
		os.Remove("6.txt")
		os.Remove("7.txt")
		os.Remove("8.txt")
	})

	t.Run("success to save file", func(t *testing.T) {
		t.Parallel()
		book := Book{
			BookConfig:  &con,
			BookModel:   model.BookModel{ID: 5, Title: "title"},
			WriterModel: model.WriterModel{Name: "writer"},
		}
		chapters := []Chapter{Chapter{Title: "chapter", Content: "content"}}
		err := book.saveChapters(chapters)
		if err != nil {
			t.Errorf("save chapter return %v", err)
		}
		result, err := os.ReadFile(book.location())
		if err != nil || string(result) != book.content()+chapters[0].content() {
			t.Errorf("read file return error: %v", err)
			t.Errorf("saved content is %v", string(result))
		}
	})

	t.Run("fail to create file", func(t *testing.T) {
		t.Parallel()
		tempCon := config.BookConfig{Storage: "./unknown_dir/"}
		book := Book{BookConfig: &tempCon}
		err := book.saveChapters(nil)
		if err == nil {
			t.Errorf("save chapter return %v", err)
		}
	})

	t.Run("create the existing file", func(t *testing.T) {
		t.Parallel()
		book := Book{
			BookConfig:  &con,
			BookModel:   model.BookModel{ID: 6, Title: "title"},
			WriterModel: model.WriterModel{Name: "writer"},
		}
		chapters := []Chapter{Chapter{Title: "chapter", Content: "content"}}
		book.saveChapters(chapters)
		err := book.saveChapters(nil)
		if err != nil {
			t.Errorf("save chapter return %v", err)
		}
		result, err := os.ReadFile(book.location())
		if err != nil || string(result) != book.content() {
			t.Errorf("read file return error: %v", err)
			t.Errorf("saved content is %v", string(result))
		}
	})

	t.Run("write empty chapters", func(t *testing.T) {
		t.Parallel()
		book := Book{
			BookConfig:  &con,
			BookModel:   model.BookModel{ID: 7, Title: "title"},
			WriterModel: model.WriterModel{Name: "writer"},
		}
		err := book.saveChapters(nil)
		if err != nil {
			t.Errorf("save chapter return %v", err)
		}
		result, err := os.ReadFile(book.location())
		if err != nil || string(result) != book.content() {
			t.Errorf("read file return error: %v", err)
			t.Errorf("saved content is %v", string(result))
		}
	})

	t.Run("too many failed chapters", func(t *testing.T) {
		t.Parallel()
		book := Book{
			BookConfig:  &con,
			BookModel:   model.BookModel{ID: 8, Title: "title"},
			WriterModel: model.WriterModel{Name: "writer"},
		}
		chapters := []Chapter{Chapter{Title: "1", Content: "failed"}}
		err := book.saveChapters(chapters)
		if err == nil || err.Error() != "too many fail chapter" {
			t.Errorf("save chapter return %v", err)
		}
		result, err := os.ReadFile(book.location())
		if err != nil || string(result) != book.content()+chapters[0].content() {
			t.Errorf("read file return error: %v", err)
			t.Errorf("saved content is %v", string(result))
		}
	})
}

func TestBook_Download(t *testing.T) {
	t.Parallel()

	server := mock.MockBookDownloadServer()
	t.Cleanup(func() {
		server.Close()
		os.Remove("10.txt")
	})

	con := config.BookConfig{
		URLConfig: config.URLConfig{
			Download:      server.URL + "/list_success_chapters/%v",
			ChapterPrefix: server.URL + "/chapter/success",
		},
		SourceKey: "test_book",
	}

	client := client.CircuitBreakerClient{}
	client.Init(10)

	book := Book{
		BookModel:            model.BookModel{ID: 10, Title: "title"},
		WriterModel:          model.WriterModel{Name: "writer"},
		BookConfig:           &con,
		CircuitBreakerClient: &client,
	}

	err := book.Download()

	if err != nil {
		t.Errorf("download return %v", err)
	}

	expected := `0
	--------------------
	success/0
	--------------------
	1
	--------------------
	success/1
	--------------------
	2
	--------------------
	success/2
	--------------------
	`
	expected = strings.ReplaceAll(expected, "\t", "")

	content, err := os.ReadFile("10.txt")
	if err != nil || string(content) != book.content()+expected {
		t.Errorf("read file return content: %v, err: %v", string(content), err)
	}
}

func TestBook_location(t *testing.T) {
	t.Parallel()
	con := config.BookConfig{Storage: "/test"}

	t.Run("hash code is 0", func(t *testing.T) {
		book := Book{
			BookModel:  model.BookModel{Site: "test", ID: 1, HashCode: 0},
			BookConfig: &con,
		}
		result := book.location()
		if result != "/test/1.txt" {
			t.Errorf("location return %v", result)
		}
	})

	t.Run("hash code is 10", func(t *testing.T) {
		book := Book{
			BookModel:  model.BookModel{Site: "test", ID: 1, HashCode: 10},
			BookConfig: &con,
		}
		result := book.location()
		if result != "/test/1-a.txt" {
			t.Errorf("location return %v", result)
		}
	})
}

func TestBook_Content(t *testing.T) {
	t.Parallel()
	filename := "./1.txt"
	os.WriteFile(filename, []byte("content"), 0644)

	t.Cleanup(func() {
		os.Remove(filename)
	})

	con := config.BookConfig{Storage: "."}

	t.Run("return existing content", func(t *testing.T) {
		book := Book{
			BookModel:  model.BookModel{ID: 1, HashCode: 0},
			BookConfig: &con,
		}
		result := book.Content()
		if result != "content" {
			t.Errorf("content return %v", result)
		}
	})

	t.Run("return non existant content", func(t *testing.T) {
		book := Book{
			BookModel:  model.BookModel{ID: 1, HashCode: 10},
			BookConfig: &con,
		}
		result := book.Content()
		if result != "" {
			t.Errorf("content return %v", result)
		}
	})
}

func TestBook_LoadChapters(t *testing.T) {
	t.Parallel()
	t.Skip()
	//TODO: load chapter from content
}
