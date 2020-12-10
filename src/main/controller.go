package main

import (
	"os"
	"strings"
	"fmt"
	"../model"
	"runtime"
	_ "net/http/pprof"
	"log"
	"net/http"
)
/*
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
*/

func help() {
	fmt.Println("Command: ")
    fmt.Println("help" + strings.Repeat(" ", 14) + "show the functin list avaliable")
    fmt.Println("download" + strings.Repeat(" ", 10) + "download books")
    fmt.Println("update" + strings.Repeat(" ", 12) + "update books information")
    fmt.Println("explore" + strings.Repeat(" ", 11) + "explore new books in internet")
    fmt.Println("check" + strings.Repeat(" ", 13) + "check database or storage error")
    fmt.Println("checkend" + strings.Repeat(" ", 10) + "check recorded books finished")
    fmt.Println("error" + strings.Repeat(" ", 13) + "update all website may have error")
    fmt.Println("backup" + strings.Repeat(" ", 12) + "backup the current database by the current date and time")
    fmt.Println("regular" + strings.Repeat(" ", 11) + "do the default operation (explore->update->download->check)")
	fmt.Println("fix" + strings.Repeat(" ", 15) + "fix the error in database and storage in the site")
	fmt.Println("\n")
}

func download(sites map[string]model.Site) {
	for name, site := range sites {
		if name == "80txt" {
			continue
		}
		fmt.Println(name + "\tdownload")
		site.Download()
		runtime.GC()
	}
}
func update(sites map[string]model.Site) {
	for name, site := range sites {
		fmt.Println(name + "\tupdate")
		site.Update()
		runtime.GC()
	}
}
func explore(sites map[string]model.Site, maxError int) {
	for name, site := range sites {
		fmt.Println(name + "\texplore")
		site.Explore(maxError)
		runtime.GC()
	}
}
func updateError(sites map[string]model.Site) {
	for name, site := range sites {
		fmt.Println(name + "\tupdate error")
		site.UpdateError()
		runtime.GC()
	}
}
func info(sites map[string]model.Site) {
	for name, site := range sites {
		fmt.Println(name + "\tinfo")
		fmt.Println(strings.Repeat("- ", 20))
		site.Info()
		fmt.Println(strings.Repeat("- ", 20))
	}
}
func check(sites map[string]model.Site) {
	for name, site := range sites {
		fmt.Println(name + "\tcheck")
		fmt.Println(strings.Repeat("- ", 20))
		site.Check()
		fmt.Println(strings.Repeat("- ", 20))
		runtime.GC()
	}
}
func checkEnd(sites map[string]model.Site) {
	for name, site := range sites {
		fmt.Println(name + "\tcheck end")
		fmt.Println(strings.Repeat("- ", 20))
		site.CheckEnd()
		fmt.Println(strings.Repeat("- ", 20))
		runtime.GC()
	}
}
func backup(sites map[string]model.Site) {
	for name, site := range sites {
		fmt.Println(name + "\tbackup")
		site.Backup()
		runtime.GC()
	}
}
func fix(sites map[string]model.Site) {
	for name, site := range sites {
		fmt.Println(name + "\tfix")
		site.Fix()
		runtime.GC()
	}
}
func random(sites map[string]model.Site) {
	for name, site := range sites {
		fmt.Println(name + "\trandom")
		results := site.Random(5)
		for _, result := range results {
			fmt.Println(result.String()+"\n")
		}
		runtime.GC()
	}
}
func test(sites map[string]model.Site) {
	site := sites["hjwzw"]
	site.Download()
	//site.Explore(1000)
	//site.Update()
}

func main() {
	go func() { log.Fatal(http.ListenAndServe(":4000", nil)) }()
	fmt.Println("test (v0.0.0) - - - - - - - - - -")
	if (len(os.Args) < 2) {
		help()
		fmt.Println("No arguements")
		return
	}
	/*
	big5Decoder := traditionalchinese.Big5.NewDecoder()

	sites := make(map[string]model.Site)
	sites["ck101"] = model.NewSite("ck101", big5Decoder, "./book-config/ck101-desktop.json", "./database/ck101.db", "./")
	*/

	sites := model.LoadSites("./config/config.json")

	switch operation := strings.ToUpper(os.Args[1]); operation {
	case "UPDATE":
		update(sites)
	case "EXPLORE":
		explore(sites, 1000)
	case "DOWNLOAD":
		download(sites)
	case "ERROR":
		updateError(sites)
	case "INFO":
		info(sites)
	case "CHECK":
		check(sites)
	case "CHECKEND":
		checkEnd(sites)
	case "BACKUP":
		backup(sites)
	case "FIX":
		fix(sites)
	case "RANDOM":
		random(sites)
	case "TEST":
		test(sites)
	default:
		help()
		fmt.Println("Invalid rguement")
	}
}
