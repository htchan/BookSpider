package router

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/htchan/BookSpider/internal/service/book"
	"github.com/htchan/BookSpider/internal/service/site"
)

func GeneralInfoAPIHandler(sites map[string]*site.Site) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		stInfo := make(map[string]repo.Summary)
		for _, site := range sites {
			stInfo[site.Name] = site.Info()
		}
		json.NewEncoder(res).Encode(stInfo)
	}
}

func SiteInfoAPIHandler(res http.ResponseWriter, req *http.Request) {
	st := req.Context().Value("site").(*site.Site)
	json.NewEncoder(res).Encode(st.Info())
}

func BookSearchAPIHandler(res http.ResponseWriter, req *http.Request) {
	st := req.Context().Value("site").(*site.Site)
	title := req.Context().Value("title").(string)
	writer := req.Context().Value("writer").(string)
	limit := req.Context().Value("limit").(int)
	offset := req.Context().Value("offset").(int)

	bks, err := st.QueryBooks(title, writer, limit, offset)
	if err != nil {
		json.NewEncoder(res).Encode(map[string]string{"error": err.Error()})
	} else {
		json.NewEncoder(res).Encode(map[string][]model.Book{"books": bks})
	}
}

func BookRandomAPIHandler(res http.ResponseWriter, req *http.Request) {
	st := req.Context().Value("site").(*site.Site)
	limit := req.Context().Value("limit").(int)

	bks, err := st.RandomBooks(limit)
	if err != nil {
		json.NewEncoder(res).Encode(map[string]string{"error": err.Error()})
	} else {
		json.NewEncoder(res).Encode(map[string][]model.Book{"books": bks})
	}
}

func BookInfoAPIHandler(res http.ResponseWriter, req *http.Request) {
	bk := req.Context().Value("book").(*model.Book)
	json.NewEncoder(res).Encode(bk)
}

func BookDownloadAPIHandler(res http.ResponseWriter, req *http.Request) {
	st := req.Context().Value("site").(*site.Site)
	bk := req.Context().Value("book").(*model.Book)
	content, err := book.Content(bk, st.StConf)
	if err != nil {
		json.NewEncoder(res).Encode(map[string]string{"error": err.Error()})
	} else {
		fileName := fmt.Sprintf("%s-%s.txt", bk.Title, bk.Writer.Name)
		res.Header().Set("Content-Type", "text/txt; charset=utf-8")
		res.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
		fmt.Fprint(res, content)
	}
}

func DBStatsAPIHandler(sites map[string]*site.Site) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		stats := make([]sql.DBStats, 0, len(sites))
		for _, st := range sites {
			stats = append(stats, st.Stat())
		}
		json.NewEncoder(res).Encode(map[string][]sql.DBStats{"stats": stats})
	}
}
