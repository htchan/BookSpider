package main

import (
	"os"
	"strings"
	"fmt"
	"github.com/htchan/BookSpider/model"
	"github.com/htchan/BookSpider/helper"
	"runtime"
	_ "net/http/pprof"
	
	"golang.org/x/sync/semaphore"
	
	"encoding/json"
	"io/ioutil"

	"sync"
	"flag"
)

var stageFileName string

type Flags struct {
	site *string
	id *int
}

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

func writeStage(s string) {
	file, err := os.OpenFile(stageFileName, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0664)
	if err != nil { fmt.Println(err, "\n", stageFileName) }
	file.WriteString(s + "\n")
    file.Close()
}

func download(sites map[string]model.Site, config model.Config, flags Flags) {
	writeStage("stage: download start")
	var wg sync.WaitGroup
	for name, site := range sites {
		if *flags.site != "" && name != *flags.site { continue }
		wg.Add(1)
		go func(name string, site model.Site) {
			defer wg.Done()
			writeStage("sub_stage: " + name + " start")
			if name == "80txt" { return }
			fmt.Println(name + "\tdownload")
			if *flags.id != -1 {
				book := site.Book(*flags.id, -1)
				book.Download(".", site.MAX_THREAD_COUNT)
			} else {
				site.Download()
			}
			runtime.GC()
			writeStage("sub_stage: " + name + " finish")
		} (name, site)
	}
	wg.Wait()
	writeStage("stage: download finish")
}
func update(sites map[string]model.Site, config model.Config, flags Flags) {
	writeStage("stage: update start")
	var wg sync.WaitGroup
	var s = semaphore.NewWeighted(int64(config.MaxThreads))
	for name, site := range sites {
		if *flags.site != "" && name != *flags.site { continue }
		wg.Add(1)
		go func(name string, site model.Site) {
			writeStage("sub_stage: " + name + " start")
			fmt.Println(name + "\tupdate")
			site.Update(s)
			runtime.GC()
			writeStage("sub_stage: " + name + " finish")
			wg.Done()
		} (name, site)
	}
	wg.Wait()
	writeStage("stage: update finish")
}
func explore(sites map[string]model.Site, config model.Config, flags Flags) {
	writeStage("stage: explore start")
	var wg sync.WaitGroup
	var s = semaphore.NewWeighted(int64(config.MaxThreads))
	for name, site := range sites {
		if *flags.site != "" && name != *flags.site { continue }
		wg.Add(1)
		go func(name string, site model.Site) {
			writeStage("sub_stage: " + name + " start")
			fmt.Println(name + "\texplore")
			site.Explore(config.MaxExploreError, s)
			runtime.GC()
			writeStage("sub_stage: " + name + " finish")
			wg.Done()
		} (name, site)
	}
	wg.Wait()
	writeStage("stage: explore finish")
}
func updateError(sites map[string]model.Site, config model.Config, flags Flags) {
	writeStage("stage: update error start")
	var wg sync.WaitGroup
	var s = semaphore.NewWeighted(int64(config.MaxThreads))
	for name, site := range sites {
		if *flags.site != "" && name != *flags.site { continue }
		wg.Add(1)
		go func(name string, site model.Site) {
			writeStage("sub_stage: " + name + " start")
			fmt.Println(name + "\tupdate error")
			site.UpdateError(s)
			runtime.GC()
			writeStage("sub_stage: " + name + " finish")
			wg.Done()
		} (name, site)
	}
	wg.Wait()
	writeStage("stage: update error finish")
}
func info(sites map[string]model.Site, config model.Config, flags Flags) {
	for name, site := range sites {
		if *flags.site != "" && name != *flags.site { continue }
		fmt.Println(name + "\tinfo")
		fmt.Println(strings.Repeat("- ", 20))
		site.Info()
		fmt.Println(strings.Repeat("- ", 20))
	}
}
func check(sites map[string]model.Site, config model.Config, flags Flags) {
	writeStage("stage: check error start")
	for name, site := range sites {
		if *flags.site != "" && name != *flags.site { continue }
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
func checkEnd(sites map[string]model.Site, config model.Config, flags Flags) {
	writeStage("stage: check end start")
	var wg sync.WaitGroup
	for name, site := range sites {
		if *flags.site != "" && name != *flags.site { continue }
		wg.Add(1)
		go func(name string, site model.Site) {
			writeStage("sub_stage: " + name + " start")
			fmt.Println(name + "\tcheck end")
			fmt.Println(strings.Repeat("- ", 20))
			site.CheckEnd()
			fmt.Println(strings.Repeat("- ", 20))
			runtime.GC()
			writeStage("sub_stage: " + name + " finish")
			wg.Done()
		} (name, site)
	}
	wg.Wait()
	writeStage("stage: check end finish")
}
func fix(sites map[string]model.Site, config model.Config, flags Flags) {
	writeStage("stage: fix start")
	var wg sync.WaitGroup
	var s = semaphore.NewWeighted(int64(config.MaxThreads))
	for name, site := range sites {
		if *flags.site != "" && name != *flags.site { continue }
		wg.Add(1)
		go func(name string, site model.Site) {
			writeStage("sub_stage: " + name + " start")
			fmt.Println(name + "\tfix")
			site.Fix(s)
			runtime.GC()
			writeStage("sub_stage: " + name + " finish")
			wg.Done()
		} (name, site)
	}
	wg.Wait()
	writeStage("stage: fix finish")
}
func random(sites map[string]model.Site, config model.Config, flags Flags) {
	for name, site := range sites {
		if *flags.site != "" && name != *flags.site { continue }
		fmt.Println(name + "\trandom")
		results := site.RandomSuggestBook(5)
		for _, result := range results { fmt.Println(result.String()+"\n") }
		runtime.GC()
	}
}
func validate(sites map[string]model.Site, config model.Config, flags Flags) {
	exploreResult := make(map[string]float64)
	downloadResult := make(map[string]float64)
	for name, site := range sites {
		if *flags.site != "" && name != *flags.site { continue }
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
func schedule(sites map[string]model.Site, config model.Config, flags Flags) {
	validate(sites, config, flags)
	update(sites, config, flags)
	explore(sites, config, flags)
	updateError(sites, config, flags)
	check(sites, config, flags)
	fix(sites, config, flags)
	checkEnd(sites, config, flags)
	download(sites, config, flags)
}
func test(sites map[string]model.Site, config model.Config, flags Flags) {
	site := sites["hjwzw"]
	fmt.Println(site.SiteName)
	fmt.Println()
	book := site.Book(36814, -1)
	book.Download("./validate-download/", 1000)
	//site.Validate()
	//site.Explore(1000)
	//site.Update()
}

func main() {
	fmt.Println("test (v0.0.0) - - - - - - - - - -")
	if (len(os.Args) < 2) {
		help()
		fmt.Println("No arguements")
		return
	}

	var flags Flags

	flags.site = flag.String("site", "", "specific site to operate")
	flags.id = flag.Int("id", -1, "specific id to operate")
	flag.Parse()

	config := model.LoadYaml("./config/config.yaml")
	sites := model.LoadSitesYaml(config)

	stageFileName = config.Backend.StageFile

	os.Remove(config.Backend.StageFile)
	os.Create(config.Backend.StageFile)

	functionMap := map[string]func(map[string]model.Site, model.Config, Flags){
		"UPDATE": update,
		"EXPLORE": explore,
		"DOWNLOAD": download,
		"ERROR": updateError,
		"INFO": info,
		"CHECK": check,
		"CHECKEND": checkEnd,
		// "BACKUP": backup,
		// "BACKUPSTRING": backupString,
		"FIX": fix,
		"RANDOM": random,
		"VALIDATE": validate,
		"TEST": test,
	}

	function, exist := functionMap[strings.ToUpper(os.Args[1])]
	if !exist {
		help()
		fmt.Println("Invalid rguement")
	} else {
		function(sites, config, flags)
	}
}
