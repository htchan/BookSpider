package services

import (
	"log"
	"sync"
	"runtime"

	"golang.org/x/sync/semaphore"

	"github.com/htchan/BookSpider/pkg/flags"
	"github.com/htchan/BookSpider/pkg/sites"
	"github.com/htchan/BookSpider/internal/utils"
)

func Update(siteMap map[string]sites.Site, flags flags.Flags) {
	utils.WriteStage("stage: update start")
	var wg sync.WaitGroup
	var s = semaphore.NewWeighted(int64(*flags.MaxThreads))
	for name, site := range siteMap {
		if *flags.Site != "" && name != *flags.Site { continue }
		wg.Add(1)
		go func(name string, site sites.Site) {
			utils.WriteStage("sub_stage: " + name + " start")
			log.Println(name + "\tupdate")
			site.Update(s)
			runtime.GC()
			utils.WriteStage("sub_stage: " + name + " finish")
			wg.Done()
		} (name, site)
	}
	wg.Wait()
	utils.WriteStage("stage: update finish")
}

func UpdateError(siteMap map[string]sites.Site, flags flags.Flags) {
	utils.WriteStage("stage: update error start")
	var wg sync.WaitGroup
	var s = semaphore.NewWeighted(int64(*flags.MaxThreads))
	for name, site := range siteMap {
		if *flags.Site != "" && name != *flags.Site { continue }
		wg.Add(1)
		go func(name string, site sites.Site) {
			utils.WriteStage("sub_stage: " + name + " start")
			log.Println(name + "\tupdate error")
			site.UpdateError(s)
			runtime.GC()
			utils.WriteStage("sub_stage: " + name + " finish")
			wg.Done()
		} (name, site)
	}
	wg.Wait()
	utils.WriteStage("stage: update error finish")
}