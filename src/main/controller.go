package main

import (
	"os"
	"strings"
	"fmt"
	"../model"
	"../helper"
	"runtime"
	_ "net/http/pprof"
	"path/filepath"
	
	"encoding/json"
	"io/ioutil"

	"sync"
)

var stageFileName string

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

func writeStage(s string) () {
	file, err := os.OpenFile(stageFileName, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0664)
	if err != nil {
		fmt.Println(err)
		fmt.Println(stageFileName)
	}
	file.WriteString(s + "\n")
    file.Close()
}

func download(sites map[string]model.Site) {
	writeStage("stage: download start")
	var wg sync.WaitGroup
	for name, site := range sites {
		wg.Add(1)
		go func() {
			writeStage("sub_stage: " + name + " start")
			if name == "80txt" {
				wg.Done()
				return
			}
			fmt.Println(name + "\tdownload")
			site.Download()
			runtime.GC()
			writeStage("sub_stage: " + name + " finish")
			wg.Done()
		} ()
	}
	writeStage("stage: download finish")
}
func update(sites map[string]model.Site) {
	writeStage("stage: update start")
	var wg sync.WaitGroup
	for name, site := range sites {
		wg.Add(1)
		go func() {
			writeStage("sub_stage: " + name + " start")
			fmt.Println(name + "\tupdate")
			site.Update()
			runtime.GC()
			writeStage("sub_stage: " + name + " finish")
			wg.Done()
		} ()
	}
	wg.Wait()
	writeStage("stage: update finish")
}
func explore(sites map[string]model.Site, maxError int) {
	writeStage("stage: explore start")
	var wg sync.WaitGroup
	for name, site := range sites {
		wg.Add(1)
		go func() {
			writeStage("sub_stage: " + name + " start")
			fmt.Println(name + "\texplore")
			site.Explore(maxError)
			runtime.GC()
			writeStage("sub_stage: " + name + " finish")
			wg.Done()
		} ()
	}
	writeStage("stage: explore finish")
}
func updateError(sites map[string]model.Site) {
	writeStage("stage: update error start")
	var wg sync.WaitGroup
	for name, site := range sites {
		wg.Add(1)
		go func() {
			writeStage("sub_stage: " + name + " start")
			fmt.Println(name + "\tupdate error")
			site.UpdateError()
			runtime.GC()
			writeStage("sub_stage: " + name + " finish")
			wg.Done()
		} ()
	}
	writeStage("stage: update error finish")
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
	writeStage("stage: check error start")
	for name, site := range sites {
		writeStage("sub_stage: " + name + " start")
		fmt.Println(name + "\tcheck")
		fmt.Println(strings.Repeat("- ", 20))
		site.Check()
		fmt.Println(strings.Repeat("- ", 20))
		runtime.GC()
		writeStage("sub_stage: " + name + " finish")
	}
	writeStage("stage: check error finish")
}
func checkEnd(sites map[string]model.Site) {
	writeStage("stage: check end start")
	var wg sync.WaitGroup
	for name, site := range sites {
		wg.Add(1)
		go func() {
			writeStage("sub_stage: " + name + " start")
			fmt.Println(name + "\tcheck end")
			fmt.Println(strings.Repeat("- ", 20))
			site.CheckEnd()
			fmt.Println(strings.Repeat("- ", 20))
			runtime.GC()
			writeStage("sub_stage: " + name + " finish")
			wg.Done()
		} ()
	}
	writeStage("stage: check end finish")
}
func backup(sites map[string]model.Site) {
	writeStage("stage: backup start")
	var wg sync.WaitGroup
	for name, site := range sites {
		wg.Add(1)
		go func() {
			writeStage("sub_stage: " + name + " start")
			fmt.Println(name + "\tbackup")
			site.Backup()
			runtime.GC()
			writeStage("sub_stage: " + name + " finish")
			wg.Done()
		} ()
	}
	writeStage("stage: backup finish")
}
func backupString(sites map[string]model.Site) {
	writeStage("stage: backup string start")
	for name, site := range sites {
		writeStage("sub_stage: " + name + " start")
		fmt.Println(name + "\tbackup")
		site.Backup()
		runtime.GC()
		writeStage("sub_stage: " + name + " finish")
	}
	writeStage("stage: backup string finish")
}
func fix(sites map[string]model.Site) {
	writeStage("stage: fix start")
	var wg sync.WaitGroup
	for name, site := range sites {
		wg.Add(1)
		go func() {
			writeStage("sub_stage: " + name + " start")
			fmt.Println(name + "\tfix")
			site.Fix()
			runtime.GC()
			writeStage("sub_stage: " + name + " finish")
			wg.Done()
		} ()
	}
	writeStage("stage: fix finish")
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
func validate(sites map[string]model.Site) {
	exploreResult := make(map[string]float64)
	downloadResult := make(map[string]float64)
	for name, site := range sites {
		fmt.Println(name + "\tvalidate explore")
		exploreResult[name] = site.Validate()
		fmt.Println(name + "\tvalidate download")
		downloadResult[name] = site.ValidateDownload()
	}
	b, err := json.Marshal(exploreResult)
	helper.CheckError(err)
	err = ioutil.WriteFile("./validate.json", b, 0644)
	helper.CheckError(err)
	b, err = json.Marshal(downloadResult)
	helper.CheckError(err)
	err = ioutil.WriteFile("./validate-download.json", b, 0644)
	helper.CheckError(err)

}
func schedule(sites map[string]model.Site) {
	validate(sites)
	backup(sites)
	update(sites)
	explore(sites, 1000)
	updateError(sites)
	check(sites)
	fix(sites)
	checkEnd(sites)
	download(sites)
}
func test(sites map[string]model.Site) {
	site := sites["hjwzw"]
	book := site.Book(36458, -1)
	book.Download("./validate-download/", 1000)
	//site.Validate()
	//site.Explore(1000)
	//site.Update()
}

func main() {
	//go func() { log.Fatal(http.ListenAndServe(":4000", nil)) }()
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
	config := model.LoadYaml("./config/config.yaml")
	sites := model.LoadSitesYaml(config)

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	stageFileName = dir + "/log/stage.txt"

	os.Remove(stageFileName)
	os.Create(stageFileName)
	
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
	case "BACKUPSTRING":
		backupString(sites)
	case "FIX":
		fix(sites)
	case "RANDOM":
		random(sites)
	case "VALIDATE":
		validate(sites)
	case "TEST":
		test(sites)
	default:
		help()
		fmt.Println("Invalid rguement")
	}
}