package main

import (
	"context"
	"log"
	"sync"

	config "github.com/htchan/BookSpider/internal/config_new"
	repo "github.com/htchan/BookSpider/internal/repo/psql"
	service_new "github.com/htchan/BookSpider/internal/service_new"
	"golang.org/x/sync/semaphore"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("load backend config: %v", err)
		return
	}

	validErr := conf.Validate()
	if validErr != nil {
		log.Fatalf("validate config fail: %v", validErr)
	}

	ctx := context.Background()
	services := make(map[string]service_new.Service)
	for _, siteName := range conf.BatchConfig.AvailableSiteNames {
		migrateDB, err := repo.OpenDatabase(siteName)
		if err != nil {
			log.Fatalf("load db Fail. site: %v; err: %v", siteName, err)
		}

		repo.Migrate(migrateDB)

		db, err := repo.OpenDatabase(siteName)
		if err != nil {
			log.Error().Err(err).Str("site", siteName).Msg("load db fail")
			return
		}

		defer db.Close()

		sema := semaphore.NewWeighted(int64(conf.SiteConfigs[siteName].MaxThreads))

		serv, err := service_new.LoadService(
			siteName, conf.SiteConfigs[siteName], db, sema, &ctx,
		)
		if err != nil {
			log.Fatalf("load service fail. site: %v, err: %v", siteName, err)
		}

		services[siteName] = serv
	}

	// loop all sites by calling process
	var wg sync.WaitGroup
	log.Println("start regular batch process")

	for _, serv := range services {
		serv := serv
		wg.Add(1)
		go func(serv service_new.Service) {
			defer wg.Done()
			err := serv.Process()
			if err != nil {
				log.Printf("Process fail: %v\n", err)
			}
		}(serv)
	}

	wg.Wait()
	log.Println("completed regular batch process")
}
