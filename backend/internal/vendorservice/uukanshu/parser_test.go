package uukanshu

import (
	"testing"

	vendor "github.com/htchan/BookSpider/internal/vendorservice"
	"github.com/stretchr/testify/assert"
)

func TestParser_ParseBook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		body      string
		want      *vendor.BookInfo
		wantError error
	}{
		{
			name: "happy flow",
			body: `<data>
				<div class="xiaoshuo_content"><dl class="jieshao"><dd class="jieshao_content">
					<h1><a title="book name最新章节"></a></h1>
					<h2><a>author</a></h2>
					<div class="shijian">d a	t
					e</div>
				</dd></dl></div>
				<div class="weizhi"><div class="path"><a></a><a>type</a></div></div>
				<div class="zhangjie"><ul id="chapterList"><li><a>chapter name</a></li></ul></div>
			</data>`,
			want: &vendor.BookInfo{
				Title: "book name", Writer: "author", Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter name",
			},
			wantError: nil,
		},
		{
			name: "title not found",
			body: `<data>
				<div class="xiaoshuo_content"><dl class="jieshao"><dd class="jieshao_content">
					<h2><a>author</a></h2>
					<div class="shijian">d a	t
					e</div>
				</dd></dl></div>
				<div class="weizhi"><div class="path"><a></a><a>type</a></div></div>
				<div class="zhangjie"><ul id="chapterList"><li><a>chapter name</a></li></ul></div>
			</data>`,
			want: &vendor.BookInfo{
				Writer: "author", Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter name",
			},
			wantError: vendor.ErrBookTitleNotFound,
		},
		{
			name: "writer not found",
			body: `<data>
				<div class="xiaoshuo_content"><dl class="jieshao"><dd class="jieshao_content">
					<h1><a title="book name最新章节"></a></h1>
					<h2><a></a></h2>
					<div class="shijian">d a	t
					e</div>
				</dd></dl></div>
				<div class="weizhi"><div class="path"><a></a><a>type</a></div></div>
				<div class="zhangjie"><ul id="chapterList"><li><a>chapter name</a></li></ul></div>
			</data>`,
			want: &vendor.BookInfo{
				Title: "book name", Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter name",
			},
			wantError: vendor.ErrBookWriterNotFound,
		},
		{
			name: "type not found",
			body: `<data>
			<div class="xiaoshuo_content"><dl class="jieshao"><dd class="jieshao_content">
				<h1><a title="book name最新章节"></a></h1>
				<h2><a>author</a></h2>
				<div class="shijian">d a	t
				e</div>
			</dd></dl></div>
			<div class="zhangjie"><ul id="chapterList"><li><a>chapter name</a></li></ul></div>
			</data>`,
			want: &vendor.BookInfo{
				Title: "book name", Writer: "author",
				UpdateDate: "date", UpdateChapter: "chapter name",
			},
			wantError: vendor.ErrBookTypeNotFound,
		},
		{
			name: "date not found",
			body: `<data>
			<div class="xiaoshuo_content"><dl class="jieshao"><dd class="jieshao_content">
				<h1><a title="book name最新章节"></a></h1>
				<h2><a>author</a></h2>
			</dd></dl></div>
			<div class="weizhi"><div class="path"><a></a><a>type</a></div></div>
			<div class="zhangjie"><ul id="chapterList"><li><a>chapter name</a></li></ul></div>
			</data>`,
			want: &vendor.BookInfo{
				Title: "book name", Writer: "author", Type: "type",
				UpdateChapter: "chapter name",
			},
			wantError: vendor.ErrBookDateNotFound,
		},
		{
			name: "chapter not found",
			body: `<data>
			<div class="xiaoshuo_content"><dl class="jieshao"><dd class="jieshao_content">
				<h1><a title="book name最新章节"></a></h1>
				<h2><a>author</a></h2>
				<div class="shijian">d a	t
				e</div>
			</dd></dl></div>
			<div class="weizhi"><div class="path"><a></a><a>type</a></div></div>
			</data>`,
			want: &vendor.BookInfo{
				Title: "book name", Writer: "author", Type: "type",
				UpdateDate: "date",
			},
			wantError: vendor.ErrBookChapterNotFound,
		},
		{
			name:      "all fields not found",
			body:      "<data></data>",
			want:      &vendor.BookInfo{},
			wantError: vendor.ErrFieldsNotFound,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			p := VendorService{}
			got, err := p.ParseBook(test.body)
			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}

func TestParser_ParseChapterList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		body      string
		want      vendor.ChapterList
		wantError error
	}{
		{
			name: "happy flow",
			body: `<data>
				<div class="zhangjie"><ul id="chapterList">
					<li><a href="chapter url 1">chapter name 1</a></li>
					<li><a href="chapter url 2">chapter name 2</a></li>
					<li><a href="chapter url 3">chapter name 3</a></li>
					<li><a href="chapter url 4">chapter name 4</a></li>
				</ul></div>
			</data>`,
			want: vendor.ChapterList{
				{URL: "chapter url 1", Title: "chapter name 1"},
				{URL: "chapter url 2", Title: "chapter name 2"},
				{URL: "chapter url 3", Title: "chapter name 3"},
				{URL: "chapter url 4", Title: "chapter name 4"},
			},
			wantError: nil,
		},
		{
			name: "2nd chapter missing href",
			body: `<data>
				<div class="zhangjie"><ul id="chapterList">
					<li><a href="chapter url 1">chapter name 1</a></li>
					<li><a href="">chapter name 2</a></li>
					<li><a href="chapter url 3">chapter name 3</a></li>
					<li><a href="chapter url 4">chapter name 4</a></li>
				</ul></div>
			</data>`,
			want: vendor.ChapterList{
				{URL: "chapter url 1", Title: "chapter name 1"},
				{URL: "", Title: "chapter name 2"},
				{URL: "chapter url 3", Title: "chapter name 3"},
				{URL: "chapter url 4", Title: "chapter name 4"},
			},
			wantError: vendor.ErrChapterListUrlNotFound,
		},
		{
			name: "3nd chapter missing title",
			body: `<data>
				<div class="zhangjie"><ul id="chapterList">
					<li><a href="chapter url 1">chapter name 1</a></li>
					<li><a href="chapter url 2">chapter name 2</a></li>
					<li><a href="chapter url 3"></a></li>
					<li><a href="chapter url 4">chapter name 4</a></li>
				</ul></div>
			</data>`,
			want: vendor.ChapterList{
				{URL: "chapter url 1", Title: "chapter name 1"},
				{URL: "chapter url 2", Title: "chapter name 2"},
				{URL: "chapter url 3", Title: ""},
				{URL: "chapter url 4", Title: "chapter name 4"},
			},
			wantError: vendor.ErrChapterListTitleNotFound,
		},
		{
			name:      "no chapters found",
			body:      `<data></data>`,
			want:      nil,
			wantError: vendor.ErrChapterListEmpty,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			p := VendorService{}
			got, err := p.ParseChapterList(test.body)
			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}

func TestParser_ParseChapter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		body      string
		want      *vendor.ChapterInfo
		wantError error
	}{
		{
			name: "happy flow",
			body: `<data>
				<div class="zhengwen_box"><div class="box_left"><div class="w_main">
					<div class="h1title"><h1 id="timu">chapter name</h1></div>
					<div class="contentbox"><div id="contentbox">chapter content</div></div>
				</div></div></div>
			</data>`,
			want: &vendor.ChapterInfo{
				Title: "chapter name", Body: "chapter content",
			},
			wantError: nil,
		},
		{
			name: "title not found",
			body: `<data>
				<div class="zhengwen_box"><div class="box_left"><div class="w_main">
					<div class="contentbox"><div id="contentbox">chapter content</div></div>
				</div></div></div>
			</data>`,
			want: &vendor.ChapterInfo{
				Title: "", Body: "chapter content",
			},
			wantError: vendor.ErrChapterTitleNotFound,
		},
		{
			name: "body not found",
			body: `<data>
				<div class="zhengwen_box"><div class="box_left"><div class="w_main">
					<div class="h1title"><h1 id="timu">chapter name</h1></div>
				</div></div></div>
			</data>`,
			want: &vendor.ChapterInfo{
				Title: "chapter name", Body: "",
			},
			wantError: vendor.ErrChapterContentNotFound,
		},
		{
			name:      "all fields not found",
			body:      "<data></data>",
			want:      &vendor.ChapterInfo{},
			wantError: vendor.ErrFieldsNotFound,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			p := VendorService{}
			got, err := p.ParseChapter(test.body)
			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, err, test.wantError)
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

			p := VendorService{}
			got := p.IsAvailable(test.body)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestParser_FindMissingIds(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ids  []int
		want []int
	}{
		{
			name: "no missing ids",
			ids:  []int{4, 2, 3, 1, 5},
			want: nil,
		},
		{
			name: "some id is missing",
			ids:  []int{3, 5, 1},
			want: []int{2, 4},
		},
		{
			name: "input ids contains negative",
			ids:  []int{3, -1},
			want: []int{1, 2},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			p := VendorService{}
			got := p.FindMissingIds(test.ids)
			assert.Equal(t, test.want, got)
		})
	}

}
