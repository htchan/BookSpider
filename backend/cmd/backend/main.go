package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	config "github.com/htchan/BookSpider/internal/config_new"
	repo "github.com/htchan/BookSpider/internal/repo/psql"
	"github.com/htchan/BookSpider/internal/router"
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
	for _, siteName := range conf.APIConfig.AvailableSiteNames {
		db, err := repo.OpenDatabase(siteName)
		if err != nil {
			log.Fatalf("load db Fail. site: %v; err: %v", siteName, err)
		}

		sema := semaphore.NewWeighted(int64(conf.SiteConfigs[siteName].MaxThreads))

		serv, err := service_new.LoadService(
			siteName, conf.SiteConfigs[siteName], db, sema, &ctx,
		)
		if err != nil {
			log.Fatalf("load service fail. site: %v, err: %v", siteName, err)
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
	log.Println("start http server")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("backend stopped: %v", err)
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
