package model

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/siongui/gojianfan"
	"go.opentelemetry.io/otel/attribute"
)

type Error struct {
	Err error
}

func (e Error) MarshalJSON() ([]byte, error) {
	if e.Err != nil {
		return []byte(`"` + e.Err.Error() + `"`), nil
	}

	return []byte(`""`), nil
}

func (e *Error) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	if str == "" {
		e.Err = nil
	} else {
		e.Err = errors.New(str)
	}
	return nil
}

type Book struct {
	Site          string     `json:"site"`
	ID            int        `json:"id"`
	HashCode      int        `json:"hash_code"`
	Title         string     `json:"title"`
	Type          string     `json:"type"`
	UpdateDate    string     `json:"update_date"`
	UpdateChapter string     `json:"update_chapter"`
	Status        StatusCode `json:"status"`
	IsDownloaded  bool       `json:"is_downloaded"`

	Writer Writer `json:"writer"`
	Error  Error  `json:"error"`
}

type BookGroup []Book

func NewBook(site string, id int) Book {
	return Book{
		Site:     site,
		ID:       id,
		HashCode: GenerateHash(),
	}
}

func GenerateHash() int {
	return int(time.Now().Unix())
}

func (bk *Book) HeaderInfo() string {
	return bk.Title + "\n" + bk.Writer.Name + "\n" + CONTENT_SEP + "\n\n"
}

func (bk Book) String() string {
	result := fmt.Sprintf("%v-%v", bk.Site, bk.ID)
	if bk.HashCode > 0 {
		result += fmt.Sprintf("-%v", bk.HashCode)
	}
	return result
}

func simplified(s string) string {
	return gojianfan.T2S(s)
}

func strToShortHex(s string) string {
	b := []byte(s)
	encoded := base64.StdEncoding.EncodeToString(b)
	return encoded
}

func (bk Book) Checksum() string {
	if len(bk.Title) > 100 {
		log.
			Error().
			Err(errors.New("book title is too long")).
			Str("book", bk.String()).
			Str("title", bk.Title).
			Msg("generate checksum failed")
		return ""
	}

	title := strings.ReplaceAll(bk.Title, " ", "")
	return strToShortHex(simplified(title))
}

func (bk *Book) FormatHashCode() string {
	return strconv.FormatInt(int64(bk.HashCode), 36)
}

func (bk *Book) IsEnd() bool {
	//TODO: fetch all chapter
	//hint: use book.generateEmptyChapters
	//TODO: check last n chapter to see if they contains any end keywords
	//hint: use len(chapters) and the n should come from book config
	if bk.UpdateDate < strconv.Itoa(time.Now().Year()-1) {
		return true
	}

	chapter := strings.ReplaceAll(bk.UpdateChapter, " ", "")
	for _, keyword := range ChapterEndKeywords {
		if strings.Contains(chapter, keyword) {
			return true
		}
	}
	return false
}

func (bk *Book) OtelAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("book_site", bk.Site),
		attribute.String("book_id", bk.String()),
		attribute.String("book_hash_code", bk.FormatHashCode()),
		attribute.String("book_title", bk.Title),
		attribute.String("book_writer", bk.Writer.Name),
		attribute.String("book_status", bk.Status.String()),
		attribute.Bool("book_is_downloaded", bk.IsDownloaded),
	}
}
