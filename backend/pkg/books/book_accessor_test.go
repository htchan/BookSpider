package books

import (
	"os"
	// "io"
	"testing"
	"errors"
	// "github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/pkg/configs"
	"github.com/htchan/BookSpider/internal/database"
	"github.com/htchan/BookSpider/internal/database/sqlite"
)

var accessorConfig = configs.LoadConfigYaml(os.Getenv("ASSETS_LOCATION") + "/test-data/config.yml").SiteConfigs["test"].BookMeta

func Test_Books_Book_GetInfo(t *testing.T) {
	db := sqlite.NewSqliteDB("./book_test.db")
	defer db.Close()
	
	t.Run("success", func(t *testing.T) {
		book := LoadBook(db, "test", 1, 100, accessorConfig)
		site, id, hashCode := book.GetInfo()
		if site != "test" || id != 1 || hashCode != 100 {
			t.Fatalf(
				"book.GetInfo return wrong value - site: %v, id: %v, hash code: %v",
				site, id, hashCode)
		}
	})
}

func Test_Books_Book_GetTitle(t *testing.T) {
	db := sqlite.NewSqliteDB("./book_test.db")
	defer db.Close()
	
	t.Run("success", func(t *testing.T) {
		book := LoadBook(db, "test", 1, 100, accessorConfig)
		title := book.GetTitle()
		if title != "title-1" {
			t.Fatalf("book.GetTitle return wrong value: %v", title)
		}
	})
}

func Test_Books_Book_SetTitle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		book := NewBook("test", 1, 100, accessorConfig)
		book.SetTitle("test-title")
		if book.bookRecord.Title != "test-title" {
			t.Fatalf("book.SetTitle update wrong value: %v", book.bookRecord.Title)
		}
	})
}

func Test_Books_Book_GetWriter(t *testing.T) {
	db := sqlite.NewSqliteDB("./book_test.db")
	defer db.Close()
	
	t.Run("success", func(t *testing.T) {
		book := LoadBook(db, "test", 1, 100, accessorConfig)
		writer := book.GetWriter()
		if writer != "writer-1" {
			t.Fatalf("book.GetWriter return wrong value: %v", writer)
		}
	})
}

func Test_Books_Book_SetWriter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		book := NewBook("test", 1, 100, accessorConfig)
		book.SetWriter("test-writer")
		if book.bookRecord.WriterId != -1 || book.writerRecord.Id != -1 ||
			book.writerRecord.Name != "test-writer" {
			t.Fatalf("book.SetWriter update wrong value: %v", book.writerRecord.Name)
		}
	})
}

func Test_Books_Book_GetType(t *testing.T) {
	db := sqlite.NewSqliteDB("./book_test.db")
	defer db.Close()
	
	t.Run("success", func(t *testing.T) {
		book := LoadBook(db, "test", 1, 100, accessorConfig)
		typeString := book.GetType()
		if typeString != "type-1" {
			t.Fatalf("book.GetType return wrong value: %v", typeString)
		}
	})
}

func Test_Books_Book_SetType(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		book := NewBook("test", 1, 100, accessorConfig)
		book.SetType("test-type")
		if book.bookRecord.Type != "test-type" {
			t.Fatalf("book.SetType update wrong value: %v", book.bookRecord.Type)
		}
	})
}

func Test_Books_Book_GetUpdateDate(t *testing.T) {
	db := sqlite.NewSqliteDB("./book_test.db")
	defer db.Close()
	
	t.Run("success", func(t *testing.T) {
		book := LoadBook(db, "test", 1, 100, accessorConfig)
		updateDate := book.GetUpdateDate()
		if updateDate != "104" {
			t.Fatalf("book.GetUpdateDate return wrong value: %v", updateDate)
		}
	})
}

func Test_Books_Book_SetUpdateDate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		book := NewBook("test", 1, 100, accessorConfig)
		book.SetUpdateDate("test-date")
		if book.bookRecord.UpdateDate != "test-date" {
			t.Fatalf("book.SetUpdateDate update wrong value: %v", book.bookRecord.UpdateDate)
		}
	})
}

func Test_Books_Book_GetUpdateChapter(t *testing.T) {
	db := sqlite.NewSqliteDB("./book_test.db")
	defer db.Close()
	
	t.Run("success", func(t *testing.T) {
		book := LoadBook(db, "test", 1, 100, accessorConfig)
		updateChapter := book.GetUpdateChapter()
		if updateChapter != "chapter-1" {
			t.Fatalf("book.GetUpdateChapter return wrong value: %v", updateChapter)
		}
	})
}

func Test_Books_Book_SetUpdateChapter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		book := NewBook("test", 1, 100, accessorConfig)
		book.SetUpdateChapter("test-chapter")
		if book.bookRecord.UpdateChapter != "test-chapter" {
			t.Fatalf("book.SetUpdateChapter update wrong value: %v", book.bookRecord.UpdateChapter)
		}
	})
}

func Test_Books_Book_GetStatus(t *testing.T) {
	db := sqlite.NewSqliteDB("./book_test.db")
	defer db.Close()
	
	t.Run("success", func(t *testing.T) {
		book := LoadBook(db, "test", 1, 100, accessorConfig)
		status := book.GetStatus()
		if status != database.InProgress {
			t.Fatalf("book.GetStatus return wrong value: %v", status)
		}
	})
}

func Test_Books_Book_SetStatus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		book := NewBook("test", 1, 100, accessorConfig)
		book.SetStatus(database.Error)
		if book.bookRecord.Status != database.Error {
			t.Fatalf("book.SetStatus update wrong value: %v", book.bookRecord.Status)
		}
	})
}

func Test_Books_Book_GetError(t *testing.T) {
	db := sqlite.NewSqliteDB("./book_test.db")
	defer db.Close()
	
	t.Run("success", func(t *testing.T) {
		book := LoadBook(db, "test", 1, 100, accessorConfig)
		err := book.GetError()
		if err != nil {
			t.Fatalf("book.GetError return wrong value: %v", err)
		}
	})
}

func Test_Books_Book_SetError(t *testing.T) {
	t.Run("success with specific error", func(t *testing.T) {
		book := NewBook("test", 1, 100, accessorConfig)
		book.SetError(errors.New("test-error"))
		if book.errorRecord == nil || book.errorRecord.Site != "test" ||
			book.errorRecord.Id != 1 || book.errorRecord.Error.Error() != "test-error" {
			t.Fatalf("book.SetError update wrong value: %v", book.errorRecord)
		}
	})
	t.Run("success with nil error", func(t *testing.T) {
		book := NewBook("test", 1, 100, accessorConfig)
		book.SetError(errors.New("test-error"))
		book.SetError(nil)
		if book.errorRecord != nil {
			t.Fatalf("book.SetError update wrong value: %v", book.errorRecord)
		}
	})
}

func Test_Books_Book_getContentLocation(t *testing.T) {
	db := sqlite.NewSqliteDB("./book_test.db")
	defer db.Close()
	
	t.Run("success", func(t *testing.T) {
		book := LoadBook(db, "test", 1, 100, accessorConfig)
		contentLocation := book.getContentLocation()
		if contentLocation != os.Getenv("ASSETS_LOCATION") + "/test-data/storage/1-v100.txt" {
			t.Fatalf("book.getContentLocation return wrong value: %v", contentLocation)
		}
	})
}

func Test_Books_Book_HasContent(t *testing.T) {
	t.Run("return true because path exist", func(t *testing.T) {
		book := NewBook("test", 1, 100, accessorConfig)
		result := book.HasContent()

		if !result {
			t.Fatalf("book does not have content for exist file")
		}
	})

	t.Run("return false because path not exist", func(t *testing.T) {
		book := NewBook("test", 1, 200, accessorConfig)
		result := book.HasContent()

		if result {
			t.Fatalf("book has content for not exist file")
		}
	})
}

func Test_Books_Book_GetContent(t *testing.T) {
	db := sqlite.NewSqliteDB("./book_test.db")
	defer db.Close()
	
	t.Run("success", func(t *testing.T) {
		book := LoadBook(db, "test", 1, 100, accessorConfig)
		content := book.GetContent()
		if content != "hello" {
			t.Fatalf("book.GetContent return wrong value: %v", content)
		}
	})
	
	t.Run("return empty string if book has no content", func(t *testing.T) {
		book := NewBook("test", 1, 200, accessorConfig)
		content := book.GetContent()
		if content != "" {
			t.Fatalf("non download book GetContent return some value: %v", content)
		}
	})
}