package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/router"
	"github.com/htchan/BookSpider/internal/service/site"
)

func main() {
	configLocation := os.Getenv("ASSETS_LOCATION") + "/config"
	var err error

	// TODO: load backend config
	backendConfig, err := config.LoadBackendConfig(configLocation)
	if err != nil {
		fmt.Printf("load backend config: %v", err)
		return
	}

	sites, err := site.LoadSitesFromConfigDirectory(configLocation, backendConfig.EnabledSiteNames)
	if err != nil {
		log.Fatal(err)
	}

	// load routes
	r := chi.NewRouter()
	if backendConfig.ContainsRoute(config.RouteAPIKey) {
		router.AddAPIRoutes(r, sites)
	}

	if backendConfig.ContainsRoute(config.RouteLiteKey) {
		router.AddLiteRoutes(r, sites)
	}

	server := http.Server{
		Addr:         ":9105",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
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
