package services

import (
	"log"
	"github.com/htchan/BookSpider/models"
	"github.com/htchan/BookSpider/helper"
	"sync"
	"runtime"
)

func Download(sites map[string]models.Site, flags models.Flags) {
	helper.WriteStage("stage: download start")
	var wg sync.WaitGroup
	for name, site := range sites {
		if *flags.Site != "" && name != *flags.Site { continue }
		wg.Add(1)
		go func(name string, site models.Site) {
			defer wg.Done()
			helper.WriteStage("sub_stage: " + name + " start")
			if name == "80txt" { return }
			log.Println(name + "\tdownload")
			if *flags.Id != -1 {
				book := site.Book(*flags.Id, -1)
				book.Download(site.DownloadLocation, site.MAX_THREAD_COUNT)
			} else {
				site.Download()
			}
			runtime.GC()
			helper.WriteStage("sub_stage: " + name + " finish")
		} (name, site)
	}
	wg.Wait()
	helper.WriteStage("stage: download finish")
}