package books

import (
	"testing"
	"errors"
	"github.com/htchan/BookSpider/internal/mock"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/pkg/configs"
)

var updateConfig = configs.LoadConfigYaml("../../assets/test-data/config.yml").SiteConfigs["test"].BookMeta

func Test_Books_Book_fetchInfo(t *testing.T) {
	book := NewBook("test", 1, -1, updateConfig)

	t.Run("Success", func(t *testing.T) {
		getWeb = mock.UpdateGetWebSuccess
		
		title, writer, typeString, updateDate, updateChapter, err := book.fetchInfo()
		
		if err != nil || title != "title-regex" ||
			writer != "writer-regex" || typeString != "type-regex" ||
			updateDate != "last-update-regex" ||
			updateChapter != "last-chapter-regex" {
				t.Fatalf("book fetch info failed - book: %v, err: %v", book.bookRecord, err)
		}
	})

	t.Run("fail because of invalid html", func(t *testing.T) {
		getWeb = mock.UpdateGetWebEmpty
		
		title, writer, typeString, updateDate, updateChapter, err := book.fetchInfo()
		
		if err == nil || title != "" || writer != "" || typeString != "" ||
			updateDate != "" || updateChapter != "" {
				t.Fatalf("book fetch info success for invalid html - book: %v, err: %v", book.bookRecord, err)
		}
	})

	t.Run("fail because some regex failed", func(t *testing.T) {
		getWeb = mock.UpdateGetWebPartialFail
		
		title, writer, typeString, updateDate, updateChapter, err := book.fetchInfo()
		
		if err == nil || title != "" || writer != "" || typeString != "" ||
			updateDate != "" || updateChapter != "" {
				t.Fatalf("book fetch info success for invalid html - book: %v, err: %v", book.bookRecord, err)
		}
	})
}

func Test_Books_Book_Update(t *testing.T) {
	t.Run("update if last update is new", func(T *testing.T) {
		getWeb = mock.UpdateGetWebSuccess

		book := NewBook("testing", 1, -1, updateConfig)
		book.SetTitle("title-regex")
		book.SetWriter("writer-regex")
		book.SetType("type-regex")
		book.SetUpdateDate("old-update-date")
		book.SetUpdateChapter("last-chapter-regex")
		book.SetStatus(database.InProgress)
		result := book.Update()

		if !result || book.GetTitle() != "title-regex" ||
			book.GetWriter() != "writer-regex" || book.GetType() != "type-regex" ||
			book.GetError() != nil || book.GetStatus() != database.InProgress ||
			book.GetUpdateDate() != "last-update-regex" ||
			book.GetUpdateChapter() != "last-chapter-regex" {
				t.Fatalf("book update success with update date different: %v", book.bookRecord)
		}
	})

	t.Run("update if last chapter is new", func(T *testing.T) {
		getWeb = mock.UpdateGetWebSuccess

		book := NewBook("testing", 1, -1, updateConfig)
		book.SetTitle("title-regex")
		book.SetWriter("writer-regex")
		book.SetType("type-regex")
		book.SetUpdateDate("last-update-regex")
		book.SetUpdateChapter("old-chapter")
		book.SetStatus(database.Download)
		result := book.Update()

		if !result || book.GetTitle() != "title-regex" ||
			book.GetWriter() != "writer-regex" || book.GetType() != "type-regex" ||
			book.GetError() != nil || book.GetStatus() != database.InProgress ||
			book.GetUpdateDate() != "last-update-regex" ||
			book.GetUpdateChapter() != "last-chapter-regex" {
				t.Fatalf("book update success with update chapter different: %v", book.bookRecord)
		}
	})

	t.Run("update error book if title is different", func(T *testing.T) {
		getWeb = mock.UpdateGetWebSuccess

		book := NewBook("testing", 1, 10, updateConfig)
		book.SetError(errors.New("test-error"))
		result := book.Update()

		if !result || book.bookRecord.HashCode != 10 || book.GetTitle() != "title-regex" ||
			book.GetWriter() != "writer-regex" || book.GetType() != "type-regex" ||
			book.GetError() != nil || book.GetStatus() != database.InProgress ||
			book.GetUpdateDate() != "last-update-regex" ||
			book.GetUpdateChapter() != "last-chapter-regex" {
				t.Fatalf("error book update success with fetch info: %v", book.bookRecord)
		}
	})

	t.Run("update normal book hash if title is different", func(T *testing.T) {
		getWeb = mock.UpdateGetWebSuccess

		book := NewBook("testing", 1, 10, updateConfig)
		book.SetTitle("old-title-regex")
		book.SetWriter("writer-regex")
		book.SetType("type-regex")
		book.SetUpdateDate("last-update-regex")
		book.SetUpdateChapter("last-chapter-regex")
		book.SetStatus(database.InProgress)
		result := book.Update()

		if !result || book.bookRecord.HashCode == 10 || book.GetTitle() != "title-regex" ||
			book.GetWriter() != "writer-regex" || book.GetType() != "type-regex" ||
			book.GetError() != nil || book.GetStatus() != database.InProgress ||
			book.GetUpdateDate() != "last-update-regex" ||
			book.GetUpdateChapter() != "last-chapter-regex" {
				t.Fatalf("book update success even partial fail: %v", book.bookRecord)
		}
	})

	t.Run("set error if error book fetch info failed", func(T *testing.T) {
		getWeb = mock.UpdateGetWebPartialFail

		book := NewBook("testing", 1, -1, updateConfig)
		result := book.Update()

		if result || book.GetTitle() != "" ||
			book.GetWriter() != "" || book.GetType() != "" ||
			book.GetError().Error() == "" || book.GetStatus() != database.Error ||
			book.GetUpdateDate() != "" || book.GetUpdateChapter() != "" {
				t.Fatalf("book update success even partial fail: %v", book.bookRecord)
		}
	})

	t.Run("do nothing if normal book fetch info failed", func(T *testing.T) {
		getWeb = mock.UpdateGetWebPartialFail

		book := NewBook("testing", 1, -1, updateConfig)
		book.SetStatus(database.InProgress)
		result := book.Update()

		if result || book.GetTitle() != "" ||
			book.GetWriter() != "" || book.GetType() != "" ||
			book.GetError() != nil || book.GetStatus() == database.Error ||
			book.GetUpdateDate() != "" || book.GetUpdateChapter() != "" {
				t.Fatalf("book update success even partial fail: %v", book.bookRecord)
		}
	})
}