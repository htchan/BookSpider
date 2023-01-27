package router

import (
	"os"

	"github.com/go-chi/chi/v5"
	service_new "github.com/htchan/BookSpider/internal/service_new"
)

func AddLiteRoutes(router chi.Router, services map[string]service_new.Service) {
	api_route_prefix := os.Getenv("BOOK_SPIDER_LITE_ROUTE_PREFIX")
	if api_route_prefix == "" {
		api_route_prefix = "/lite/novel"
	}

	router.Route(api_route_prefix, func(router chi.Router) {
		router.Route("/sites/{siteName}", func(router chi.Router) {
			router.Use(GetSiteMiddleware(services))
			router.Get("/", SiteLiteHandlerfunc)

			router.With(GetSearchParamsMiddleware).With(GetPageParamsMiddleware).Get("/search", SearchLiteHandler)
			router.With(GetPageParamsMiddleware).Get("/random", RandomLiteHandler)

			router.Route("/books", func(router chi.Router) {
				router.Route("/{idHash:\\d+(-[\\w]+)?}", func(router chi.Router) {
					// idHash format is <id>-<hash>
					router.Use(GetBookMiddleware)
					router.Get("/", BookLiteHandler)
					router.Get("/download", DownloadLiteHandler)
				})
			})
		})

		router.Get("/", GeneralLiteHandler(services))
	})
}
