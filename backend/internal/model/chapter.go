package model

import (
	"errors"
	"fmt"
	"strings"
)

var ChapterEndKeywords = []string{
	// traditional chinese
	"番外", "結局", "新書", "完結", "尾聲", "感言", "後記", "完本",
	"全書完", "全文完", "全文終", "全文結", "劇終", "（完）", "終章",
	"外傳", "結尾",
	// simplified chinese
	"番外", "结局", "新书", "完结", "尾声", "感言", "后记", "完本",
	"全书完", "全文完", "全文终", "全文结", "剧终", "（完）", "终章",
	"外传", "结尾",
}

var (
	ErrCannotParseContent = errors.New("cannot parse content")
)

type Chapter struct {
	Index   int
	URL     string
	Title   string
	Content string
	Error   error
}

type Chapters []Chapter

var CONTENT_SEP = strings.Repeat("-", 20)

func NewChapter(i int, url, title string) Chapter {
	return Chapter{Index: i, URL: url, Title: title}
}

func (c *Chapter) OptimizeContent() {
	lines := strings.Split(c.Content, "\n")
	lines = removeEmptyLines(lines)

	c.Content = strings.Join(lines, "\n\n")
}

func (c *Chapter) ContentString() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n", c.Title, CONTENT_SEP, c.Content, CONTENT_SEP)
}

func removeEmptyLines(lines []string) []string {
	result := make([]string, 0)

	for _, line := range lines {
		data := strings.TrimSpace(line)
		if len(data) > 0 {
			result = append(result, data)
		}
	}

	return result
}

func StringToChapters(content string) (Chapters, error) {
	splitedContent := strings.Split(content, CONTENT_SEP)
	if len(splitedContent) < 3 {
		return nil, ErrCannotParseContent
	}
	chapters := make(Chapters, 0)
	chapter := Chapter{
		Index:   0,
		Title:   strings.Join(removeEmptyLines(strings.Split(splitedContent[1], "\n")), "\n\n"),
		Content: strings.Join(removeEmptyLines(strings.Split(splitedContent[2], "\n")), "\n\n"),
	}
	chapter.OptimizeContent()
	chapters = append(chapters, chapter)
	chapter = Chapter{}
	index := 1

	for i, con := range splitedContent[3:] {
		lines := removeEmptyLines(strings.Split(con, "\n"))
		if i == len(splitedContent)-4 && len(lines) == 0 {
			continue
		}
		if len(lines) == 1 && chapter.Title == "" {
			chapter.Title = lines[0]
		} else {
			if chapter.Title == "" {
				lastChapterLines := strings.Split(chapters[len(chapters)-1].Content, "\n")
				chapter.Title = lastChapterLines[len(lastChapterLines)-1]
				chapters[len(chapters)-1].Content = strings.Join(lastChapterLines[:len(lastChapterLines)-2], "\n")
			}
			chapter.Content = strings.Join(lines, "\n\n")
			chapter.Index = index
			index++
			chapters = append(chapters, chapter)
			chapter = Chapter{}
		}
	}

	return chapters, nil
}
