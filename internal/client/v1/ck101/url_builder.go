package ck101

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

func bookURL(bookID string) string {
	return fmt.Sprintf(bookURLTemplate, bookID)
}

func chapterListURL(bookID string) string {
	return fmt.Sprintf(chapterListURLTemplate, bookID)
}

func chapterURL(uri string) string {
	if uri == "" {
		return ""
	} else if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
		return uri
	} else if strings.HasPrefix(uri, "/") {
		return vendorProtocol + "://" + vendorHost + uri
	}

	log.Error().
		Str("vendor", Host).
		Str("url", uri).
		Msg("unexpected resources for building chapter url")

	return uri
}

func availabilityURL() string {
	return vendorProtocol + "://" + vendorHost
}
