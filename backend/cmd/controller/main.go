package main

import (
	"os"
	"strings"
	"log"
	"github.com/htchan/BookSpider/pkg/configs"
	"github.com/htchan/BookSpider/pkg/sites"
	"github.com/htchan/BookSpider/pkg/flags"
	"github.com/htchan/BookSpider/pkg/services"
	"github.com/htchan/BookSpider/internal/utils"
	_ "net/http/pprof"

	// "flag"
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

func info(siteMap map[string]sites.Site, flags flags.Flags) {
	for name, site := range siteMap {
		if *flags.Site != "" && name != *flags.Site { continue }
		log.Println(name + "\tinfo")
		log.Println(strings.Repeat("- ", 20))
		site.Info()
		log.Println(strings.Repeat("- ", 20))
	}
}
// func schedule(siteMap map[string]sites.Site, config configs.Config, flags flags.Flags) {
// 	validate(sites, config, flags)
// 	update(sites, config, flags)
// 	explore(sites, config, flags)
// 	updateError(sites, config, flags)
// 	check(sites, config, flags)
// 	fix(sites, config, flags)
// 	checkEnd(sites, config, flags)
// 	services.Download(sites, flags)
// }
func test(siteMap map[string]sites.Site, flags flags.Flags) {
	site := siteMap["hjwzw"]
	log.Println(site.SiteName)
	log.Println()
	// book := books.NewBook(site.SiteName, 36814, -1, site.meta, site.decoder, site.bookTx)
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

	flag := flags.NewFlags()

	config := configs.LoadConfigYaml("./configs/config.yaml")
	siteMap := configs.LoadSitesYaml(config)

	flag.Load(config.MaxThreads)


	utils.StageFileName = config.Backend.StageFile

	os.Remove(config.Backend.StageFile)
	os.Create(config.Backend.StageFile)

	functionMap := map[string]func(map[string]sites.Site, flags.Flags){
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
	function, exist := functionMap[strings.ToUpper(*flag.Operation)]
	if !exist {
		help()
		log.Println("Invalid rguement")
	} else {
		function(siteMap, *flag)
	}
}
