package main

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/htchan/BookSpider/internal/common"
	"github.com/htchan/BookSpider/internal/config/v2"
	repo "github.com/htchan/BookSpider/internal/repo/sqlc"
	"github.com/htchan/BookSpider/internal/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func CalculateNextRunTime() time.Time {
	result := time.Now().UTC().Truncate(24*time.Hour).AddDate(0, 1, 0)
	result = time.Date(result.Year(), result.Month(), 1, 20, 0, 0, 0, time.UTC)

	if result.Weekday() != time.Friday {
		nDaysLater := int(time.Friday - result.Weekday())
		if nDaysLater < 0 {
			nDaysLater += 7
		}

		result = result.AddDate(0, 0, nDaysLater)
	}

	return result
}

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

	conf, confErr := config.LoadWorkerConfig()
	if confErr != nil {
		log.Error().Err(confErr).Msg("load backend config")
		return
	}

	validErr := conf.Validate()
	if validErr != nil {
		log.Error().Err(validErr).Msg("validate config fail")
		return
	}

	repo.Migrate(conf.DatabaseConfig, "/migrations")

	db, dbErr := repo.OpenDatabaseByConfig(conf.DatabaseConfig)
	if dbErr != nil {
		log.Error().Err(dbErr).Msg("load db fail")
		return
	}

	defer db.Close()

	services := common.LoadServices(conf.AvailableSiteNames, db, conf.SiteConfigs, int64(conf.MaxWorkingThreads))

	// loop all sites by calling process
	var wg sync.WaitGroup

	for true {
		time.Sleep(time.Until(CalculateNextRunTime()))
		log.Log().Msg("start regular batch process")

		for _, serv := range services {
			serv := serv
			wg.Add(1)
			go func(serv service.Service) {
				defer wg.Done()
				ctx := log.Logger.WithContext(context.Background())
				processErr := serv.Process(ctx)
				if processErr != nil {
					log.Error().Err(processErr).Str("site", serv.Name()).Msg("process failed")
				}
			}(serv)
		}

		wg.Wait()
		log.Log().Msg("completed regular batch process")
	}
}
