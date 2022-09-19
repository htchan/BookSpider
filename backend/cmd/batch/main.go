package main

import (
	"log"
	"os"
	"sync"

	"github.com/htchan/ApiParser"
	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/service/site"
)

func main() {
	configLocation := os.Getenv("ASSETS_LOCATION") + "/config"
	var err error

	// TODO: load backend config
	batchConfig, err := config.LoadBatchConfig(configLocation)
	if err != nil {
		log.Fatalf("load backend config: %v", err)
		return
	}

	sites, err := site.LoadSitesFromConfigDirectory(configLocation, batchConfig.EnabledSites)
	if err != nil {
		log.Fatal(err)
	}

	ApiParser.SetDefault(
		ApiParser.FromDirectory(os.Getenv("ASSETS_LOCATION") + "/api_parser"))

	// loop all sites by calling process
	var wg sync.WaitGroup
	log.Println("start")
	for _, st := range sites {
		st := st
		wg.Add(1)
		go func(st *site.Site) {
			defer wg.Done()
			err := site.Process(st)
			if err != nil {
				log.Printf("Process fail: %v\n", err)
			}
		}(st)
	}
	wg.Wait()
	log.Println("completed")
}
