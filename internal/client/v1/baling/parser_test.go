package baling

import (
	"testing"
	"time"

	"github.com/htchan/BookSpider/internal/client/v1"
	"github.com/stretchr/testify/assert"
)

func TestParser_ParseBook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		body    string
		want    *client.BookInfo
		wantErr string
	}{
		{
			name: "happy flow",
			body: `<data>
				<meta property="og:novel:book_name" content="book name" />
				<meta property="og:novel:author" content="author" />
				<meta property="og:novel:category" content="type" />
				<div><div class="txt_info"></div><div class="txt_info"></div>
				<div class="txt_info"></div><div class="txt_info">2000-01-01</div></div>
				<div class="yulan"><a>chapter name</a></div>
				</div>
			</data>`,
			want: &client.BookInfo{
				Title: "book name", Author: "author", Type: "type",
				UpdateDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), UpdateChapter: "chapter name",
			},
			wantErr: "",
		},
		{
			name: "title not found",
			body: `<data>
				<meta property="og:novel:author" content="author" />
				<meta property="og:novel:category" content="type" />
				<div><div class="txt_info"></div><div class="txt_info"></div>
				<div class="txt_info"></div><div class="txt_info">更新时间：2000-01-01</div></div>
				<div class="yulan"><a>chapter name</a></div>
			</data>`,
			want: &client.BookInfo{
				Author: "author", Type: "type",
				UpdateDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), UpdateChapter: "chapter name",
			},
			wantErr: client.ErrBookTitleNotFound.Error(),
		},
		{
			name: "writer not found",
			body: `<data>
				<meta property="og:novel:book_name" content="book name" />
				<meta property="og:novel:category" content="type" />
				<div><div class="txt_info"></div><div class="txt_info"></div>
				<div class="txt_info"></div><div class="txt_info">更新时间：2000-01-01</div></div>
				<div class="yulan"><a>chapter name</a></div>
			</data>`,
			want: &client.BookInfo{
				Title: "book name", Type: "type",
				UpdateDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), UpdateChapter: "chapter name",
			},
			wantErr: client.ErrBookWriterNotFound.Error(),
		},
		{
			name: "type not found",
			body: `<data>
				<meta property="og:novel:book_name" content="book name" />
				<meta property="og:novel:author" content="author" />
				<div><div class="txt_info"></div><div class="txt_info"></div>
				<div class="txt_info"></div><div class="txt_info">更新时间：2000-01-01</div></div>
				<div class="yulan"><a>chapter name</a></div>
			</data>`,
			want: &client.BookInfo{
				Title: "book name", Author: "author",
				UpdateDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), UpdateChapter: "chapter name",
			},
			wantErr: client.ErrBookTypeNotFound.Error(),
		},
		{
			name: "date not found",
			body: `<data>
				<meta property="og:novel:book_name" content="book name" />
				<meta property="og:novel:author" content="author" />
				<meta property="og:novel:category" content="type" />
				<div class="yulan"><a>chapter name</a></div>
			</data>`,
			want: &client.BookInfo{
				Title: "book name", Author: "author", Type: "type",
				UpdateDate: time.Now().UTC().Truncate(time.Second), UpdateChapter: "chapter name",
			},
			wantErr: client.ErrBookDateNotFound.Error(),
		},
		{
			name: "chapter not found",
			body: `<data>
				<meta property="og:novel:book_name" content="book name" />
				<meta property="og:novel:author" content="author" />
				<meta property="og:novel:category" content="type" />
				<div><div class="txt_info"></div><div class="txt_info"></div>
				<div class="txt_info"></div><div class="txt_info">更新时间：2000-01-01</div></div>
			</data>`,
			want: &client.BookInfo{
				Title: "book name", Author: "author", Type: "type",
				UpdateDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: client.ErrBookChapterNotFound.Error(),
		},
		{
			name: "all fields not found",
			body: "<data></data>",
			want: &client.BookInfo{
				UpdateDate: time.Now().UTC().Truncate(time.Second),
			},
			wantErr: client.ErrFieldsNotFound.Error(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseBook(test.body)
			assert.Equal(t, test.want, got)
			if test.wantErr != "" {
				assert.ErrorContains(t, err, test.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParser_ParseChapterList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		body    string
		want    client.ChapterEntryList
		wantErr string
	}{
		{
			name: "happy flow",
			body: `<data>
				<div class="yulan">
					<li><a href="chapter url 1">chapter name 1</a></li>
					<li><a href="chapter url 2">chapter name 2</a></li>
					<li><a href="chapter url 3">chapter name 3</a></li>
					<li><a href="chapter url 4">chapter name 4</a></li>
				</div>
			</data>`,
			want: client.ChapterEntryList{
				{URL: "chapter url 1", Title: "chapter name 1"},
				{URL: "chapter url 2", Title: "chapter name 2"},
				{URL: "chapter url 3", Title: "chapter name 3"},
				{URL: "chapter url 4", Title: "chapter name 4"},
			},
			wantErr: "",
		},
		{
			name: "2nd chapter missing href",
			body: `<data>
				<div class="yulan">
					<li><a href="chapter url 1">chapter name 1</a></li>
					<li><a href="">chapter name 2</a></li>
					<li><a href="chapter url 3">chapter name 3</a></li>
					<li><a href="chapter url 4">chapter name 4</a></li>
				</div>
			</data>`,
			want: client.ChapterEntryList{
				{URL: "chapter url 1", Title: "chapter name 1"},
				{URL: "", Title: "chapter name 2"},
				{URL: "chapter url 3", Title: "chapter name 3"},
				{URL: "chapter url 4", Title: "chapter name 4"},
			},
			wantErr: client.ErrChapterListUrlNotFound.Error(),
		},
		{
			name: "3nd chapter missing title",
			body: `<data>
				<div class="yulan">
					<li><a href="chapter url 1">chapter name 1</a></li>
					<li><a href="chapter url 2">chapter name 2</a></li>
					<li><a href="chapter url 3"></a></li>
					<li><a href="chapter url 4">chapter name 4</a></li>
				</div>
			</data>`,
			want: client.ChapterEntryList{
				{URL: "chapter url 1", Title: "chapter name 1"},
				{URL: "chapter url 2", Title: "chapter name 2"},
				{URL: "chapter url 3", Title: ""},
				{URL: "chapter url 4", Title: "chapter name 4"},
			},
			wantErr: client.ErrChapterListTitleNotFound.Error(),
		},
		{
			name:    "no chapters found",
			body:    `<data></data>`,
			want:    nil,
			wantErr: client.ErrChapterListEmpty.Error(),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseChapterList(test.body)
			assert.Equal(t, test.want, got)
			if test.wantErr != "" {
				assert.ErrorContains(t, err, test.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParser_ParseChapter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		body    string
		want    *client.ChapterContent
		wantErr string
	}{
		{
			name: "happy flow",
			body: `<data>
				<div class="date"><h1>chapter name</h1></div>
				<div class="book_content">chapter content</div>
			</data>`,
			want: &client.ChapterContent{
				Title: "chapter name", Body: "chapter content",
			},
			wantErr: "",
		},
		{
			name: "title empty",
			body: `<data>
			<div class="date"><h1></h1></div>
			<div class="book_content">chapter content</div>
			</data>`,
			want: &client.ChapterContent{
				Title: "", Body: "chapter content",
			},
			wantErr: client.ErrChapterTitleNotFound.Error(),
		},
		{
			name: "body empty",
			body: `<data>
			<div class="date"><h1>chapter name</h1></div>
			<div class="book_content"></div>
			</data>`,
			want: &client.ChapterContent{
				Title: "chapter name", Body: "",
			},
			wantErr: client.ErrChapterContentNotFound.Error(),
		},
		{
			name:    "all fields not found",
			body:    "<data></data>",
			want:    &client.ChapterContent{},
			wantErr: client.ErrFieldsNotFound.Error(),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseChapter(test.body)
			assert.Equal(t, test.want, got)
			if test.wantErr != "" {
				assert.ErrorContains(t, err, test.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParser_IsAvailable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
		want bool
	}{
		{
			name: "return true",
			body: "80txt",
			want: true,
		},
		{
			name: "return false",
			body: "",
			want: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := isAvailable(test.body)
			assert.Equal(t, test.want, got)
		})
	}
}
