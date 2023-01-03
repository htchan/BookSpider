package parse

type Parser interface {
	ParseBook(html string) (*ParsedBookFields, error)
	ParseChapterList(html string) (*ParsedChapterList, error)
	ParseChapter(html string) (*ParsedChapterFields, error)
}
