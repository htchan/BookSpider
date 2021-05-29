package model

import (
	"gopkg.in/yaml.v2"
	"encoding/json"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/traditionalchinese"
	"io/ioutil"
	"strconv"

	"github.com/htchan/BookSpider/helper"
)

type Config struct {
	Sites map[string]map[string]string `yaml:"sites"`
	MaxThreads int `yaml:"maxThreads"`
	MaxExploreError int `yaml:"maxExploreError"`
	Backend struct {
		Api []string `yaml:"api"`
		StageFile string `yaml:"stageFile"`
		LogFile string `yaml:"logFile"`
	}
}

func NewSiteYaml(siteName string, decoder *encoding.Decoder, configFileLocation string, 
	databaseLocation string, downloadLocation string, MAX_THREAD_COUNT int) Site {
	//database, err := sql.Open("sqlite3", databaseLocation)
	//helper.CheckError(err);
	//database.SetMaxIdleConns(10);
	//database.SetMaxOpenConns(99999);
	data, err := ioutil.ReadFile(configFileLocation)
	helper.CheckError(err);
	var info map[string]string;
	helper.CheckError(yaml.Unmarshal(data, &info))
	site := Site{
		SiteName: siteName,
		database: nil,
		metaBaseUrl: info["metaBaseUrl"],
		metaDownloadUrl: info["metaDownloadUrl"],
		metaChapterUrl: info["metaChapterUrl"],
		chapterPattern: info["chapterPattern"],
		decoder: decoder,
		titleRegex: info["titleRegex"],
		writerRegex: info["writerRegex"],
		typeRegex: info["typeRegex"],
		lastUpdateRegex: info["lastUpdateRegex"],
		lastChapterRegex: info["lastChapterRegex"],
		chapterUrlRegex: info["chapterUrlRegex"],
		chapterTitleRegex: info["chapterTitleRegex"],
		chapterContentRegex: info["chapterContentRegex"],
		databaseLocation: databaseLocation,
		downloadLocation: downloadLocation,
		MAX_THREAD_COUNT: MAX_THREAD_COUNT};
	return site;
}

func LoadSiteYaml(siteName string, config map[string]string) Site {
	var decoder *encoding.Decoder
	if (config["decode"] == "big5") {
		decoder = traditionalchinese.Big5.NewDecoder()
	} else {
		decoder = nil
	}
	MAX_THREAD_COUNT, _ := strconv.Atoi(config["threadsCount"])
	site := NewSiteYaml(siteName, decoder,
					config["configLocation"],
					config["databaseLocation"],
					config["downloadLocation"],
					MAX_THREAD_COUNT)
	return site
}

func LoadSitesYaml(config Config) map[string]Site {
	sites := make(map[string]Site)
	for siteName, siteConfig := range config.Sites {
		sites[siteName] = LoadSiteYaml(siteName, siteConfig)
	}
	return sites
}

func LoadYaml(configFileLocation string) Config {
	var config Config
	data, err := ioutil.ReadFile(configFileLocation)
	helper.CheckError(err)
	helper.CheckError(yaml.Unmarshal(data, &config))
	return config
}

func NewSiteJson(siteName string, decoder *encoding.Decoder, configFileLocation string, 
	databaseLocation string, downloadLocation string) Site {
	//database, err := sql.Open("sqlite3", databaseLocation)
	//helper.CheckError(err);
	//database.SetMaxIdleConns(10);
	//database.SetMaxOpenConns(99999);
	data, err := ioutil.ReadFile(configFileLocation)
	helper.CheckError(err);
	var info map[string]string;
	helper.CheckError(json.Unmarshal(data, &info))
	site := Site{
		SiteName: siteName,
		//database: database,
		database: nil,
		metaBaseUrl: info["metaBaseUrl"],
		metaDownloadUrl: info["metaDownloadUrl"],
		metaChapterUrl: info["metaChapterUrl"],
		chapterPattern: info["chapterPattern"],
		decoder: decoder,
		titleRegex: info["titleRegex"],
		writerRegex: info["writerRegex"],
		typeRegex: info["typeRegex"],
		lastUpdateRegex: info["lastUpdateRegex"],
		lastChapterRegex: info["lastChapterRegex"],
		chapterUrlRegex: info["chapterUrlRegex"],
		chapterTitleRegex: info["chapterTitleRegex"],
		chapterContentRegex: info["chapterContentRegex"],
		databaseLocation: databaseLocation,
		downloadLocation: downloadLocation};
	return site;
}

func LoadSiteJson(siteName string, config map[string]string) Site {
	var decoder *encoding.Decoder
	if (config["decode"] == "big5") {
		decoder = traditionalchinese.Big5.NewDecoder()
	} else {
		decoder = nil
	}
	site :=NewSiteJson(siteName, decoder,
					config["configLocation"],
					config["databaseLocation"],
					config["downloadLocation"])
	return site
}
func LoadSitesJson(config Config) map[string]Site {
	sites := make(map[string]Site)
	for siteName, config := range config.Sites {
		sites[siteName] = LoadSiteJson(siteName, config)
	}
	return sites
}

func LoadJson(configFileLocation string) Config {
	var config Config
	data, err := ioutil.ReadFile(configFileLocation)
	helper.CheckError(err)
	helper.CheckError(json.Unmarshal(data, &config))
	return config
}