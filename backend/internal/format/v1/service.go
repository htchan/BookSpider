package format

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/htchan/BookSpider/internal/format"
	"github.com/htchan/BookSpider/internal/model"
)

type serviceImpl struct{}

var _ format.Service = (*serviceImpl)(nil)

func NewService() format.Service {
	return &serviceImpl{}
}

func (serv *serviceImpl) ChaptersFromTxt(ctx context.Context, reader io.Reader) (model.Chapters, error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("fail to read file: %w", err)
	}

	content := strings.Split(string(bytes), "\n")

	// split it into title, writer, chapters
	//title := content[0]
	//writer := content[1]

	chapters := make(model.Chapters, 0)
	// look the remaining content from index 1 to find content
	for i := 3; i < len(content); i++ {
		if content[i] == model.CONTENT_SEP {
			if content[i+1] == model.CONTENT_SEP && len(chapters) > 0 {
				chapters[len(chapters)-1].Content += "\n"
			}
			continue
		} else if i < len(content)-3 && content[i+1] == model.CONTENT_SEP {
			if len(chapters) > 0 && chapters[len(chapters)-1].Content == "" {
				chapters[len(chapters)-1].Content += content[i] + "\n"
			} else if len(chapters) > 0 && content[i+3] == model.CONTENT_SEP && content[i-1] != model.CONTENT_SEP {
				chapters[len(chapters)-1].Content += content[i] + "\n"
			} else {
				chapters = append(chapters, model.Chapter{Title: content[i], Index: len(chapters)})
			}
		} else if len(chapters) > 0 {
			chapters[len(chapters)-1].Content += content[i] + "\n"
		}
	}

	return chapters, nil
}
