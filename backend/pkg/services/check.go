package services

import (
	"log"
	"github.com/htchan/BookSpider/pkg/flags"
	"github.com/htchan/BookSpider/pkg/sites"
	"github.com/htchan/BookSpider/internal/utils"
	"sync"
	"strings"
	"runtime"
)

func Check(siteMap map[string]sites.Site, flags flags.Flags) {
	utils.WriteStage("stage: check error start")
	for name, site := range siteMap {
		if *flags.Site != "" && name != *flags.Site { continue }
		utils.WriteStage("sub_stage: " + name + " start")
		log.Println(name + "\tcheck")
		log.Println(strings.Repeat("- ", 20))
		site.Check()
		log.Println(strings.Repeat("- ", 20))
		runtime.GC()
		utils.WriteStage("sub_stage: " + name + " finish")
	}
	utils.WriteStage("stage: check error finish")
}

func CheckEnd(siteMap map[string]sites.Site, flags flags.Flags) {
	utils.WriteStage("stage: check end start")
	var wg sync.WaitGroup
	for name, site := range siteMap {
		if *flags.Site != "" && name != *flags.Site { continue }
		wg.Add(1)
		go func(name string, site sites.Site) {
			utils.WriteStage("sub_stage: " + name + " start")
			log.Println(name + "\tcheck end")
			log.Println(strings.Repeat("- ", 20))
			site.CheckEnd()
			log.Println(strings.Repeat("- ", 20))
			runtime.GC()
			utils.WriteStage("sub_stage: " + name + " finish")
			wg.Done()
		} (name, site)
	}
	wg.Wait()
	utils.WriteStage("stage: check end finish")
}