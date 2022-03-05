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
	logging.Info("Command: ")
    logging.Info("help" + strings.Repeat(" ", 14) + "show the functin list avaliable")
    logging.Info("download" + strings.Repeat(" ", 10) + "download books")
    logging.Info("update" + strings.Repeat(" ", 12) + "update books information")
    logging.Info("explore" + strings.Repeat(" ", 11) + "explore new books in internet")
    logging.Info("check" + strings.Repeat(" ", 13) + "check database or storage error")
    logging.Info("checkend" + strings.Repeat(" ", 10) + "check recorded books finished")
    logging.Info("error" + strings.Repeat(" ", 13) + "update all website may have error")
    logging.Info("backup" + strings.Repeat(" ", 12) + "backup the current database by the current date and time")
    logging.Info("regular" + strings.Repeat(" ", 11) + "do the default operation (explore->update->download->check)")
	logging.Info("fix" + strings.Repeat(" ", 15) + "fix the error in database and storage in the site")
	logging.Info("")
}

func info(site *sites.Site, flags *flags.Flags) error {
	if *flags.Site != "" && site.Name != *flags.Site { return nil }
	logging.Info(site.Name + "\tinfo")
	logging.Info(strings.Repeat("- ", 20))
	logging.Info("%v", site.Map())
	logging.Info(strings.Repeat("- ", 20))
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
	logging.Info(site.Name)
	logging.Info("")
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
				logging.Error("%v", err)
				logging.Error("'%v' '%v' '%v' '%v' '%v'", *f.Operation, *f.Site, *f.Id, *f.HashCode, *f.MaxThreads)
			}
			runtime.GC()
			utils.WriteStage("sub_stage: " + name + " finish")
			wg.Done()
		} (name, site)
	}
	wg.Wait()
}

func main() {
	logging.SetLogLevel(0)
	logging.UseFile(os.Getenv("ASSETS_LOCATION") + "/log/controller.log")
	// runtime.GOMAXPROCS(3)
	logging.Info("test (v0.0.0) - - - - - - - - - -")
	if (len(os.Args) < 2) {
		help()
		logging.Info("No arguements")
		return
	}

	flag := flags.NewFlags()

	config := configs.LoadConfigYaml(os.Getenv("ASSETS_LOCATION") + "/configs/config.yaml")
	siteMap := make(map[string]*sites.Site)
	for key, siteConfig := range config.SiteConfigs {
		siteMap[key] = sites.NewSite(key, siteConfig)
		siteMap[key].OpenDatabase()
	}

	flag.Load(config.MaxThreads)


	utils.StageFileName = os.Getenv("ASSETS_LOCATION") + config.Backend.StageFile

	os.Remove(config.Backend.StageFile)
	os.Create(config.Backend.StageFile)

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
		logging.Info("Invalid rguement")
	} else {
		utils.WriteStage("stage: " + *flag.Operation + " start")
		execute(function, siteMap, flag)
		utils.WriteStage("stage: " + *flag.Operation + " finish")
	}
}
