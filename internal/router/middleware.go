package router

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ContextKey string

const (
	ContextKeyReqID        ContextKey = "req_id"
	ContextKeySiteName     ContextKey = "site_name"
	ContextKeyReadDataServ ContextKey = "read_data_serv"
	ContextKeyBook         ContextKey = "book"
	ContextKeyBookGroup    ContextKey = "book_group"
	ContextKeyTitle        ContextKey = "title"
	ContextKeyWriter       ContextKey = "writer"
	ContextKeyPage         ContextKey = "page"
	ContextKeyPerPage      ContextKey = "per_page"
	ContextKeyLimit        ContextKey = "limit"
	ContextKeyOffset       ContextKey = "offset"
	ContextKeyUriPrefix    ContextKey = "uri_prefix"
	ContextKeyFormat       ContextKey = "format"
)

func getTracer() trace.Tracer {
	tracer := otel.Tracer("htchan/BookSpider/api")
	return tracer
}

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			ctx, span := getTracer().Start(req.Context(), fmt.Sprintf("%s %s", req.Method, req.RequestURI))
			defer span.End()

			next.ServeHTTP(res, req.WithContext(ctx))
		},
	)
}

func GetReadDataServiceMiddleware(readDataServ service.ReadDataService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(res http.ResponseWriter, req *http.Request) {
				logger := zerolog.Ctx(req.Context())
				_, span := getTracer().Start(req.Context(), "get read data service middleware")
				defer span.End()

				if readDataServ == nil {
					span.SetStatus(codes.Error, "read data service not initialized")
					span.RecordError(errors.New("read data service not initialized"))

					logger.Error().Msg("read data service not initialized")
					writeError(res, http.StatusInternalServerError, errors.New("read data service not initialized"))
					return
				}

				ctx := context.WithValue(req.Context(), ContextKeyReadDataServ, readDataServ)
				next.ServeHTTP(res, req.WithContext(ctx))
			},
		)
	}
}
func GetSiteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			siteName := chi.URLParam(req, "siteName")

			ctx := context.WithValue(req.Context(), ContextKeySiteName, siteName)
			next.ServeHTTP(res, req.WithContext(ctx))
		},
	)
}
func GetBookMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			logger := zerolog.Ctx(req.Context())
			idHash := chi.URLParam(req, "idHash")
			site := req.Context().Value(ContextKeySiteName).(string)
			serv := req.Context().Value(ContextKeyReadDataServ).(service.ReadDataService)
			var (
				bk    *model.Book
				group *model.BookGroup
				err   error
			)

			_, span := getTracer().Start(req.Context(), "get book middleware")
			defer span.End()

			span.SetAttributes(
				attribute.String("id_hash", idHash),
			)

			idHashArray := strings.Split(idHash, "-")
			if len(idHashArray) == 1 {
				id := idHashArray[0]
				bk, group, err = serv.BookGroup(req.Context(), site, id, "")
			} else if len(idHashArray) == 2 {
				id, hash := idHashArray[0], idHashArray[1]
				bk, group, err = serv.BookGroup(req.Context(), site, id, hash)
			}
			if err != nil {
				span.SetStatus(codes.Error, "get book failed")
				span.RecordError(err)

				logger.
					Error().
					Err(err).
					Str("site", site).
					Str("id-hash", idHash).
					Msg("get book middleware failed")
				writeError(res, http.StatusNotFound, errors.New("book not found"))
				return
			}

			span.End()

			ctx := context.WithValue(req.Context(), ContextKeyBook, bk)
			ctx = context.WithValue(ctx, ContextKeyBookGroup, group)
			next.ServeHTTP(res, req.WithContext(ctx))
		},
	)
}
func GetSearchParamsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			title := req.URL.Query().Get("title")
			ctx := context.WithValue(req.Context(), ContextKeyTitle, title)

			writer := req.URL.Query().Get("writer")
			ctx = context.WithValue(ctx, ContextKeyWriter, writer)

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

			ctx := context.WithValue(req.Context(), ContextKeyLimit, perPage)
			ctx = context.WithValue(ctx, ContextKeyOffset, offset)
			ctx = context.WithValue(ctx, ContextKeyPage, page)
			ctx = context.WithValue(ctx, ContextKeyPerPage, perPage)

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
			ctx := context.WithValue(req.Context(), ContextKeyFormat, format)

			next.ServeHTTP(res, req.WithContext(ctx))
		},
	)
}
func logRequest() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(res http.ResponseWriter, req *http.Request) {
				requestID := uuid.New()

				ctx := context.WithValue(req.Context(), ContextKeyReqID, requestID)
				logger := log.With().
					Str("request_id", requestID.String()).
					Logger()

				start := time.Now().UTC().Truncate(5 * time.Second)
				next.ServeHTTP(res, req.WithContext(logger.WithContext(ctx)))

				logger.Info().
					Str("path", req.URL.String()).
					Str("duration", time.Since(start).String()).
					Msg("request handled")
			},
		)
	}
}

func SetUriPrefixMiddleware(uriPrefix string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(res http.ResponseWriter, req *http.Request) {
				ctx := context.WithValue(req.Context(), ContextKeyUriPrefix, uriPrefix)
				next.ServeHTTP(res, req.WithContext(ctx))
			},
		)
	}
}
