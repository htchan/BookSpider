package books

import (
	"testing"
	"github.com/htchan/BookSpider/internal/mock"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/pkg/configs"

	"sync"
	"os"
	"strconv"
	"fmt"
)

var downloadConfig = configs.LoadSourceConfigs(os.Getenv("ASSETS_LOCATION") + "/test-data/configs")["test_source_key"]

func TestBooks_Book_Download(t *testing.T) {
	downloadConfig.ConstSleep = 0
	downloadConfig.SourceKey = "test_source_key"
	server := mock.DownloadServer()
	defer server.Close()

	t.Run("func getEmptyChapters", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			book := NewBook("test", 1, -1, downloadConfig)
			book.config.DownloadUrl = server.URL + "/content/success/1"

			chapters, err := book.getEmptyChapters()
			if err != nil || len(chapters) != 4 {
				t.Fatalf("book getEmptyChapters return err: %v, chapters: %v", err, chapters)
			}

			for i, chapter := range chapters {
				s := strconv.Itoa(i + 1)
				if chapter.Url != fmt.Sprintf(book.config.ChapterUrl, "/" + s) || chapter.Title != s {
					t.Fatalf("book getEmptyChapters return wrong result at position %v, chatper: %v", i, chapter)
				}
			}
		})

		t.Run("fail when response is empty", func(t *testing.T) {
			t.Parallel()
			book := NewBook("test", 1, -1, downloadConfig)
			book.config.DownloadUrl = server.URL + "/content/empty"
			book.config.ConstSleep = 0

			_, err := book.getEmptyChapters()
			if err == nil {
				t.Fatalf("book getEmptyChapters return nil error for invalid response")
			}
		})

		t.Run("fail when no url found", func(t *testing.T) {
			t.Parallel()
			book := NewBook("test", 1, -1, downloadConfig)
			book.config.DownloadUrl = server.URL + "/content/no_url"

			_, err := book.getEmptyChapters()
			if err == nil {
				t.Fatalf("book getEmptyChapters return nil error for invalid response")
			}
		})
	})

	t.Run("func downloadChapters", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			book := NewBook("test", 1, -1, downloadConfig)
			book.config.ConstSleep = 0
			book.config.DownloadUrl = server.URL + "/chapter/success/"
			urls := []string { "1", "2", "3" }
			titles := []string { "title-1", "title-2", "title-3" }
			chapters := []Chapter{
				NewChapter(1, urls[0], titles[0], &book.config),
				NewChapter(2, urls[1], titles[1], &book.config),
				NewChapter(3, urls[2], titles[2], &book.config),
			}

			result := book.downloadChapters(chapters, 1)

			if len(result) != 3 {
				t.Fatalf("book download Chapters has wrong chapter count: %v", len(result))
			}

			for i, chapter := range result {
				if chapter.Url != book.config.DownloadUrl + urls[i] || 
					chapter.Title != titles[i] ||
					chapter.Content != strconv.Itoa(i + 1) {
						t.Fatalf("result at position %v incorrect: %v", chapter.Index, chapter)
				}
			} 
		})
	})
	
	t.Run("func saveContent", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			book := NewBook("test", 1, 1, downloadConfig)
			book.config.ConstSleep = 0
			book.SetTitle("book-title")
			book.SetWriter("book-writer")

			chapters := []Chapter{
				Chapter{ Index: 2, Title: "title-2", Url: "url-2", Content: "content-2" },
				Chapter{ Index: 1, Title: "title-1", Url: "url-1", Content: "content-1" },
				Chapter{ Index: 3, Title: "title-3", Url: "url-3", Content: "content-3" },
			}
			book.saveContent("/test-data/storage", chapters)

			b, err := os.ReadFile(os.Getenv("ASSETS_LOCATION") + "/test-data/storage/1-v1.txt")
			utils.CheckError(err)
			reference, err := os.ReadFile(os.Getenv("ASSETS_LOCATION") + "/test-data/storage/1-v1-reference.txt")
			utils.CheckError(err)

			if string(b) != string(reference){
				t.Fatalf("book saveContent save such content: %v", string(b))
			}
		})
	})
	
	t.Run("func Download", func(t *testing.T) {
		var mutex sync.Mutex

		t.Run("success", func(t *testing.T) {
			book := NewBook("test", 1, 10, downloadConfig)
			book.config.DownloadUrl = server.URL + "/content/success/1"
			book.config.ChapterUrl = server.URL + "/chapter/success/%v"
			book.config.ConstSleep = 0
			book.SetTitle("book-title")
			book.SetWriter("book-writer")
			result := book.Download("/test-data/storage", 10, &mutex)

			if !result {
				t.Fatalf("book download failed")
			}

			b, err := os.ReadFile(os.Getenv("ASSETS_LOCATION") + "/test-data/storage/1-v10.txt")
			utils.CheckError(err)
			reference, err := os.ReadFile(os.Getenv("ASSETS_LOCATION") + "/test-data/storage/1-v10-reference.txt")
			utils.CheckError(err)

			if string(b) != string(reference){
				t.Fatalf("book saveContent save such content: %v", string(b))
			}
		})

		t.Run("fail if too much chapters return error", func(t *testing.T) {
			book := NewBook("test", 1, 10, downloadConfig)
			book.config.DownloadUrl = server.URL + "/content/success/1"
			book.config.ChapterUrl = server.URL + "/chapter/invalid/%v"
			book.config.ConstSleep = 0
			book.SetTitle("book-title")
			book.SetWriter("book-writer")
			result := book.Download("/test-data/storage", 10, &mutex)

			if result {
				t.Fatalf("book download success")
			}
		})

		t.Run("fail because of wrong download url", func(t *testing.T) {
			book := NewBook("test", 1, 10, downloadConfig)
			book.config.DownloadUrl = server.URL + "/content/invalid/1"
			book.config.ConstSleep = 0
			book.SetTitle("book-title")
			book.SetWriter("book-writer")
			result := book.Download("/test-data/storage", 10, &mutex)

			if result {
				t.Fatalf("book download success")
			}
		})
	})
}
