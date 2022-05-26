package books

import (
	"os"
	"testing"
	"errors"
	"github.com/htchan/BookSpider/pkg/configs"
	"github.com/htchan/BookSpider/internal/database"
)

var accessorConfig = configs.LoadSourceConfigs(os.Getenv("ASSETS_LOCATION") + "/test-data/configs")["test_source_key"]

func TestBooks_Book_Accessor(t *testing.T) {
	book := Book{
		bookRecord: &database.BookRecord{
			Site: "test",
			Id: 1,
			HashCode: 100,
			Title: "title-1",
			WriterId: 1,
			Type: "type-1",
			UpdateDate: "104",
			UpdateChapter: "chapter-1",
			Status: database.InProgress,
		},
		writerRecord: &database.WriterRecord{
			Id: 1,
			Name: "writer-1",
		},
		errorRecord: nil,
		config: accessorConfig.Populate(1),
	}
	
	t.Run("func GetInfo/success", func(t *testing.T) {
		t.Parallel()
		site, id, hashCode := book.GetInfo()
		if site != "test" || id != 1 || hashCode != 100 {
			t.Errorf(
				"book.GetInfo return wrong value - site: %v, id: %v, hash code: %v",
				site, id, hashCode)
		}
	})

	t.Run("func GetTitle/success", func(t *testing.T) {
		t.Parallel()
		title := book.GetTitle()
		if title != "title-1" {
			t.Errorf("book.GetTitle return wrong value: %v", title)
		}
	})

	t.Run("func SetTitle/success", func(t *testing.T) {
		t.Parallel()
		book := NewBook("test", 1, 100, accessorConfig)
		book.SetTitle("test-title")
		if book.bookRecord.Title != "test-title" {
			t.Errorf("book.SetTitle update wrong value: %v", book.bookRecord.Title)
		}
	})

	t.Run("func GetWriter/success", func(t *testing.T) {
		t.Parallel()
		writer := book.GetWriter()
		if writer != "writer-1" {
			t.Errorf("book.GetWriter return wrong value: %v", writer)
		}
	})

	t.Run("func SetWriter/success", func(t *testing.T) {
		t.Parallel()
		book := NewBook("test", 1, 100, accessorConfig)
		book.SetWriter("test-writer")
		if book.bookRecord.WriterId != -1 || book.writerRecord.Id != -1 ||
			book.writerRecord.Name != "test-writer" {
			t.Errorf("book.SetWriter update wrong value: %v", book.writerRecord.Name)
		}
	})

	t.Run("func GetType/success", func(t *testing.T) {
		t.Parallel()
		typeString := book.GetType()
		if typeString != "type-1" {
			t.Errorf("book.GetType return wrong value: %v", typeString)
		}
	})

	t.Run("func SetType/success", func(t *testing.T) {
		t.Parallel()
		book := NewBook("test", 1, 100, accessorConfig)
		book.SetType("test-type")
		if book.bookRecord.Type != "test-type" {
			t.Errorf("book.SetType update wrong value: %v", book.bookRecord.Type)
		}
	})

	t.Run("func GetUpdateDate/success", func(t *testing.T) {
		t.Parallel()
		updateDate := book.GetUpdateDate()
		if updateDate != "104" {
			t.Errorf("book.GetUpdateDate return wrong value: %v", updateDate)
		}
	})

	t.Run("func SetUpdateDate/success", func(t *testing.T) {
		t.Parallel()
		book := NewBook("test", 1, 100, accessorConfig)
		book.SetUpdateDate("test-date")
		if book.bookRecord.UpdateDate != "test-date" {
			t.Errorf("book.SetUpdateDate update wrong value: %v", book.bookRecord.UpdateDate)
		}
	})

	t.Run("func GetUpdateChapter/success", func(t *testing.T) {
		t.Parallel()
		updateChapter := book.GetUpdateChapter()
		if updateChapter != "chapter-1" {
			t.Errorf("book.GetUpdateChapter return wrong value: %v", updateChapter)
		}
	})

	t.Run("func SetUpdateChapter/success", func(t *testing.T) {
		t.Parallel()
		book := NewBook("test", 1, 100, accessorConfig)
		book.SetUpdateChapter("test-chapter")
		if book.bookRecord.UpdateChapter != "test-chapter" {
			t.Errorf("book.SetUpdateChapter update wrong value: %v", book.bookRecord.UpdateChapter)
		}
	})

	t.Run("func GetStatus/success", func(t *testing.T) {
		t.Parallel()
		status := book.GetStatus()
		if status != database.InProgress {
			t.Errorf("book.GetStatus return wrong value: %v", status)
		}
	})

	t.Run("func SetStatus/success", func(t *testing.T) {
		t.Parallel()
		book := NewBook("test", 1, 100, accessorConfig)
		book.SetStatus(database.Error)
		if book.bookRecord.Status != database.Error {
			t.Errorf("book.SetStatus update wrong value: %v", book.bookRecord.Status)
		}
	})

	t.Run("func GetError", func(t *testing.T) {
		t.Parallel()
		t.Run("success for nil error", func(t *testing.T) {
			err := book.GetError()
			if err != nil {
				t.Errorf("book.GetError return wrong value: %v", err)
			}
		})

		t.Run("success for existing error", func(t *testing.T) {
			book.errorRecord = &database.ErrorRecord{Error: errors.New("testing")}
			err := book.GetError()
			if err == nil || err.Error() != "testing" {
				t.Errorf("book.GetError return wrong value: %v", err)
			}
		})
	})

	t.Run("func SetError", func(t *testing.T) {
		t.Parallel()
		t.Run("success with specific error", func(t *testing.T) {
			t.Parallel()
			book := NewBook("test", 1, 100, accessorConfig)
			book.SetError(errors.New("test-error"))
			if book.errorRecord == nil || book.errorRecord.Site != "test" ||
				book.errorRecord.Id != 1 || book.errorRecord.Error.Error() != "test-error" {
				t.Errorf("book.SetError update wrong value: %v", book.errorRecord)
			}
		})

		t.Run("success with nil error", func(t *testing.T) {
			t.Parallel()
			book := NewBook("test", 1, 100, accessorConfig)
			book.SetError(errors.New("test-error"))
			book.SetError(nil)
			if book.errorRecord != nil {
				t.Errorf("book.SetError update wrong value: %v", book.errorRecord)
			}
		})
	})

	t.Run("func getContentLocation/success", func(t *testing.T) {
		t.Parallel()
		contentLocation := book.getContentLocation("/test-data/storage")
		if contentLocation != os.Getenv("ASSETS_LOCATION") + "/test-data/storage/1-v100.txt" {
			t.Errorf("book.getContentLocation return wrong value: %v", contentLocation)
		}
	})

	t.Run("func HasContent", func(t *testing.T) {
		t.Parallel()
		t.Run("return true because path exist", func(t *testing.T) {
			t.Parallel()
			book := NewBook("test", 1, 100, accessorConfig)
			result := book.HasContent("/test-data/storage")

			if !result {
				t.Errorf("book does not have content for exist file")
			}
		})

		t.Run("return false because path not exist", func(t *testing.T) {
			t.Parallel()
			book := NewBook("test", 1, 200, accessorConfig)
			result := book.HasContent("/test-data/storage")

			if result {
				t.Errorf("book has content for not exist file")
			}
		})
	})

	t.Run("func GetContent", func(t *testing.T) {
		t.Parallel()
		t.Run("success", func(t *testing.T) {
			t.Parallel()
			content := book.GetContent("/test-data/storage")
			if content != "hello" {
				t.Errorf("book.GetContent return wrong value: %v", content)
			}
		})
		
		t.Run("return empty string if book has no content", func(t *testing.T) {
			t.Parallel()
			book := NewBook("test", 1, 200, accessorConfig)
			content := book.GetContent("/test-data/storage")
			if content != "" {
				t.Errorf("non download book GetContent return some value: %v", content)
			}
		})
	})

	t.Run("func Map/success", func(t *testing.T) {
		result := book.Map()

		if result["site"] != "test" ||
			result["id"] != 1 ||
			result["hash"] != "2s" ||
			result["status"] != "in_progress" ||
			result["title"] != "title-1" ||
			result["writer"] != "writer-1" ||
			result["type"] != "type-1" ||
			result["updateChapter"] != "chapter-1" ||
			result["updateDate"] != "104" {
			t.Errorf("wrong map: %v", result)
		}
	})

	t.Run("func String", func(t *testing.T) {
		result := book.String()

		if result != "test-1-2s" {
			t.Errorf("wrong string: %v", result)
		}
	})
}