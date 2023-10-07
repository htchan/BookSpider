package main

import (
	"context"
	"os"
	"sync"

	config "github.com/htchan/BookSpider/internal/config_new"
	repo "github.com/htchan/BookSpider/internal/repo/sqlc"
	service_new "github.com/htchan/BookSpider/internal/service_new"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/semaphore"
)

func main() {
	outputPath := os.Getenv("OUTPUT_PATH")
	if outputPath != "" {
		writer, err := os.OpenFile(outputPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err == nil {
			log.Logger = log.Logger.Output(writer)
			defer writer.Close()
		} else {
			log.Fatal().
				Err(err).
				Str("output_path", outputPath).
				Msg("set logger output failed")
		}
	}

	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.99999Z07:00"

	conf, confErr := config.LoadConfig()
	if confErr != nil {
		log.Error().Err(confErr).Msg("load backend config")
		return
	}

	validErr := conf.Validate()
	if validErr != nil {
		log.Error().Err(validErr).Msg("validate config fail")
		return
	}

	ctx := context.Background()
	publicSema := semaphore.NewWeighted(int64(conf.BatchConfig.MaxWorkingThreads))
	services := make(map[string]service_new.Service)
	for _, siteName := range conf.BatchConfig.AvailableSiteNames {
		migrateDB, migrateDBErr := repo.OpenDatabase(siteName)
		if migrateDBErr != nil {
			log.Error().Err(migrateDBErr).Str("site", siteName).Msg("load db for migration Fail")
			return
		}

		migrateErr := repo.Migrate(migrateDB)
		if migrateErr != nil {
			log.Error().Err(migrateErr).Str("site", siteName).Msg("migrate fail")
		}

		db, dbErr := repo.OpenDatabase(siteName)
		if dbErr != nil {
			log.Error().Err(dbErr).Str("site", siteName).Msg("load db fail")
			return
		}

		defer db.Close()

		serv, loadServErr := service_new.LoadService(
			siteName, conf.SiteConfigs[siteName], db, ctx, publicSema,
		)
		if loadServErr != nil {
			log.Error().Err(loadServErr).Str("site", siteName).Msg("load service fail")
			return
		}

		services[siteName] = serv
	}

	// loop all sites by calling process
	var wg sync.WaitGroup
	log.Log().Msg("start regular batch process")

	for _, serv := range services {
		serv := serv
		wg.Add(1)
		go func(serv service_new.Service) {
			defer wg.Done()
			processErr := serv.Process()
			if processErr != nil {
				log.Error().Err(processErr).Str("site", serv.Name()).Msg("process failed")
			}
		}(serv)
	}

	wg.Wait()
	log.Log().Msg("completed regular batch process")
}
