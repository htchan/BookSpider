package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/htchan/BookSpider/internal/config/v1"
	"github.com/htchan/BookSpider/internal/service/v1"
)

func AddLiteRoutes(router chi.Router, conf *config.BackendConfig, readDataServices service.ReadDataService) {
	router.Route(conf.LiteRoutePrefix, func(router chi.Router) {
		router.Use(logRequest())
		router.Use(TraceMiddleware)
		router.Use(SetUriPrefixMiddleware(conf.LiteRoutePrefix))
		router.Use(GetReadDataServiceMiddleware(readDataServices))

		router.Route("/sites/{siteName}", func(router chi.Router) {
			router.Use(GetSiteMiddleware)
			router.Get("/", SiteLiteHandlerfunc)

			router.With(GetSearchParamsMiddleware).With(GetPageParamsMiddleware).Get("/search", SearchLiteHandler)
			router.With(GetPageParamsMiddleware).Get("/random", RandomLiteHandler)

			router.Route("/books", func(router chi.Router) {
				router.Route("/{idHash:\\d+(-[\\w]+)?}", func(router chi.Router) {
					// idHash format is <id>-<hash>
					router.Use(GetBookMiddleware)
					router.Get("/", BookLiteHandler)
					router.With(GetDownloadParamsMiddleware).Get("/download", DownloadLiteHandler)
				})
			})
		})

		router.Get("/", GeneralLiteHandler(conf.AvailableSites))
		router.With(GetSearchParamsMiddleware).With(GetPageParamsMiddleware).Get("/search", SearchLiteHandler)
		router.With(GetPageParamsMiddleware).Get("/random", RandomLiteHandler)

	})
}
