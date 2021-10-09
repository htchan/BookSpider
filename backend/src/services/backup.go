package services

import (
	"log"
	"github.com/htchan/BookSpider/models"
	"github.com/htchan/BookSpider/helper"
	"sync"
	"runtime"
)

func BackupSql(sites map[string]models.Site, flags models.Flags) {
	helper.WriteStage("stage: backup as sql start")
	var wg sync.WaitGroup
	for name, site := range sites {
		if *flags.Site != "" && name != *flags.Site { continue }
		wg.Add(1)
		go func(name string, site models.Site) {
			helper.WriteStage("sub_stage: " + name + " start")
			log.Println(name + "\tbackup as sql")
			site.BackupSql("/backup")
			runtime.GC()
			helper.WriteStage("sub_stage: " + name + " finish")
			wg.Done()
		} (name, site)
	}
	wg.Wait()
	helper.WriteStage("stage: backup as sql finish")
}