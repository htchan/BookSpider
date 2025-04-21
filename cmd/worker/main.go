package main

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/htchan/BookSpider/internal/common"
	"github.com/htchan/BookSpider/internal/config/v2"
	repo "github.com/htchan/BookSpider/internal/repo/sqlc"
	"github.com/htchan/BookSpider/internal/service"
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

func CalculateNextRunTime(conf *config.ScheduleConfig) time.Time {
	result := time.Now().UTC().Truncate(24 * time.Hour)

	for true {
		result = time.Date(result.Year(), result.Month(), conf.InitDate, conf.InitHour, conf.InitMinute, 0, 0, time.UTC)
		if result.Weekday() != conf.MatchWeekday {
			nDaysLater := int(conf.MatchWeekday - result.Weekday())
			if nDaysLater < 0 {
				nDaysLater += 7
			}

			result = result.AddDate(0, 0, nDaysLater)
		}

		if time.Now().Before(result) {
			return result
		}

		result = result.AddDate(0, conf.IntervalMonth, conf.IntervalDay)
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

	tp, err := otelProvider(conf.TraceConfig)
	if err != nil {
		log.Error().Err(err).Msg("init tracer failed")
	}
	defer tp.Shutdown(context.Background())

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
		until := CalculateNextRunTime(&conf.ScheduleConfig)
		log.Log().Time("scheduled_at", until).Msg("start sleep")
		time.Sleep(time.Until(until))
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
