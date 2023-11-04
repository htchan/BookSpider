package router

import (
	"embed"
	"fmt"
	"net/http"
	"text/template"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/htchan/BookSpider/internal/service"
	"github.com/rs/zerolog/log"
)

//go:embed templates/*
var files embed.FS

var customTemplateFunc = template.FuncMap{
	"arr": func(eles ...any) []any { return eles },
}

func GeneralLiteHandler(services map[string]service.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		uriPrefix := req.Context().Value(URI_PREFIX_KEY).(string)
		t, err :=
			new(template.Template).
				Funcs(customTemplateFunc).
				ParseFS(files, "templates/sites.html", "templates/components/site-card.html")
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Msg("general lite handler parse fs fail")
			res.Write([]byte(err.Error()))
			return
		}
		execErr := t.ExecuteTemplate(res, "sites.html", struct {
			Services  map[string]service.Service
			UriPrefix string
		}{Services: services, UriPrefix: uriPrefix})
		if execErr != nil {
			res.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(execErr).Msg("compute response failed")
		}
	}
}

func SiteLiteHandlerfunc(res http.ResponseWriter, req *http.Request) {
	uriPrefix := req.Context().Value(URI_PREFIX_KEY).(string)
	t, err := new(template.Template).
		Funcs(customTemplateFunc).
		ParseFS(files, "templates/site.html")
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		log.Error().Err(err).Msg("site lite handler parse fs fail")
		res.Write([]byte(err.Error()))
		return
	}

	serv := req.Context().Value(SERV_KEY).(service.Service)
	execErr := t.ExecuteTemplate(res, "site.html", struct {
		Name      string
		UriPrefix string
		Summary   repo.Summary
	}{
		Name:      serv.Name(),
		UriPrefix: uriPrefix,
		Summary:   serv.Stats(req.Context()),
	})
	if execErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Error().Err(execErr).Msg("compute response failed")
	}
}

func SearchLiteHandler(res http.ResponseWriter, req *http.Request) {
	uriPrefix := req.Context().Value(URI_PREFIX_KEY).(string)
	t, err := new(template.Template).
		Funcs(customTemplateFunc).
		ParseFS(files, "templates/result.html", "templates/components/book-card.html")
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		log.Error().Err(err).Msg("search lite handler parse fs fail")
		return
	}

	serv := req.Context().Value(SERV_KEY).(service.Service)
	title := req.Context().Value(TITLE_KEY).(string)
	writer := req.Context().Value(WRITER_KEY).(string)
	limit := req.Context().Value(LIMIT_KEY).(int)
	offset := req.Context().Value(OFFSET_KEY).(int)
	if limit == 0 {
		limit = 10
	}

	bks, err := serv.QueryBooks(req.Context(), title, writer, limit, offset)

	if err != nil {
		res.WriteHeader(404)
		fmt.Fprint(res, "books not found")
		return
	}
	pageNo := (offset / limit)

	execErr := t.ExecuteTemplate(res, "result.html", struct {
		Name       string
		UriPrefix  string
		Books      []model.Book
		lastPageNo int
		nextPageNo int
	}{
		Name:       serv.Name(),
		UriPrefix:  uriPrefix,
		Books:      bks,
		lastPageNo: pageNo - 1,
		nextPageNo: pageNo + 1,
	})
	if execErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Error().Err(execErr).Msg("compute response failed")
		res.Write([]byte(execErr.Error()))
	}
}

func RandomLiteHandler(res http.ResponseWriter, req *http.Request) {
	uriPrefix := req.Context().Value(URI_PREFIX_KEY).(string)
	t, err := new(template.Template).
		Funcs(customTemplateFunc).
		ParseFS(files, "templates/result.html", "templates/components/book-card.html")
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		log.Error().Err(err).Msg("random lite handler parse fs fail")
		return
	}

	serv := req.Context().Value(SERV_KEY).(service.Service)
	limit := req.Context().Value(LIMIT_KEY).(int)
	if limit == 0 {
		limit = 10
	}

	bks, err := serv.RandomBooks(req.Context(), limit)

	if err != nil {
		res.WriteHeader(404)
		fmt.Fprint(res, "books not found")
		return
	}

	execErr := t.ExecuteTemplate(res, "result.html", struct {
		Name      string
		UriPrefix string
		Books     []model.Book
	}{
		Name:      serv.Name(),
		UriPrefix: uriPrefix,
		Books:     bks,
	})
	if execErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Error().Err(execErr).Msg("compute response failed")
	}
}

func BookLiteHandler(res http.ResponseWriter, req *http.Request) {
	uriPrefix := req.Context().Value(URI_PREFIX_KEY).(string)
	t, err := new(template.Template).
		Funcs(customTemplateFunc).
		ParseFS(files, "templates/book.html", "templates/components/book-card.html")
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		log.Error().Err(err).Msg("book lite handler parse fs fail")
		return
	}

	serv := req.Context().Value(SERV_KEY).(service.Service)
	bk := req.Context().Value(BOOK_KEY).(*model.Book)
	group := req.Context().Value(BOOK_GROUP_KEY).(*model.BookGroup)

	execErr := t.ExecuteTemplate(res, "book.html", struct {
		Name      string
		UriPrefix string
		Book      *model.Book
		Group     *model.BookGroup
	}{
		Name:      serv.Name(),
		UriPrefix: uriPrefix,
		Book:      bk,
		Group:     group,
	})
	if execErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Error().Err(execErr).Msg("compute response failed")
	}
}

func DownloadLiteHandler(res http.ResponseWriter, req *http.Request) {
	serv := req.Context().Value(SERV_KEY).(service.Service)
	bk := req.Context().Value(BOOK_KEY).(*model.Book)
	content, err := serv.BookContent(req.Context(), bk)
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
