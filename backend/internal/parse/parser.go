package parse

//go:generate mockgen -source=./$GOFILE -destination=../mock/$GOFILE -package=mock
type Parser interface {
	ParseBook(html string) (*ParsedBookFields, error)
	ParseChapterList(html string) (*ParsedChapterList, error)
	ParseChapter(html string) (*ParsedChapterFields, error)
}
