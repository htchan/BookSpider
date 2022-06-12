package mock

import (
	"fmt"
	"net/http/httptest"
	"net/http"
	"strings"
	"regexp"
)

func DownloadServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			if strings.HasPrefix(req.URL.Path, "/content/imbalance_url_title") {
				fmt.Fprintf(res, "chapter-url-regex-/1 chapter-title-regex-1" +
				"chapter-url-regex-/2 chapter-title-regex-2" +
				"chapter-url-regex-/3 chapter-title-regex-3" +
				"chapter-url-regex-/4 ")
			} else if strings.HasPrefix(req.URL.Path, "/content/empty") {
				fmt.Fprintf(res, "")
			} else if strings.HasPrefix(req.URL.Path, "/content/no_url") {
				fmt.Fprintf(res, "hello")
			} else if strings.HasPrefix(req.URL.Path, "/content/success") {
				fmt.Fprintf(res, "chapter-url-regex-/1 chapter-title-regex-1 " +
				"chapter-url-regex-/2 chapter-title-regex-2 " +
				"chapter-url-regex-/3 chapter-title-regex-3 " +
				"chapter-url-regex-/4 chapter-title-regex-4 ")
			} else if strings.HasPrefix(req.URL.Path, "/chapter/success") {
				s := strings.ReplaceAll(req.URL.Path, "success", "")
				s = strings.ReplaceAll(s, "content", "")
				s = strings.ReplaceAll(s, "chapter", "")
				s = strings.ReplaceAll(s, "/", "")
				fmt.Fprintf(res, "chapter-content-" + s + "-content-regex")
			} else if strings.HasPrefix(req.URL.Path, "/chapter/replace") {
				fmt.Fprintf(res, "chapter-content-url-hello<br />-content-regex")
			} else if strings.HasPrefix(req.URL.Path, "/chapter/invalid") {

			}
		}))
}

func MockBookDownloadServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			if strings.Contains(req.URL.Path, "/chapter/success") {
				reg := regexp.MustCompile(".*/chapter/")
				data := reg.ReplaceAllString(req.URL.Path, "")
				fmt.Fprintf(res, fmt.Sprintf("chapter-content-%v", data))
			} else if strings.HasPrefix(req.URL.Path, "/chapter/400") {
				res.WriteHeader(400)
			} else if strings.HasPrefix(req.URL.Path, "/chapter/unknown") {
				fmt.Fprintf(res, "unknown content")
			} else if strings.HasPrefix(req.URL.Path, "/list_chapters") {
				fmt.Fprintf(res, "url-0 title-0 url-1 title-1 url-2 title-2 ")
			} else if strings.HasPrefix(req.URL.Path, "/list_success_chapters") {
				fmt.Fprintf(res, "url-/0 title-0 url-/1 title-1 url-/2 title-2 ")
			} else if strings.HasPrefix(req.URL.Path, "/400") {
				res.WriteHeader(400)
			} else if strings.HasPrefix(req.URL.Path, "/unknown") {
				fmt.Fprintf(res, "unknown content")
			}
		}))
}
