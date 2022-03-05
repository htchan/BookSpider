package main

import (
	"os"
	"time"
	"io/ioutil"
	"strings"
	"net/http"
	// "log"
	
	"encoding/json"

	"github.com/htchan/BookSpider/pkg/sites"
	"github.com/htchan/BookSpider/pkg/configs"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/logging"
)

type Logs struct {
	logLocation string
	Logs []string
	MemoryLastUpdate, FileLastUpdate time.Time
	size int64
}

var stageFileName string
const readLogLen int64 = 10000
const RECORD_PER_PAGE = 50
var apiFunc = map[string]func() {
	"info": func() { http.HandleFunc("/api/novel/info", GeneralInfo) },
	"siteInfo": func() { http.HandleFunc("/api/novel/sites/", SiteInfo) },
	"bookInfo": func() { http.HandleFunc("/api/novel/books/", BookInfo) },

	"download": func() { for name := range siteMap { http.HandleFunc("/api/novel/download/"+name+"/", BookDownload) } },
	"search": func() { for name := range siteMap { http.HandleFunc("/api/novel/search/"+name+"", BookSearch) } },
	"random": func() { for name := range siteMap { http.HandleFunc("/api/novel/random/"+name, BookRandom) } },

	"process": func() { http.HandleFunc("/api/novel/process", ProcessState) },
}

func (logs *Logs) update() {
	if time.Now().Unix() - logs.MemoryLastUpdate.Unix() < 60 { return }

	logFileStat, err := os.Stat(logs.logLocation)
	if err != nil { return }
	fileSize := logFileStat.Size()
	logs.FileLastUpdate = logFileStat.ModTime()
	if fileSize == logs.size { return }

	file, err := os.Open(logs.logLocation)
	utils.CheckError(err)
	defer file.Close()
	logs.size = fileSize
	offset := fileSize - readLogLen
	if offset < readLogLen { offset = 0 }
	
	b := make([]byte, readLogLen)
	file.ReadAt(b, offset)
	logs.Logs = strings.Split(string(b), "\n")
	logs.Logs = logs.Logs[1:]
	for i, _ := range logs.Logs {
		logs.Logs[i] = strings.ReplaceAll(logs.Logs[i], "\"", "\\\"")
	}
	logs.MemoryLastUpdate = time.Now()
}
func setHeader(res http.ResponseWriter) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
}
func response(res http.ResponseWriter, data map[string]interface{}) {
	encoder := json.NewEncoder(res)
	utils.CheckError(encoder.Encode(data))
}
func error(res http.ResponseWriter, code int, msg string) {
	res.WriteHeader(code)
	response(res, map[string]interface{} {
		"code": code,
		"message": msg,
	})
}
//TODO
func ProcessState(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	
	logs.update()
	data, err := ioutil.ReadFile(stageFileName)
	var stageStr []string
	if err != nil {
		stageStr = append(stageStr, err.Error())
	} else {
		stageStr = append(stageStr, strings.Split(string(data), "\n")...)
	}

	response(res, map[string]interface{} {
		"time": logs.FileLastUpdate.Unix(),
		"stage": stageStr,
		"logs": logs.Logs,
	})
}

var currentProcess string
var logs Logs
var config *configs.Config
var siteMap map[string]*sites.Site

func setup(configFileLocation string) {
	currentProcess = ""
	config = configs.LoadConfigYaml(configFileLocation)
	stageFileName = config.Backend.StageFile
	siteMap = make(map[string]*sites.Site)
	for key, siteConfig := range config.SiteConfigs {
		siteMap[key] = sites.NewSite(key, siteConfig)
		siteMap[key].OpenDatabase()
		//TODO: deploy a thread to close the database if it is not opened
	}
}
func startServer(addr string) {
	for _, api := range config.Backend.Api { apiFunc[api]() }
	logging.Info("started")
	logging.Error("%v", http.ListenAndServe(addr, nil))
}
func main() {
	setup(os.Getenv("ASSETS_LOCATION") + "/configs/config.yaml")
	logs = Logs{
		logLocation: config.Backend.LogFile, 
		Logs: make([]string, 100), 
		MemoryLastUpdate: time.Unix(0, 0), 
		FileLastUpdate: time.Unix(0, 0)}

	startServer(":9427")
}

/*
// get info of server
host:port

// get currentProcess info
host:port/currentProcess

// get info of site
host:port/info/<site>

// get info of book (all versions)
host:port/<site>/<id>
host:port/info/<site>/<id>

// get info of book
host:port/<site>/<id>/<version>
host:port/info/<site>/<id>/<version>

// download book
host:port/download/<site>/<id>

// search books
host:port/search?title=???&writer=???
*/
