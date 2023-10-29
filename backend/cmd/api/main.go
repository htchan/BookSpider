package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/htchan/BookSpider/internal/common"
	config "github.com/htchan/BookSpider/internal/config_new"
	repo "github.com/htchan/BookSpider/internal/repo/sqlc"
	"github.com/htchan/BookSpider/internal/router"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

	repo.Migrate(conf.DatabaseConfig, "/migrations")

	db, dbErr := repo.OpenDatabaseByConfig(conf.DatabaseConfig)
	if dbErr != nil {
		log.Error().Err(dbErr).Msg("load db fail")
		return
	}

	defer db.Close()

	// ctx := context.Background()
	// publicSema := semaphore.NewWeighted(int64(conf.BatchConfig.MaxWorkingThreads))
	// services := make(map[string]service_new.Service)
	// for _, siteName := range conf.APIConfig.AvailableSiteNames {
	// 	serv, loadServErr := service_new.LoadService(
	// 		siteName, conf.SiteConfigs[siteName], db, ctx, publicSema,
	// 	)
	// 	if loadServErr != nil {
	// 		log.Error().Err(loadServErr).Str("site", siteName).Msg("load service fail")
	// 		return
	// 	}

	// 	services[siteName] = serv
	// }
	services := common.LoadServices(conf.APIConfig.AvailableSiteNames, db, conf)

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
