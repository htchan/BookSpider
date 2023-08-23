package router

import (
	"embed"
	"fmt"
	"net/http"
	"text/template"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	service_new "github.com/htchan/BookSpider/internal/service_new"
	"github.com/rs/zerolog/log"
)

//go:embed templates/*
var files embed.FS

func GeneralLiteHandler(services map[string]service_new.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFS(files, "templates/sites.html", "templates/components/site-card.html")
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			log.Error().Err(err).Msg("general lite handler parse fs fail")
			return
		}
		t.Execute(res, services)
	}
}

func SiteLiteHandlerfunc(res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFS(files, "templates/site.html")
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		log.Error().Err(err).Msg("site lite handler parse fs fail")
		return
	}

	serv := req.Context().Value(SERV_KEY).(service_new.Service)
	t.Execute(res, struct {
		Name    string
		Summary repo.Summary
	}{
		Name:    serv.Name(),
		Summary: serv.Stats(),
	})
}

func SearchLiteHandler(res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFS(files, "templates/result.html", "templates/components/book-card.html")
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		log.Error().Err(err).Msg("search lite handler parse fs fail")
		return
	}

	serv := req.Context().Value(SERV_KEY).(service_new.Service)
	title := req.Context().Value(TITLE_KEY).(string)
	writer := req.Context().Value(WRITER_KEY).(string)
	limit := req.Context().Value(LIMIT_KEY).(int)
	offset := req.Context().Value(OFFSET_KEY).(int)
	if limit == 0 {
		limit = 10
	}

	bks, err := serv.QueryBooks(title, writer, limit, offset)

	if err != nil {
		res.WriteHeader(404)
		fmt.Fprint(res, "books not found")
		return
	}
	pageNo := (offset / limit)

	t.Execute(res, struct {
		Name       string
		Books      []model.Book
		lastPageNo int
		nextPageNo int
	}{
		Name:       serv.Name(),
		Books:      bks,
		lastPageNo: pageNo - 1,
		nextPageNo: pageNo + 1,
	})
}

func RandomLiteHandler(res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFS(files, "templates/result.html", "templates/components/book-card.html")
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		log.Error().Err(err).Msg("random lite handler parse fs fail")
		return
	}

	serv := req.Context().Value(SERV_KEY).(service_new.Service)
	limit := req.Context().Value(LIMIT_KEY).(int)
	if limit == 0 {
		limit = 10
	}

	bks, err := serv.RandomBooks(limit)

	if err != nil {
		res.WriteHeader(404)
		fmt.Fprint(res, "books not found")
		return
	}

	t.Execute(res, struct {
		Name  string
		Books []model.Book
	}{
		Name:  serv.Name(),
		Books: bks,
	})
}

func BookLiteHandler(res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFS(files, "templates/book.html", "templates/components/book-card.html")
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		log.Error().Err(err).Msg("book lite handler parse fs fail")
		return
	}

	serv := req.Context().Value(SERV_KEY).(service_new.Service)
	bk := req.Context().Value(BOOK_KEY).(*model.Book)
	group := req.Context().Value(BOOK_GROUP_KEY).(*model.BookGroup)

	t.Execute(res, struct {
		Name  string
		Book  *model.Book
		Group *model.BookGroup
	}{
		Name:  serv.Name(),
		Book:  bk,
		Group: group,
	})
}

func DownloadLiteHandler(res http.ResponseWriter, req *http.Request) {
	serv := req.Context().Value(SERV_KEY).(service_new.Service)
	bk := req.Context().Value(BOOK_KEY).(*model.Book)
	content, err := serv.BookContent(bk)
	if err != nil {
		res.WriteHeader(500)
		log.Error().Err(err).Str("book", bk.String()).Msg("download lite handler failed")
		return
	} else {
		fileName := fmt.Sprintf("%s-%s.txt", bk.Title, bk.Writer.Name)
		res.Header().Set("Content-Type", "text/txt; charset=utf-8")
		res.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
		fmt.Fprint(res, content)
	}
}
