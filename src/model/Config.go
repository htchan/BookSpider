package model

import (
	"gopkg.in/yaml.v2"
	"encoding/json"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/traditionalchinese"
	"io/ioutil"

	"../helper"
)

type Config struct {
	Sites map[string]map[string]string
	Api []string
}

func NewSiteYaml(siteName string, decoder *encoding.Decoder, configFileLocation string, databaseLocation string, downloadLocation string) (Site) {
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
		MetaBaseUrl: info["metaBaseUrl"],
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

func LoadSiteYaml(config map[string]string) (Site) {
	var decoder *encoding.Decoder
	if (config["decode"] == "big5") {
		decoder = traditionalchinese.Big5.NewDecoder()
	} else {
		decoder = nil
	}
	site := NewSiteYaml(config["name"], decoder,
					config["configLocation"],
					config["databaseLocation"],
					config["downloadLocation"])
	return site
}

func LoadSitesYaml(config Config) (map[string]Site) {
	sites := make(map[string]Site)
	for siteName, siteConfig := range config.Sites {
		sites[siteName] = LoadSiteYaml(siteConfig)
	}
	return sites
}

func LoadYaml(configFileLocation string) (Config) {
	var config Config
	data, err := ioutil.ReadFile(configFileLocation)
	helper.CheckError(err)
	helper.CheckError(yaml.Unmarshal(data, &config))
	return config
}

func NewSiteJson(siteName string, decoder *encoding.Decoder, configFileLocation string, databaseLocation string, downloadLocation string) (Site) {
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
		MetaBaseUrl: info["metaBaseUrl"],
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

func LoadSiteJson(config map[string]string) (Site) {
	var decoder *encoding.Decoder
	if (config["decode"] == "big5") {
		decoder = traditionalchinese.Big5.NewDecoder()
	} else {
		decoder = nil
	}
	site :=NewSiteJson(config["name"], decoder,
					config["configLocation"],
					config["databaseLocation"],
					config["downloadLocation"])
	return site
}
func LoadSitesJson(config Config) (map[string]Site) {
	sites := make(map[string]Site)
	for siteName, config := range config.Sites {
		sites[siteName] = LoadSiteJson(config)
	}
	return sites
}

func LoadJson(configFileLocation string) (Config) {
	var config Config
	data, err := ioutil.ReadFile(configFileLocation)
	helper.CheckError(err)
	helper.CheckError(json.Unmarshal(data, &config))
	return config
}