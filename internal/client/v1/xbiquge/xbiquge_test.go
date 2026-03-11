package xbiquge

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"go.uber.org/goleak"
)

var serv *httptest.Server

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "check for memory leaks")
	flag.Parse()

	serv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "success") {
			if strings.Contains(r.URL.Path, "book_info") {
				w.Write([]byte(`<data>
					<meta property="og:novel:book_name" content="book name" />
					<meta property="og:novel:author" content="author" />
					<meta property="og:novel:category" content="type" />
					<meta property="og:novel:update_time" content="2000-01-01" />
					<meta property="og:novel:latest_chapter_name" content="chapter name" />
				</data>`))
			} else if strings.Contains(r.URL.Path, "chapter_list") {
				w.Write([]byte(`<data>
					<div>
						<dd><a href="chapter url 1">chapter name 1</a></dd>
						<dd><a href="chapter url 2">chapter name 2</a></dd>
						<dd><a href="chapter url 3">chapter name 3</a></dd>
						<dd><a href="chapter url 4">chapter name 4</a></dd>
					</div>
				</data>`))
			} else if strings.Contains(r.URL.Path, "chapter") {
				w.Write([]byte(`<data>
					<div class="bookname"><h1>chapter name</h1></div>
					<div id="content">chapter content</div>
				</data>`))
			} else if strings.Contains(r.URL.Path, "available") {
				w.Write([]byte(`笔趣阁`))
			}
		} else if strings.Contains(r.URL.Path, "not_found") {
			w.WriteHeader(http.StatusNotFound)
		} else if strings.Contains(r.URL.Path, "error") {
			w.WriteHeader(http.StatusBadRequest)
		} else if strings.Contains(r.URL.Path, "forbidden") {
			w.WriteHeader(http.StatusForbidden)
		} else if strings.Contains(r.URL.Path, "timeout") {
			time.Sleep(100 * time.Millisecond)
		}
	}))

	if *leak {
		goleak.VerifyTestMain(
			m,
			goleak.Cleanup(func(exitCode int) { serv.Close() }),
		)
	} else {
		code := m.Run()
		serv.Close()
		os.Exit(code)
	}
}
