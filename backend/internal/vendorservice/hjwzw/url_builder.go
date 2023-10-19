package hjwzw

import (
	"fmt"
	"strings"

	vendor "github.com/htchan/BookSpider/internal/vendorservice"
	"github.com/rs/zerolog/log"
)

type BookURLBuilder struct{}

var _ vendor.BookURLBuilder = (*BookURLBuilder)(nil)

func NewURLBuilder() *BookURLBuilder {
	return &BookURLBuilder{}
}

func (b *BookURLBuilder) BookURL(bookID string) string {
	return fmt.Sprintf(bookURLTemplate, bookID)
}

func (b *BookURLBuilder) ChapterListURL(bookID string) string {
	return fmt.Sprintf(chapterListURLTemplate, bookID)
}

func (b *BookURLBuilder) ChapterURL(resources ...string) string {
	if len(resources) == 1 {
		res := resources[0]
		if strings.HasPrefix(res, "http") {
			return res
		} else if strings.HasPrefix(res, "/") {
			return vendorProtocol + "://" + vendorHost + res
		}
	}

	log.Error().
		Str("vendor", "hjwzw").
		Strs("resources", resources).
		Msg("unexpected resources for building chapter url")

	return ""
}

func (b *BookURLBuilder) AvailabilityURL() string {
	return vendorProtocol + "://" + vendorHost
}
