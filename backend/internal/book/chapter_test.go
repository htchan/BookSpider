package book

import (
	"github.com/google/go-cmp/cmp"
	"github.com/htchan/BookSpider/internal/book/model"
	"github.com/htchan/BookSpider/internal/client"
	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/mock"
	"testing"
)

func TestChapter_NewChapter(t *testing.T) {
	t.Parallel()

	bk := Book{}
	result := NewChapter(1, "url", "title", &bk)
	if result.Index != 1 || result.URL != "url" || result.Title != "title" ||
		result.Content != "" || !cmp.Equal(result.Book, &bk) {
		t.Errorf("new chapter return: %v", result)
	}
}

func TestChapter_chapterURL(t *testing.T) {
	t.Parallel()

	bk := Book{
		BookModel: model.BookModel{ID: 1},
		BookConfig: &config.BookConfig{
			URLConfig: config.URLConfig{ChapterPrefix: "chapter", Download: "download/%v/"},
		},
	}

	t.Run("append to Chapter URL if url is start with /", func(t *testing.T) {
		t.Parallel()
		chapter := Chapter{Book: &bk, URL: "/hello"}
		result := chapter.chapterURL()
		if result != "chapter/hello" {
			t.Errorf("chapter url return: %v", result)
		}
	})

	t.Run("append to Chapter URL if url is start with http", func(t *testing.T) {
		t.Parallel()
		chapter := Chapter{Book: &bk, URL: "http://hello"}
		result := chapter.chapterURL()
		if result != "chapterhttp://hello" {
			t.Errorf("chapter url return: %v", result)
		}
	})

	t.Run("append to downlaod URL if url is not start with / or http", func(t *testing.T) {
		t.Parallel()
		chapter := Chapter{Book: &bk, URL: "hello"}
		result := chapter.chapterURL()
		if result != "download/1/hello" {
			t.Errorf("chapter url return: %v", result)
		}
	})
}

func TestChapter_generateIndex(t *testing.T) {
	t.Skip("future not available now")
}

func TestChapter_optimizeContent(t *testing.T) {
	t.Parallel()

	t.Run("remove all unwanted content", func(t *testing.T) {
		t.Parallel()
		chapter := Chapter{Content: "&nbsp;<b></b><p></p>                "}
		chapter.optimizeContent()
		if chapter.Content != "" {
			t.Errorf("optimize content update content to %v", chapter.Content)
		}
	})

	t.Run("replace unwanted content to wanted", func(t *testing.T) {
		t.Parallel()
		chapter := Chapter{Content: "<p/><br />"}
		chapter.optimizeContent()
		if chapter.Content != "\n\n" {
			t.Errorf("optimize content update content to %v", chapter.Content)
		}
	})
}

func TestChapter_sortChapters(t *testing.T) {
	t.Parallel()

	chapters := []Chapter{
		Chapter{Index: 5}, Chapter{Index: 4}, Chapter{Index: 3},
		Chapter{Index: 2}, Chapter{Index: 1}, Chapter{Index: 0},
	}
	sortChapters(chapters)
	for i, chapter := range chapters {
		if chapter.Index != i {
			t.Errorf("index of chapter at position %v is %v", i, chapter.Index)
		}
	}
}

func TestChapter_optimizeChapters(t *testing.T) {
	t.Skip("future not available now")
}

func TestChapter_Fetch(t *testing.T) {
	t.Parallel()

	server := mock.MockBookDownloadServer()

	t.Cleanup(func() {
		server.Close()
	})

	client := client.CircuitBreakerClient{}
	client.Init(0)

	con := config.BookConfig{
		URLConfig: config.URLConfig{ChapterPrefix: server.URL},
		SourceKey: "test_book",
	}

	t.Run("success to fetch and parse data", func(t *testing.T) {
		t.Parallel()
		chapter := Chapter{
			Book: &Book{CircuitBreakerClient: &client, BookConfig: &con},
			URL:  "/chapter/success",
		}

		chapter.Fetch()
		if chapter.Content != "success" {
			t.Errorf("fetch update content to %v", chapter.Content)
		}
	})

	t.Run("fail to fetch data", func(t *testing.T) {
		t.Parallel()
		chapter := Chapter{
			Book: &Book{CircuitBreakerClient: &client, BookConfig: &con},
			URL:  "/chapter/400",
		}

		chapter.Fetch()
		if chapter.Content != "load html failed - code 400" {
			t.Errorf("fetch update content to %v", chapter.Content)
		}
	})

	t.Run("fail to parse data", func(t *testing.T) {
		t.Parallel()
		chapter := Chapter{
			Book: &Book{CircuitBreakerClient: &client, BookConfig: &con},
			URL:  "/chapter/unknown",
		}

		chapter.Fetch()
		if chapter.Content != "recognize html fail\nunknown content" {
			t.Errorf("fetch update content to %v", chapter.Content)
		}
	})
}

func TestChapter_content(t *testing.T) {
	t.Parallel()

	chapter := Chapter{
		Title:   "title",
		Content: "content",
	}

	if chapter.content() != "title\n--------------------\ncontent\n--------------------\n" {
		t.Errorf("content return %v", chapter.content())
	}
}
