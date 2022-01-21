package books

import (
	"testing"
	"errors"
	"github.com/htchan/BookSpider/internal/mock"
	"github.com/htchan/BookSpider/pkg/configs"
)

func Test_Books_NewChapter(t *testing.T) {
	t.Run("success with appending download url", func(t *testing.T) {
		config := configs.BookConfig{ DownloadUrl: "test/" }
		chapter := NewChapter(0, "0", "title-0", &config)
		if chapter.Index != 0 || chapter.Url != "test/0" ||
			chapter.Title != "title-0" || chapter.Content != "" {
				t.Fatalf("NewChatper return wrong result: %v", chapter)
		}
	})

	t.Run("success with formating download url", func(t *testing.T) {
		config := configs.BookConfig{ ChapterUrl: "test%v" }
		chapter := NewChapter(0, "/0", "title-0", &config)
		if chapter.Index != 0 || chapter.Url != "test/0" ||
			chapter.Title != "title-0" || chapter.Content != "" {
				t.Fatalf("NewChatper return wrong result: %v", chapter)
		}
	})
}

func Test_Books_Chapter_generateIndex(t *testing.T) {
	t.Run("success without any replace", func(t *testing.T) {
		chapter := Chapter{ Title: "abc421def" }
		chapter.generateIndex()
		if chapter.Index != 4210 {
			t.Fatalf("chapter generate index is %v, not 4210", chapter.Index)
		}
	})

	t.Run("success with simplified number", func(t *testing.T) {
		chapter := Chapter{ Title: "abc四二一def" }
		chapter.generateIndex()
		if chapter.Index != 4210 {
			t.Fatalf("chapter generate index is %v, not 4210", chapter.Index)
		}
	})

	t.Run("success with simplified number and replace word", func(t *testing.T) {
		chapter := Chapter{ Title: "abc四百二十一def" }
		chapter.generateIndex()
		if chapter.Index != 4210 {
			t.Fatalf("chapter generate index is %v, not 4210", chapter.Index)
		}
	})

	t.Run("success with traditional number", func(t *testing.T) {
		chapter := Chapter{ Title: "abc肆貳壹def" }
		chapter.generateIndex()
		if chapter.Index != 4210 {
			t.Fatalf("chapter generate index is %v, not 4210", chapter.Index)
		}
	})

	t.Run("success with traditional number and replace word", func(t *testing.T) {
		chapter := Chapter{ Title: "abc肆佰貳拾壹def" }
		chapter.generateIndex()
		if chapter.Index != 4210 {
			t.Fatalf("chapter generate index is %v, not 4210", chapter.Index)
		}
	})

	t.Run("success with 上 in title", func(t *testing.T) {
		chapter := Chapter{ Title: "abc肆佰貳拾壹def上" }
		chapter.generateIndex()
		if chapter.Index != 4212 {
			t.Fatalf("chapter generate index is %v, not 4212", chapter.Index)
		}
	})

	t.Run("fail without any number in title", func(t *testing.T) {
		chapter := Chapter{ Title: "abcdef" }
		chapter.generateIndex()
		if chapter.Index != 9999990 {
			t.Fatalf("chapter generate index is %v, not 9999990", chapter.Index)
		}
	})
}

func Test_Books_Book_optimizedContent(t *testing.T) {
	t.Run("success to replace string", func(t *testing.T) {
		chapter := Chapter{
			Content: "a<br />b&nbsp;c<b>d</b>e<p>f</p>g                h<p/>i",
		}
		chapter.optimizeContent()

		if chapter.Content != "abcdefgh\ni" {
			t.Fatalf("chapter optimize content generate %v, but not %v",
				chapter.Content, "abcdefgh\ni")
		}
	})
}

func Test_Books_Chapter_sortChapters(t *testing.T) {
	t.Run("success to sort by index ascending order", func(t *testing.T) {
		chapters := []Chapter{
			Chapter{ Index: 5, Url: "1"},
			Chapter{ Index: 1, Url: "2"},
		}

		sortChapters(chapters)

		if chapters[0].Index != 1 || chapters[1].Index != 5 {
			t.Fatalf("sort chatpers return wrong order: %v", chapters)
		}
	})

	t.Run("not change order if index are same", func(t *testing.T) {
		chapters := []Chapter{
			Chapter{ Index: 5, Url: "5"},
			Chapter{ Index: 1, Url: "1"},
			Chapter{ Index: 3, Url: "3"},
			Chapter{ Index: 6, Url: "6"},
			Chapter{ Index: 1, Url: "2"},
			Chapter{ Index: 3, Url: "4"},
		}

		sortChapters(chapters)

		if chapters[0].Url != "1" || chapters[1].Url != "2" ||
			chapters[2].Url != "3" || chapters[3].Url != "4" ||
			chapters[4].Index != 5 || chapters[5].Index != 6 {
			t.Fatalf("sort chatpers return wrong order: %v", chapters)
		}
	})
}

func Test_Books_Chapter_OptimizeChapters(t *testing.T) {

}

func Test_Books_Chatper_Download(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		getWeb = mock.ChapterGetWebSuccess
		config := configs.BookConfig{
			DownloadUrl: "",
			ChapterContentRegex: "chapter-content-(\\d)-content-regex",
		}
		chapter := NewChapter(0, "0", "title-0", &config)
		validHTML := func(_ string)error { return nil }
		chapter.Download(&config, validHTML)

		if chapter.Content != "0" {
			t.Fatalf("chapter Download fail with content: %v", chapter.Content)
		}
	})

	t.Run("fail with invalid HTML", func(t *testing.T) {
		getWeb = mock.ChapterGetWebSuccess
		config := configs.BookConfig{
			DownloadUrl: "",
			ChapterContentRegex: "chapter-content-(\\d)-content-regex",
		}
		chapter := NewChapter(0, "0", "title-0", &config)
		validHTML := func(_ string)error { return errors.New("test error") }
		chapter.Download(&config, validHTML)

		if chapter.Content != "load html fail" {
			t.Fatalf("chapter Download fail with content: %v", chapter.Content)
		}
	})

	t.Run("fail with not recognized HTML", func(t *testing.T) {
		getWeb = mock.ChapterGetWebSuccess
		config := configs.BookConfig{
			DownloadUrl: "",
		}
		chapter := NewChapter(0, "0", "title-0", &config)
		validHTML := func(_ string)error { return nil }
		chapter.Download(&config, validHTML)

		if chapter.Content != "recognize html fail\nchapter-content-0-content-regex" {
			t.Fatalf("chapter Download fail with content: %v", chapter.Content)
		}
	})
}