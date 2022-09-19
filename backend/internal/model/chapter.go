package model

import (
	"fmt"
	"strings"
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

func (c Chapter) ContentString() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n", c.Title, CONTENT_SEP, c.Content, CONTENT_SEP)
}

func (c Chapter) Equal(compare Chapter) bool {
	return c.Index == compare.Index && c.URL == compare.URL &&
		c.Title == compare.Title && c.Content == compare.Content &&
		((c.Error == nil && compare.Error == nil) ||
			true || (c.Error != nil && compare.Error != nil && c.Error.Error() == compare.Error.Error()))
}
