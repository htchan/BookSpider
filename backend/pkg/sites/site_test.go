package sites

import (
	"testing"
	"log"
	"io/ioutil"
	"os"
	"io"
	"github.com/htchan/BookSpider/pkg/books"
	"github.com/htchan/BookSpider/internal/utils"
)

var meta, _ = books.NewMetaInfo(map[string]string{
	"baseUrl": "https://www.ck101.org/book/%v.html",
	"downloadUrl": "https://www.ck101.org/0/%v/",
	"chapterUrl": "https://www.ck101.org%v",
	"chapterUrlPattern": "/.*?\\.html",
	"titleRegex": "<h1><a.*?>(.*?)</a></h1>",
	"writerRegex": "作者︰<a.*?>(.*?)</a>",
	"typeRegex": " &gt; (.*?) &gt; ",
	"lastUpdateRegex": "最新章節\\((\\d{4}-\\d{2}-\\d{2})\\)",
	"lastChapterRegex": "<a.*?id=\"newchapter\".*>(.*?)./a>",
	"chapterUrlRegex": "<dd><a href=\"(.*?)\">.*?</a></dd>",
	"chapterTitleRegex": "<dd><a href.*?>(.*?)</a></dd>",
	"chapterContentRegex": "(?s)<div.*?yuedu_zhengwen.*?>(.*?)</div>",
})

func init() {
	log.SetOutput(ioutil.Discard)
	source, err := os.Open("../../test/site-test-data/ck101_template.db")
	utils.CheckError(err)
	destination, err := os.Create("../../test/site-test-data/ck101.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func mapEqual(map1, map2 map[string]interface{}) (bool) {
	for key:= range map1 {
		if _, exist:= map2[key]; !exist {
			return false
		}
	}
	for key:= range map2 {
		if val, exist:= map2[key]; !exist || val != map2[key] {
			return false
		}
	}
	return true
}

func TestLoadSite(t *testing.T) {
	t.Skip("I'm lazy")
}

func TestInfo(t *testing.T) {
	t.Skip("I'm lazy")
}

func TestRandomSuggestBook(t *testing.T) {
	t.Skip("I'm lazy")
}

func TestSearch(t *testing.T) {
	t.Skip("I'm lazy")
}

func TestBackupSql(t *testing.T) {
	t.Skip("I'm lazy")
}

func TestMap(t *testing.T) {
	t.Skip("I'm lazy")
}
