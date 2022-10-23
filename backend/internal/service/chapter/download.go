package chapter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/htchan/ApiParser"
	"github.com/htchan/BookSpider/internal/client"
	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/model"
)

func chapterURL(bookID int, chapter *model.Chapter, bkConf *config.BookConfig) string {
	if strings.HasPrefix(chapter.URL, "http") {
		return chapter.URL
	} else if strings.HasPrefix(chapter.URL, "/") {
		return bkConf.URLConfig.ChapterPrefix + chapter.URL
	}
	downloadURL := fmt.Sprintf(bkConf.URLConfig.Download, bookID)
	if !strings.HasSuffix(downloadURL, "/") {
		downloadURL = downloadURL + "/"
	}
	return downloadURL + chapter.URL
}

func Download(bookID int, chapter *model.Chapter, bkConf *config.BookConfig, stConf *config.SiteConfig, c *client.CircuitBreakerClient) error {
	html, err := c.Get(chapterURL(bookID, chapter, bkConf))
	if err != nil {
		chapter.Error = fmt.Errorf("download chapter error: %w", err)
		return chapter.Error
	}
	for _, r := range bkConf.UnwantContent {
		re, err := regexp.Compile(r)
		if err != nil {
			continue
		}
		html = re.ReplaceAllString(html, "")
	}

	responseApi := ApiParser.Parse(stConf.BookKey+".chapter_content", html)
	content, ok := responseApi.Data["ChapterContent"]

	if !ok {
		chapter.Error = fmt.Errorf("chapter content not found")
	} else {
		chapter.Content = content
	}
	chapter.OptimizeContent()

	return chapter.Error
}
