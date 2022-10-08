package router

import (
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/htchan/BookSpider/internal/service/site"
)

func AddLiteRoutes(router chi.Router, sites map[string]*site.Site) {
	api_route_prefix := os.Getenv("BOOK_SPIDER_LITE_ROUTE_PREFIX")
	if api_route_prefix == "" {
		api_route_prefix = "/lite/novel"
	}

	router.Route(api_route_prefix, func(router chi.Router) {
		router.Route("/sites/{siteName}", func(router chi.Router) {
			router.Use(GetSite(sites))
			router.Get("/", SiteLiteHandlerfunc)

			router.With(GetSearchParams).With(GetPageParams).Get("/search", SearchLiteHandler)
			router.With(GetPageParams).Get("/random", RandomLiteHandler)

			router.Route("/books", func(router chi.Router) {
				router.Route("/{idHash:\\d+(-[\\w]+)?}", func(router chi.Router) {
					// idHash format is <id>-<hash>
					router.Use(GetBook)
					router.Get("/", BookLiteHandler)
					router.Get("/download", DownloadLiteHandler)
				})
			})
		})

		router.Get("/", GeneralLiteHandler(sites))
	})
}
