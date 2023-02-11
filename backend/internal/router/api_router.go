package router

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	config "github.com/htchan/BookSpider/internal/config_new"
	service_new "github.com/htchan/BookSpider/internal/service_new"
)

var UnauthorizedError = errors.New("unauthorized")
var InvalidParamsError = errors.New("invalid params")
var RecordNotFoundError = errors.New("record not found")

func writeError(res http.ResponseWriter, statusCode int, err error) {
	res.WriteHeader(statusCode)
	fmt.Fprintln(res, fmt.Sprintf(`{ "error": "%v" }`, err))
}

func AddAPIRoutes(router chi.Router, conf config.APIConfig, services map[string]service_new.Service) {
	router.Route(conf.APIRoutePrefix, func(router chi.Router) {
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

		router.Get("/info", GeneralInfoAPIHandler(services))

		router.Route("/sites/{siteName}", func(router chi.Router) {
			router.Use(GetSiteMiddleware(services))
			router.Get("/", SiteInfoAPIHandler)

			router.Route("/books", func(router chi.Router) {
				router.With(GetSearchParamsMiddleware).With(GetPageParamsMiddleware).Get("/search", BookSearchAPIHandler)
				router.With(GetPageParamsMiddleware).Get("/random", BookRandomAPIHandler)

				router.Route("/{idHash:\\d+(-[\\w]+)?}", func(router chi.Router) {
					// idHash format is <id>-<hash>
					router.Use(GetBookMiddleware)
					router.With().Get("/", BookInfoAPIHandler)
					router.Get("/download", BookDownloadAPIHandler)
				})
			})
		})

		router.Get("/db-stats", DBStatsAPIHandler(services))
	})
}
