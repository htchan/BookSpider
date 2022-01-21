package mock

import (
	"golang.org/x/text/encoding"
	"fmt"
)

func DownloadGetWebSuccess(_ string, _ int, _ *encoding.Decoder, _ int) (html string, i int) {
	return "chapter-url-regex-1 chapter-title-regex-1 " +
	"chapter-url-regex-2 chapter-title-regex-2 " +
	"chapter-url-regex-3 chapter-title-regex-3 " +
	"chapter-url-regex-4 chapter-title-regex-4 ", 0
}

func DownloadGetWebImbalanceUrlTitle(_ string, _ int, _ *encoding.Decoder, _ int) (html string, i int) {
	return "chapter-url-regex-1 chapter-title-regex-1" +
	"chapter-url-regex-2 chapter-title-regex-2" +
	"chapter-url-regex-3 chapter-title-regex-3" +
	"chapter-url-regex-4", 0
}

func DownloadGetWebEmpty(_ string, _ int, _ *encoding.Decoder, _ int) (html string, i int) {
	return "", 10
}

func DownloadGetWebNoUrl(_ string, _ int, _ *encoding.Decoder, _ int) (html string, i int) {
	return "hello", 10
}

func DownloadGetWebFullSuccess(url string, _ int, _ *encoding.Decoder, _ int) (html string, i int) {
	return "chapter-url-regex-1 chapter-title-regex-1 " +
		"chapter-url-regex-2 chapter-title-regex-2 " +
		"chapter-url-regex-3 chapter-title-regex-3 " +
		"chapter-url-regex-4 chapter-title-regex-4 " +
		fmt.Sprintf("chapter-content-%v-content-regex", url), 0
}