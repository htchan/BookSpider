package main

import (
	"net/http"
	"strings"
	"strconv"
	"github.com/htchan/BookSpider/internal/database"
	"fmt"
)

func BookDownload(res http.ResponseWriter, req *http.Request) {
	setHeader(res)
	uri := strings.Split(req.URL.Path, "/")
	if len(uri) < 6 {
		error(res, http.StatusBadRequest, "not enough parameter")
	}
	siteName := uri[4]
	site, ok := siteMap[siteName]
	id, err := strconv.Atoi(uri[5])
	if !ok {
		error(res, http.StatusNotFound, "site <" + siteName + "> not found")
		return
	} else if err != nil {
		error(res, http.StatusBadRequest, "id <" + uri[5] + "> is not a number")
		return
	}
	hashCode := "";
	if len(uri) > 6 {
		hashCode = uri[6]
	}
	book := site.SearchByIdHash(id, hashCode)
	if book == nil || book.GetTitle() == "" {
		error(res, http.StatusNotFound, "book <" + strconv.Itoa(id) + ">, " +
				"hash <" + hashCode + "> in site <" + siteName + "> not found")
		return
	} else if book.GetStatus() != database.Download {
		error(res, http.StatusNotAcceptable, "book <" + uri[5] + "> not download yet")
		return
	}
	fileName := book.GetTitle() + "-" + book.GetWriter()
	_, _, hashCodeInt := book.GetInfo()
	if hashCodeInt > 0 {
		fileName += "-v" + strconv.FormatInt(int64(hashCodeInt), 36)
	}
	content := book.GetContent(site.StorageDirectory())
	res.Header().Set("Content-Type", "text/txt; charset=utf-8")
	res.Header().Set("Content-Disposition", "attachment; filename=\"" + fileName + ".txt\"")
	fmt.Fprintf(res, content)
}

func BookSearch(res http.ResponseWriter, req *http.Request) {
	setHeader(res)
	uri := strings.Split(req.URL.Path, "/")
	siteName := uri[4]
	site, ok := siteMap[siteName]
	title := req.URL.Query().Get("title")
	writer := req.URL.Query().Get("writer")
	pageStr := req.URL.Query().Get("page")
	bookResults := site.SearchByTitleWriter(title, writer)
	page, err := strconv.Atoi(pageStr)
	if !ok {
		error(res, http.StatusNotFound, "site <" + siteName + "> not found")
		return
	} else if err != nil {
		page = 0
	}
	booksArray := make([]map[string]interface{}, 0)
	for i := page * RECORD_PER_PAGE; i < (page + 1) * RECORD_PER_PAGE && i < len(bookResults); i++ {
		booksArray = append(booksArray, bookResults[i].Map())
	}
	response(res, map[string]interface{} {
		"books": booksArray,
	})
}

func BookRandom(res http.ResponseWriter, req *http.Request) {
	setHeader(res)
	uri := strings.Split(req.URL.Path, "/")
	siteName := uri[4]
	site, ok := siteMap[siteName]
	if !ok {
		error(res, http.StatusNotFound, "site <" + siteName + "> not found")
		return
	}
	count, err := strconv.Atoi(req.URL.Query().Get("count"))
	if (err != nil) { count = 20 }
	if (count > 50) { count = 50 }
	status, ok := database.StatusCodeMap[strings.ToUpper(req.URL.Query().Get("status"))]
	if !ok { status = database.Download }
	bookResults := site.RandomSuggestBook(count, status)
	booksArray := make([]map[string]interface{}, 0)
	for _, book := range bookResults {
		booksArray = append(booksArray, book.Map())
	}
	response(res, map[string]interface{} {
		"books": booksArray,
	})
}