package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/htchan/BookSpider/internal/config/v2"
	mockclient "github.com/htchan/BookSpider/internal/mock/client/v2"
	mockparser "github.com/htchan/BookSpider/internal/mock/parser"
	mockrepo "github.com/htchan/BookSpider/internal/mock/repo"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/parse"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
)

func TestServiceImp_baseURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		serv ServiceImp
		bk   *model.Book
		want string
	}{
		{
			name: "happy flow",
			serv: ServiceImp{conf: config.SiteConfig{
				URL: config.URLConfig{Base: "http://test.com/books/%v"},
			}},
			bk:   &model.Book{Site: "test", ID: 1, HashCode: 1},
			want: "http://test.com/books/1",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := test.serv.baseURL(test.bk)

			assert.Equal(t, test.want, got)
		})
	}
}

func TestServiceImp_UpdateBook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupServ    func(ctrl *gomock.Controller) ServiceImp
		bk           *model.Book
		wantBook     *model.Book
		wantError    bool
		wantErrorStr string
	}{
		{
			name: "happy flow without any changes",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mockrepo.NewMockRepository(ctrl)

				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/books/1").Return("content", nil)

				p := mockparser.NewMockParser(ctrl)
				p.EXPECT().ParseBook("content").Return(parse.NewParsedBookFields(
					"title", "writer", "type", "date", "chapter",
				), nil)

				return ServiceImp{
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(1),
					rpo:    rpo,
					client: c,
					parser: p,
					conf: config.SiteConfig{URL: config.URLConfig{
						Base: "http://test.com/books/%v",
					}},
				}
			},
			bk: &model.Book{
				ID: 1, Status: model.StatusInProgress,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
			},
			wantBook: &model.Book{
				ID: 1, Status: model.StatusInProgress,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
			},
			wantError:    false,
			wantErrorStr: "",
		},
		{
			name: "happy flow for new book",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().CreateBook(&model.Book{
					ID: 1, HashCode: model.GenerateHash(), Status: model.StatusInProgress,
					Title: "title new", Writer: model.Writer{ID: 10, Name: "writer new"}, Type: "type new",
					UpdateDate: "date new", UpdateChapter: "chapter new",
				}).Return(nil)
				rpo.EXPECT().SaveWriter(&model.Writer{ID: 10, Name: "writer new"}).Return(nil)

				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/books/1").Return("content", nil)

				p := mockparser.NewMockParser(ctrl)
				p.EXPECT().ParseBook("content").Return(parse.NewParsedBookFields(
					"title new", "writer new", "type new", "date new", "chapter new",
				), nil)

				return ServiceImp{
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(1),
					rpo:    rpo,
					client: c,
					parser: p,
					conf: config.SiteConfig{URL: config.URLConfig{
						Base: "http://test.com/books/%v",
					}},
				}
			},
			bk: &model.Book{
				ID: 1, Status: model.StatusEnd,
				Title: "title", Writer: model.Writer{ID: 10, Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
			},
			wantBook: &model.Book{
				ID: 1, HashCode: model.GenerateHash(), Status: model.StatusInProgress,
				Title: "title new", Writer: model.Writer{ID: 10, Name: "writer new"}, Type: "type new",
				UpdateDate: "date new", UpdateChapter: "chapter new",
			},
			wantError:    false,
			wantErrorStr: "",
		},
		{
			name: "happy flow for updated existing book",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().UpdateBook(&model.Book{
					ID: 1, Status: model.StatusInProgress,
					Title: "title", Writer: model.Writer{ID: 10, Name: "writer"}, Type: "type",
					UpdateDate: "date new", UpdateChapter: "chapter new",
				}).Return(nil)
				rpo.EXPECT().SaveWriter(&model.Writer{ID: 10, Name: "writer"}).Return(nil)

				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/books/1").Return("content", nil)

				p := mockparser.NewMockParser(ctrl)
				p.EXPECT().ParseBook("content").Return(parse.NewParsedBookFields(
					"title", "writer", "type", "date new", "chapter new",
				), nil)

				return ServiceImp{
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(1),
					rpo:    rpo,
					client: c,
					parser: p,
					conf: config.SiteConfig{URL: config.URLConfig{
						Base: "http://test.com/books/%v",
					}},
				}
			},
			bk: &model.Book{
				ID: 1, Status: model.StatusEnd,
				Title: "title", Writer: model.Writer{ID: 10, Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
			},
			wantBook: &model.Book{
				ID: 1, Status: model.StatusInProgress,
				Title: "title", Writer: model.Writer{ID: 10, Name: "writer"}, Type: "type",
				UpdateDate: "date new", UpdateChapter: "chapter new",
			},
			wantError:    false,
			wantErrorStr: "",
		},
		{
			name: "get website fail",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mockrepo.NewMockRepository(ctrl)

				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/books/1").Return("", errors.New("get web error"))

				p := mockparser.NewMockParser(ctrl)

				return ServiceImp{
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(1),
					rpo:    rpo,
					client: c,
					parser: p,
					conf: config.SiteConfig{URL: config.URLConfig{
						Base: "http://test.com/books/%v",
					}},
				}
			},
			bk: &model.Book{
				ID: 1, Status: model.StatusInProgress,
				Title: "title", Writer: model.Writer{ID: 10, Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
			},
			wantBook: &model.Book{
				ID: 1, Status: model.StatusInProgress,
				Title: "title", Writer: model.Writer{ID: 10, Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
			},
			wantError:    true,
			wantErrorStr: "fetch html fail: get web error",
		},
		{
			name: "parse content fail",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mockrepo.NewMockRepository(ctrl)

				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/books/1").Return("content", nil)

				p := mockparser.NewMockParser(ctrl)
				p.EXPECT().ParseBook("content").Return(nil, parse.ErrParseBookFieldsNotFound)

				return ServiceImp{
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(1),
					rpo:    rpo,
					client: c,
					parser: p,
					conf: config.SiteConfig{URL: config.URLConfig{
						Base: "http://test.com/books/%v",
					}},
				}
			},
			bk: &model.Book{
				ID: 1, Status: model.StatusInProgress,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
			},
			wantBook: &model.Book{
				ID: 1, Status: model.StatusInProgress,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
			},
			wantError:    true,
			wantErrorStr: "parse html fail: parse book fail: fields not found",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := test.setupServ(ctrl)
			err := serv.UpdateBook(test.bk)

			if test.wantError {
				assert.Error(t, err)
			}

			if err != nil {
				assert.EqualError(t, err, test.wantErrorStr)
			}

			assert.Equal(t, test.wantBook, test.bk)
		})
	}
}

func TestServiceImp_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupServ    func(ctrl *gomock.Controller) ServiceImp
		wantError    bool
		wantErrorStr string
	}{
		{
			name: "happy flow that some book update success and some fail",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mockrepo.NewMockRepository(ctrl)

				bookChan := make(chan model.Book, 2)
				bookChan <- model.Book{
					ID: 1, Status: model.StatusInProgress, Type: "type",
					Title: "title", Writer: model.Writer{Name: "writer"},
					UpdateDate: "date", UpdateChapter: "chapter",
				}
				bookChan <- model.Book{
					ID: 2, Status: model.StatusInProgress, Type: "type",
					Title: "title", Writer: model.Writer{Name: "writer"},
					UpdateDate: "date", UpdateChapter: "chapter",
				}
				close(bookChan)
				rpo.EXPECT().FindBooksForUpdate().Return(bookChan, nil)
				rpo.EXPECT().UpdateBook(&model.Book{
					ID: 1, Status: model.StatusInProgress, Type: "type",
					Title: "title", Writer: model.Writer{Name: "writer"},
					UpdateDate: "date new", UpdateChapter: "chapter new",
				}).Return(nil)
				rpo.EXPECT().SaveWriter(&model.Writer{Name: "writer"}).Return(nil)

				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/books/1").Return("content 1", nil)
				c.EXPECT().Get(gomock.Any(), "http://test.com/books/2").Return("", errors.New("get web fail"))

				p := mockparser.NewMockParser(ctrl)
				p.EXPECT().ParseBook("content 1").Return(parse.NewParsedBookFields(
					"title", "writer", "type", "date new", "chapter new",
				), nil)

				return ServiceImp{
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(1),
					rpo:    rpo,
					client: c,
					parser: p,
					conf: config.SiteConfig{URL: config.URLConfig{
						Base: "http://test.com/books/%v",
					}},
				}
			},
			wantError:    false,
			wantErrorStr: "",
		},
		{
			name: "load book from DB got error",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBooksForUpdate().Return(nil, errors.New("some error"))
				return ServiceImp{rpo: rpo}
			},
			wantError:    true,
			wantErrorStr: "fail to load books from DB: some error",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := test.setupServ(ctrl)
			err := serv.Update()

			if test.wantError {
				assert.Error(t, err)
			}

			if err != nil {
				assert.EqualError(t, err, test.wantErrorStr)
			}
		})
	}
}
