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
	"path/filepath"
	//"encoding/json"
	//"golang.org/x/text/encoding/traditionalchinese"
	//"golang.org/x/text/encoding"
	"encoding/json"

	//"../helper"
	"../model"
	"../helper"
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
	if time.Now().Unix() - logs.MemoryLastUpdate.Unix() < 60 {
		return
	}
	logFileStat, err := os.Stat(logs.logLocation)
	if err != nil {
		return
	}
	fileSize := logFileStat.Size()
	logs.FileLastUpdate = logFileStat.ModTime()
	if fileSize == logs.size {
		return
	}
	file, err := os.Open(logs.logLocation)
	helper.CheckError(err)
	defer file.Close()
	logs.size = fileSize
	offset := fileSize - readLogLen
	if offset < readLogLen {
		offset = 0
	}
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

func ProcessState(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	modifyTime := logs.FileLastUpdate.String()[:19]
	// get last several line of nohup.out
	logs.update()
	// print them
	response(res, map[string]interface{} {
		"time": modifyTime,
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
	data, err := ioutil.ReadFile(stageFileName)
	var stageStr string
	if err != nil {
		stageStr = err.Error()
	} else {
		stageStr = strings.ReplaceAll(string(data), "\n", "\\n")
	}
	siteNames := make([]string, 0)
	for siteName, _ := range sites {
		siteNames = append(siteNames, siteName)
	}
	sort.Strings(siteNames)
	response(res, map[string]interface{} {
		"stage": stageStr,
		"siteNames": siteNames,
	})
}

func SiteInfo(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	uri := strings.Split(req.URL.Path, "/")
	siteName := uri[2]
	site, ok := sites[siteName]
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		// fmt.Fprintf(res, "{\"code\" : 404, \"message\" : \"site <" + siteName + "> not found\"}")
		response(res, map[string]interface{} {
			"code": 404,
			"message": "site <" + siteName + "> not found",
		})
		return
	}
	response(res, site.Map())
}

func BookInfo(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	uri := strings.Split(req.URL.Path, "/")
	siteName := uri[2]
	site, ok := sites[siteName]
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		// fmt.Fprintf(res, "{\"code\" : 404, \"message\" : \"site <" + siteName + "> not found\"}")
		response(res, map[string]interface{} {
			"code": 404,
			"message": "site <" + siteName + "> not found",
		})
		return
	}
	id, err := strconv.Atoi(uri[3])
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		// fmt.Fprintf(res, "{\"code\" : 400, \"message\" : \"id <" + uri[3] + "> is not a number\"}")
		response(res, map[string]interface{} {
			"code": 400,
			"message": "id <" + uri[3] + "> is not a number",
		})
		return
	}
	version := -1;
	if len(uri) > 4 {
		version, err =strconv.Atoi(uri[4])
		if err != nil {
			version = -1
		}
	}
	book := site.Book(id, version)
	if (book.Title == "") {
		res.WriteHeader(http.StatusNotFound)
		// fmt.Fprintf(res, "{\"code\" : 404, \"message\" : \"book <" + strconv.Itoa(id) + ">, version <" + strconv.Itoa(version) + "> in site <" + siteName + "> not found\"}")
		response(res, map[string]interface{} {
			"code": 404,
			"message": "book <" + strconv.Itoa(id) + ">, version <" + strconv.Itoa(version) + "> in site <" + siteName + "> not found",
		})
	} else {
		response(res, book.Map())
	}
}

func BookDownload(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	uri := strings.Split(req.URL.Path, "/")
	siteName := uri[2]
	site, ok := sites[siteName]
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		// fmt.Fprintf(res, "{\"code\" : 404, \"message\" : \"site <" + siteName + "> not found\"}")
		response(res, map[string]interface{} {
			"code": 404,
			"message": "site <" + siteName + "> not found",
		})
		return
	}
	id, err := strconv.Atoi(uri[3])
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		// fmt.Fprintf(res, "{\"code\" : 400, \"message\" : \"id <" + uri[3] + "> is not a number\"}")
		response(res, map[string]interface{} {
			"code": 400,
			"message": "id <" + uri[3] + "> is not a number",
		})
		return
	}
	version := -1;
	if len(uri) > 4 {
		version, err = strconv.Atoi(uri[4])
		if err != nil {
			version = -1
		}
	}
	book := site.Book(id, version)
	if (book.Title == "") {
		res.WriteHeader(http.StatusNotFound)
		// fmt.Fprintf(res, "{\"code\" : 404, \"message\" : \"book <" + strconv.Itoa(id) + ">, version <" + strconv.Itoa(version) + "> in site <" + siteName + "> not found\"}")
		response(res, map[string]interface{} {
			"code": 404,
			"message": "book <" + strconv.Itoa(id) + ">, version <" + strconv.Itoa(version) + "> in site <" + siteName + "> not found",
		})
		return
	}
	if !book.DownloadFlag {
		res.WriteHeader(http.StatusNotAcceptable)
		// fmt.Fprintf(res, "{\"code\" : 406, \"message\" : \"book <" + uri[3] + "> not download yet\"}")
		response(res, map[string]interface{} {
			"code": 406,
			"message": "book <" + uri[3] + "> not download yet",
		})
		return
	}
	fileName := book.Title + "-" + book.Writer
	if book.Version > 0 {
		fileName += "-v" + strconv.Itoa(book.Version)
	}
	content := site.BookContent(book)
	res.Header().Set("Content-Type", "text/txt; charset=utf-8")
	res.Header().Set("Content-Disposition", "attachment; filename=\"" + fileName + ".txt\"")
	fmt.Fprintf(res, content)
}

func BookSearch(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	uri := strings.Split(req.URL.Path, "/")
	siteName := uri[2]
	site, ok := sites[siteName]
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		// fmt.Fprintf(res, "{\"code\" : 404, \"message\" : \"site <" + siteName + "> not found\"}")
		response(res, map[string]interface{} {
			"code": 404,
			"message": "site <" + siteName + "> not found",
		})
		return
	}
	title := strings.ReplaceAll(req.URL.Query().Get("title"), "*", "%")
	writer := strings.ReplaceAll(req.URL.Query().Get("writer"), "*", "%")
	pageStr := req.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
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
	siteName := uri[2]
	site, ok := sites[siteName]
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		// fmt.Fprintf(res, "{\"code\" : 404, \"message\" : \"site <" + siteName + "> not found\"}")
		response(res, map[string]interface{} {
			"code": 404,
			"message": "site <" + siteName + "> not found",
		})
		return
	}
	num, err := strconv.Atoi(req.URL.Query().Get("num"))
	if (err != nil) {
		num = 20;
	}
	if (num > 50) { num = 50 }
	books := site.Random(num)
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
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	stageFileName = dir + "/log/stage.txt"
	logs = Logs{
		logLocation: dir + "/nohup.out", 
		Logs: make([]string, 100), 
		MemoryLastUpdate: time.Unix(0, 0), 
		FileLastUpdate: time.Unix(0, 0)}
	config := model.LoadYaml("./config/config.yaml")
	sites = model.LoadSitesYaml(config)
	apiFunc := make(map[string]func())
	apiFunc["search"] = func() { for name := range sites { http.HandleFunc("/search/"+name+"", BookSearch) } }
	apiFunc["download"] = func() { for name := range sites { http.HandleFunc("/download/"+name+"/", BookDownload) } }
	apiFunc["siteInfo"] = func() { for name := range sites { http.HandleFunc("/info/"+name, SiteInfo) } }
	apiFunc["bookInfo"] = func() { for name := range sites { http.HandleFunc("/info/"+name+"/", BookInfo) } }
	apiFunc["random"] = func() { for name := range sites { http.HandleFunc("/random/"+name, BookRandom) } }
	apiFunc["process"] = func() { http.HandleFunc("/process", ProcessState) }
	apiFunc["info"] = func() { http.HandleFunc("/info", GeneralInfo) }
	apiFunc["validate"] = func() { http.HandleFunc("/validate", ValidateState) }

	for _, api := range config.Api { apiFunc[api]() }
	fmt.Println("started")
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
