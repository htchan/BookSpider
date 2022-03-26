package configs

import (
	"gopkg.in/yaml.v2"
	// "encoding/json"
	"io/ioutil"
	"golang.org/x/text/encoding"
	"fmt"

	"github.com/htchan/BookSpider/internal/utils"
)

type BookConfig struct {
	BaseUrl string `yaml:"baseUrl"`
    DownloadUrl string `yaml:"downloadUrl"`
    ChapterUrl string `yaml:"chapterUrl"`
    ChapterUrlPattern string `yaml:"chapterUrlPattern"`
	TitleRegex string `yaml:"titleRegex"`
    WriterRegex string `yaml:"writerRegex"`
    TypeRegex string `yaml:"typeRegex"`
    LastUpdateRegex string `yaml:"lastUpdateRegex"`
    LastChapterRegex string `yaml:"lastChapterRegex"`
	ChapterUrlRegex string `yaml:"chapterUrlRegex"`
    ChapterTitleRegex string `yaml:"chapterTitleRegex"`
    ChapterContentRegex string `yaml:"chapterContentRegex"`
	Decoder *encoding.Decoder
	CONST_SLEEP int
	StorageDirectory string
	UseRequestInterval bool
}

func LoadBookConfigYaml(bookConfigLocation string) (bookConfig *BookConfig) {
	defer utils.Recover(func() { bookConfig = nil })
	bookConfig = new(BookConfig)

	data, err := ioutil.ReadFile(bookConfigLocation)
	utils.CheckError(err)
	utils.CheckError(yaml.Unmarshal(data, bookConfig))
	return
}

func LoadBookConfigJson(bookConfigLocation string) *BookConfig {
	return new(BookConfig)
}

func (config BookConfig)Populate(id int) BookConfig {
	config.BaseUrl = fmt.Sprintf(config.BaseUrl, id)
	config.DownloadUrl = fmt.Sprintf(config.DownloadUrl, id)
	return config
}