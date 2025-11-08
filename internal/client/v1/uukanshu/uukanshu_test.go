package uukanshu

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
					<div class="xiaoshuo_content"><dl class="jieshao"><dd class="jieshao_content">
						<h1><a title="book name最新章节"></a></h1>
						<h2><a>author</a></h2>
						<div class="shijian">5月</div>
					</dd></dl></div>
					<div class="weizhi"><div class="path"><a></a><a>type</a></div></div>
					<div class="zhangjie"><ul id="chapterList"><li><a>chapter name</a></li></ul></div>
				</data>`))
			} else if strings.Contains(r.URL.Path, "chapter_list") {
				w.Write([]byte(`<data>
					<div class="zhangjie"><ul id="chapterList">
						<li><a href="chapter url 4">chapter name 4</a></li>
						<li><a href="chapter url 3">chapter name 3</a></li>
						<li><a href="chapter url 2">chapter name 2</a></li>
						<li><a href="chapter url 1">chapter name 1</a></li>
					</ul></div>
				</data>`))
			} else if strings.Contains(r.URL.Path, "chapter") {
				w.Write([]byte(`<data>
					<div class="zhengwen_box"><div class="box_left"><div class="w_main">
						<div class="h1title"><h1 id="timu">chapter name</h1></div>
						<div class="contentbox"><div id="contentbox">chapter content</div></div>
					</div></div></div>
				</data>`))
			} else if strings.Contains(r.URL.Path, "available") {
				w.Write([]byte(`UU看书`))
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
