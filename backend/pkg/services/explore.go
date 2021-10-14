package services

import (
	"log"
	"github.com/htchan/BookSpider/pkg/flags"
	"github.com/htchan/BookSpider/pkg/sites"
	"github.com/htchan/BookSpider/internal/utils"
	"sync"
	"runtime"
	"golang.org/x/sync/semaphore"
)

func Explore(siteMap map[string]sites.Site, flags flags.Flags) {
	utils.WriteStage("stage: explore start")
	var wg sync.WaitGroup
	var s = semaphore.NewWeighted(int64(*flags.MaxThreads))
	for name, site := range siteMap {
		if *flags.Site != "" && name != *flags.Site { continue }
		wg.Add(1)
		go func(name string, site sites.Site) {
			utils.WriteStage("sub_stage: " + name + " start")
			log.Println(name + "\texplore")
			site.Explore(*flags.MaxThreads, s)
			runtime.GC()
			utils.WriteStage("sub_stage: " + name + " finish")
			wg.Done()
		} (name, site)
	}
	wg.Wait()
	utils.WriteStage("stage: explore finish")
}