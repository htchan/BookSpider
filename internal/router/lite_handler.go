package router

import (
	"embed"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/htchan/BookSpider/internal/format/v1"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/htchan/BookSpider/internal/service"
	"github.com/rs/zerolog"
)

//go:embed templates/*
var files embed.FS

var customTemplateFunc = template.FuncMap{
	"arr": func(eles ...any) []any { return eles },
}

// @Summary		Home page
// @description	home page
// @Tags			book-spider-lite
// @Produce		html
// @Success		200	{string}	string
// @Router			/lite/book-spider/ [get]
func GeneralLiteHandler(services map[string]service.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		logger := zerolog.Ctx(req.Context())
		uriPrefix := req.Context().Value(ContextKeyUriPrefix).(string)
		t, err :=
			new(template.Template).
				Funcs(customTemplateFunc).
				ParseFS(
					files,
					"templates/sites.html",
					"templates/components/site-card.html",
					"templates/styles/site-button.html",
				)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			logger.Error().Err(err).Msg("general lite handler parse fs fail")
			res.Write([]byte(err.Error()))
			return
		}
		execErr := t.ExecuteTemplate(res, "sites.html", struct {
			Services  map[string]service.Service
			UriPrefix string
		}{Services: services, UriPrefix: uriPrefix})
		if execErr != nil {
			res.WriteHeader(http.StatusInternalServerError)
			logger.Error().Err(execErr).Msg("compute response failed")
		}
	}
}

// @Summary		site info page
// @description	site info page
// @Tags			book-spider-lite
// @Produce		html
// @Param			siteName	path		string	true	"site name"
// @Success		200			{string}	string
// @Router			/lite/book-spider/sites/{siteName} [get]
func SiteLiteHandlerfunc(res http.ResponseWriter, req *http.Request) {
	logger := zerolog.Ctx(req.Context())
	uriPrefix := req.Context().Value(ContextKeyUriPrefix).(string)
	t, err := new(template.Template).
		Funcs(customTemplateFunc).
		ParseFS(files, "templates/site.html")
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		logger.Error().Err(err).Msg("site lite handler parse fs fail")
		res.Write([]byte(err.Error()))
		return
	}

	siteName := req.Context().Value(ContextKeySiteName).(string)
	serv := req.Context().Value(ContextKeyReadDataServ).(service.ReadDataService)
	execErr := t.ExecuteTemplate(res, "site.html", struct {
		Name      string
		UriPrefix string
		Summary   repo.Summary
	}{
		Name:      siteName,
		UriPrefix: uriPrefix,
		Summary:   serv.Stats(req.Context(), siteName),
	})
	if execErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Error().Err(execErr).Msg("compute response failed")
	}
}

// @Summary		Search result page
// @description	search result page
// @Tags			book-spider-lite
// @Produce		html
// @Param			siteName	path		string	true	"site name"
// @Success		200			{string}	string
// @Router			/lite/book-spider/sites/{siteName}/search [get]
func SearchLiteHandler(res http.ResponseWriter, req *http.Request) {
	logger := zerolog.Ctx(req.Context())
	uriPrefix := req.Context().Value(ContextKeyUriPrefix).(string)
	t, err := new(template.Template).
		Funcs(customTemplateFunc).
		ParseFS(
			files,
			"templates/result.html",
			"templates/components/book-card.html",
			"templates/components/pagination.html",
			"templates/styles/book-box.html",
			"templates/styles/pagination.html",
		)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	serv := req.Context().Value(ContextKeyReadDataServ).(service.ReadDataService)
	title := req.Context().Value(ContextKeyTitle).(string)
	writer := req.Context().Value(ContextKeyWriter).(string)
	page := req.Context().Value(ContextKeyPage).(int)
	perPage := req.Context().Value(ContextKeyPerPage).(int)
	limit := req.Context().Value(ContextKeyLimit).(int)
	offset := req.Context().Value(ContextKeyOffset).(int)
	if limit == 0 {
		limit = 10
	}

	bks, err := serv.SearchBooks(req.Context(), title, writer, limit, offset)

	if err != nil {
		res.WriteHeader(404)
		fmt.Fprint(res, "books not found")
		return
	}

	execErr := t.ExecuteTemplate(res, "result.html", struct {
		Name           string
		UriPrefix      string
		Books          []model.Book
		Title          string
		Writer         string
		PreviousPage   int
		NextPage       int
		PerPage        int
		ShowPagination bool
	}{
		Name:           "Search Result",
		UriPrefix:      uriPrefix,
		Books:          bks,
		Title:          title,
		Writer:         writer,
		PreviousPage:   page - 1,
		NextPage:       page + 1,
		PerPage:        perPage,
		ShowPagination: true,
	})
	if execErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Error().Err(execErr).Msg("compute response failed")
		res.Write([]byte(execErr.Error()))
	}
}

// @Summary		Random result page
// @description	random result page
// @Tags			book-spider-lite
// @Produce		html
// @Param			siteName	path		string	true	"site name"
// @Success		200			{string}	string
// @Router			/lite/book-spider/sites/{siteName}/random [get]
func RandomLiteHandler(res http.ResponseWriter, req *http.Request) {
	logger := zerolog.Ctx(req.Context())
	uriPrefix := req.Context().Value(ContextKeyUriPrefix).(string)
	t, err := new(template.Template).
		Funcs(customTemplateFunc).
		ParseFS(
			files,
			"templates/result.html",
			"templates/components/book-card.html",
			"templates/components/pagination.html",
			"templates/styles/book-box.html",
			"templates/styles/pagination.html",
		)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		logger.Error().Err(err).Msg("random lite handler parse fs fail")
		return
	}

	serv := req.Context().Value(ContextKeyReadDataServ).(service.ReadDataService)
	limit := req.Context().Value(ContextKeyLimit).(int)
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
		Name           string
		UriPrefix      string
		Books          []model.Book
		ShowPagination bool
	}{
		Name:           "Random",
		UriPrefix:      uriPrefix,
		Books:          bks,
		ShowPagination: false,
	})
	if execErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Error().Err(execErr).Msg("compute response failed")
	}
}

// @Summary		Book info page
// @description	book info page
// @Tags			book-spider-lite
// @Produce		html
// @Param			siteName	path		string	true	"site name"
// @Param			idHash		path		string	true	"id and hash in format <id>[-<hash>]. -<hash is optional"
// @Success		200			{string}	string
// @Router			/lite/book-spider/sites/{siteName}/books/{idHash} [get]
func BookLiteHandler(res http.ResponseWriter, req *http.Request) {
	logger := zerolog.Ctx(req.Context())
	uriPrefix := req.Context().Value(ContextKeyUriPrefix).(string)
	t, err := new(template.Template).
		Funcs(customTemplateFunc).
		ParseFS(
			files,
			"templates/book.html",
			"templates/components/book-card.html",
			"templates/styles/book-box.html",
		)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		logger.Error().Err(err).Msg("book lite handler parse fs fail")
		return
	}

	bk := req.Context().Value(ContextKeyBook).(*model.Book)
	group := req.Context().Value(ContextKeyBookGroup).(*model.BookGroup)

	execErr := t.ExecuteTemplate(res, "book.html", struct {
		UriPrefix string
		Book      *model.Book
		Group     *model.BookGroup
	}{
		UriPrefix: uriPrefix,
		Book:      bk,
		Group:     group,
	})
	if execErr != nil {
		res.WriteHeader(http.StatusInternalServerError)
		logger.Error().Err(execErr).Msg("compute response failed")
	}
}

// @Summary		Book download page
// @description	book download page
// @Tags			book-spider-lite
// @Produce		html
// @Param			siteName	path		string	true	"site name"
// @Param			idHash		path		string	true	"id and hash in format <id>[-<hash>]. -<hash is optional"
// @Param			format		query		string	true	"txt (default) or epub"
// @Success		200			{string}	string
// @Router			/lite/book-spider/sites/{siteName}/books/{idHash}/download [get]
func DownloadLiteHandler(res http.ResponseWriter, req *http.Request) {
	logger := zerolog.Ctx(req.Context())
	serv := req.Context().Value(ContextKeyReadDataServ).(service.ReadDataService)
	bk := req.Context().Value(ContextKeyBook).(*model.Book)
	formatStr := req.Context().Value(ContextKeyFormat).(string)

	content, err := serv.BookContent(req.Context(), bk)
	if err != nil {
		res.WriteHeader(500)
		logger.Error().Err(err).Str("book", bk.String()).Msg("download lite handler failed")
		return
	}

	switch formatStr {
	case "epub":
		formatServ := format.NewService()
		chapters, err := formatServ.ChaptersFromTxt(req.Context(), strings.NewReader(content))
		if err != nil {
			res.WriteHeader(500)
			logger.Error().Err(err).Str("book", bk.String()).Msg("download lite handler failed")
			return
		}
		fileName := fmt.Sprintf("%s-%s.epub", bk.Title, bk.Writer.Name)
		res.Header().Set("Content-Type", "application/epub+zip; charset=utf-8")
		res.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
		formatServ.WriteBookEpub(req.Context(), bk, chapters, res)
	default:
		fileName := fmt.Sprintf("%s-%s.txt", bk.Title, bk.Writer.Name)
		res.Header().Set("Content-Type", "text/txt; charset=utf-8")
		res.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
		fmt.Fprint(res, content)
	}
}
