package main

import (
	"os"
	"strings"
	"log"
	"github.com/htchan/BookSpider/models"
	"github.com/htchan/BookSpider/services"
	"github.com/htchan/BookSpider/helper"
	_ "net/http/pprof"

	"flag"
)

var stageFileName string

func help() {
	log.Println("Command: ")
    log.Println("help" + strings.Repeat(" ", 14) + "show the functin list avaliable")
    log.Println("download" + strings.Repeat(" ", 10) + "download books")
    log.Println("update" + strings.Repeat(" ", 12) + "update books information")
    log.Println("explore" + strings.Repeat(" ", 11) + "explore new books in internet")
    log.Println("check" + strings.Repeat(" ", 13) + "check database or storage error")
    log.Println("checkend" + strings.Repeat(" ", 10) + "check recorded books finished")
    log.Println("error" + strings.Repeat(" ", 13) + "update all website may have error")
    log.Println("backup" + strings.Repeat(" ", 12) + "backup the current database by the current date and time")
    log.Println("regular" + strings.Repeat(" ", 11) + "do the default operation (explore->update->download->check)")
	log.Println("fix" + strings.Repeat(" ", 15) + "fix the error in database and storage in the site")
	log.Println("\n")
}

func info(sites map[string]models.Site, flags models.Flags) {
	for name, site := range sites {
		if *flags.Site != "" && name != *flags.Site { continue }
		log.Println(name + "\tinfo")
		log.Println(strings.Repeat("- ", 20))
		site.Info()
		log.Println(strings.Repeat("- ", 20))
	}
}
// func schedule(sites map[string]models.Site, config models.Config, flags models.Flags) {
// 	validate(sites, config, flags)
// 	update(sites, config, flags)
// 	explore(sites, config, flags)
// 	updateError(sites, config, flags)
// 	check(sites, config, flags)
// 	fix(sites, config, flags)
// 	checkEnd(sites, config, flags)
// 	services.Download(sites, flags)
// }
func test(sites map[string]models.Site, flags models.Flags) {
	site := sites["hjwzw"]
	log.Println(site.SiteName)
	log.Println()
	// book := site.Book(36814, -1)
	// book.Download("./validate-download/", 1000)
	site.BackupSql("/backup")
	//site.Validate()
	//site.Explore(1000)
	//site.Update()
}

func main() {
	// runtime.GOMAXPROCS(3)
	log.Println("test (v0.0.0) - - - - - - - - - -")
	if (len(os.Args) < 2) {
		help()
		log.Println("No arguements")
		return
	}

	var flags models.Flags

	flags.Operation = flag.String("operation", "", "the operation to work on")
	flags.Site = flag.String("site", "", "specific site to operate")
	flags.Id = flag.Int("id", -1, "specific id to operate")
	flags.MaxThreads = flag.Int("max-threads", -1, "maximum number of threads to carry the process")
	flag.Parse()
	// log.Println(flags.Site, *flags.Site)

	config := models.LoadYaml("./config/config.yaml")
	sites := models.LoadSitesYaml(config)

	if *flags.MaxThreads <= 0 {
		*flags.MaxThreads = config.MaxThreads
	} else if *flags.MaxThreads <= 100 {
		*flags.MaxThreads = 100
	}

	helper.StageFileName = config.Backend.StageFile

	os.Remove(config.Backend.StageFile)
	os.Create(config.Backend.StageFile)

	functionMap := map[string]func(map[string]models.Site, models.Flags){
		"UPDATE": services.Update,
		"EXPLORE": services.Explore,
		"DOWNLOAD": services.Download,
		"ERROR": services.UpdateError,
		"INFO": info,
		"CHECK": services.Check,
		"CHECKEND": services.CheckEnd,
		// "BACKUP": backup,
		"BACKUP": services.BackupSql,
		"FIX": services.Fix,
		"RANDOM": nil,
		"VALIDATE": services.Validate,
		"TEST": test,
	}

	// function, exist := functionMap[strings.ToUpper(os.Args[1])]
	function, exist := functionMap[strings.ToUpper(*flags.Operation)]
	if !exist {
		help()
		log.Println("Invalid rguement")
	} else {
		function(sites, flags)
	}
}
