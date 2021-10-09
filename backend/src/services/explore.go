package services

import (
	"log"
	"github.com/htchan/BookSpider/models"
	"github.com/htchan/BookSpider/helper"
	"sync"
	"runtime"
	"golang.org/x/sync/semaphore"
)

func Explore(sites map[string]models.Site, flags models.Flags) {
	helper.WriteStage("stage: explore start")
	var wg sync.WaitGroup
	var s = semaphore.NewWeighted(int64(*flags.MaxThreads))
	for name, site := range sites {
		if *flags.Site != "" && name != *flags.Site { continue }
		wg.Add(1)
		go func(name string, site models.Site) {
			helper.WriteStage("sub_stage: " + name + " start")
			log.Println(name + "\texplore")
			site.Explore(*flags.MaxThreads, s)
			runtime.GC()
			helper.WriteStage("sub_stage: " + name + " finish")
			wg.Done()
		} (name, site)
	}
	wg.Wait()
	helper.WriteStage("stage: explore finish")
}