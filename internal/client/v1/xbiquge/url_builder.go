package xbiquge

import (
	"fmt"
	"strings"
)

func bookURL(bookID string) string {
	return fmt.Sprintf(bookURLTemplate, bookID)
}

func chapterListURL(bookID string) string {
	return fmt.Sprintf(chapterListURLTemplate, bookID)
}

func chapterURL(uri string, bookID string) string {
	if uri == "" {
		return ""
	} else if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
		return uri
	} else if strings.HasPrefix(uri, "/") {
		return vendorProtocol + "://" + vendorHost + uri
	} else {
		return fmt.Sprintf(chapterURLTemplate, bookID, uri)
	}
}

func availabilityURL() string {
	return vendorProtocol + "://" + vendorHost
}
