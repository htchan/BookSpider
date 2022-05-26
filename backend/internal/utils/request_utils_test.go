package utils

import (
	"testing"
	"time"
	"github.com/htchan/BookSpider/internal/mock"
	"golang.org/x/text/encoding/traditionalchinese"
)

func TestUtils_Request_getWeb(t *testing.T) {
	server := mock.UtilsServer("stub utils server")
	defer server.Close()
	encodeServer := mock.UtilsEncoderServer("一二三")
	defer encodeServer.Close()

	t.Run("func getWeb", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			response := getWeb(server.URL + "/testing")
			if response != "stub utils server" {
				t.Errorf("get Web return \"%v\"", response)
			}
		})

		t.Run("failed if status code is not 200", func(t *testing.T) {
			response := getWeb(server.URL + "/error")
			if response != "400" {
				t.Errorf("get Web return \"%v\" for non 200 response", response)
			}
		})
	})
	
	t.Run("func GetWeb", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			response, trial := GetWeb(server.URL + "/testing", 1, nil, 0)
			if response != "stub utils server" || trial != 0 {
				t.Errorf("utils.GetWeb return \"%v\", %v for success case",
					response, trial)
			}
		})

		t.Run("success for specific encoder", func(t *testing.T) {
			decoder := traditionalchinese.Big5.NewDecoder()
			response, trial := GetWeb(encodeServer.URL + "/testing", 1, decoder, 0)
			if response != "一二三" || trial != 0 {
				t.Errorf("utils.GetWeb return \"%v\", %v for success case",
					response, trial)
			}
		})

		t.Run("fail if decoder not match", func(t *testing.T) {
			response, trial := GetWeb(encodeServer.URL + "/testing", 1, nil, 0)
			if response == "一二三" || trial != 0 {
				t.Errorf("utils.GetWeb return \"%v\", %v for success case",
					response, trial)
			}
		})
	})

	t.Run("func RequestInterval", func(t *testing.T) {
		t.Run("success to wait 1 second", func(t *testing.T) {
			SlowRequest = true
			before := time.Now().Unix()
			RequestInterval()
			after := time.Now().Unix()
			if after - before < 1 {
				t.Errorf("it wait for %v unix", after - before)
			}
		})

		t.Run("not wait any time if SlowRequest is false", func(t *testing.T) {
			SlowRequest = false
			before := time.Now().Unix()
			RequestInterval()
			after := time.Now().Unix()
			if after - before > 0 {
				t.Errorf("it wait for %v unix", after - before)
			}
		})
	})
}
