package router

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/htchan/BookSpider/internal/service"
	"github.com/rs/zerolog"
)

// @Summary		Get all sites info
// @description	get all sites info
// @Tags			book-spider-api
// @Accept			json
// @Produce		json
// @Success		200	{object}	map[string]repo.Summary
// @Router			/api/book-spider/info [get]
func GeneralInfoAPIHandler(services map[string]service.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		servInfo := make(map[string]repo.Summary)
		for _, serv := range services {
			servInfo[serv.Name()] = serv.Stats(req.Context())
		}
		json.NewEncoder(res).Encode(servInfo)
	}
}

// @Summary		Get site info
// @description	get site info
// @Tags			book-spider-api
// @Accept			json
// @Produce		json
// @Param			siteName	path		string	true	"site name"
// @Success		200			{object}	repo.Summary
// @Failure		404			{object}	errResp
// @Router			/api/book-spider/sites/{siteName} [get]
func SiteInfoAPIHandler(res http.ResponseWriter, req *http.Request) {
	site := req.Context().Value(ContextKeySiteName).(string)
	serv := req.Context().Value(ContextKeyReadDataServ).(service.ReadDataService)
	json.NewEncoder(res).Encode(serv.Stats(req.Context(), site))
}

// @Summary		Search books
// @description	search books
// @Tags			book-spider-api
// @Accept			json
// @Produce		json
// @Param			siteName	path		string	true	"site name"
// @Success		200			{object}	booksResp
// @Failure		400			{object}	errResp
// @Router			/api/book-spider/sites/{siteName}/books/search [get]
func BookSearchAPIHandler(res http.ResponseWriter, req *http.Request) {
	logger := zerolog.Ctx(req.Context())
	serv := req.Context().Value(ContextKeyReadDataServ).(service.ReadDataService)
	title := req.Context().Value(ContextKeyTitle).(string)
	writer := req.Context().Value(ContextKeyWriter).(string)
	limit := req.Context().Value(ContextKeyLimit).(int)
	offset := req.Context().Value(ContextKeyOffset).(int)

	bks, err := serv.SearchBooks(req.Context(), title, writer, limit, offset)
	if err != nil {
		logger.Error().Err(err).Msg("query books failed")
		writeError(res, 400, err)
	} else {
		json.NewEncoder(res).Encode(booksResp{bks})
	}
}

// @Summary		List random books
// @description	list random books
// @Tags			book-spider-api
// @Accept			json
// @Produce		json
// @Param			siteName	path		string	true	"site name"
// @Success		200			{object}	booksResp
// @Failure		400			{object}	errResp
// @Router			/api/book-spider/sites/{siteName}/books/random [get]
func BookRandomAPIHandler(res http.ResponseWriter, req *http.Request) {
	logger := zerolog.Ctx(req.Context())
	serv := req.Context().Value(ContextKeyReadDataServ).(service.ReadDataService)
	limit := req.Context().Value(ContextKeyLimit).(int)

	bks, err := serv.RandomBooks(req.Context(), limit)
	if err != nil {
		logger.Error().Err(err).Msg("random books railed")
		writeError(res, 400, err)
	} else {
		json.NewEncoder(res).Encode(booksResp{bks})
	}
}

// @Summary		Get book info
// @description	get book info
// @Tags			book-spider-api
// @Accept			json
// @Produce		json
// @Param			siteName	path		string	true	"site name"
// @Param			idHash		path		string	true	"id and hash in format <id>[-<hash>]. -<hash is optional"
// @Success		200			{object}	model.Book
// @Failure		400			{object}	errResp
// @Router			/api/book-spider/sites/{siteName}/books/{idHash} [get]
func BookInfoAPIHandler(res http.ResponseWriter, req *http.Request) {
	bk := req.Context().Value(ContextKeyBook).(*model.Book)
	json.NewEncoder(res).Encode(bk)
}

// @Summary		Download book
// @description	download book in txt format
// @Tags			book-spider-api
// @Accept			json
// @Produce		json
// @Param			siteName	path		string	true	"site name"
// @Param			idHash		path		string	true	"id and hash in format <id>[-<hash>]. -<hash is optional"
// @Success		200			{string}	string "the book content"
// @Failure		400			{object}	errResp
// @Router			/api/book-spider/sites/{siteName}/books/{idHash}/download [get]
func BookDownloadAPIHandler(res http.ResponseWriter, req *http.Request) {
	logger := zerolog.Ctx(req.Context())
	serv := req.Context().Value(ContextKeyReadDataServ).(service.ReadDataService)
	bk := req.Context().Value(ContextKeyBook).(*model.Book)
	content, err := serv.BookContent(req.Context(), bk)
	if err != nil {
		logger.Error().Err(err).Msg("book content failed")
		writeError(res, 400, err)
	} else {
		fileName := fmt.Sprintf("%s-%s.txt", bk.Title, bk.Writer.Name)
		res.Header().Set("Content-Type", "text/txt; charset=utf-8")
		res.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
		fmt.Fprint(res, content)
	}
}

// @Summary		DB stats
// @description	db stats
// @Tags			book-spider-api
// @Accept			json
// @Produce		json
// @Param			siteName	path		string	true	"site name"
// @Param			idHash		path		string	true	"id and hash in format <id>[-<hash>]. -<hash is optional"
// @Success		200			{object}	map[string]interface{}
// @Failure		400			{object}	errResp
// @Router			/api/book-spider/db-stats [get]
func DBStatsAPIHandler(service service.ReadDataService) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		json.NewEncoder(res).Encode(dbStatsResp{[]sql.DBStats{service.DBStats(req.Context())}})
	}
}
