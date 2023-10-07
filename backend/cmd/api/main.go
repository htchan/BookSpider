package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	config "github.com/htchan/BookSpider/internal/config_new"
	repo "github.com/htchan/BookSpider/internal/repo/sqlc"
	"github.com/htchan/BookSpider/internal/router"
	service_new "github.com/htchan/BookSpider/internal/service_new"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/semaphore"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339Nano})

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
	for _, siteName := range conf.APIConfig.AvailableSiteNames {
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

	// load routes
	r := chi.NewRouter()
	// if conf.APIConfig.ContainsRoute(config.RouteAPIKey) {
	router.AddAPIRoutes(r, conf.APIConfig, services)
	// }

	// if backendConfig.ContainsRoute(config.RouteLiteKey) {
	router.AddLiteRoutes(r, conf.APIConfig, services)
	// }

	server := http.Server{
		Addr:         ":9427",
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 300 * time.Second,
		IdleTimeout:  300 * time.Second,
	}
	// go func() {
	log.Info().Msg("start http server")

	if httpErr := server.ListenAndServe(); httpErr != nil {
		log.Error().Err(httpErr).Msg("backend stopped")
		return
	}
	// }()

	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, os.Interrupt)
	// <-sigChan
	// log.Println("received interrupt signal")

	// // Setup graceful shutdown
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()
	// server.Shutdown(ctx)
}
