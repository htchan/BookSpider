package client

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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
