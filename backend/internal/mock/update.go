package mock

import (

	"golang.org/x/text/encoding"
)

func UpdateGetWebSuccess(_ string, _ int, _ *encoding.Decoder, _ int) (html string, i int) {
	return "title-regex writer-regex type-regex last-update-regex " +
		"last-chapter-regex", 0
}

func UpdateGetWebPartialFail(_ string, _ int, _ *encoding.Decoder, _ int) (html string, i int) {
	return "title-regex writer-regex type-regex last-update-regex ", 0
}

func UpdateGetWebEmpty(_ string, _ int, _ *encoding.Decoder, _ int) (html string, i int) {
	return "", 10
}

func UpdateGetWebNumber(_ string, _ int, _ *encoding.Decoder, _ int) (html string, i int) {
	return "200", 10
}