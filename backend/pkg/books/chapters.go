package books

import (
	"strings"
	"strconv"
	"sort"
	"fmt"

	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/pkg/configs"
	"github.com/htchan/ApiParser"
)

type Chapter struct {
	Index int
	Url, Title, Content string
}

func NewChapter(i int, url, title string, config *configs.SourceConfig) (chapter Chapter) {
	chapter.Index = i
	//TODO: check the reason of adding the http here
	if strings.HasPrefix(url, "/") || strings.HasPrefix(url, "http") {
		url = fmt.Sprintf(config.ChapterUrl, url)
	} else {
		url = config.DownloadUrl + url
	}
	chapter.Url = url
	chapter.Title = title
	return
}

func (chapter *Chapter)generateIndex() {
	numberMap := map[string]string {
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
	} else
	if strings.Contains(chapter.Title, "中") {
		chapter.Index += 5
	} else
	if strings.Contains(chapter.Title, "下") {
		chapter.Index += 8
	}
}

func (chapter *Chapter)optimizeContent() {
	replaceItems := []struct{
		old, new string
	} {
		{"<br />", ""},
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

}

func (chapter *Chapter)Download(config *configs.SourceConfig, validHTML func(string)error) {
	// get chapter resource
	html, _ := utils.GetWeb(chapter.Url, 10, config.Decoder, config.ConstSleep)
	if err := validHTML(html); err != nil {
		chapter.Content = "load html fail"
		return
	}
	// extract chapter
	responseApi := ApiParser.Parse(html, config.SourceKey + ".chapter_content")
	content, ok := responseApi.Data["Content"]

	// content, err := utils.Search(html, config.ChapterContentRegex)
	if !ok {
		chapter.Content = "recognize html fail\n" + html
	} else {
		chapter.Content = content
		chapter.optimizeContent()
	}
}