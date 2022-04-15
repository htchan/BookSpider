package main

import (
	"github.com/htchan/BookSpider/pkg/configs"
	"github.com/htchan/BookSpider/pkg/books"
	"github.com/htchan/BookSpider/internal/utils"

	"log"
	"os"
	"encoding/json"
	"sync"
	"time"
	"runtime"
)

func download_books_chapters_title(site string, config *configs.BookConfig, start, end int) (titles []string) {
	titles = make([]string, 0)
	for i := start; i < end; i++ {
		book := books.NewBook(site, i, -1, config)
		t, err := book.GetChapterTitles()
		if err != nil {
			site, id, _ := book.GetInfo()
			log.Printf("%v-%v: download titles: %v", site, id, err)
		}
		titles = append(titles, t ...)
		log.Printf("%v completed", i - start)
	}
	return
}
func download_books_chapters_title_thread(site string, config *configs.BookConfig, start, end int) (titles []string) {
	titles = make([]string, 0)
	// log.Println(config)
	ch := make(chan []string, 1)
	go func() {
		for i := start; i < end; i++ {
			log.Printf("%v completed", i - start)
			titles = append(titles, <-ch ...)
		}
	}()
	var wg sync.WaitGroup
	for i := start; i < end; i++ {
		wg.Add(1)
		time.Sleep(1000 * time.Millisecond)
		go func(i int) {
			defer wg.Done()
			book := books.NewBook(site, i, -1, config)
			t, err := book.GetChapterTitles()
			if err != nil {
				site, id, _ := book.GetInfo()
				log.Printf("%v-%v: download titles: %v", site, id, err)
				ch <- []string{}
				return
			}
			ch <- t
		} (i)
	}
	wg.Wait()
	return
}

func writeJson(location string, obj interface{}) {
	b, err := json.Marshal(obj)
	utils.CheckError(err)
	os.WriteFile(location, b, 0644)
}

func main() {
	runtime.GOMAXPROCS(3)
	configs := []struct{
		name string
		config *configs.BookConfig
	} {
		// { "ck101", configs.LoadBookConfigYaml("../configs/site-config/ck101-desktop.yaml") },
		{ "xqishu", configs.LoadBookConfigYaml("../configs/site-config/xqishu-desktop.yaml") },
		// { "hjwzw", configs.LoadBookConfigYaml("../configs/site-config/hjwzw-desktop.yaml") },
	}
	titles := make([]string, 0)
	s := time.Now()
	for _, config := range configs {
		log.Println(config.config.CONST_SLEEP)
		titles = append(titles, download_books_chapters_title(config.name, config.config, 1, 100)...)
	}
	log.Println("time used %s", time.Since(s))

	writeJson("./titles.json", titles)
}