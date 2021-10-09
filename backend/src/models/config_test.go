package models

import (
	"testing"
	"log"
	"io/ioutil"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func (site Site) equal(target Site) (bool) {
	return site == target
}

func (config Config) equal(target Config) (bool) {
	if config.MaxThreads != target.MaxThreads || config.MaxExploreError != target.MaxExploreError ||
	config.Backend.StageFile != target.Backend.StageFile || config.Backend.LogFile != config.Backend.LogFile {
		return false
	}
	if len(config.Backend.Api) != len(target.Backend.Api) {
		return false
	}
	for i := range config.Backend.Api {
		if config.Backend.Api[i] != target.Backend.Api[i] {
			return false
		}
	}
	for siteName := range config.Sites {
		if _, exist := target.Sites[siteName]; !exist {
			return false
		}
		for key := range config.Sites[siteName] {
			if val, exist := target.Sites[siteName][key]; !exist || val != config.Sites[siteName][key] {
				return false
			}
		}
	}
	return true
}

func TestConfig_NewSiteYaml(t *testing.T) {
	actual := NewSiteYaml("testing", nil, "./test_res/test.yaml", "./", "./", 100)
	expected := Site{
		SiteName: "testing",
		database: nil,
		metaBaseUrl: "https://test/book/%v.html",
		metaDownloadUrl: "https://test/%v/",
		metaChapterUrl: "https://test/%v",
		chapterPattern: "/.*?\\.html",
		decoder: nil,
		titleRegex: "title regex",
		writerRegex: "writer regex",
		typeRegex: "type regex",
		lastUpdateRegex: "last update regex",
		lastChapterRegex: "last chapter regex",
		chapterUrlRegex: "chapter url regex",
		chapterTitleRegex: "chapter title regex",
		chapterContentRegex: "chapter content regex",
		databaseLocation: "./",
		DownloadLocation: "./",
		MAX_THREAD_COUNT: 100};
	if !actual.equal(expected) {
		t.Fatalf("helper.NewSiteYaml(\"testing\", nil, \"./test_res/test.yaml\"," +
			"\"./\", \"./\", 100) result gives\n\"%v\", but not\n\"%v\"\n",
			actual, expected)
	}
}

func TestConfig_LoadSiteYaml(t *testing.T) {
	config := map[string]string {
		"decode": "",
		"threadsCount": "100",
		"configLocation": "./test_res/test.yaml",
		"databaseLocation": "./",
		"downloadLocation": "./",
	}
	actual := LoadSiteYaml("testing", config)
	expected := Site{
		SiteName: "testing",
		database: nil,
		metaBaseUrl: "https://test/book/%v.html",
		metaDownloadUrl: "https://test/%v/",
		metaChapterUrl: "https://test/%v",
		chapterPattern: "/.*?\\.html",
		decoder: nil,
		titleRegex: "title regex",
		writerRegex: "writer regex",
		typeRegex: "type regex",
		lastUpdateRegex: "last update regex",
		lastChapterRegex: "last chapter regex",
		chapterUrlRegex: "chapter url regex",
		chapterTitleRegex: "chapter title regex",
		chapterContentRegex: "chapter content regex",
		databaseLocation: "./",
		DownloadLocation: "./",
		MAX_THREAD_COUNT: 100};
		if !actual.equal(expected) {
			t.Fatalf("helper.NewSiteYaml(\"testing\", nil, \"./test_res/test.yaml\"," +
				"\"./\", \"./\", 100) result gives\n\"%v\", but not\n\"%v\"\n",
				actual, expected)
		}
}

func TestConfig_LoadSitesYaml(t *testing.T) {
	config := Config{
		Sites: map[string]map[string]string{
			"testing": map[string]string{
				"decode" : "",
				"configLocation" : "./test_res/test.yaml",
				"databaseLocation" : "./",
				"downloadLocation" : "./",
				"threadsCount" : "100",
			},
		},
		MaxThreads: 100,
		MaxExploreError: 100,
		Backend: struct {
			Api []string `yaml:"api"`
			StageFile string `yaml:"stageFile"`
			LogFile string `yaml:"logFile"`
		} {
			Api: []string{"api1", "api2", "api3"},
			StageFile: "stage file",
			LogFile: "log file",
		},
	}
	actual := LoadSitesYaml(config)
	expected := map[string]Site{
		"testing":Site{
			SiteName: "testing",
			database: nil,
			metaBaseUrl: "https://test/book/%v.html",
			metaDownloadUrl: "https://test/%v/",
			metaChapterUrl: "https://test/%v",
			chapterPattern: "/.*?\\.html",
			decoder: nil,
			titleRegex: "title regex",
			writerRegex: "writer regex",
			typeRegex: "type regex",
			lastUpdateRegex: "last update regex",
			lastChapterRegex: "last chapter regex",
			chapterUrlRegex: "chapter url regex",
			chapterTitleRegex: "chapter title regex",
			chapterContentRegex: "chapter content regex",
			databaseLocation: "./",
			DownloadLocation: "./",
			MAX_THREAD_COUNT: 100},
	}
	for key := range expected {
		if _, ok := actual[key]; !ok {
			t.Fatalf("helper.LoadSitesYaml(%v) result gives\n\"%v\", but not\n\"%v\"\n",
				config, actual[key], expected[key])
		}
	}
	for key := range actual {
		if !actual[key].equal(expected[key]) {
			t.Fatalf("helper.LoadSitesYaml(%v) result gives\n\"%v\", but not\n\"%v\"\n",
				config, actual[key], expected[key])
		}
	}
}

func TestConfig_LoadYaml(t *testing.T) {
	actual := LoadYaml("./test_res/test_config.yaml")
	expected := Config{
		Sites: map[string]map[string]string{
			"testing": map[string]string{
				"decode" : "big5",
				"configLocation" : "./test_res/test.yaml",
				"databaseLocation" : "./",
				"downloadLocation" : "./",
				"threadsCount" : "100",
			},
		},
		MaxThreads: 1000,
		MaxExploreError: 1000,
		Backend: struct {
			Api []string `yaml:"api"`
			StageFile string `yaml:"stageFile"`
			LogFile string `yaml:"logFile"`
		} {
			Api: []string{"search", "download", "bookInfo", "siteInfo"},
			StageFile: "/log/stage.txt",
			LogFile: "/log/controller.log",
		},
	}
	if !actual.equal(expected) {
		t.Fatalf("helper.LoadYaml(\"%v\") result gives\n\"%v\", but not\n\"%v\"\n",
			"./test_res/test_config.yaml", actual, expected)
	}
}