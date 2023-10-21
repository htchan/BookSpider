package vendor

import (
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type BookInfo struct {
	Title         string
	Writer        string
	Type          string
	UpdateDate    string
	UpdateChapter string
	// TODO: move update_date to udpate_datetime so we can save specific data to DB instead of a string
	UpdateDateTime time.Time
}

type ChapterListInfo struct {
	URL   string
	Title string
}
type ChapterList []ChapterListInfo

type ChapterInfo struct {
	Title string
	Body  string
}

//go:generate mockgen -destination=../mock/vendorservice/vendor_service.go -package=mockvendorservice . VendorService
type VendorService interface {
	// url builder
	BookURL(bookID string) string
	ChapterListURL(bookID string) string
	ChapterURL(resources ...string) string
	AvailabilityURL() string

	// content parser
	ParseBook(body string) (*BookInfo, error)
	ParseChapterList(body string) (ChapterList, error)
	ParseChapter(body string) (*ChapterInfo, error)
	IsAvailable(body string) bool
	FindMissingIds(ids []int) []int
}

func GetGoqueryContent(s *goquery.Selection) string {
	return strings.TrimSpace(s.Clone().Children().Remove().End().Text())
}
