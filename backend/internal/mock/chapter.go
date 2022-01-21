package mock

import (
	"golang.org/x/text/encoding"
	"fmt"
)

func ChapterGetWebSuccess(url string, _ int, _ *encoding.Decoder, _ int) (html string, i int) {
	return fmt.Sprintf("chapter-content-%v-content-regex", url), 0
}

func ChapterGetWebInvalid(_ string, _ int, _ *encoding.Decoder, _ int) (html string, i int) {
	return "", 0
}

func ChapterGetWebNoContent(_ string, _ int, _ *encoding.Decoder, _ int) (html string, i int) {
	return "hello", 0
}