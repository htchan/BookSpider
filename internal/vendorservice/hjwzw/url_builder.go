package hjwzw

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

func (b *VendorService) BookURL(bookID string) string {
	return fmt.Sprintf(bookURLTemplate, bookID)
}

func (b *VendorService) ChapterListURL(bookID string) string {
	return fmt.Sprintf(chapterListURLTemplate, bookID)
}

func (b *VendorService) ChapterURL(resources ...string) string {
	if len(resources) == 0 {
		return ""
	}

	uri := resources[0]
	if uri == "" {
		return ""
	} else if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
		return uri
	} else if strings.HasPrefix(uri, "/") {
		return vendorProtocol + "://" + vendorHost + uri
	}

	log.Error().
		Str("vendor", Host).
		Strs("resources", resources).
		Msg("unexpected resources for building chapter url")

	return uri
}

func (b *VendorService) AvailabilityURL() string {
	return vendorProtocol + "://" + vendorHost
}
