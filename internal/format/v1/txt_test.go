package format

import (
	"bytes"
	"context"
	"testing"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestWriteBookTxt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		serv        *serviceImpl
		bk          *model.Book
		chapters    model.Chapters
		wantContent string
		wantError   error
	}{
		{
			name: "happy flow",
			serv: &serviceImpl{},
			bk:   &model.Book{Title: "title", Writer: model.Writer{Name: "writer"}},
			chapters: model.Chapters{
				{Title: "chapter 1", Content: "content 1\ncontent 1"},
				{Title: "chapter 2", Content: "content 2\ncontent 2"},
			},
			wantContent: `title
writer
--------------------

chapter 1
--------------------
content 1
content 1
--------------------
chapter 2
--------------------
content 2
content 2
--------------------
`,
			wantError: nil,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var buffer bytes.Buffer

			err := test.serv.WriteBookTxt(context.Background(), test.bk, test.chapters, &buffer)
			assert.Equal(t, test.wantContent, buffer.String())
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}
