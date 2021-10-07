package main

import (
	"os"
	"fmt"
	"time"
	"sort"
	"io/ioutil"
	"strings"
	"strconv"
	"net/http"
	"log"
	//"encoding/json"
	//"golang.org/x/text/encoding/traditionalchinese"
	//"golang.org/x/text/encoding"
	"encoding/json"

	//"../helper"
	"github.com/htchan/BookSpider/models"
	"github.com/htchan/BookSpider/helper"
)

type Logs struct {
	logLocation string
	Logs []string
	MemoryLastUpdate, FileLastUpdate time.Time
	size int64
}

var stageFileName string
const readLogLen int64 = 10000

func (logs *Logs) update() {
	if time.Now().Unix() - logs.MemoryLastUpdate.Unix() < 60 { return }

	logFileStat, err := os.Stat(logs.logLocation)
	if err != nil { return }
	fileSize := logFileStat.Size()
	logs.FileLastUpdate = logFileStat.ModTime()
	if fileSize == logs.size { return }

	file, err := os.Open(logs.logLocation)
	helper.CheckError(err)
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

func response(res http.ResponseWriter, data map[string]interface{}) {
	dataByte, err := json.Marshal(data)
	helper.CheckError(err)
	fmt.Fprintln(res, string(dataByte))
}
func error(res http.ResponseWriter, code int, msg string) {
	res.WriteHeader(code)
	response(res, map[string]interface{} {
		"code": code,
		"message": msg,
	})
}

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

func ValidateState(res http.ResponseWriter, req *http.Request) {
	//TODO: conside to turn this json to a yaml array to show what site is working
	//TODO: maybe also put it into config.yaml
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	b, err := ioutil.ReadFile("./validate.json")
	if err != nil {
		fmt.Fprintf(res, "{}")
		return
	}
	fmt.Fprintf(res, string(b))
}

func GeneralInfo(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	siteNames := make([]string, 0)
	for siteName, _ := range sites {
		siteNames = append(siteNames, siteName)
	}
	sort.Strings(siteNames)

	response(res, map[string]interface{} {
		"siteNames": siteNames,
	})
}

func SiteInfo(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	uri := strings.Split(req.URL.Path, "/")
	siteName := uri[4]
	site, ok := sites[siteName]
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		error(res, http.StatusNotFound, "site <" + siteName + "> not found")
		return
	}

	response(res, site.Map())
}

func BookInfo(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	uri := strings.Split(req.URL.Path, "/")
	if len(uri) < 6 {
		error(res, http.StatusBadRequest, "not enough parameter")
	}
	siteName := uri[4]
	site, ok := sites[siteName]
	id, err := strconv.Atoi(uri[5])
	if !ok {
		error(res, http.StatusNotFound, "site <" + siteName + "> not found")
		return
	} else if err != nil {
		error(res, http.StatusBadRequest, "id <" + uri[5] + "> is not a number")
		return
	}

	version := -1;
	if len(uri) > 6 {
		version, err = strconv.Atoi(uri[6])
		if err != nil { version = -1 }
	}

	book := site.Book(id, version)
	if (book.Title == "") {
		error(res, http.StatusNotFound, "book <" + strconv.Itoa(id) + ">, " +
				"version <" + strconv.Itoa(version) + "> in site <" + siteName + "> not found")
	} else {
		response(res, book.Map())
	}
}

func BookDownload(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	uri := strings.Split(req.URL.Path, "/")
	if len(uri) < 6 {
		error(res, http.StatusBadRequest, "not enough parameter")
	}
	siteName := uri[4]
	site, ok := sites[siteName]
	id, err := strconv.Atoi(uri[5])
	if !ok {
		error(res, http.StatusNotFound, "site <" + siteName + "> not found")
		return
	} else if err != nil {
		error(res, http.StatusBadRequest, "id <" + uri[5] + "> is not a number")
		return
	}
	version := -1;
	if len(uri) > 6 {
		version, err = strconv.Atoi(uri[6])
		if err != nil {
			version = -1
		}
	}
	book := site.Book(id, version)
	if (book.Title == "") {
		error(res, http.StatusNotFound, "book <" + strconv.Itoa(id) + ">, " +
				"version <" + strconv.Itoa(version) + "> in site <" + siteName + "> not found")
		return
	} else if !book.DownloadFlag {
		error(res, http.StatusNotAcceptable, "book <" + uri[5] + "> not download yet")
		return
	}
	fileName := book.Title + "-" + book.Writer
	if book.Version > 0 {
		fileName += "-v" + strconv.Itoa(book.Version)
	}
	content := book.Content(site.DownloadLocation)
	res.Header().Set("Content-Type", "text/txt; charset=utf-8")
	res.Header().Set("Content-Disposition", "attachment; filename=\"" + fileName + ".txt\"")
	fmt.Fprintf(res, content)
}

func BookSearch(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	uri := strings.Split(req.URL.Path, "/")
	siteName := uri[4]
	site, ok := sites[siteName]
	title := strings.ReplaceAll(req.URL.Query().Get("title"), "*", "%")
	writer := strings.ReplaceAll(req.URL.Query().Get("writer"), "*", "%")
	pageStr := req.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if !ok {
		error(res, http.StatusNotFound, "site <" + siteName + "> not found")
		return
	} else if err != nil {
		page = 0
	}
	books := site.Search(title, writer, page)
	booksArray := make([]map[string]interface{}, 0)
	for _, book := range books {
		booksArray = append(booksArray, book.Map())
	}
	response(res, map[string]interface{} {
		"books": booksArray,
	})
}

func BookRandom(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	uri := strings.Split(req.URL.Path, "/")
	siteName := uri[4]
	site, ok := sites[siteName]
	if !ok {
		error(res, http.StatusNotFound, "site <" + siteName + "> not found")
		return
	}
	num, err := strconv.Atoi(req.URL.Query().Get("num"))
	if (err != nil) { num = 20 }
	if (num > 50) { num = 50 }
	books := site.RandomSuggestBook(num)
	booksArray := make([]map[string]interface{}, 0)
	for _, book := range books {
		booksArray = append(booksArray, book.Map())
	}
	response(res, map[string]interface{} {
		"books": booksArray,
	})
}

var currentProcess string
var logs Logs
var sites map[string]model.Site

func main() {
	currentProcess = ""
	config := model.LoadYaml("./config/config.yaml")
	stageFileName = config.Backend.StageFile
	logs = Logs{
		logLocation: config.Backend.LogFile, 
		Logs: make([]string, 100), 
		MemoryLastUpdate: time.Unix(0, 0), 
		FileLastUpdate: time.Unix(0, 0)}
	sites = model.LoadSitesYaml(config)
	apiFunc := make(map[string]func())
	apiFunc["search"] = func() { for name := range sites { http.HandleFunc("/api/novel/search/"+name+"", BookSearch) } }
	apiFunc["download"] = func() { for name := range sites { http.HandleFunc("/api/novel/download/"+name+"/", BookDownload) } }
	apiFunc["siteInfo"] = func() { for name := range sites { http.HandleFunc("/api/novel/sites/"+name, SiteInfo) } }
	apiFunc["bookInfo"] = func() { for name := range sites { http.HandleFunc("/api/novel/books/"+name+"/", BookInfo) } }
	apiFunc["random"] = func() { for name := range sites { http.HandleFunc("/api/novel/random/"+name, BookRandom) } }
	apiFunc["process"] = func() { http.HandleFunc("/api/novel/process", ProcessState) }
	apiFunc["info"] = func() { http.HandleFunc("/api/novel/info", GeneralInfo) }
	apiFunc["validate"] = func() { http.HandleFunc("/api/novel/validate", ValidateState) }

	for _, api := range config.Backend.Api { apiFunc[api]() }
	log.Println("started")
	log.Fatal(http.ListenAndServe(":9427", nil))
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
