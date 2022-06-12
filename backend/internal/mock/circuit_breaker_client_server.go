package mock

import (
	"fmt"
	"net/http/httptest"
	"net/http"
	"time"
)

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
			} else if req.URL.Path == "/timeout" {
				time.Sleep(time.Duration(timeout) * time.Second)
				fmt.Fprintf(res, "timeout")
			} else if req.URL.Path == "/empty" {
			}
		}))
}