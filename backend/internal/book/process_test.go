package book

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/htchan/BookSpider/internal/book/model"
	"github.com/htchan/BookSpider/internal/client"
	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/mock"
	"os"
	"testing"
	"time"
)

func TestBook_isEnd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		book   Book
		result bool
	}{
		{name: "keyword in last chapter", book: Book{BookModel: model.BookModel{UpdateChapter: "abc后记def"}}, result: true},
		{name: "keyword not in last chapter", book: Book{BookModel: model.BookModel{UpdateChapter: "abcdef"}}, result: false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result := test.book.isEnd()
			if result != test.result {
				t.Errorf("book.isEnd() return %v, expect %v", result, test.result)
			}
		})
	}
}

func TestBook_Process(t *testing.T) {
	t.Parallel()

	dateLayout := "2006-01-02T15:04:05"
	oldDate := (time.Now().Add(-1.5 * 365 * 24 * time.Hour)).Format(dateLayout)
	newDate := (time.Now().Add(-10 * 24 * time.Hour)).Format(dateLayout)
	server := mock.MockBookServer(oldDate, newDate)
	client := &client.CircuitBreakerClient{
		CircuitBreakerClientConfig: config.CircuitBreakerClientConfig{
			Timeout: 1,
		},
	}
	client.Init(10)

	t.Cleanup(func() {
		server.Close()
		os.Remove("./11.txt")
		os.Remove("./12.txt")
		os.Remove("./13.txt")
	})

	tests := []struct {
		name               string
		book               Book
		result             bool
		expectedBookStatus model.StatusCode
		downloaded         bool
		downloadContent    string
	}{
		{
			name: "success update, end and download",
			book: Book{
				BookModel: model.BookModel{
					ID: 11, Title: "title", Type: "type",
					UpdateChapter: "chapter", UpdateDate: oldDate,
				},
				WriterModel: model.WriterModel{Name: "writer"},
				BookConfig: &config.BookConfig{
					URLConfig: config.URLConfig{
						Base:          server.URL + "/updated_chapter_end/%v",
						Download:      server.URL + "/download/%v",
						ChapterPrefix: server.URL + "/chapter_success",
					},
					UpdateDateLayout: dateLayout,
					Storage:          "",
					SourceKey:        "test_book",
				},
				CircuitBreakerClient: client,
			},
			result:             true,
			expectedBookStatus: model.Download,
			downloaded:         true,
			downloadContent: fmt.Sprintf(
				"title_new\nwriter_new\n%v\n\ntitle_1\n%v\ncontent1\n%v\n",
				CONTENT_SEP, CONTENT_SEP, CONTENT_SEP,
			),
		},
		{
			name: "updated but not ended",
			book: Book{
				BookModel: model.BookModel{
					ID: 12, Title: "title", Type: "type",
					UpdateChapter: "chapter_old", UpdateDate: oldDate,
				},
				WriterModel: model.WriterModel{Name: "writer"},
				BookConfig: &config.BookConfig{
					URLConfig: config.URLConfig{
						Base: server.URL + "/updated_chapter_not_end/%v",
					},
					UpdateDateLayout: dateLayout,
					Storage:          "",
					SourceKey:        "test_book",
				},
				CircuitBreakerClient: client,
			},
			result:             true,
			expectedBookStatus: model.InProgress,
			downloaded:         false,
			downloadContent:    "",
		},
		{
			name: "not updated but update date is 1 year before",
			book: Book{
				BookModel: model.BookModel{
					ID: 13, Title: "title", Type: "type",
					UpdateChapter: "chapter", UpdateDate: oldDate,
				},
				WriterModel: model.WriterModel{Name: "writer"},
				BookConfig: &config.BookConfig{
					URLConfig: config.URLConfig{
						Base:          server.URL + "/no_updated/%v",
						Download:      server.URL + "/download/%v",
						ChapterPrefix: server.URL + "/chapter_success",
					},
					UpdateDateLayout: dateLayout,
					Storage:          "",
					SourceKey:        "test_book",
				},
				CircuitBreakerClient: client,
			},
			result:             true,
			expectedBookStatus: model.Download,
			downloaded:         true,
			downloadContent: fmt.Sprintf(
				"title\nwriter\n%v\n\ntitle_1\n%v\ncontent1\n%v\n",
				CONTENT_SEP, CONTENT_SEP, CONTENT_SEP,
			),
		},
		{
			name: "not updated and update date within 1 year",
			book: Book{
				BookModel: model.BookModel{
					UpdateDate: newDate,
				},
				BookConfig: &config.BookConfig{
					URLConfig: config.URLConfig{
						Base:     server.URL + "/error/%v",
						Download: server.URL + "/error/%v",
					},
					UpdateDateLayout: dateLayout,
					Storage:          "",
					SourceKey:        "test_book",
				},
				CircuitBreakerClient: client,
			},
			result:             false,
			expectedBookStatus: model.Error,
			downloaded:         false,
			downloadContent:    "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			result := test.book.Process()

			if result != test.result {
				t.Errorf("book.Process() return %v, want %v", result, test.result)
			}

			if test.book.Status != test.expectedBookStatus {
				t.Errorf("book.Process() set book status to %v, want %v", test.book.Status, test.expectedBookStatus)
			}

			if test.downloaded {
				content := test.book.Content()
				if !cmp.Equal(string(content), test.downloadContent) {
					t.Errorf("book.Process() download wrong content")
					t.Error(cmp.Diff(string(content), test.downloadContent))
				}
			}
		})
	}
}
