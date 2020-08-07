package main

import (
	"os"
	"fmt"
	"time"
	"io/ioutil"
	"runtime"
	"strings"
	"strconv"
	"net/http"
	"log"
	"encoding/json"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding"

	"./helper"
	"./model"
)

type Logs struct {
	Logs []string
	LastUpdate time.Time
}

func LoadSites(configLocation string) (map[string]model.Site) {
	sites := make(map[string]model.Site)
	data, err := ioutil.ReadFile(configLocation)
	helper.CheckError(err)
	var info []map[string]interface{}
	if err = json.Unmarshal(data, &info); err != nil {
        panic(err)
	}
	for _, config := range info {
		var decoder *encoding.Decoder
		if (config["decode"] == "big5") {
			decoder = traditionalchinese.Big5.NewDecoder()
		} else {
			decoder = nil
		}
		sites[config["name"].(string)] =
			model.NewSite(config["name"].(string), decoder,
							config["configLocation"].(string),
							config["databaseLocation"].(string),
							config["downloadLocation"].(string))
	}
	return sites
}

func download() {
	for name, site := range sites {
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
	fmt.Fprintf(res, "{\"code\" : 202, \"message\" : \"<" + operation + " is put to process\"}")
}

func (logs *Logs) update() {
	if (time.Now().Unix() - logs.LastUpdate.Unix() > 600) {
		data, err := ioutil.ReadFile("./nohup.out")
		if err != nil {
			return
		}
		dataStrings := strings.Split(string(data), "\n")
		min := len(logs.Logs)
		if (len(logs.Logs) > len(data)) {
			min = len(data)
		}
		for i := 1; i < min; i += 1 {
			logs.Logs[i] = dataStrings[len(dataStrings) - min + i]
		}
		logs.LastUpdate = time.Now()
	}
}

func ProcessState(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	f, err := os.Stat("nohup.out")
	// get create time, last update time of nohup.out
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(res, "{\"code\" : 404, \"message\" : \"no log found\"}")
		return
	}
	modifyTime := f.ModTime().String()[:19]
	// get last several line of nohup.out
	logs.update()
	// print them
	fmt.Fprintf(res, "{")
	fmt.Fprintf(res, "\"time\" : \"" + modifyTime + "\", ")
	fmt.Fprintf(res, "\"logs\" : [\n")
	for i, log := range logs.Logs {
		fmt.Fprintf(res, "\"" + log + "\"")
		if i < len(logs.Logs) {
			fmt.Fprintf(res, ",\n")
		}
	}
	fmt.Fprintf(res, "]}")
}

func GeneralInfo(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(res, "{\"currentProcess\" : \""+currentProcess+"\", ")
	fmt.Fprintf(res, "\"siteNames\" : [")
	var i int
	for siteName, _ := range sites {
		fmt.Fprintf(res, "\""+siteName+"\"")
		if i < len(sites)-1 {
			fmt.Fprintf(res, ", ")
		}
		i += 1
	}
	fmt.Fprintf(res, "]}")
}

func SiteInfo(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
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
	fmt.Println(id)
	book := site.Book(id)
	if (book.Title == "") {
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(res, "{\"code\" : 404, \"message\" : \"book <" + strconv.Itoa(id) + "> in site <" + siteName + "> not found\"}")
	} else {
		fmt.Fprintf(res, book.JsonString())
	}
}

func BookDownload(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
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
	book := site.Book(id)
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

func bookSearch(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
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

var currentProcess string
var logs Logs
var sites map[string]model.Site

func main() () {
	currentProcess = ""
	logs = Logs{Logs: make([]string, 100), LastUpdate: time.Unix(0, 0)}
	sites = LoadSites("./config/config.json")
	http.HandleFunc("/start", Start)
	for name, _ := range sites {
		http.HandleFunc("/search/"+name+"", bookSearch)
		http.HandleFunc("/download/"+name+"/", BookDownload)
		http.HandleFunc("/info/"+name+"/", BookInfo)
		http.HandleFunc("/info/"+name, SiteInfo)
	}
	http.HandleFunc("/process", ProcessState)
	http.HandleFunc("/info", GeneralInfo)
	fmt.Println("started")
	log.Fatal(http.ListenAndServe(":9001", nil))
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
