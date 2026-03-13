package main

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	shutdown "github.com/htchan/goshutdown"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/htchan/BookSpider/internal/common"
	"github.com/htchan/BookSpider/internal/config/v2"
	intOtel "github.com/htchan/BookSpider/internal/otel"
	repo "github.com/htchan/BookSpider/internal/repo/sqlc"
	"github.com/htchan/BookSpider/internal/router"
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

	conf, confErr := config.LoadAPIConfig()
	if confErr != nil {
		log.Error().Err(confErr).Msg("load backend config")
		return
	}

	validErr := conf.Validate()
	if validErr != nil {
		log.Error().Err(validErr).Msg("validate config fail")
		return
	}

	tp, err := intOtel.NewProvider(conf.TraceConfig)
	if err != nil {
		log.Error().Err(err).Msg("init tracer failed")
	}

	repo.Migrate(conf.DatabaseConfig, "/migrations")

	db, dbErr := repo.OpenDatabaseByConfig(conf.DatabaseConfig)
	if dbErr != nil {
		log.Error().Err(dbErr).Msg("load db fail")
		return
	}

	defer db.Close()

	services := common.LoadServices(conf.AvailableSiteNames, db, conf.SiteConfigs, 1)
	readDataService := common.LoadReadDataService(db, conf.SiteConfigs)

	shutdown.LogEnabled = true
	shutdownHandler := shutdown.New(syscall.SIGINT, syscall.SIGTERM)

	// load routes
	r := chi.NewRouter()
	router.AddAPIRoutes(r, conf, services, readDataService)
	router.AddLiteRoutes(r, conf, services, readDataService)

	server := http.Server{
		Addr:         ":9427",
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 300 * time.Second,
		IdleTimeout:  300 * time.Second,
	}
	go func() {
		log.Info().Msg("start http server")

		if httpErr := server.ListenAndServe(); httpErr != nil {
			log.Error().Err(httpErr).Msg("backend stopped")
			return
		}
	}()

	shutdownHandler.Register("api server", func() error {
		server.Shutdown(context.Background())

		return nil
	})
	shutdownHandler.Register("database", db.Close)
	shutdownHandler.Register("tracer", func() error {
		return tp.Shutdown(context.Background())
	})

	shutdownHandler.Listen(60 * time.Second)
}
