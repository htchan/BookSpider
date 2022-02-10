package mock

import (
	"net/http/httptest"
	"net/http"
	"fmt"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

func UtilsServer(expectResponse string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			if (req.URL.Path == "/error") {
				res.WriteHeader(400)
			} else {
				fmt.Fprintf(res, expectResponse)
			}
		}))
}

func UtilsEncoderServer(expectResponse string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			if (req.URL.Path == "/error") {
				res.WriteHeader(400)
			} else {
				encoder := traditionalchinese.Big5.NewEncoder()
				response, _, _ := transform.String(encoder, expectResponse)
				fmt.Fprintf(res, response)
			}
		}))
}