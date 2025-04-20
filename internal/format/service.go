package format

import (
	"context"
	"io"

	"github.com/htchan/BookSpider/internal/model"
)

type Service interface {
	ChaptersFromTxt(context.Context, io.Reader) (model.Chapters, error)

	WriteBookTxt(context.Context, *model.Book, model.Chapters, io.Writer) error
	WriteBookEpub(context.Context, *model.Book, model.Chapters, io.Writer) error
}
