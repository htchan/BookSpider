package services

import (
	"log"
	"github.com/htchan/BookSpider/models"
	"github.com/htchan/BookSpider/helper"
	"sync"
	"strings"
	"runtime"
)

func Check(sites map[string]models.Site, flags models.Flags) {
	helper.WriteStage("stage: check error start")
	for name, site := range sites {
		if *flags.Site != "" && name != *flags.Site { continue }
		helper.WriteStage("sub_stage: " + name + " start")
		log.Println(name + "\tcheck")
		log.Println(strings.Repeat("- ", 20))
		site.Check()
		log.Println(strings.Repeat("- ", 20))
		runtime.GC()
		helper.WriteStage("sub_stage: " + name + " finish")
	}
	helper.WriteStage("stage: check error finish")
}

func CheckEnd(sites map[string]models.Site, flags models.Flags) {
	helper.WriteStage("stage: check end start")
	var wg sync.WaitGroup
	for name, site := range sites {
		if *flags.Site != "" && name != *flags.Site { continue }
		wg.Add(1)
		go func(name string, site models.Site) {
			helper.WriteStage("sub_stage: " + name + " start")
			log.Println(name + "\tcheck end")
			log.Println(strings.Repeat("- ", 20))
			site.CheckEnd()
			log.Println(strings.Repeat("- ", 20))
			runtime.GC()
			helper.WriteStage("sub_stage: " + name + " finish")
			wg.Done()
		} (name, site)
	}
	wg.Wait()
	helper.WriteStage("stage: check end finish")
}