package router

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/service/site"
)

func GetSite(sites map[string]*site.Site) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(res http.ResponseWriter, req *http.Request) {
				siteName := chi.URLParam(req, "siteName")
				st, ok := sites[siteName]
				if !ok {
					fmt.Println(sites, siteName, st, ok)
					fmt.Fprint(res, `{"error": "site not found"}`)
					return
				}
				ctx := context.WithValue(req.Context(), "site", st)
				next.ServeHTTP(res, req.WithContext(ctx))
			},
		)
	}
}
func GetBook(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			idHash := chi.URLParam(req, "idHash")
			st := req.Context().Value("site").(*site.Site)
			var (
				bk  *model.Book
				err error
			)

			idHashArray := strings.Split(idHash, "-")
			if len(idHashArray) == 1 {
				id, _ := strconv.Atoi(idHashArray[0])
				bk, err = st.BookFromID(id)
			} else if len(idHashArray) == 2 {
				id, _ := strconv.Atoi(idHashArray[0])
				hash := idHashArray[1]
				bk, err = st.BookFromIDHash(id, hash)
			}
			if err != nil {
				fmt.Fprintf(res, `{"error": "book not found"}`)
				return
			}

			ctx := context.WithValue(req.Context(), "book", bk)
			next.ServeHTTP(res, req.WithContext(ctx))
		},
	)
}
func GetSearchParams(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			title := req.URL.Query().Get("title")
			ctx := context.WithValue(req.Context(), "title", title)

			writer := req.URL.Query().Get("writer")
			ctx = context.WithValue(ctx, "writer", writer)

			next.ServeHTTP(res, req.WithContext(ctx))
		},
	)
}
func GetPageParams(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			page, _ := strconv.Atoi(req.URL.Query().Get("page"))
			perPage, _ := strconv.Atoi(req.URL.Query().Get("per_page"))
			offset := page * perPage

			ctx := context.WithValue(req.Context(), "limit", perPage)
			ctx = context.WithValue(ctx, "offset", offset)

			next.ServeHTTP(res, req.WithContext(ctx))
		},
	)
}
