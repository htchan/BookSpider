package router

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	service_new "github.com/htchan/BookSpider/internal/service_new"
)

func GeneralInfoAPIHandler(services map[string]service_new.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		servInfo := make(map[string]repo.Summary)
		for _, serv := range services {
			servInfo[serv.Name()] = serv.Stats()
		}
		json.NewEncoder(res).Encode(servInfo)
	}
}

func SiteInfoAPIHandler(res http.ResponseWriter, req *http.Request) {
	serv := req.Context().Value(SERV_KEY).(service_new.Service)
	json.NewEncoder(res).Encode(serv.Stats())
}

func BookSearchAPIHandler(res http.ResponseWriter, req *http.Request) {
	serv := req.Context().Value(SERV_KEY).(service_new.Service)
	title := req.Context().Value(TITLE_KEY).(string)
	writer := req.Context().Value(WRITER_KEY).(string)
	limit := req.Context().Value(LIMIT_KEY).(int)
	offset := req.Context().Value(OFFSET_KEY).(int)

	bks, err := serv.QueryBooks(title, writer, limit, offset)
	if err != nil {
		json.NewEncoder(res).Encode(map[string]string{"error": err.Error()})
	} else {
		json.NewEncoder(res).Encode(map[string][]model.Book{"books": bks})
	}
}

func BookRandomAPIHandler(res http.ResponseWriter, req *http.Request) {
	serv := req.Context().Value(SERV_KEY).(service_new.Service)
	limit := req.Context().Value(LIMIT_KEY).(int)

	bks, err := serv.RandomBooks(limit)
	if err != nil {
		json.NewEncoder(res).Encode(map[string]string{"error": err.Error()})
	} else {
		json.NewEncoder(res).Encode(map[string][]model.Book{"books": bks})
	}
}

func BookInfoAPIHandler(res http.ResponseWriter, req *http.Request) {
	bk := req.Context().Value(BOOK_KEY).(*model.Book)
	json.NewEncoder(res).Encode(bk)
}

func BookDownloadAPIHandler(res http.ResponseWriter, req *http.Request) {
	serv := req.Context().Value(SERV_KEY).(service_new.Service)
	bk := req.Context().Value(BOOK_KEY).(*model.Book)
	content, err := serv.BookContent(bk)
	if err != nil {
		json.NewEncoder(res).Encode(map[string]string{"error": err.Error()})
	} else {
		fileName := fmt.Sprintf("%s-%s.txt", bk.Title, bk.Writer.Name)
		res.Header().Set("Content-Type", "text/txt; charset=utf-8")
		res.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
		fmt.Fprint(res, content)
	}
}

func DBStatsAPIHandler(services map[string]service_new.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		stats := make([]sql.DBStats, 0, len(services))
		for _, serv := range services {
			stats = append(stats, serv.DBStats())
		}
		json.NewEncoder(res).Encode(map[string][]sql.DBStats{"stats": stats})
	}
}
