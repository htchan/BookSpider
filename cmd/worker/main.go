package main

import (
	"context"
	"os"
	"syscall"
	"time"

	"github.com/htchan/BookSpider/internal/common"
	"github.com/htchan/BookSpider/internal/config/v1"
	repo "github.com/htchan/BookSpider/internal/repo/sqlc"
	"github.com/htchan/BookSpider/internal/service/v1"
	bookprocess "github.com/htchan/BookSpider/internal/tasks/nats/book_process"
	"github.com/nats-io/nats.go/jetstream"

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

	conf, confErr := config.LoadWorkerConfig()
	if confErr != nil {
		log.Error().Err(confErr).Msg("load worker config")
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
	clients, err := common.LoadClients(context.Background(), conf.Clients)
	if err != nil {
		log.Error().Err(err).Msg("load clients fail")
	}

	bookService := service.NewBookService(clients, rpo, conf.Common.StoragePath)

	shutdown.LogEnabled = true
	shutdownHandler := shutdown.New(syscall.SIGINT, syscall.SIGTERM)

	nc, err := common.ConnectNatsQueue(&conf.Nats)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to nats server")
	}

	processTasks := make([]jetstream.ConsumeContext, 0, len(conf.AvailableSites))
	bookProcessTasks := bookprocess.NewTaskSet(nc, bookService, conf.AvailableSites)
	for _, task := range bookProcessTasks {
		consumer, err := task.Subscribe(context.Background())
		if err != nil {
			log.Fatal().Err(err).
				Str("task", "book-process").
				Msg("failed to subscribe to nats server")
		}

		processTasks = append(processTasks, consumer)
	}

	// TODO: register batch task

	// TODO: register shutdown handler for batch task

	for _, task := range processTasks {
		shutdownHandler.Register("process task", func() error {
			task.Stop()

			return nil
		})
	}

	shutdownHandler.Register("nats connect", func() error {
		nc.Close()

		return nil
	})
	shutdownHandler.Register("database", db.Close)
	shutdownHandler.Register("tracer", func() error {
		return tp.Shutdown(context.Background())
	})

	shutdownHandler.Listen(60 * time.Second)
}
