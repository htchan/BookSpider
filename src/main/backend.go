package main

import (
	"os"
	"fmt"
	"time"
	"sort"
	"io/ioutil"
	"runtime"
	"strings"
	"strconv"
	"net/http"
	"log"
	"path/filepath"
	//"encoding/json"
	//"golang.org/x/text/encoding/traditionalchinese"
	//"golang.org/x/text/encoding"

	//"../helper"
	"../model"
)

type Logs struct {
	logLocation string
	Logs []string
	MemoryLastUpdate, FileLastUpdate time.Time
	size int64
}

var stageFileName string

func download() {
	for name, site := range sites {
		if name == "80txt" {
			continue
		}
		currentProcess = name + "\tdownload"
		fmt.Println(name + "\tdownload")
		site.Download()
		runtime.GC()
	}
	currentProcess = ""
}
func update() {
	for name, site := range sites {
		currentProcess = name + "\tupdate"
		fmt.Println(name + "\tupdate")
		site.Update()
		runtime.GC()
	}
	currentProcess = ""
}
func explore(maxError int) {
	for name, site := range sites {
		currentProcess = name + "\texplore"
		fmt.Println(name + "\texplore")
		site.Explore(maxError)
		runtime.GC()
	}
	currentProcess = ""
}
func updateError() {
	for name, site := range sites {
		currentProcess = name + "\tupdate error"
		fmt.Println(name + "\tupdate error")
		site.UpdateError()
		runtime.GC()
	}
	currentProcess = ""
}
func check() {
	for name, site := range sites {
		currentProcess = name + "\tcheck"
		fmt.Println(name + "\tcheck")
		site.Check()
		runtime.GC()
	}
	currentProcess = ""
}
func checkEnd() {
	for name, site := range sites {
		currentProcess = name + "\tcheck end"
		fmt.Println(name + "\tcheck end")
		site.CheckEnd()
		runtime.GC()
	}
	currentProcess = ""
}
func backup() {
	for name, site := range sites {
		currentProcess = name + "\tbackup"
		fmt.Println(name + "\tbackup")
		site.Backup()
		runtime.GC()
	}
	currentProcess = ""
}
func fix() {
	for name, site := range sites {
		currentProcess = name + "\tfix"
		fmt.Println(name + "\tfix")
		site.Fix()
		runtime.GC()
	}
	currentProcess = ""
}

func Start(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	operation := req.URL.Query().Get("operation")
	if currentProcess != "" {
		res.WriteHeader(http.StatusLocked)
		fmt.Fprintf(res, "{\"code\" : 423, \"message\" : \"locked - <" + currentProcess + "> is operating, please wait until it finish\"}")
		return
	}
	switch strings.ToUpper(operation) {
	case "UPDATE":
		go update()
	case "EXPLORE":
		go explore(1000)
	case "DOWNLOAD":
		go download()
	case "ERROR":
		go updateError()
	case "CHECK":
		go check()
	case "CHECKEND":
		go checkEnd()
	case "BACKUP":
		go backup()
	case "FIX":
		go fix()
	default:
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(res, "{\"code\" : 404, \"message\" : \"not found - operation <"+operation+"> not found\"}")
		return
	}
	res.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(res, "{\"code\" : 202, \"message\" : \"<" + operation + "> is put to process\"}")
}

func (logs *Logs) update() {
	logFileStat, err := os.Stat(logs.logLocation)
	if err != nil {
		return
	}
	fileSize := logFileStat.Size()
	logs.FileLastUpdate = logFileStat.ModTime()
	if (fileSize - logs.size < 1000) && (time.Now().Sub(logs.MemoryLastUpdate).Seconds() > 3600) {
		return
	}
	logs.size = fileSize
	data, err := ioutil.ReadFile(logs.logLocation)
	if err != nil {
		return
	}
	dataStrings := strings.Split(string(data), "\n")
	min := len(logs.Logs)
	if (len(logs.Logs) > len(dataStrings)) {
		min = len(dataStrings)
	}
	for i := 1; i < min; i += 1 {
		logs.Logs[i] = dataStrings[len(dataStrings) - min + i]
	}
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
	fmt.Fprintf(res, site.JsonString())
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
		fmt.Fprintf(res, book.JsonString())
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
		fmt.Fprintf(res, book.JsonString())
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
		fmt.Fprintf(res, book.JsonString())
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
	logs = Logs{logLocation: dir + "/nohup.out", Logs: make([]string, 100), MemoryLastUpdate: time.Unix(0,0), FileLastUpdate: time.Unix(0, 0)}
	config := model.LoadYaml("./config/config.yaml")
	sites = model.LoadSitesYaml(config)
	apiFunc := make(map[string]func())
	apiFunc["control"] = func() { http.HandleFunc("/start", Start) }
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

<operation> := download, update, explore, update error, check, backup
// start a currentProcess
host:port/start?operation=<operation>
** if current currentProcess is not "", currentProcess cannot be start and return 403

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
