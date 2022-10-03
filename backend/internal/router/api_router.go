package router

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/htchan/BookSpider/internal/service/site"
)

var UnauthorizedError = errors.New("unauthorized")
var InvalidParamsError = errors.New("invalid params")
var RecordNotFoundError = errors.New("record not found")

func writeError(res http.ResponseWriter, statusCode int, err error) {
	res.WriteHeader(statusCode)
	fmt.Fprintln(res, fmt.Sprintf(`{ "error": "%v" }`, err))
}

func AddAPIRoutes(router chi.Router, sites map[string]*site.Site) {
	api_route_prefix := os.Getenv("BOOK_SPIDER_API_ROUTE_PREFIX")
	if api_route_prefix == "" {
		api_route_prefix = "/api/novel"
	}

	router.Route(api_route_prefix, func(router chi.Router) {
		router.Use(
			cors.Handler(
				cors.Options{
					AllowedOrigins: []string{"*"},
					AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
					AllowedHeaders: []string{"*"},
					MaxAge:         300, // Maximum value not ignored by any of major browsers
				},
			),
		)

		router.Get("/info", GeneralInfoAPIHandler(sites))

		router.Route("/sites/{siteName}", func(router chi.Router) {
			router.Use(GetSite(sites))
			router.Get("/", SiteInfoAPIHandler)

			router.Route("/books", func(router chi.Router) {
				router.With(GetSearchParams).With(GetPageParams).Get("/search", BookSearchAPIHandler)
				router.With(GetPageParams).Get("/random", BookRandomAPIHandler)

				router.Route("/{idHash:\\d+(-[\\w]+)?}", func(router chi.Router) {
					// idHash format is <id>-<hash>
					router.Use(GetBook)
					router.With().Get("/", BookInfoAPIHandler)
					router.Get("/download", BookDownloadAPIHandler)
				})
			})
		})

		router.Get("/db-stats", DBStatsAPIHandler(sites))
	})
}
