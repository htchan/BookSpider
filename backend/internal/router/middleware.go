package router

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ContextKey string

const (
	SERV_KEY       ContextKey = "serv"
	BOOK_KEY       ContextKey = "book"
	BOOK_GROUP_KEY ContextKey = "book_group"
	TITLE_KEY      ContextKey = "title"
	WRITER_KEY     ContextKey = "writer"
	LIMIT_KEY      ContextKey = "limit"
	OFFSET_KEY     ContextKey = "offset"
	URI_PREFIX_KEY ContextKey = "uri_prefix"
	FORMAT_KEY     ContextKey = "format"
)

func GetSiteMiddleware(services map[string]service.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(res http.ResponseWriter, req *http.Request) {
				logger := zerolog.Ctx(req.Context())
				siteName := chi.URLParam(req, "siteName")
				availableSites := make([]string, 0, len(services))
				for key := range services {
					availableSites = append(availableSites, key)
				}
				serv, ok := services[siteName]
				if !ok {
					logger.
						Error().
						Err(errors.New("site not found")).
						Str("site", siteName).
						Strs("available sites", availableSites).
						Msg("get site middleware failed")
					writeError(res, http.StatusNotFound, errors.New("site not found"))
					return
				}
				ctx := context.WithValue(req.Context(), SERV_KEY, serv)
				next.ServeHTTP(res, req.WithContext(ctx))
			},
		)
	}
}
func GetBookMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			logger := zerolog.Ctx(req.Context())
			idHash := chi.URLParam(req, "idHash")
			serv := req.Context().Value(SERV_KEY).(service.Service)
			var (
				bk    *model.Book
				group *model.BookGroup
				err   error
			)

			idHashArray := strings.Split(idHash, "-")
			if len(idHashArray) == 1 {
				id := idHashArray[0]
				bk, group, err = serv.BookGroup(req.Context(), id, "")
			} else if len(idHashArray) == 2 {
				id, hash := idHashArray[0], idHashArray[1]
				bk, group, err = serv.BookGroup(req.Context(), id, hash)
			}
			if err != nil {
				logger.
					Error().
					Err(err).
					Str("site", serv.Name()).
					Str("id-hash", idHash).
					Msg("get book middleware failed")
				writeError(res, http.StatusNotFound, errors.New("book not found"))
				return
			}

			ctx := context.WithValue(req.Context(), BOOK_KEY, bk)
			ctx = context.WithValue(ctx, BOOK_GROUP_KEY, group)
			next.ServeHTTP(res, req.WithContext(ctx))
		},
	)
}
func GetSearchParamsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			title := req.URL.Query().Get("title")
			ctx := context.WithValue(req.Context(), TITLE_KEY, title)

			writer := req.URL.Query().Get("writer")
			ctx = context.WithValue(ctx, WRITER_KEY, writer)

			next.ServeHTTP(res, req.WithContext(ctx))
		},
	)
}
func GetPageParamsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			page, _ := strconv.Atoi(req.URL.Query().Get("page"))
			perPage, _ := strconv.Atoi(req.URL.Query().Get("per_page"))
			offset := page * perPage

			ctx := context.WithValue(req.Context(), LIMIT_KEY, perPage)
			ctx = context.WithValue(ctx, OFFSET_KEY, offset)

			next.ServeHTTP(res, req.WithContext(ctx))
		},
	)
}
func GetDownloadParamsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			format := req.URL.Query().Get("format")
			if format == "" {
				format = "txt"
			}
			ctx := context.WithValue(req.Context(), FORMAT_KEY, format)

			next.ServeHTTP(res, req.WithContext(ctx))
		},
	)
}

func ZerologMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			requestUUID := uuid.New()
			logger := log.With().
				Str("request_uuid", requestUUID.String()).
				Str("method", req.Method).
				Str("endpoint", req.RequestURI).Logger()
			logger.Info().Msg("API called")

			next.ServeHTTP(res, req.WithContext(logger.WithContext(req.Context())))
		},
	)
}

func SetUriPrefixMiddleware(uriPrefix string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(res http.ResponseWriter, req *http.Request) {
				ctx := context.WithValue(req.Context(), URI_PREFIX_KEY, uriPrefix)
				next.ServeHTTP(res, req.WithContext(ctx))
			},
		)
	}
}
