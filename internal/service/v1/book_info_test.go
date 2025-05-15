package service

import (
	"context"
	"os"
	"testing"

	"github.com/htchan/BookSpider/internal/config/v2"
	mockrepo "github.com/htchan/BookSpider/internal/mock/repo"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	serv "github.com/htchan/BookSpider/internal/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestServiceImpl_BookInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		bk   *model.Book
		want string
	}{
		{
			name: "happy flow",
			bk:   &model.Book{HashCode: 0},
			want: `{"site":"","id":0,"hash_code":"0","title":"","writer":"","type":"","update_date":"","update_chapter":"","status":"ERROR","is_downloaded":false,"error":""}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := new(ServiceImpl).BookInfo(context.Background(), test.bk)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestServiceImpl_BookContent(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() {
		assert.NoError(t, os.RemoveAll("./book-content"))
	})

	if !assert.NoError(t, os.Mkdir("./book-content", os.ModePerm)) ||
		!assert.NoError(t, os.WriteFile("./book-content/123.txt", []byte("test"), 0644)) ||
		!assert.NoError(t, os.WriteFile("./book-content/123-va.txt", []byte("test v2"), 0644)) {
		return
	}

	tests := []struct {
		name      string
		serv      *ServiceImpl
		bk        *model.Book
		want      string
		wantError error
	}{
		{
			name:      "book without hashcode content exist",
			serv:      &ServiceImpl{conf: config.SiteConfig{Storage: "./book-content"}},
			bk:        &model.Book{ID: 123, IsDownloaded: true},
			want:      "test",
			wantError: nil,
		},
		{
			name:      "book with hashcode content exist",
			serv:      &ServiceImpl{conf: config.SiteConfig{Storage: "./book-content"}},
			bk:        &model.Book{ID: 123, HashCode: 10, IsDownloaded: true},
			want:      "test v2",
			wantError: nil,
		},
		{
			name:      "book not downloaded",
			serv:      &ServiceImpl{conf: config.SiteConfig{Storage: "./book-content"}},
			bk:        &model.Book{ID: 123, HashCode: 10, IsDownloaded: false},
			want:      "",
			wantError: serv.ErrBookNotDownload,
		},
		{
			name:      "book content not exist",
			serv:      &ServiceImpl{conf: config.SiteConfig{Storage: "./book-content"}},
			bk:        &model.Book{ID: 456, IsDownloaded: true},
			want:      "",
			wantError: serv.ErrBookFileNotFound,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := test.serv.BookContent(context.Background(), test.bk)
			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}

func TestServiceImpl_BookChapters(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() {
		assert.NoError(t, os.RemoveAll("./book-chapters"))
	})

	if !assert.NoError(t, os.Mkdir("./book-chapters", os.ModePerm)) ||
		!assert.NoError(t, os.WriteFile("./book-chapters/123.txt", []byte(
			"test"+model.CONTENT_SEP+
				"title 1"+model.CONTENT_SEP+"content 1"+model.CONTENT_SEP+
				"title 2"+model.CONTENT_SEP+"content 2"+model.CONTENT_SEP,
		), 0644)) ||
		!assert.NoError(t, os.WriteFile("./book-chapters/123-va.txt", []byte("test"), 0644)) {
		return
	}

	tests := []struct {
		name      string
		serv      *ServiceImpl
		bk        *model.Book
		want      model.Chapters
		wantError error
	}{
		{
			name: "successfully read and parse content to chapters",
			serv: &ServiceImpl{conf: config.SiteConfig{Storage: "./book-chapters"}},
			bk:   &model.Book{ID: 123, IsDownloaded: true},
			want: model.Chapters{
				{Index: 0, Title: "title 1", Content: "content 1"},
				{Index: 1, Title: "title 2", Content: "content 2"},
			},
			wantError: nil,
		},
		{
			name:      "fail to parse content to chapters",
			serv:      &ServiceImpl{conf: config.SiteConfig{Storage: "./book-chapters"}},
			bk:        &model.Book{ID: 123, HashCode: 10, IsDownloaded: true},
			want:      nil,
			wantError: model.ErrCannotParseContent,
		},
		{
			name:      "fail to read content ",
			serv:      &ServiceImpl{conf: config.SiteConfig{Storage: "./book-chapters"}},
			bk:        &model.Book{ID: 456, HashCode: 0, IsDownloaded: true},
			want:      nil,
			wantError: serv.ErrBookFileNotFound,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := test.serv.BookChapters(context.Background(), test.bk)
			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}

func TestServiceImpl_Book(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		getService func(*gomock.Controller) *ServiceImpl
		id         string
		hash       string
		want       *model.Book
		wantError  error
	}{
		{
			name: "book found with pure ID",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBookById(gomock.Any(), 123).Return(&model.Book{ID: 123, HashCode: 0}, nil)

				return &ServiceImpl{rpo: rpo}
			},
			id:        "123",
			hash:      "",
			want:      &model.Book{ID: 123, HashCode: 0},
			wantError: nil,
		},
		{
			name: "book found with ID and Hashcode",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBookByIdHash(gomock.Any(), 123, 10).Return(&model.Book{ID: 123, HashCode: 10}, nil)

				return &ServiceImpl{rpo: rpo}
			},
			id:        "123",
			hash:      "A",
			want:      &model.Book{ID: 123, HashCode: 10},
			wantError: nil,
		},
		{
			name: "invalid id",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)

				return &ServiceImpl{rpo: rpo}
			},
			id:        "abc",
			hash:      "",
			want:      nil,
			wantError: serv.ErrInvalidBookID,
		},
		{
			name: "invalid hashcode",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)

				return &ServiceImpl{rpo: rpo}
			},
			id:        "123",
			hash:      "abc-def",
			want:      nil,
			wantError: serv.ErrInvalidHashCode,
		},
		{
			name: "book not found",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBookById(gomock.Any(), 123).Return(nil, repo.ErrBookNotExist)

				return &ServiceImpl{rpo: rpo}
			},
			id:        "123",
			hash:      "",
			want:      nil,
			wantError: repo.ErrBookNotExist,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			got, err := test.getService(ctrl).Book(context.Background(), test.id, test.hash)
			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}

func TestServiceImpl_BookGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		getService    func(*gomock.Controller) *ServiceImpl
		id            string
		hash          string
		wantBook      *model.Book
		wantBookGroup *model.BookGroup
		wantError     error
	}{
		{
			name: "book group found with pure ID",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBookGroupByID(gomock.Any(), 123).Return(
					model.BookGroup{{ID: 123, HashCode: 0}, {ID: 456, HashCode: 0}},
					nil,
				)

				return &ServiceImpl{rpo: rpo}
			},
			id:            "123",
			hash:          "",
			wantBook:      &model.Book{ID: 123, HashCode: 0},
			wantBookGroup: &model.BookGroup{{ID: 456, HashCode: 0}},
			wantError:     nil,
		},
		{
			name: "book group found with ID and Hashcode",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBookGroupByIDHash(gomock.Any(), 123, 10).Return(
					model.BookGroup{{ID: 123, HashCode: 10}, {ID: 456, HashCode: 0}},
					nil,
				)

				return &ServiceImpl{rpo: rpo}
			},
			id:            "123",
			hash:          "A",
			wantBook:      &model.Book{ID: 123, HashCode: 10},
			wantBookGroup: &model.BookGroup{{ID: 456, HashCode: 0}},
			wantError:     nil,
		},
		{
			name: "invalid id",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)

				return &ServiceImpl{rpo: rpo}
			},
			id:            "abc",
			hash:          "",
			wantBook:      nil,
			wantBookGroup: nil,
			wantError:     serv.ErrInvalidBookID,
		},
		{
			name: "invalid hashcode",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)

				return &ServiceImpl{rpo: rpo}
			},
			id:            "123",
			hash:          "abc-def",
			wantBook:      nil,
			wantBookGroup: nil,
			wantError:     serv.ErrInvalidHashCode,
		},
		{
			name: "book not found",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBookGroupByID(gomock.Any(), 123).Return(nil, repo.ErrBookNotExist)

				return &ServiceImpl{rpo: rpo}
			},
			id:            "123",
			hash:          "",
			wantBook:      nil,
			wantBookGroup: nil,
			wantError:     repo.ErrBookNotExist,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gotBk, gotBkGroup, err := test.getService(ctrl).BookGroup(context.Background(), test.id, test.hash)
			assert.Equal(t, test.wantBook, gotBk)
			assert.Equal(t, test.wantBookGroup, gotBkGroup)
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}

func TestServiceImpl_QueryBooks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		getService func(*gomock.Controller) *ServiceImpl
		title      string
		writer     string
		limit      int
		offset     int
		want       []model.Book
		wantError  error
	}{
		{
			name: "happy flow with books",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBooksByTitleWriter(gomock.Any(), "title", "writer", 10, 0).
					Return([]model.Book{{ID: 123, HashCode: 0}}, nil)

				return &ServiceImpl{rpo: rpo}
			},
			title:     "title",
			writer:    "writer",
			limit:     10,
			offset:    0,
			want:      []model.Book{{ID: 123, HashCode: 0}},
			wantError: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := test.getService(ctrl)

			got, err := svc.QueryBooks(
				context.Background(), test.title, test.writer, test.limit, test.offset,
			)
			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}

func TestServiceImpl_RandomBooks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		getService func(*gomock.Controller) *ServiceImpl
		limit      int
		want       []model.Book
		wantError  error
	}{
		{
			name: "happy flow",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBooksByRandom(gomock.Any(), 10).Return([]model.Book{{ID: 123, HashCode: 0}}, nil)

				return &ServiceImpl{rpo: rpo}
			},
			limit:     10,
			want:      []model.Book{{ID: 123, HashCode: 0}},
			wantError: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := test.getService(ctrl)

			got, err := svc.RandomBooks(context.Background(), test.limit)
			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}
