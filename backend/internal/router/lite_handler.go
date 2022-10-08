package router

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/htchan/BookSpider/internal/service/site"
)

//go:embed templates/*
var files embed.FS

func GeneralLiteHandler(sites map[string]*site.Site) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFS(files, "templates/sites.html")
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			log.Println(err)
			return
		}
		t.Execute(res, sites)
	}
}

func SiteLiteHandlerfunc(res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFS(files, "templates/site.html")
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}

	st := req.Context().Value("site").(*site.Site)
	t.Execute(res, struct {
		Name    string
		Summary repo.Summary
	}{
		Name:    st.Name,
		Summary: st.Info(),
	})
}

func SearchLiteHandler(res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFS(files, "templates/result.html")
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}

	st := req.Context().Value("site").(*site.Site)
	title := req.Context().Value("title").(string)
	writer := req.Context().Value("writer").(string)
	limit := req.Context().Value("limit").(int)
	offset := req.Context().Value("offset").(int)
	if limit == 0 {
		limit = 10
	}

	bks, err := st.QueryBooks(title, writer, limit, offset)

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
		Name:       st.Name,
		Books:      bks,
		lastPageNo: pageNo - 1,
		nextPageNo: pageNo + 1,
	})
}

func RandomLiteHandler(res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFS(files, "templates/result.html")
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}

	st := req.Context().Value("site").(*site.Site)
	limit := req.Context().Value("limit").(int)
	if limit == 0 {
		limit = 10
	}

	bks, err := st.RandomBooks(limit)

	if err != nil {
		res.WriteHeader(404)
		fmt.Fprint(res, "books not found")
		return
	}

	t.Execute(res, struct {
		Name  string
		Books []model.Book
	}{
		Name:  st.Name,
		Books: bks,
	})
}

func BookLiteHandler(res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFS(files, "templates/book.html")
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}

	st := req.Context().Value("site").(*site.Site)
	bk := req.Context().Value("book").(*model.Book)

	t.Execute(res, struct {
		Name string
		Book *model.Book
	}{
		Name: st.Name,
		Book: bk,
	})
}

func DownloadLiteHandler(res http.ResponseWriter, req *http.Request) {
	st := req.Context().Value("site").(*site.Site)
	bk := req.Context().Value("book").(*model.Book)
	content, err := site.Content(st, bk)
	if err != nil {
		res.WriteHeader(500)
		fmt.Println(err)
		return
	} else {
		fileName := fmt.Sprintf("%s-%s.txt", bk.Title, bk.Writer.Name)
		res.Header().Set("Content-Type", "text/txt; charset=utf-8")
		res.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
		fmt.Fprint(res, content)
	}
}
