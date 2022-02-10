package mock

import (
	"fmt"
	"net/http/httptest"
	"net/http"
)

func ChapterServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/unrecognize" {
				fmt.Fprintf(res, "hello")
			} else if req.URL.Path == "/invalid" {
				fmt.Fprintf(res, "")
			} else {
				fmt.Fprintf(res, "chapter-content-success-content-regex")
			}
		}))
}