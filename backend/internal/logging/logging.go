package logging

import (
	"log"
	"fmt"
)

func logInfo(header string, data interface{}) {
	if data != nil {
		log.Printf("%v %v", header, data)
	} else {
		log.Println(header)
	}
}

func LogBookEvent(book, action, event string, data interface{}) {
	logInfo(fmt.Sprintf("book-spider.book.%v.%v.%v", book, action, event), data)
}

func LogSiteEvent(site, action, event string, data interface{}) {
	logInfo(fmt.Sprintf("book-spider.site.%v.%v.%v", site, action, event), data)
}

func LogRequestEvent(action, event string, data interface{}) {
	logInfo(fmt.Sprintf("book-spider.request.%v.%v", action, event), data)
}

func LogEvent(area, event string, data interface{}) {
	logInfo(fmt.Sprintf("book-spider.%v.%v", area, event), data)
}