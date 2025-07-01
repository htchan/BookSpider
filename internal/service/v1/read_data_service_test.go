package service

import (
	"context"
	"database/sql"
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

func TestNewReadDataService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		rpo   repo.Repository
		confs map[string]config.SiteConfig
		want  *ReadDataServiceImpl
	}{
		{
			name: "happy flow",
			want: &ReadDataServiceImpl{},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := NewReadDataService(test.rpo, test.confs)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestReadDataReadDataServiceImpl_bookFileLocation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		serv *ReadDataServiceImpl
		bk   *model.Book
		want string
	}{
		{
			name: "return book location for book with id only",
			serv: &ReadDataServiceImpl{confs: map[string]config.SiteConfig{"test": {Storage: "/data"}}},
			bk:   &model.Book{Site: "test", ID: 123, HashCode: 0},
			want: "/data/123.txt",
		},
		{
			name: "return book location for book with id and hash code",
			serv: &ReadDataServiceImpl{confs: map[string]config.SiteConfig{"test": {Storage: "/data"}}},
			bk:   &model.Book{Site: "test", ID: 123, HashCode: 456},
			want: "/data/123-vco.txt",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.serv.bookFileLocation(test.bk)
			assert.Equal(t, test.want, got)
		})
	}

}

func TestReadDataReadDataServiceImpl_Book(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		getService func(*gomock.Controller) *ReadDataServiceImpl
		site       string
		id         string
		hash       string
		want       *model.Book
		wantError  error
	}{
		{
			name: "book found with pure ID",
			getService: func(ctrl *gomock.Controller) *ReadDataServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBookById(gomock.Any(), "test", 123).Return(&model.Book{ID: 123, HashCode: 0}, nil)

				return &ReadDataServiceImpl{rpo: rpo}
			},
			site:      "test",
			id:        "123",
			hash:      "",
			want:      &model.Book{ID: 123, HashCode: 0},
			wantError: nil,
		},
		{
			name: "book found with ID and Hashcode",
			getService: func(ctrl *gomock.Controller) *ReadDataServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBookByIdHash(gomock.Any(), "test", 123, 10).Return(&model.Book{ID: 123, HashCode: 10}, nil)

				return &ReadDataServiceImpl{rpo: rpo}
			},
			site:      "test",
			id:        "123",
			hash:      "A",
			want:      &model.Book{ID: 123, HashCode: 10},
			wantError: nil,
		},
		{
			name: "invalid id",
			getService: func(ctrl *gomock.Controller) *ReadDataServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)

				return &ReadDataServiceImpl{rpo: rpo}
			},
			site:      "test",
			id:        "abc",
			hash:      "",
			want:      nil,
			wantError: serv.ErrInvalidBookID,
		},
		{
			name: "invalid hashcode",
			getService: func(ctrl *gomock.Controller) *ReadDataServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)

				return &ReadDataServiceImpl{rpo: rpo}
			},
			site:      "test",
			id:        "123",
			hash:      "abc-def",
			want:      nil,
			wantError: serv.ErrInvalidHashCode,
		},
		{
			name: "book not found",
			getService: func(ctrl *gomock.Controller) *ReadDataServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBookById(gomock.Any(), "test", 123).Return(nil, repo.ErrBookNotExist)

				return &ReadDataServiceImpl{rpo: rpo}
			},
			site:      "test",
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

			got, err := test.getService(ctrl).Book(context.Background(), test.site, test.id, test.hash)
			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}

func TestReadDataReadDataServiceImpl_BookContent(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() {
		assert.NoError(t, os.RemoveAll("./book-content-read"))
	})

	if !assert.NoError(t, os.Mkdir("./book-content-read", os.ModePerm)) ||
		!assert.NoError(t, os.WriteFile("./book-content-read/123.txt", []byte("test"), 0644)) ||
		!assert.NoError(t, os.WriteFile("./book-content-read/123-va.txt", []byte("test v2"), 0644)) {
		return
	}

	tests := []struct {
		name      string
		serv      *ReadDataServiceImpl
		bk        *model.Book
		want      string
		wantError error
	}{
		{
			name:      "book without hashcode content exist",
			serv:      &ReadDataServiceImpl{confs: map[string]config.SiteConfig{"test": {Storage: "./book-content-read"}}},
			bk:        &model.Book{Site: "test", ID: 123, IsDownloaded: true},
			want:      "test",
			wantError: nil,
		},
		{
			name:      "book with hashcode content exist",
			serv:      &ReadDataServiceImpl{confs: map[string]config.SiteConfig{"test": {Storage: "./book-content-read"}}},
			bk:        &model.Book{Site: "test", ID: 123, HashCode: 10, IsDownloaded: true},
			want:      "test v2",
			wantError: nil,
		},
		{
			name:      "book not downloaded",
			serv:      &ReadDataServiceImpl{confs: map[string]config.SiteConfig{"test": {Storage: "./book-content-read"}}},
			bk:        &model.Book{Site: "test", ID: 123, HashCode: 10, IsDownloaded: false},
			want:      "",
			wantError: serv.ErrBookNotDownload,
		},
		{
			name:      "book content not exist",
			serv:      &ReadDataServiceImpl{confs: map[string]config.SiteConfig{"test": {Storage: "./book-content-read"}}},
			bk:        &model.Book{Site: "test", ID: 456, IsDownloaded: true},
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

func TestReadDataReadDataServiceImpl_BookChapters(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() {
		assert.NoError(t, os.RemoveAll("./book-chapters-read"))
	})

	if !assert.NoError(t, os.Mkdir("./book-chapters-read", os.ModePerm)) ||
		!assert.NoError(t, os.WriteFile("./book-chapters-read/123.txt", []byte(
			"test"+model.CONTENT_SEP+
				"title 1"+model.CONTENT_SEP+"content 1"+model.CONTENT_SEP+
				"title 2"+model.CONTENT_SEP+"content 2"+model.CONTENT_SEP,
		), 0644)) ||
		!assert.NoError(t, os.WriteFile("./book-chapters-read/123-va.txt", []byte("test"), 0644)) {
		return
	}

	tests := []struct {
		name      string
		serv      *ReadDataServiceImpl
		bk        *model.Book
		want      model.Chapters
		wantError error
	}{
		{
			name: "successfully read and parse content to chapters",
			serv: &ReadDataServiceImpl{confs: map[string]config.SiteConfig{"test": {Storage: "./book-chapters-read"}}},
			bk:   &model.Book{Site: "test", ID: 123, IsDownloaded: true},
			want: model.Chapters{
				{Index: 0, Title: "title 1", Content: "content 1"},
				{Index: 1, Title: "title 2", Content: "content 2"},
			},
			wantError: nil,
		},
		{
			name:      "fail to parse content to chapters",
			serv:      &ReadDataServiceImpl{confs: map[string]config.SiteConfig{"test": {Storage: "./book-chapters-read"}}},
			bk:        &model.Book{Site: "test", ID: 123, HashCode: 10, IsDownloaded: true},
			want:      nil,
			wantError: model.ErrCannotParseContent,
		},
		{
			name:      "fail to read content ",
			serv:      &ReadDataServiceImpl{confs: map[string]config.SiteConfig{"test": {Storage: "./book-chapters-read"}}},
			bk:        &model.Book{Site: "test", ID: 456, HashCode: 0, IsDownloaded: true},
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

func TestReadDataReadDataServiceImpl_BookGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		getService    func(*gomock.Controller) *ReadDataServiceImpl
		site          string
		id            string
		hash          string
		wantBook      *model.Book
		wantBookGroup *model.BookGroup
		wantError     error
	}{
		{
			name: "book group found with pure ID",
			getService: func(ctrl *gomock.Controller) *ReadDataServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBookGroupByID(gomock.Any(), "test", 123).Return(
					model.BookGroup{{Site: "test", ID: 123, HashCode: 0}, {Site: "test", ID: 456, HashCode: 0}},
					nil,
				)

				return &ReadDataServiceImpl{rpo: rpo}
			},
			site:          "test",
			id:            "123",
			hash:          "",
			wantBook:      &model.Book{Site: "test", ID: 123, HashCode: 0},
			wantBookGroup: &model.BookGroup{{Site: "test", ID: 456, HashCode: 0}},
			wantError:     nil,
		},
		{
			name: "book group found with ID and Hashcode",
			getService: func(ctrl *gomock.Controller) *ReadDataServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBookGroupByIDHash(gomock.Any(), "test", 123, 10).Return(
					model.BookGroup{{Site: "test", ID: 123, HashCode: 10}, {Site: "test", ID: 456, HashCode: 0}},
					nil,
				)

				return &ReadDataServiceImpl{rpo: rpo}
			},
			site:          "test",
			id:            "123",
			hash:          "A",
			wantBook:      &model.Book{Site: "test", ID: 123, HashCode: 10},
			wantBookGroup: &model.BookGroup{{Site: "test", ID: 456, HashCode: 0}},
			wantError:     nil,
		},
		{
			name: "invalid id",
			getService: func(ctrl *gomock.Controller) *ReadDataServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)

				return &ReadDataServiceImpl{rpo: rpo}
			},
			site:          "test",
			id:            "abc",
			hash:          "",
			wantBook:      nil,
			wantBookGroup: nil,
			wantError:     serv.ErrInvalidBookID,
		},
		{
			name: "invalid hashcode",
			getService: func(ctrl *gomock.Controller) *ReadDataServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)

				return &ReadDataServiceImpl{rpo: rpo}
			},
			site:          "test",
			id:            "123",
			hash:          "abc-def",
			wantBook:      nil,
			wantBookGroup: nil,
			wantError:     serv.ErrInvalidHashCode,
		},
		{
			name: "book not found",
			getService: func(ctrl *gomock.Controller) *ReadDataServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBookGroupByID(gomock.Any(), "test", 123).Return(nil, repo.ErrBookNotExist)

				return &ReadDataServiceImpl{rpo: rpo}
			},
			site:          "test",
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

			gotBk, gotBkGroup, err := test.getService(ctrl).BookGroup(context.Background(), test.site, test.id, test.hash)
			assert.Equal(t, test.wantBook, gotBk)
			assert.Equal(t, test.wantBookGroup, gotBkGroup)
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}

func TestReadDataReadDataServiceImpl_SearchBook(t *testing.T) {

	t.Parallel()

	tests := []struct {
		name       string
		getService func(*gomock.Controller) *ReadDataServiceImpl
		title      string
		writer     string
		limit      int
		offset     int
		want       []model.Book
		wantError  error
	}{
		{
			name: "happy flow with books",
			getService: func(ctrl *gomock.Controller) *ReadDataServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBooksByTitleWriter(gomock.Any(), "title", "writer", 10, 0).
					Return([]model.Book{{ID: 123, HashCode: 0}}, nil)

				return &ReadDataServiceImpl{rpo: rpo}
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

			got, err := svc.SearchBooks(
				context.Background(), test.title, test.writer, test.limit, test.offset,
			)
			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}

func TestReadDataReadDataServiceImpl_RandomBook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		getService func(*gomock.Controller) *ReadDataServiceImpl
		limit      int
		want       []model.Book
		wantError  error
	}{
		{
			name: "happy flow",
			getService: func(ctrl *gomock.Controller) *ReadDataServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBooksByRandom(gomock.Any(), 10).Return([]model.Book{{ID: 123, HashCode: 0}}, nil)

				return &ReadDataServiceImpl{rpo: rpo}
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

func TestReadDataReadDataServiceImpl_Stats(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		getService func(*gomock.Controller) *ReadDataServiceImpl
		want       repo.Summary
	}{
		{
			name: "happy flow",
			getService: func(ctrl *gomock.Controller) *ReadDataServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().Stats(gomock.Any(), "test").Return(repo.Summary{})

				return &ReadDataServiceImpl{rpo: rpo}
			},
			want: repo.Summary{},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := test.getService(ctrl)

			got := svc.Stats(t.Context(), "test")
			assert.Equal(t, test.want, got)
		})
	}

}

func TestReadDataReadDataServiceImpl_DBStats(t *testing.T) {

	t.Parallel()

	tests := []struct {
		name       string
		getService func(*gomock.Controller) *ReadDataServiceImpl
		want       sql.DBStats
	}{
		{
			name: "happy flow",
			getService: func(ctrl *gomock.Controller) *ReadDataServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().DBStats(gomock.Any()).Return(sql.DBStats{})

				return &ReadDataServiceImpl{rpo: rpo}
			},
			want: sql.DBStats{},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := test.getService(ctrl)

			got := svc.DBStats(t.Context())
			assert.Equal(t, test.want, got)
		})
	}
}
