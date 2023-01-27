package router

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	service_new "github.com/htchan/BookSpider/internal/service_new"
)

//go:embed templates/*
var files embed.FS

func GeneralLiteHandler(services map[string]service_new.Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		t, err := template.ParseFS(files, "templates/sites.html")
		if err != nil {
			res.WriteHeader(http.StatusNotFound)
			log.Println(err)
			return
		}
		t.Execute(res, services)
	}
}

func SiteLiteHandlerfunc(res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFS(files, "templates/site.html")
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		log.Println(err)
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
	t, err := template.ParseFS(files, "templates/result.html")
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		log.Println(err)
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
	t, err := template.ParseFS(files, "templates/result.html")
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		log.Println(err)
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
	t, err := template.ParseFS(files, "templates/book.html")
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}

	serv := req.Context().Value(SERV_KEY).(service_new.Service)
	bk := req.Context().Value(BOOK_KEY).(*model.Book)

	t.Execute(res, struct {
		Name string
		Book *model.Book
	}{
		Name: serv.Name(),
		Book: bk,
	})
}

func DownloadLiteHandler(res http.ResponseWriter, req *http.Request) {
	serv := req.Context().Value(SERV_KEY).(service_new.Service)
	bk := req.Context().Value(BOOK_KEY).(*model.Book)
	content, err := serv.BookContent(bk)
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
