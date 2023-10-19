package vendor

//go:generate mockgen -destination=../mock/vendorservice/url_builder.go -package=mockvendorservice . BookURLBuilder
type BookURLBuilder interface {
	BookURL(bookID string) string
	ChapterListURL(bookID string) string
	ChapterURL(resources ...string) string
	AvailabilityURL() string
}
