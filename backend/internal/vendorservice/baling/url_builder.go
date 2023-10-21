package baling

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
	if len(resources) == 1 {
		res := resources[0]
		if strings.HasPrefix(res, "http") {
			return res
		} else if strings.HasPrefix(res, "/") {
			return vendorProtocol + "://" + vendorHost + res
		}
	}

	log.Error().
		Str("vendor", Host).
		Strs("resources", resources).
		Msg("unexpected resources for building chapter url")

	return ""
}

func (b *VendorService) AvailabilityURL() string {
	return vendorProtocol + "://" + vendorHost
}
