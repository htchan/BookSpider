package mock

import (
	"fmt"
	"net/http/httptest"
	"net/http"
)

func UpdateServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			if (req.URL.Path == "/partial_fail") {
				fmt.Fprintf(res, "title-regex writer-regex type-regex last-update-regex ")
			} else if (req.URL.Path == "/empty") {
				fmt.Fprintf(res, "")
			} else if (req.URL.Path == "/number") {
				fmt.Fprintf(res, "200")
			} else {
				fmt.Fprintf(res, "title-regex writer-regex type-regex last-update-regex " +
				"last-chapter-regex")
			}
		}))
}