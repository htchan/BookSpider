package services

import (
	"log"
	"github.com/htchan/BookSpider/pkg/flags"
	"github.com/htchan/BookSpider/pkg/sites"
	"github.com/htchan/BookSpider/internal/utils"
	"sync"
	"runtime"
)

func Download(siteMap map[string]sites.Site, flags flags.Flags) {
	utils.WriteStage("stage: download start")
	var wg sync.WaitGroup
	for name, site := range siteMap {
		if *flags.Site != "" && name != *flags.Site { continue }
		wg.Add(1)
		go func(name string, site sites.Site) {
			defer wg.Done()
			defer runtime.GC()
			utils.WriteStage("sub_stage: " + name + " start")
			defer utils.WriteStage("sub_stage: " + name + " finish")
			if name == "80txt" { return }
			log.Println(name + "\tdownload")
			if *flags.Id > 0 {
				book, err := site.Book(*flags.Id, -1)
				if err != nil {
					book.Log(map[string]interface{}{
						"error": err.Error(), "stage": "download",
					})
				}
				book.Download(site.DownloadLocation, site.MAX_THREAD_COUNT)
			} else {
				site.Download()
			}
		} (name, site)
	}
	wg.Wait()
	utils.WriteStage("stage: download finish")
}