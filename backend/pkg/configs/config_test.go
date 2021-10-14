package configs

import (
	"testing"
	"log"
	"io/ioutil"

	"github.com/htchan/BookSpider/pkg/sites"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func mapEqual(map1, map2 map[string]sites.Site) bool {
	for key := range map1 {
		if _, exist := map2[key]; !exist {
			return false
		}
	}
	for key := range map2 {
		if val, exist := map2[key]; !exist || equalSite(val, map1[key]) {
			return false
		}
	}
	return true
}

func equalSite(site1, site2 sites.Site) (bool) {
	return site1 == site2
}

func equalConfig(config1, config2 Config) (bool) {
	if config1.MaxThreads != config2.MaxThreads ||
		config1.MaxExploreError != config2.MaxExploreError ||
		config1.Backend.StageFile != config2.Backend.StageFile ||
		config1.Backend.LogFile != config1.Backend.LogFile {
		return false
	}
	if len(config1.Backend.Api) != len(config2.Backend.Api) {
		return false
	}
	for i := range config1.Backend.Api {
		if config1.Backend.Api[i] != config2.Backend.Api[i] {
			return false
		}
	}
	for siteName := range config1.Sites {
		if _, exist := config2.Sites[siteName]; !exist {
			return false
		}
		for key := range config1.Sites[siteName] {
			if val, exist := config2.Sites[siteName][key]; !exist || val != config1.Sites[siteName][key] {
				return false
			}
		}
	}
	return true
}

var expectMeta = map[string]string{
	"baseUrl": "https://test/book/%v.html",
	"downloadUrl": "https://test/%v/",
	"chapterUrl": "https://test/%v",
	"chapterUrlPattern": "/.*?\\.html",
	"titleRegex": "title regex",
	"writerRegex": "writer regex",
	"typeRegex": "type regex",
	"lastUpdateRegex": "last update regex",
	"lastChapterRegex": "last chapter regex",
	"chapterUrlRegex": "chapter url regex",
	"chapterTitleRegex": "chapter title regex",
	"chapterContentRegex": "chapter content regex",
}

func TestLoadSiteYaml(t *testing.T) {
	config := map[string]string {
		"decode": "",
		"threadsCount": "100",
		"configFileLocation": "../../test/config-test-data/test.yaml",
		"databaseLocation": "../../test/config-test-data",
		"downloadLocation": "../../test/config-test-data",
	}
	actual := LoadSiteYaml("testing", config)
	expected, _ := sites.LoadSite("testing", config, expectMeta)
	if !equalSite(actual, *expected) {
		t.Fatalf("utils.NewSiteYaml(\"testing\", nil, \"./test_res/test.yaml\"," +
			"\"./\", \"./\", 100) result gives\n\"%v\", but not\n\"%v\"\n",
			actual, *expected)
	}
}

func TestLoadConfigYaml(t *testing.T) {
	config := LoadConfigYaml("../../test/config-test-data/test_config.yaml")
	expectConfig := Config{
		Sites: map[string]map[string]string{
			"testing": map[string]string{
				"decode" : "big5",
				"configFileLocation" : "../../test/config-test-data/test.yaml",
				"databaseLocation" : "../../test/config-test-data/",
				"downloadLocation" : "../../test/config-test-data/",
				"threadsCount" : "100",
			},
		},
		MaxThreads: 1000,
		MaxExploreError: 1000 ,
		Backend: struct {
			Api []string `yaml:"api"`
			StageFile string `yaml:"stageFile"`
			LogFile string `yaml:"logFile"`
		} {
			Api: []string {"search", "download", "bookInfo", "siteInfo"},
			StageFile: "/log/stage.txt",
			LogFile: "/log/controller.log",
		},
	}
	if !equalConfig(config, expectConfig) {
		t.Fatalf("LoadConfigYaml failed, return\n%v but not\n%v", config, expectConfig)
	}
}

func TestLoadSitesYaml(t *testing.T) {
	siteConfig := map[string]string {
		"decode": "big5",
		"threadsCount": "100",
		"configFileLocation": "../../test/config-test-data/test.yaml",
		"databaseLocation": "../../test/config-test-data",
		"downloadLocation": "../../test/config-test-data",
	}
	config := LoadConfigYaml("../../test/config-test-data/test_config.yaml")
	actualSites := LoadSitesYaml(config)
	expectSite, _ := sites.LoadSite("testing", siteConfig, expectMeta)
	expectSites := map[string]sites.Site{
		"testing": *expectSite,
	}
	if !mapEqual(actualSites, expectSites) {
		t.Fatalf("LoadSitesYaml fail return\n%v but not\n%v", actualSites, expectSites)
	}
}
