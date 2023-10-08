package vendor

//go:generate mockgen -destination=../mock/vendor/url_builder.go -package=mockvendor . BookURLBuilder
type BookURLBuilder interface {
	BookURL(bookID string) string
	ChapterListURL(bookID string) string
	ChapterURL(resources ...string) string
	AvailabilityURL() string
}
