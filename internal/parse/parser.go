package parse

//go:generate go tool mockgen -destination=../mock/parser/parser.go -package=mockparser . Parser
type Parser interface {
	ParseBook(html string) (*ParsedBookFields, error)
	ParseChapterList(html string) (*ParsedChapterList, error)
	ParseChapter(html string) (*ParsedChapterFields, error)
}
