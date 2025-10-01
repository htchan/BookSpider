package client

import (
	"context"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type BookInfo struct {
	Title         string
	Author        string
	Type          string
	UpdateChapter string
	UpdateDate    time.Time
}

type ChapterEntry struct {
	Title string
	URL   string
}

type ChapterEntryList []ChapterEntry

type ChapterContent struct {
	Title string
	Body  string
}

type Client interface {
	GetBookInfo(ctx context.Context, bookID string) (*BookInfo, error)
	GetBookChapterList(ctx context.Context, bookID string) (ChapterEntryList, error)
	GetChapterContent(ctx context.Context, chapter ChapterEntry) (*ChapterContent, error)
	Available(ctx context.Context) bool
}

func ParseDoc(body string) (*goquery.Document, error) {
	return goquery.NewDocumentFromReader(strings.NewReader(body))
}

func GetGoqueryContentWithoutChildren(s *goquery.Selection) string {
	html, err := s.Html()
	if err == nil {
		replaceItems := []struct {
			old, new string
		}{
			{"<br/>", "\n"},
			{"&nbsp;", ""},
			{"\u00a0", ""},
			{"<b>", ""},
			{"</b>", ""},
			{"<p>", ""},
			{"</p>", "\n"},
			{"                ", ""},
			{"<p/>", "\n"},
		}
		for _, replaceItem := range replaceItems {
			html = strings.ReplaceAll(
				html, replaceItem.old, replaceItem.new)
		}

		s.SetHtml(html)
	}

	return strings.TrimSpace(s.Clone().Children().Remove().End().Text())
}

func GetGoqueryContentWithChildren(s *goquery.Selection) string {
	html, err := s.Html()
	if err == nil {
		replaceItems := []struct {
			old, new string
		}{
			{"<br/>", "\n"},
			{"&nbsp;", ""},
			{"\u00a0", ""},
			{"<b>", ""},
			{"</b>", ""},
			{"<p>", ""},
			{"</p>", "\n"},
			{"                ", ""},
			{"<p/>", "\n"},
		}
		for _, replaceItem := range replaceItems {
			html = strings.ReplaceAll(
				html, replaceItem.old, replaceItem.new)
		}

		s.SetHtml(html)
	}

	return strings.TrimSpace(s.Clone().Text())
}
