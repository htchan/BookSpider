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
	offset := fileSize - 1000
	if offset < 1000 {
		offset = 0
	}
	b := make([]byte, 1000)
	file.ReadAt(b, offset)
	logs.Logs = strings.Split(string(b), "\n")
	logs.MemoryLastUpdate = time.Now()
}

func ProcessState(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	/*
	f, err := os.Stat("nohup.out")
	// get create time, last update time of nohup.out
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(res, "{\"code\" : 404, \"message\" : \"no log found\"}")
		return
	}
	modifyTime := f.ModTime().String()[:19]
	*/
	modifyTime := logs.FileLastUpdate.String()[:19]
	// get last several line of nohup.out
	logs.update()
	// print them
	fmt.Fprintf(res, "{")
	fmt.Fprintf(res, "\"time\" : \"" + modifyTime + "\", ")
	fmt.Fprintf(res, "\"logs\" : [\n")
	for i, log := range logs.Logs {
		fmt.Fprintf(res, "\"" + strings.ReplaceAll(log, "\"", "\\\"") + "\"")
		if i < len(logs.Logs) - 1 {
			fmt.Fprintf(res, ",\n")
		}
	}
	fmt.Fprintf(res, "]}")
}

func ValidateState(res http.ResponseWriter, req *http.Request) {
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
	if err != nil {
		fmt.Fprintf(res, "{\"stage\" : \""+err.Error()+"\", ")
	} else {
		fmt.Fprintf(res, "{\"stage\" : \""+strings.ReplaceAll(string(data), "\n", "\\n")+"\", ")
	}
	siteNames := make([]string, 0)
	for siteName, _ := range sites {
		siteNames = append(siteNames, siteName)
	}
	sort.Strings(siteNames)
	fmt.Fprintf(res, "\"siteNames\" : [")
	for i, siteName := range siteNames {
		fmt.Fprintf(res, "\"" + siteName + "\"")
		if i < len(siteNames) - 1 {
			fmt.Fprintf(res, ", ")
		}
	}
	fmt.Fprintf(res, "]}")
}

func SiteInfo(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	uri := strings.Split(req.URL.Path, "/")
	siteName := uri[2]
	site, ok := sites[siteName]
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(res, "{\"code\" : 404, \"message\" : \"site <" + siteName + "> not found\"}")
		return
	}
	siteByte, err := json.Marshal(site.Map())
	helper.CheckError(err)
	fmt.Fprintf(res, string(siteByte))
}

func BookInfo(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	uri := strings.Split(req.URL.Path, "/")
	siteName := uri[2]
	site, ok := sites[siteName]
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(res, "{\"code\" : 404, \"message\" : \"site <" + siteName + "> not found\"}")
		return
	}
	id, err := strconv.Atoi(uri[3])
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(res, "{\"code\" : 400, \"message\" : \"id <" + uri[3] + "> is not a number\"}")
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
		fmt.Fprintf(res, "{\"code\" : 404, \"message\" : \"book <" + strconv.Itoa(id) + ">, version <" + strconv.Itoa(version) + "> in site <" + siteName + "> not found\"}")
	} else {
		bookByte, err := json.Marshal(book.Map())
		helper.CheckError(err)
		fmt.Fprintf(res, string(bookByte))
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
		fmt.Fprintf(res, "{\"code\" : 404, \"message\" : \"site <" + siteName + "> not found\"}")
		return
	}
	id, err := strconv.Atoi(uri[3])
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(res, "{\"code\" : 400, \"message\" : \"id <" + uri[3] + "> is not a number\"}")
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
		fmt.Fprintf(res, "{\"code\" : 404, \"message\" : \"book <" + strconv.Itoa(id) + ">, version <" + strconv.Itoa(version) + "> in site <" + siteName + "> not found\"}")
		return
	}
	if !book.DownloadFlag {
		res.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprintf(res, "{\"code\" : 406, \"message\" : \"book <" + uri[3] + "> not download yet\"}")
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
		fmt.Fprintf(res, "{\"code\" : 404, \"message\" : \"site <" + siteName + "> not found\"}")
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
	fmt.Fprintf(res, "{\"books\" : [")
	for i, book := range books {
		bookByte, err := json.Marshal(book.Map())
		helper.CheckError(err)
		fmt.Fprintf(res, string(bookByte))
		if i < len(books)-1 {
			fmt.Fprintf(res, ",")
		}
	}
	fmt.Fprintf(res, "]}")
}

func BookRandom(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	uri := strings.Split(req.URL.Path, "/")
	siteName := uri[2]
	site, ok := sites[siteName]
	if !ok {
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(res, "{\"code\" : 404, \"message\" : \"site <" + siteName + "> not found\"}")
		return
	}
	num, err := strconv.Atoi(req.URL.Query().Get("num"))
	if (err != nil) {
		num = 20;
	}
	if (num > 50) { num = 50 }
	books := site.Random(num)
	fmt.Fprintf(res, "{\"books\" : [")
	for i, book := range books {
		bookByte, err := json.Marshal(book.Map())
		helper.CheckError(err)
		fmt.Fprintf(res, string(bookByte))
		if i < len(books)-1 {
			fmt.Fprintf(res, ",")
		}
	}
	fmt.Fprintf(res, "]}")
}

var currentProcess string
var logs Logs
var sites map[string]model.Site

func main() () {
	currentProcess = ""
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	stageFileName = dir + "/log/stage.txt"
	logs = Logs{
		logLocation: dir + "/nohup.out", 
		Logs: make([]string, 100), 
		MemoryLastUpdate: time.Unix(0,0), 
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
