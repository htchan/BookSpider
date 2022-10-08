package mock

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"
)

var MaxExplore = 5

func MockSiteServer() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(
			func(res http.ResponseWriter, req *http.Request) {
				if strings.HasPrefix(req.URL.Path, "/no-update-book") {
					fmt.Fprintf(res, "title writer type date-1234 chapter ")
				} else if strings.HasPrefix(req.URL.Path, "/explore") {
					id, err := strconv.Atoi(strings.ReplaceAll(req.URL.Path, "/explore/", ""))
					if err != nil || id > MaxExplore {
						fmt.Fprintln(res, "error")
					} else {
						fmt.Fprintf(res, "title-%v writer-%v type-%v date-%v chapter-%v ", id, id, id, id, id)
					}
				} else if strings.HasPrefix(req.URL.Path, "/update-book/chapter-end") {
					fmt.Fprintf(res, "title writer type date-1234 chapter後記 ")
				} else if strings.HasPrefix(req.URL.Path, "/update-book/chapter-not-end") {
					fmt.Fprintf(res, "title writer type date-1234 chapter-new ")
				} else if strings.HasPrefix(req.URL.Path, "/update-book/title") {
					fmt.Fprintf(res, "title-new writer type date-1234 chapter ")
				} else if strings.HasPrefix(req.URL.Path, "/error") {
					fmt.Fprintf(res, "error-content")
				} else if strings.HasPrefix(req.URL.Path, "/chapter/unrecognize") {
					fmt.Fprintf(res, "unknown")
				} else if strings.HasPrefix(req.URL.Path, "/chapter/valid") {
					fmt.Fprintf(res, "chapter-content-success-content-regex")
				} else if strings.HasPrefix(req.URL.Path, "/chapter/extra-content") {
					fmt.Fprintf(res, "chapter-content-success-<br />content&nbsp;-regex")
				} else if strings.HasPrefix(req.URL.Path, "/chapter-header/valid") {
					fmt.Fprintf(res, "url-/1 title-1 "+
						"url-/2 title-2 "+
						"url-/3 title-3 "+
						"url-/4 title-4 ")
				} else if strings.HasPrefix(req.URL.Path, "/chapter-header/empty") {
					fmt.Fprintf(res, "")
				} else if strings.HasPrefix(req.URL.Path, "/chapter-header/imbalance-url-title") {
					fmt.Fprintf(res, "url-/1 title-1 "+
						"url-/2 title-2 "+
						"url-/3 title-3 "+
						"url-/4 ")
				} else if strings.HasPrefix(req.URL.Path, "/chapter-header/url-not-found") {
					fmt.Fprintf(res, "hello")
				}
			},
		),
	)
}

func MockCircuitBreakerServer(timeout int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/200" {
				res.WriteHeader(200)
				fmt.Fprintf(res, "200")
			} else if req.URL.Path == "/400" {
				res.WriteHeader(400)
				fmt.Fprintf(res, "400")
			} else if req.URL.Path == "/503" {
				res.WriteHeader(503)
				fmt.Fprintf(res, "503")
			} else if req.URL.Path == "/502" {
				res.WriteHeader(502)
				fmt.Fprintf(res, "502")
			} else if req.URL.Path == "/timeout" {
				time.Sleep(time.Duration(timeout) * time.Second)
				fmt.Fprintf(res, "timeout")
			} else if req.URL.Path == "/empty" {
			}
		}))
}
