package hjwzw

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
				<meta property="og:novel:update_time" content="2000-01-01" />
				<meta property="og:novel:latest_chapter_name" content="chapter name" />
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
				<meta property="og:novel:update_time" content="2000-01-01" />
				<meta property="og:novel:latest_chapter_name" content="chapter name" />
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
				<meta property="og:novel:update_time" content="2000-01-01" />
				<meta property="og:novel:latest_chapter_name" content="chapter name" />
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
				<meta property="og:novel:update_time" content="2000-01-01" />
				<meta property="og:novel:latest_chapter_name" content="chapter name" />
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
				<meta property="og:novel:latest_chapter_name" content="chapter name" />
			</data>`,
			want: &client.BookInfo{
				Title: "book name", Author: "author", Type: "type",
				UpdateDate: time.Now().UTC().Truncate(time.Minute), UpdateChapter: "chapter name",
			},
			wantErr: client.ErrBookDateNotFound.Error(),
		},
		{
			name: "chapter not found",
			body: `<data>
				<meta property="og:novel:book_name" content="book name" />
				<meta property="og:novel:author" content="author" />
				<meta property="og:novel:category" content="type" />
					<meta property="og:novel:update_time" content="2000-01-01" />
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
				UpdateDate: time.Now().UTC().Truncate(time.Minute),
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
				<div id="tbchapterlist"><table><tbody><tr>
					<td><a href="chapter url 1">chapter name 1</a></td>
					<td><a href="chapter url 2">chapter name 2</a></td>
					<td><a href="chapter url 3">chapter name 3</a></td>
					<td><a href="chapter url 4">chapter name 4</a></td>
				</tr></tbody></table></div>
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
				<div id="tbchapterlist"><table><tbody><tr>
					<td><a href="chapter url 1">chapter name 1</a></td>
					<td><a href="">chapter name 2</a></td>
					<td><a href="chapter url 3">chapter name 3</a></td>
					<td><a href="chapter url 4">chapter name 4</a></td>
				</tr></tbody></table></div>
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
				<div id="tbchapterlist"><table><tbody><tr>
					<td><a href="chapter url 1">chapter name 1</a></td>
					<td><a href="chapter url 2">chapter name 2</a></td>
					<td><a href="chapter url 3"></a></td>
					<td><a href="chapter url 4">chapter name 4</a></td>
				</tr></tbody></table></div>
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
				<table><tbody><tr><td>
					<h1>chapter name</h1>
					<div></div><div></div><div></div><div></div>
					<div><p>chapter content</p></div>
				</td></tr></tbody></table>
			</data>`,
			want: &client.ChapterContent{
				Title: "chapter name", Body: "chapter content",
			},
			wantErr: "",
		},
		{
			name: "title empty",
			body: `<data>
				<table><tbody><tr><td>
					<h1></h1>
					<div></div><div></div><div></div><div></div>
					<div><p>chapter content</p></div>
				</td></tr></tbody></table>
			</data>`,
			want: &client.ChapterContent{
				Title: "", Body: "chapter content",
			},
			wantErr: client.ErrChapterTitleNotFound.Error(),
		},
		{
			name: "body empty",
			body: `<data>
				<table><tbody><tr><td>
					<h1>chapter name</h1>
					<div></div><div></div><div></div><div></div>
					<div><p></p></div>
				</td></tr></tbody></table>
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
			body: "黃金屋",
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
