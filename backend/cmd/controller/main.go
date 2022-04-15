package main

import (
	"os"
	"strings"
	"github.com/htchan/BookSpider/pkg/configs"
	"github.com/htchan/BookSpider/pkg/sites"
	"github.com/htchan/BookSpider/pkg/flags"
	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/logging"
	_ "net/http/pprof"
	"sync"
	"runtime"
)

var stageFileName string

func help() {
	logging.LogEvent("controller", "help", map[string]string{
		"help": "show the functin list avaliable",
		"download": "download books",
		"update": "update books information",
		"explore": "explore new books in internet",
		"check": "check database or storage error",
		"checkend": "check recorded books finished",
		"error": "update all website may have error",
		"backup": "backup the current database by the current date and time",
		"regular": "do the default operation (explore->update->download->check)",
		"fix": "fix the error in database and storage in the site",
	})
}

func info(site *sites.Site, flags *flags.Flags) error {
	if *flags.Site != "" && site.Name != *flags.Site { return nil }
	logging.LogEvent("controller", site.Name + ".info", site.Map())
	return nil
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
func test(site *sites.Site, flags *flags.Flags) error {
	logging.LogEvent("controller", "test", site.Name)
	// book := books.NewBook(site.SiteName, 36814, -1, site.meta, site.decoder, site.bookTx)
	// book.Download("./validate-download/", 1000)
	sites.Backup(site, flags)
	//site.Validate()
	//site.Explore(1000)
	//site.Update()
	return nil
}

func execute(function sites.SiteOperation, siteMap map[string]*sites.Site, f *flags.Flags) {
	var wg sync.WaitGroup
	for name, site := range siteMap {
		wg.Add(1)
		go func(name string, site *sites.Site) {
			utils.WriteStage("sub_stage: " + name + " start")
			err := site.OpenDatabase()
			utils.CheckError(err)
			err = function(site, f)
			site.CommitDatabase()
			site.CloseDatabase()
			if err != nil {
				logging.LogEvent("controller", "execute.error", map[string]interface{} {
					"error": err,
					"operation": *f.Operation,
					"site": *f.Site,
					"id": *f.Id,
					"hash": *f.HashCode,
					"threads": *f.MaxThreads,
				})
			}
			runtime.GC()
			utils.WriteStage("sub_stage: " + name + " finish")
			wg.Done()
		} (name, site)
	}
	wg.Wait()
}

func main() {
	// runtime.GOMAXPROCS(3)
	logging.LogEvent("controller", "start", nil)
	if (len(os.Args) < 2) {
		help()
		logging.LogEvent("controller", "error", "No Arguements")
		return
	}

	flag := flags.NewFlags()

	config := configs.LoadSystemConfigs(os.Getenv("ASSETS_LOCATION") + "/configs")
	siteMap := make(map[string]*sites.Site)
	for key, siteConfig := range config.AvailableSiteConfigs {
		siteMap[key] = sites.NewSite(key, siteConfig)
		siteMap[key].OpenDatabase()
	}

	flag.Load(10)


	utils.StageFileName = os.Getenv("ASSETS_LOCATION") + "/log/stage.txt"

	os.Remove(utils.StageFileName)
	os.Create(utils.StageFileName)

	functionMap := map[string]sites.SiteOperation{
		"UPDATE": sites.Update,
		"EXPLORE": sites.Explore,
		"DOWNLOAD": sites.Download,
		"ERROR": sites.UpdateError,
		"INFO": info,
		"CHECK": sites.Check,
		// "CHECKEND": sites.CheckEnd,
		// "BACKUP": backup,
		"BACKUP": sites.Backup,
		"FIX": sites.Fix,
		// "RANDOM": sites.Random,
		"VALIDATE": sites.Validate,
		"TEST": test,
	}

	// function, exist := functionMap[strings.ToUpper(os.Args[1])]
	function, exist := functionMap[strings.ToUpper(*flag.Operation)]
	
	if !exist {
		help()
		logging.LogEvent("controller", "error", "Invalid rguement")
	} else {
		utils.WriteStage("stage: " + *flag.Operation + " start")
		execute(function, siteMap, flag)
		utils.WriteStage("stage: " + *flag.Operation + " finish")
	}
}
