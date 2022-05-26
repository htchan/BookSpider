package main

import (
	"testing"
	"net/http"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/htchan/BookSpider/internal/utils"
	"github.com/htchan/BookSpider/internal/database"
)

func initMainTest() {
	source, err := os.Open(os.Getenv("ASSETS_LOCATION") + "/test-data/internal_database_sqlite.db")
	utils.CheckError(err)
	destination, err := os.Create("./test.db")
	utils.CheckError(err)
	io.Copy(destination, source)
	source.Close()
	destination.Close()
}

func cleanupMainTest() {
	os.Remove("test.db")
}

func Test_Backend_Api(t *testing.T) {
	setup(os.Getenv("ASSETS_LOCATION") + "/test-data/configs")
	go func() {
		startServer(":8080")
	}()
	client := &http.Client{Timeout: 5 * time.Second}
	host := "http://localhost:8080"

	t.Run("API general info", func(t *testing.T) {
		route := "/api/novel/info"
		t.Run("success", func(t *testing.T) {
			res, err := client.Get(host + route)
			utils.CheckError(err)
			var jsonData map[string][]string
			err = json.NewDecoder(res.Body).Decode(&jsonData)
			utils.CheckError(err)
			if value, ok := jsonData["siteNames"]; ok == false || len(value) != 1 || value[0] != "test" {
				t.Errorf("unexpected response data: %v", jsonData)
			}
		})
	})

	t.Run("API site info", func(t *testing.T) {
		route := "/api/novel/sites/"
		t.Run("success", func(t *testing.T) {
			res, err := client.Get(host + route + "test")
			utils.CheckError(err)
			var jsonData database.SummaryRecord
			err = json.NewDecoder(res.Body).Decode(&jsonData)
			utils.CheckError(err)
			if jsonData.BookCount != 6 || jsonData.ErrorCount != 3 ||
				jsonData.WriterCount != 3 || jsonData.UniqueBookCount != 5 ||
				jsonData.MaxBookId != 5 || jsonData.LatestSuccessId != 3 ||
				jsonData.StatusCount[database.Error] != 3 ||
				jsonData.StatusCount[database.InProgress] != 1 ||
				jsonData.StatusCount[database.End] != 1 ||
				jsonData.StatusCount[database.Download] != 1 {
				t.Errorf("unexpected response data: %v", jsonData)
			}
		})

		t.Run("failed if querying not exist site name", func(t *testing.T) {
			res, err := client.Get(host + route + "unknown/")
			utils.CheckError(err)
			var jsonData map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&jsonData)
			utils.CheckError(err)
			if jsonData["code"] != 404.0 || jsonData["message"] != "site <unknown> not found" {
				t.Errorf("unexpected response data: %v", jsonData)
			}
		})
	})

	t.Run("API book info", func(t *testing.T) {
		route := "/api/novel/books/"
		
		t.Run("success without version", func(t *testing.T) {
			res, err := client.Get(host + route + "test/3")
			utils.CheckError(err)
			var jsonData map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&jsonData)
			utils.CheckError(err)
			if jsonData["site"] != "test" || jsonData["id"] != 3.0 || jsonData["hash"] != "5k" ||
				jsonData["title"] != "title-3-new" || jsonData["writer"] != "writer-3" ||
				jsonData["type"] != "type-3-new" || jsonData["updateDate"] != "100" ||
				jsonData["updateChapter"] != "chapter-3-new" || jsonData["status"] != "end" {
					t.Errorf("unexpected response data: %v", jsonData)
				}
		})
		
		t.Run("success with version", func(t *testing.T) {
			res, err := client.Get(host + route + "test/3/2u")
			utils.CheckError(err)
			var jsonData map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&jsonData)
			utils.CheckError(err)
			if jsonData["site"] != "test" || jsonData["id"] != 3.0 || jsonData["hash"] != "2u" ||
				jsonData["title"] != "title-3" || jsonData["writer"] != "writer-2" ||
				jsonData["type"] != "type-3" || jsonData["updateDate"] != "102" ||
				jsonData["updateChapter"] != "chapter-3" || jsonData["status"] != "download" {
					t.Errorf("unexpected response data: %v", jsonData)
				}
		})

		t.Run("fail if not enough params", func(t *testing.T) {
			res, err := client.Get(host + route + "test")
			utils.CheckError(err)
			var jsonData map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&jsonData)
			utils.CheckError(err)
			if jsonData["code"] != 400.0 || jsonData["message"] != "not enough parameters" {
				t.Errorf("unexpected response data: %v", jsonData)
			}
		})

		t.Run("fail if query not exist site name", func(t *testing.T) {
			res, err := client.Get(host + route + "unknown/1")
			utils.CheckError(err)
			var jsonData map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&jsonData)
			utils.CheckError(err)
			if jsonData["code"] != 404.0 || jsonData["message"] != "site <unknown> not found" {
				t.Errorf("unexpected response data: %v", jsonData)
			}
		})

		t.Run("fail if query not exist book id", func(t *testing.T) {
			res, err := client.Get(host + route + "test/999")
			utils.CheckError(err)
			var jsonData map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&jsonData)
			utils.CheckError(err)
			if jsonData["code"] != 404.0 || jsonData["message"] != "book <999>, hash <> in site <test> not found" {
				t.Errorf("unexpected response data: %v", jsonData)
			}
		})

		t.Run("fail if query not exist book version", func(t *testing.T) {
			res, err := client.Get(host + route + "test/1/123")
			utils.CheckError(err)
			var jsonData map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&jsonData)
			utils.CheckError(err)
			if jsonData["code"] != 404.0 || jsonData["message"] != "book <1>, hash <123> in site <test> not found" {
				t.Errorf("unexpected response data: %v", jsonData)
			}
		})
	})

	t.Run("API book content", func(t *testing.T) {
		// write file to storage
		content := []byte("some test content")
		downloadFileName := os.Getenv("ASSETS_LOCATION") + "/test-data/storage/3-v102.txt"
		err := os.WriteFile(downloadFileName, content, 0644)
		utils.CheckError(err)
		defer os.Remove(downloadFileName)

		route := "/api/novel/download/"
		t.Run("fail without version", func(t *testing.T) {
			res, err := client.Get(host + route + "test/3")
			utils.CheckError(err)
			var jsonData map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&jsonData)
			utils.CheckError(err)
			if jsonData["code"] != 406.0 || jsonData["message"] != "book <3> not download yet" {
				t.Errorf("unexpected response data: %v", jsonData)
			}
		})

		t.Run("success with version", func(t *testing.T) {
			res, err := client.Get(host + route + "test/3/2u")
			utils.CheckError(err)
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			utils.CheckError(err)
			if string(data) != string(content) {
				t.Errorf("unexpected response data: %v", string(data))
			}
		})

		t.Run("fail if target book id not exist", func(t *testing.T) {
			res, err := client.Get(host + route + "test/999")
			utils.CheckError(err)
			var jsonData map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&jsonData)
			utils.CheckError(err)
			if jsonData["code"] != 404.0 || jsonData["message"] != "book <999>, hash <> in site <test> not found" {
				t.Errorf("unexpected response data: %v", jsonData)
			}
		})

		t.Run("fail if target book version not exist", func(t *testing.T) {
			res, err := client.Get(host + route + "test/1/123")
			utils.CheckError(err)
			var jsonData map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&jsonData)
			utils.CheckError(err)
			if jsonData["code"] != 404.0 || jsonData["message"] != "book <1>, hash <123> in site <test> not found" {
				t.Errorf("unexpected response data: %v", jsonData)
			}
		})
	})

	t.Run("API book search", func(t *testing.T) {
		route := "/api/novel/search"

		t.Run("success with multi results even only title match", func(t *testing.T) {
			res, err := client.Get(host + route + "/test?title=title&writer=abc")
			utils.CheckError(err)
			// t.Errorf(string(data))
			var jsonData map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&jsonData)
			utils.CheckError(err)
			if len(jsonData["books"].([]interface{})) != 3 {
				t.Errorf("unexpected response data: %v", jsonData)
			}
		})

		t.Run("return empty array if not exist", func(t *testing.T) {
			res, err := client.Get(host + route + "/test?title=abc&writer=abc")
			utils.CheckError(err)
			// t.Errorf(string(data))
			var jsonData map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&jsonData)
			utils.CheckError(err)
			if len(jsonData["books"].([]interface{})) != 0 {
				t.Errorf("unexpected response data: %v", jsonData)
			}
		})

		t.Run("fail if both title and writer queries are empty", func(t *testing.T) {
			res, err := client.Get(host + route + "/test")
			utils.CheckError(err)
			// t.Errorf(string(data))
			var jsonData map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&jsonData)
			utils.CheckError(err)
			if len(jsonData["books"].([]interface{})) != 0 {
				t.Errorf("unexpected response data: %v", jsonData)
			}
		})
	})

	t.Run("API book random", func(t *testing.T) {
		route := "/api/novel/random/"
		
		t.Run("success to return book not in order", func(t *testing.T) {
			res, err := client.Get(host + route + "test")
			utils.CheckError(err)
			var jsonData map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&jsonData)
			utils.CheckError(err)
			if len(jsonData["books"].([]interface{})) != 1 {
				t.Errorf("unexpected response data: %v", jsonData)
			}
		})

		t.Run("success to return book with specific status", func(t *testing.T) {
			res, err := client.Get(host + route + "test?status=error")
			utils.CheckError(err)
			var jsonData map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&jsonData)
			utils.CheckError(err)
			if len(jsonData["books"].([]interface{})) != 6 {
				t.Errorf("unexpected response data: %v", jsonData)
			}
		})
	})
}


func TestMain(m *testing.M) {
	initMainTest()

	code := m.Run()

	cleanupMainTest()
	os.Exit(code)
}