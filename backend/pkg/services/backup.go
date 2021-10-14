package services

import (
	"log"
	"github.com/htchan/BookSpider/pkg/flags"
	"github.com/htchan/BookSpider/pkg/sites"
	"github.com/htchan/BookSpider/internal/utils"
	"sync"
	"runtime"
)

func BackupSql(siteMap map[string]sites.Site, flags flags.Flags) {
	utils.WriteStage("stage: backup as sql start")
	var wg sync.WaitGroup
	for name, site := range siteMap {
		if *flags.Site != "" && name != *flags.Site { continue }
		wg.Add(1)
		go func(name string, site sites.Site) {
			utils.WriteStage("sub_stage: " + name + " start")
			log.Println(name + "\tbackup as sql")
			site.BackupSql("/backup")
			runtime.GC()
			utils.WriteStage("sub_stage: " + name + " finish")
			wg.Done()
		} (name, site)
	}
	wg.Wait()
	utils.WriteStage("stage: backup as sql finish")
}