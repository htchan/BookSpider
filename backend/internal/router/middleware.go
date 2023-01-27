package router

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/htchan/BookSpider/internal/model"
	service_new "github.com/htchan/BookSpider/internal/service_new"
)

const (
	SERV_KEY   = "serv"
	BOOK_KEY   = "book"
	TITLE_KEY  = "title"
	WRITER_KEY = "writer"
	LIMIT_KEY  = "limit"
	OFFSET_KEY = "offset"
)

func GetSiteMiddleware(services map[string]service_new.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(res http.ResponseWriter, req *http.Request) {
				siteName := chi.URLParam(req, "siteName")
				serv, ok := services[siteName]
				if !ok {
					fmt.Println(services, siteName, serv, ok)
					fmt.Fprint(res, `{"error": "site not found"}`)
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
			idHash := chi.URLParam(req, "idHash")
			serv := req.Context().Value(SERV_KEY).(service_new.Service)
			var (
				bk  *model.Book
				err error
			)

			idHashArray := strings.Split(idHash, "-")
			if len(idHashArray) == 1 {
				id, _ := strconv.Atoi(idHashArray[0])
				bk, err = serv.Book(id, "")
			} else if len(idHashArray) == 2 {
				id, _ := strconv.Atoi(idHashArray[0])
				hash := idHashArray[1]
				bk, err = serv.Book(id, hash)
			}
			if err != nil {
				fmt.Printf("cannot query book. site: %v; id-hash: %v; err: %v", serv.Name(), idHash, err)
				fmt.Fprintf(res, `{"error": "book not found"}`)
				return
			}

			ctx := context.WithValue(req.Context(), BOOK_KEY, bk)
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
