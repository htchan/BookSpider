package book

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/htchan/ApiParser"
	"github.com/htchan/BookSpider/internal/logging"
	"github.com/htchan/BookSpider/internal/utils"
)

type Chapter struct {
	Index               int
	URL, Title, Content string
	*Book
}

func NewChapter(i int, url, title string, book *Book) Chapter {
	return Chapter{
		Index: i,
		URL:   url,
		Title: title,
		Book:  book,
	}
}

func (chapter Chapter) chapterURL() string {
	//TODO: check the reason of adding the http here
	if strings.HasPrefix(chapter.URL, "/") || strings.HasPrefix(chapter.URL, "http") {
		return chapter.Book.BookConfig.URLConfig.ChapterPrefix + chapter.URL
	} else {
		return chapter.Book.downloadURL() + chapter.URL
	}
}

func (chapter *Chapter) generateIndex() {
	numberMap := map[string]string{
		"序": "0",
		"一": "1", "二": "2", "三": "3", "四": "4", "五": "5",
		"六": "6", "七": "7", "八": "8", "九": "9",
		"壹": "1", "貳": "2", "叁": "3", "肆": "4", "伍": "5",
		"陸": "6", "柒": "7", "捌": "8", "玖": "9",
		"后记": "9999990", "後記": "9999990", "新书": "9999990", "新書": "9999990",
		"结局": "9999990", "結局": "9999990", "感言": "9999990", "尾声": "9999990",
		"尾聲": "9999990", "终章": "9999990", "終章": "9999990", "外传": "9999990",
		"外傳": "9999990", "完本": "9999990", "结束": "9999990", "結束": "9999990",
		"完結": "9999990", "完结": "9999990", "终结": "9999990", "終結": "9999990",
		"番外": "9999990", "结尾": "9999990", "結尾": "9999990",
		"全书完": "9999990", "全書完": "9999990", "全本完": "9999990",
	}
	replaceList := "十拾百佰千仟万萬 ()"
	var numberTitle string
	for _, char := range chapter.Title {
		str := string(char)
		if i, ok := numberMap[str]; ok {
			numberTitle += i
		} else if !strings.Contains(replaceList, str) {
			numberTitle += str
		}
	}
	result, err := utils.Search(numberTitle, "(\\d+)")
	if err != nil {
		chapter.Index = 9999990
	} else {
		chapter.Index, _ = strconv.Atoi(result)
		chapter.Index *= 10
	}
	if strings.Contains(chapter.Title, "上") {
		chapter.Index += 2
	} else if strings.Contains(chapter.Title, "中") {
		chapter.Index += 5
	} else if strings.Contains(chapter.Title, "下") {
		chapter.Index += 8
	}
}

func (chapter *Chapter) optimizeContent() {
	replaceItems := []struct {
		old, new string
	}{
		{"<br />", "\n"},
		{"&nbsp;", ""},
		{"<b>", ""},
		{"</b>", ""},
		{"<p>", ""},
		{"</p>", ""},
		{"                ", ""},
		{"<p/>", "\n"},
	}
	for _, replaceItem := range replaceItems {
		chapter.Content = strings.ReplaceAll(
			chapter.Content, replaceItem.old, replaceItem.new)
	}
}

func sortChapters(chapters []Chapter) {
	sort.Slice(chapters, func(i, j int) bool {
		return chapters[i].Index < chapters[j].Index
	})
}

func optimizeChapters(chapters []Chapter) {
	// remove duplicated chapters
}

func (chapter *Chapter) Fetch() {
	// get chapter resource
	html, err := chapter.Book.Get(chapter.chapterURL())
	if err != nil {
		chapter.Content = fmt.Sprintf("load html failed - %v", err)
		return
	}
	// extract chapter
	responseApi := ApiParser.Parse(chapter.Book.BookConfig.SourceKey+".chapter_content", html)
	content, ok := responseApi.Data["ChapterContent"]

	// content, err := utils.Search(html, config.ChapterContentRegex)
	if !ok {
		chapter.Content = "recognize html fail\n" + html
		chapter.optimizeContent()
		logging.LogChapterEvent(chapter.chapterURL(), "download", "failed", "content not found")
	} else {
		chapter.Content = content
		chapter.optimizeContent()
	}
}

func (chapter Chapter) content() string {
	return chapter.Title + "\n" + CONTENT_SEP + "\n" + chapter.Content + "\n" + CONTENT_SEP + "\n"
}
