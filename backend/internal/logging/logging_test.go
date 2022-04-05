package logging

import (
	"testing"
	"log"
	"bytes"
	"os"
)

var buf bytes.Buffer

func stubLogOutput() {
    log.SetOutput(&buf)
}
func revertLogOutput() {
	log.SetOutput(os.Stderr)
}

func Test(t *testing.T) {
	stubLogOutput()

	t.Run("func logInfo", func(t *testing.T) {
		t.Run("test with data is not nil", func(t *testing.T) {
			buf.Reset()
			logInfo("test", "hello")
			result := buf.String()

			if result[20:] != "test hello\n" {
				t.Fatalf("got result: %v", result[20:])
			}
		})
		
		t.Run("test with data is nil", func(t *testing.T) {
			buf.Reset()
			logInfo("test", nil)
			result := buf.String()

			if result[20:] != "test\n" {
				t.Fatalf("got result: %v", result[20:])
			}
		})
	})

	t.Run("func LogBookEvent", func(t *testing.T) {
		buf.Reset()

		LogBookEvent("site-id-version", "action", "event", nil)
		result := buf.String()

		if result[20:] != "book-spider.book.site-id-version.action.event\n" {
			t.Fatalf("got result: %v", result[20:])
		}
	})

	t.Run("func LogSiteEvent", func(t *testing.T) {
		buf.Reset()

		LogSiteEvent("site-name", "action", "event", nil)
		result := buf.String()

		if result[20:] != "book-spider.site.site-name.action.event\n" {
			t.Fatalf("got result: %v", result[20:])
		}
	})

	t.Run("func LogRequestEvent", func(t *testing.T) {
		buf.Reset()

		LogRequestEvent("action", "event", nil)
		result := buf.String()

		if result[20:] != "book-spider.request.action.event\n" {
			t.Fatalf("got result: %v", result[20:])
		}
	})

	t.Run("func LogEvent", func(t *testing.T) {
		buf.Reset()

		LogEvent("area", "event", nil)
		result := buf.String()

		if result[20:] != "book-spider.area.event\n" {
			t.Fatalf("got result: %v", result[20:])
		}
	})

	revertLogOutput()
}