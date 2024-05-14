package format

import (
	"context"
	"io"

	"github.com/htchan/BookSpider/internal/model"
)

func (serv *serviceImpl) WriteBookTxt(ctx context.Context, bk *model.Book, chapters model.Chapters, writer io.Writer) error {
	// write title and writer
	writer.Write([]byte(bk.HeaderInfo()))
	for _, chapter := range chapters {
		writer.Write([]byte(chapter.ContentString()))
	}
	return nil
}
