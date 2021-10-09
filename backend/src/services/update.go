package services

import (
	"log"
	"sync"
	"runtime"

	"golang.org/x/sync/semaphore"

	"github.com/htchan/BookSpider/models"
	"github.com/htchan/BookSpider/helper"
)

func Update(sites map[string]models.Site, flags models.Flags) {
	helper.WriteStage("stage: update start")
	var wg sync.WaitGroup
	var s = semaphore.NewWeighted(int64(*flags.MaxThreads))
	for name, site := range sites {
		if *flags.Site != "" && name != *flags.Site { continue }
		wg.Add(1)
		go func(name string, site models.Site) {
			helper.WriteStage("sub_stage: " + name + " start")
			log.Println(name + "\tupdate")
			site.Update(s)
			runtime.GC()
			helper.WriteStage("sub_stage: " + name + " finish")
			wg.Done()
		} (name, site)
	}
	wg.Wait()
	helper.WriteStage("stage: update finish")
}

func UpdateError(sites map[string]models.Site, flags models.Flags) {
	helper.WriteStage("stage: update error start")
	var wg sync.WaitGroup
	var s = semaphore.NewWeighted(int64(*flags.MaxThreads))
	for name, site := range sites {
		if *flags.Site != "" && name != *flags.Site { continue }
		wg.Add(1)
		go func(name string, site models.Site) {
			helper.WriteStage("sub_stage: " + name + " start")
			log.Println(name + "\tupdate error")
			site.UpdateError(s)
			runtime.GC()
			helper.WriteStage("sub_stage: " + name + " finish")
			wg.Done()
		} (name, site)
	}
	wg.Wait()
	helper.WriteStage("stage: update error finish")
}