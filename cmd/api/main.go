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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/htchan/BookSpider/internal/config/v1"
	repo "github.com/htchan/BookSpider/internal/repo/sqlc"
	"github.com/htchan/BookSpider/internal/router"
	"github.com/htchan/BookSpider/internal/service/v1"
)

func otelProvider(conf config.TraceConfig) (*tracesdk.TracerProvider, error) {
	exp, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(conf.OtelURL),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(conf.OtelServiceName),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp, nil
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

	conf, confErr := config.LoadBackendConfig()
	if confErr != nil {
		log.Error().Err(confErr).Msg("load backend config")
		return
	}

	tp, err := otelProvider(conf.Trace)
	if err != nil {
		log.Error().Err(err).Msg("init tracer failed")
	}

	repo.Migrate(conf.Database, "/migrations")

	db, dbErr := repo.OpenDatabaseByConfig(conf.Database)
	if dbErr != nil {
		log.Error().Err(dbErr).Msg("load db fail")
		return
	}

	rpo := repo.NewRepo(db)
	readDataService := service.NewReadDataService(rpo, conf.Common.StoragePath)

	shutdown.LogEnabled = true
	shutdownHandler := shutdown.New(syscall.SIGINT, syscall.SIGTERM)

	// load routes
	r := chi.NewRouter()
	router.AddAPIRoutes(r, conf, readDataService)
	router.AddLiteRoutes(r, conf, readDataService)

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
