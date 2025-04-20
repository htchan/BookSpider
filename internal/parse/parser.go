package parse

//go:generate mockgen -destination=../mock/parser/parser.go -package=mockparser . Parser
type Parser interface {
	ParseBook(html string) (*ParsedBookFields, error)
	ParseChapterList(html string) (*ParsedChapterList, error)
	ParseChapter(html string) (*ParsedChapterFields, error)
}
