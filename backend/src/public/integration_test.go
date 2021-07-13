package public

import (
	"testing"
	"log"
	"io/ioutil"
	"github.com/htchan/BookSpider/model"
	"github.com/htchan/BookSpider/helper"
	"fmt"
	"os"
)

var config model.Config
var testSites map[string]model.Site

func init() {
	log.SetOutput(ioutil.Discard)
	config = model.LoadYaml("./config/config.yaml")
	for key := range config.Sites {
		config.Sites[key]["databaseLocation"] = "./test_res/template.db"
		config.Sites[key]["downloadLocation"] = "./test_res/"
	}
	testSites = model.LoadSitesYaml(config)
}


func TestIntegrate_LoadBookUpdate(t *testing.T) {
	var testcases = []struct {
		site string
		id int
		expected string
	} {
		{"ck101", 5, "ck101\t5\t1\n神經病戀愛指南\t小貓一尾\n2016-10-09\t第89章"},
		{"80txt", 6, "80txt\t6\t0\n最终进化\t卷土\n2013-07-19 14:11:54\t新书天择已上传，敬请光临。"},
		{"xqishu", 7, "xqishu\t7\t0\n神印王座\t唐家三少\n2017-07-27 12:06:10\t新书：斗罗大陆II《绝世唐门》开启"},
		{"hjwzw", 1890, "hjwzw\t1890\t0\n西游記\t吳承恩\n2013-08-02\t第一百回 徑回東土 五圣成真"},
		// {"bestory", 5, "ck101\t5\t1\n神經病戀愛指南\t小貓一尾\n2016-10-09\t第89章"},
	}
	for i := range testcases {
		t.Run(fmt.Sprintf("update_%v-%v", testcases[i].site, testcases[i].id), func(t * testing.T) {
			site, exist := testSites[testcases[i].site]
			if !exist { t.Fatalf("site %v not found", testcases[i].site) }
			book := site.Book(testcases[i].id, -1)
			book.Update()
			if book.String() != testcases[i].expected {
				t.Fatalf("site[\"%v\"].Book(%v).String() gives\n%v\nbut not\n%v",
					testcases[i].site, testcases[i].id, book.String(), testcases[i].expected)
			}
		})
	}
}
func TestIntegrate_LoadBookDownload(t *testing.T) {
	var testcases = []struct {
		site string
		id int
		expected bool
		actualFileName, expectedFileName string
	} {
		{"ck101", 8, true, "./test_res/8.txt", "./test_res/test-ck101-8.txt"},
		{"80txt", 6, true, "./test_res/6.txt", "./test_res/test-80txt-6.txt"},
		{"xqishu", 7, true, "./test_res/7.txt", "./test_res/test-xqishu-7.txt"},
		{"hjwzw", 1890, true, "./test_res/1890.txt", "./test_res/test-hjwzw-1890.txt"},
		// {"bestory", 9, true, "./test_res/9.txt", "./test_res/test-bestory-9.txt"},
	}
	for i := range testcases {
		t.Run(fmt.Sprintf("download_%v-%v", testcases[i].site, testcases[i].id), func(t * testing.T) {
			site, exist := testSites[testcases[i].site]
			if !exist { t.Fatalf("site %v not found", testcases[i].site) }
			book := site.Book(testcases[i].id, -1)
			book.Update()
			actual := book.Download("./test_res/", 1000)
			actualData, err := ioutil.ReadFile(testcases[i].actualFileName)
			helper.CheckError(err)
			expectedData, err := ioutil.ReadFile(testcases[i].expectedFileName)
			helper.CheckError(err)
			lenDiff := len(actualData) - len(expectedData)
			if actual != testcases[i].expected || lenDiff < 0 { lenDiff = -lenDiff }
			if lenDiff > 10 {
				t.Fatalf("site[\"%v\"].Book(%v).Download(\"./test_res\") gives\n" +
					"length = %v\nbut not\nlength = %v",
					testcases[i].site, testcases[i].id, len(actualData), len(expectedData))
			}
			os.Remove(testcases[i].actualFileName)
		})
	}
}