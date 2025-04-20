package format

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/htchan/BookSpider/internal/format"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want format.Service
	}{
		{
			name: "happy flow",
			want: &serviceImpl{},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := NewService()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestChaptersFromTxt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		serv         *serviceImpl
		reader       io.Reader
		wantChapters model.Chapters
		wantError    error
	}{
		{
			name: "happy flow/format 1",
			serv: &serviceImpl{},
			reader: strings.NewReader(`title
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
`),
			wantChapters: model.Chapters{
				{Index: 0, Title: "chapter 1", Content: "content 1\n\ncontent 1\n"},
				{Index: 1, Title: "chapter 2", Content: "content 2\n\ncontent 2\n\n"},
			},
			wantError: nil,
		},
		{
			name: "happy flow/format 2",
			serv: &serviceImpl{},
			reader: strings.NewReader(`title
writer
--------------------

chapter 1
--------------------
content 1

content 1

chapter 2
--------------------
content 2

content 2
`),
			wantChapters: model.Chapters{
				{Index: 0, Title: "chapter 1", Content: "content 1\n\ncontent 1\n\n"},
				{Index: 1, Title: "chapter 2", Content: "content 2\n\ncontent 2\n\n"},
			},
			wantError: nil,
		},
		{
			name:         "happy flow/invalid format",
			serv:         &serviceImpl{},
			reader:       strings.NewReader(""),
			wantChapters: model.Chapters{},
			wantError:    nil,
		},
		{
			name: "happy flow/format 1 with one line content",
			serv: &serviceImpl{},
			reader: strings.NewReader(`title
writer
--------------------

chapter 1
--------------------
content 1
--------------------
chapter 2
--------------------
content 2
--------------------
`),
			wantChapters: model.Chapters{
				{Index: 0, Title: "chapter 1", Content: "content 1\n"},
				{Index: 1, Title: "chapter 2", Content: "content 2\n\n"},
			},
			wantError: nil,
		},
		{
			name: "happy flow/format 2 with one line content",
			serv: &serviceImpl{},
			reader: strings.NewReader(`title
writer
--------------------

chapter 1
--------------------
content 1

chapter 2
--------------------
content 2
`),
			wantChapters: model.Chapters{
				{Index: 0, Title: "chapter 1", Content: "content 1\n\n"},
				{Index: 1, Title: "chapter 2", Content: "content 2\n\n"},
			},
			wantError: nil,
		},
		{
			name: "happy flow/format 1 with empty content",
			serv: &serviceImpl{},
			reader: strings.NewReader(`title
writer
--------------------

chapter 1
--------------------

--------------------
chapter 2
--------------------

--------------------
`),
			wantChapters: model.Chapters{
				{Index: 0, Title: "chapter 1", Content: "\n"},
				{Index: 1, Title: "chapter 2", Content: "\n\n"},
			},
			wantError: nil,
		},
		{
			name: "happy flow/format 2 with empty content",
			serv: &serviceImpl{},
			reader: strings.NewReader(`title
writer
--------------------

chapter 1
--------------------


chapter 2
--------------------

`),
			wantChapters: model.Chapters{
				{Index: 0, Title: "chapter 1", Content: "\n\n"},
				{Index: 1, Title: "chapter 2", Content: "\n\n"},
			},
			wantError: nil,
		},
		{
			name:         "invalid flow/invalid format",
			serv:         &serviceImpl{},
			reader:       strings.NewReader(""),
			wantChapters: model.Chapters{},
			wantError:    nil,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := test.serv.ChaptersFromTxt(context.Background(), test.reader)
			assert.Equal(t, test.wantChapters, got)
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}
