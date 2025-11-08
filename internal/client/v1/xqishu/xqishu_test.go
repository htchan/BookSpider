package xqishu

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
					<div class="tit1"><h1>book name</h1></div>
					<div class="date">
						<span>小说作者：author</span>
						<span></span>
						<span>更新时间：2000-01-01</span>
					</div>
					<div class="crumbs"><a></a><a>type</a></div>
					<a class="zx_zhang">chapter name</a>
				</data>`))
			} else if strings.Contains(r.URL.Path, "chapter_list") {
				w.Write([]byte(`<data>
					<div class="book_con_list"><ul>
						<li><a href="/chapter url 1">chapter name 1</a></li>
						<li><a href="/chapter url 2">chapter name 2</a></li>
						<li><a href="/chapter url 3">chapter name 3</a></li>
						<li><a href="/chapter url 4">chapter name 4</a></li>
					</ul></div>
				</data>`))
			} else if strings.Contains(r.URL.Path, "chapter") {
				w.Write([]byte(`<data>
					<div class="date"><h1>chapter name</h1></div>
					<div class="book_content">chapter content</div>
				</data>`))
			} else if strings.Contains(r.URL.Path, "available") {
				w.Write([]byte(`求书网`))
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
		goleak.VerifyTestMain(m)
		serv.Close()
	} else {
		code := m.Run()
		serv.Close()
		os.Exit(code)
	}
}
