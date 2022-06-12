package mock

import (
	"fmt"
	"net/http/httptest"
	"net/http"
	"strings"
	"strconv"
)

func UpdateServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			if strings.HasPrefix(req.URL.Path, "/partial_fail") {
				fmt.Fprintf(res, "title-regex writer-regex type-regex last-update-regex ")
			} else if (req.URL.Path == "/empty") {
				fmt.Fprintf(res, "")
			} else if (req.URL.Path == "/number") {
				fmt.Fprintf(res, "200")
			} else if strings.HasPrefix(req.URL.Path, "/specific_success") {
				fmt.Fprintf(res, "title-regex writer-regex type-regex 104 chapter-1")
			} else if strings.HasPrefix(req.URL.Path, "/success/") {
				i, err := strconv.Atoi(strings.ReplaceAll(req.URL.Path, "/success/", ""))
				if err == nil && i < 9 {
					fmt.Fprintf(res, "title-regex writer-regex type-regex last-update-regex " +
					"last-chapter-regex")
				} else {
					fmt.Fprintf(res, "title-regex writer-regex type-regex last-update-regex ")
				}
			} else {
				fmt.Fprintf(res, "title-regex writer-regex type-regex last-update-regex " +
				"last-chapter-regex")
			}
		}))
}

func MockBookUpdateServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			if strings.HasPrefix(req.URL.Path, "/updated") {
				fmt.Fprintf(res, "title_new writer_new type_new date_new chapter_new ")
			} else if strings.HasPrefix(req.URL.Path, "/no_updated") {
				fmt.Fprintf(res, "title writer type date chapter ")
			} else if strings.HasPrefix(req.URL.Path, "/zero_length") {
			} else if strings.HasPrefix(req.URL.Path, "/400") {
				res.WriteHeader(400)
			} else if strings.HasPrefix(req.URL.Path, "/503") {
				res.WriteHeader(503)
			} else if strings.HasPrefix(req.URL.Path, "/missing_date") {
				fmt.Fprintf(res, "title writer type chapter ")
			}
	}))
}