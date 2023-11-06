package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/htchan/BookSpider/internal/config/v2"
	mockclient "github.com/htchan/BookSpider/internal/mock/client/v2"
	mockparser "github.com/htchan/BookSpider/internal/mock/parser"
	mockrepo "github.com/htchan/BookSpider/internal/mock/repo"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/parse"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
)

func TestServiceImp_ExploreBook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setupServ      func(ctrl *gomock.Controller) ServiceImp
		bk             *model.Book
		expectBk       *model.Book
		expectedError  bool
		expectErrorStr string
	}{
		{
			name: "happy flow with new book (error == nil)",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/book/1").Return("basic book info", nil)

				p := mockparser.NewMockParser(ctrl)
				p.EXPECT().ParseBook("basic book info").Return(parse.NewParsedBookFields(
					"title",
					"writer",
					"type",
					"date",
					"chapter",
				), nil)

				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().CreateBook(&model.Book{Site: "test-explore-book", ID: 1})
				rpo.EXPECT().UpdateBook(&model.Book{
					Site:          "test-explore-book",
					ID:            1,
					Title:         "title",
					Writer:        model.Writer{Name: "writer"},
					Type:          "type",
					Status:        model.StatusInProgress,
					UpdateDate:    "date",
					UpdateChapter: "chapter",
				})
				rpo.EXPECT().SaveWriter(&model.Writer{Name: "writer"})
				rpo.EXPECT().SaveError(&model.Book{
					Site:          "test-explore-book",
					ID:            1,
					Title:         "title",
					Writer:        model.Writer{Name: "writer"},
					Type:          "type",
					Status:        model.StatusInProgress,
					UpdateDate:    "date",
					UpdateChapter: "chapter",
				}, nil)

				return ServiceImp{
					client: c,
					parser: p,
					rpo:    rpo,
					conf: config.SiteConfig{
						URL: config.URLConfig{Base: "http://test.com/book/%v"},
					},
				}
			},
			bk: &model.Book{Site: "test-explore-book", ID: 1},
			expectBk: &model.Book{
				Site: "test-explore-book", ID: 1, Status: model.StatusInProgress,
				Title: "title", Writer: model.Writer{Name: "writer"},
				Type: "type", UpdateDate: "date", UpdateChapter: "chapter",
			},
			expectedError:  false,
			expectErrorStr: "",
		},
		{
			name: "happy flow with existing book (error != nil)",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/book/2").Return("basic book info", nil)

				p := mockparser.NewMockParser(ctrl)
				p.EXPECT().ParseBook("basic book info").Return(parse.NewParsedBookFields(
					"title",
					"writer",
					"type",
					"date",
					"chapter",
				), nil)

				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().UpdateBook(&model.Book{
					Site:          "test-explore-book",
					ID:            2,
					Title:         "title",
					Writer:        model.Writer{Name: "writer"},
					Type:          "type",
					Status:        model.StatusInProgress,
					UpdateDate:    "date",
					UpdateChapter: "chapter",
				})
				rpo.EXPECT().SaveWriter(&model.Writer{Name: "writer"})
				rpo.EXPECT().SaveError(&model.Book{
					Site:          "test-explore-book",
					ID:            2,
					Title:         "title",
					Writer:        model.Writer{Name: "writer"},
					Type:          "type",
					Status:        model.StatusInProgress,
					UpdateDate:    "date",
					UpdateChapter: "chapter",
				}, nil)

				return ServiceImp{
					client: c,
					parser: p,
					rpo:    rpo,
					conf: config.SiteConfig{
						URL: config.URLConfig{Base: "http://test.com/book/%v"},
					},
				}
			},
			bk: &model.Book{Site: "test-explore-book", ID: 2, Error: errors.New("existing error")},
			expectBk: &model.Book{
				Site: "test-explore-book", ID: 2, Status: model.StatusInProgress,
				Title: "title", Writer: model.Writer{Name: "writer"},
				Type: "type", UpdateDate: "date", UpdateChapter: "chapter",
			},
			expectedError:  false,
			expectErrorStr: "",
		},
		{
			name: "fail if book status is not error",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				return ServiceImp{}
			},
			bk: &model.Book{Site: "test-explore-book", ID: 3, Status: model.StatusInProgress},
			expectBk: &model.Book{
				Site: "test-explore-book", ID: 3, Status: model.StatusInProgress},
			expectedError:  true,
			expectErrorStr: "book status is not error",
		},
		{
			name: "parse book fail",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/book/4").Return("basic book info", nil)

				p := mockparser.NewMockParser(ctrl)
				p.EXPECT().ParseBook("basic book info").Return(nil, parse.ErrParseBookFieldsNotFound)

				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().SaveError(&model.Book{
					Site:  "test-explore-book",
					ID:    4,
					Error: fmt.Errorf("parse html fail: %w", parse.ErrParseBookFieldsNotFound),
				}, fmt.Errorf("parse html fail: %w", parse.ErrParseBookFieldsNotFound))

				return ServiceImp{
					client: c,
					parser: p,
					rpo:    rpo,
					conf: config.SiteConfig{
						URL: config.URLConfig{Base: "http://test.com/book/%v"},
					},
				}
			},
			bk: &model.Book{Site: "test-explore-book", ID: 4, Error: errors.New("existing error")},
			expectBk: &model.Book{
				Site: "test-explore-book", ID: 4, Error: fmt.Errorf("parse html fail: %w", parse.ErrParseBookFieldsNotFound),
			},
			expectedError:  true,
			expectErrorStr: "explore book fail: parse html fail: parse book fail: fields not found",
		},
		{
			name: "update book fail",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/book/5").Return("basic book info", nil)

				p := mockparser.NewMockParser(ctrl)
				p.EXPECT().ParseBook("basic book info").Return(parse.NewParsedBookFields(
					"title",
					"writer",
					"type",
					"date",
					"chapter",
				), nil)

				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().SaveWriter(&model.Writer{Name: "writer"})
				rpo.EXPECT().UpdateBook(&model.Book{
					Site:          "test-explore-book",
					ID:            5,
					Title:         "title",
					Writer:        model.Writer{Name: "writer"},
					Type:          "type",
					Status:        model.StatusInProgress,
					UpdateDate:    "date",
					UpdateChapter: "chapter",
				}).Return(errors.New("update book error"))
				rpo.EXPECT().SaveError(&model.Book{
					Site:          "test-explore-book",
					ID:            5,
					Title:         "title",
					Writer:        model.Writer{Name: "writer"},
					Type:          "type",
					Status:        model.StatusInProgress,
					UpdateDate:    "date",
					UpdateChapter: "chapter",
					Error:         errors.New("update book error"),
				}, errors.New("update book error"))

				return ServiceImp{
					client: c,
					parser: p,
					rpo:    rpo,
					conf: config.SiteConfig{
						URL: config.URLConfig{Base: "http://test.com/book/%v"},
					},
				}
			},
			bk: &model.Book{Site: "test-explore-book", ID: 5, Error: errors.New("existing error")},
			expectBk: &model.Book{
				Site: "test-explore-book", ID: 5, Status: model.StatusInProgress,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
				Error: errors.New("update book error"),
			},
			expectedError:  true,
			expectErrorStr: "explore book fail: update book error",
		},
		{
			name: "update book fail and save error fail",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)
				c.EXPECT().Get(gomock.Any(), "http://test.com/book/5").Return("basic book info", nil)

				p := mockparser.NewMockParser(ctrl)
				p.EXPECT().ParseBook("basic book info").Return(parse.NewParsedBookFields(
					"title",
					"writer",
					"type",
					"date",
					"chapter",
				), nil)

				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().SaveWriter(&model.Writer{Name: "writer"})
				rpo.EXPECT().UpdateBook(&model.Book{
					Site:          "test-explore-book",
					ID:            5,
					Title:         "title",
					Writer:        model.Writer{Name: "writer"},
					Type:          "type",
					Status:        model.StatusInProgress,
					UpdateDate:    "date",
					UpdateChapter: "chapter",
				}).Return(errors.New("update book error"))
				rpo.EXPECT().SaveError(&model.Book{
					Site:          "test-explore-book",
					ID:            5,
					Title:         "title",
					Writer:        model.Writer{Name: "writer"},
					Type:          "type",
					Status:        model.StatusInProgress,
					UpdateDate:    "date",
					UpdateChapter: "chapter",
					Error:         errors.New("update book error"),
				}, errors.New("update book error")).Return(errors.New("fake fail"))

				return ServiceImp{
					client: c,
					parser: p,
					rpo:    rpo,
					conf: config.SiteConfig{
						URL: config.URLConfig{Base: "http://test.com/book/%v"},
					},
				}
			},
			bk: &model.Book{Site: "test-explore-book", ID: 5, Error: errors.New("existing error")},
			expectBk: &model.Book{
				Site: "test-explore-book", ID: 5, Status: model.StatusInProgress,
				Title: "title", Writer: model.Writer{Name: "writer"}, Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter",
				Error: errors.New("update book error"),
			},
			expectedError:  true,
			expectErrorStr: "save error fail: fake fail",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := test.setupServ(ctrl)
			err := serv.ExploreBook(test.bk)

			if test.expectedError {
				assert.Error(t, err)
			}

			if err != nil {
				assert.EqualError(t, err, test.expectErrorStr)
			}
			assert.Equal(t, test.bk, test.expectBk)
		})
	}
}

func TestServiceImp_exploreExisting(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		setupServ        func(ctrl *gomock.Controller) ServiceImp
		summary          repo.Summary
		errorCount       int
		expectErrorCount int
	}{
		{
			name: "happy flow",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)

				p := mockparser.NewMockParser(ctrl)

				rpo := mockrepo.NewMockRepository(ctrl)
				for i := 101; i <= 105; i++ {
					rpo.EXPECT().FindBookById(i).Return(&model.Book{
						Site:  "test-explore-existing",
						ID:    i,
						Error: errors.New("not found"),
					}, nil)

					c.EXPECT().Get(gomock.Any(), fmt.Sprintf("https://test.com/book/%v", i)).Return(fmt.Sprintf("content %v", i), nil)

					p.EXPECT().ParseBook(fmt.Sprintf("content %v", i)).Return(parse.NewParsedBookFields(
						"title", "writer", "type", "date", "chapter",
					), nil)

					bk := &model.Book{
						Site:          "test-explore-existing",
						ID:            i,
						Title:         "title",
						Writer:        model.Writer{Name: "writer"},
						Type:          "type",
						Status:        model.StatusInProgress,
						UpdateDate:    "date",
						UpdateChapter: "chapter",
					}

					rpo.EXPECT().UpdateBook(bk).Return(nil)
					rpo.EXPECT().SaveWriter(&bk.Writer).Return(nil)
					rpo.EXPECT().SaveError(bk, nil)
				}

				return ServiceImp{
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(1),
					client: c,
					parser: p,
					rpo:    rpo,
					conf: config.SiteConfig{
						URL:             config.URLConfig{Base: "https://test.com/book/%v"},
						MaxExploreError: 1,
					},
				}
			},
			summary: repo.Summary{
				LatestSuccessID: 100,
				MaxBookID:       105,
			},
			errorCount:       0,
			expectErrorCount: 0,
		},
		{
			name: "update error count if fail to load book",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)

				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBookById(201).Return(nil, errors.New("db error"))

				return ServiceImp{
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(1),
					client: c,
					rpo:    rpo,
					conf: config.SiteConfig{
						URL:             config.URLConfig{Base: "https://test.com/book/%v"},
						MaxExploreError: 5,
					},
				}
			},
			summary: repo.Summary{
				LatestSuccessID: 200,
				MaxBookID:       201,
			},
			errorCount:       0,
			expectErrorCount: 1,
		},
		{
			name: "update error count if fail to parse book",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)

				p := mockparser.NewMockParser(ctrl)

				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBookById(301).Return(&model.Book{
					Site:  "test-explore-existing",
					ID:    301,
					Error: errors.New("not found"),
				}, nil)

				c.EXPECT().Get(gomock.Any(), "https://test.com/book/301").Return("content", nil)

				p.EXPECT().ParseBook("content").Return(nil, parse.ErrParseBookFieldsNotFound)

				err := fmt.Errorf("parse html fail: %w", parse.ErrParseBookFieldsNotFound)
				rpo.EXPECT().SaveError(&model.Book{
					Site:   "test-explore-existing",
					ID:     301,
					Status: model.StatusError,
					Error:  err,
				}, err).Return(nil)

				return ServiceImp{
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(1),
					client: c,
					parser: p,
					rpo:    rpo,
					conf: config.SiteConfig{
						URL:             config.URLConfig{Base: "https://test.com/book/%v"},
						MaxExploreError: 1,
					},
				}
			},
			summary: repo.Summary{
				LatestSuccessID: 300,
				MaxBookID:       301,
			},
			errorCount:       0,
			expectErrorCount: 1,
		},
		{
			name: "clear error count if explore book success",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)

				p := mockparser.NewMockParser(ctrl)

				rpo := mockrepo.NewMockRepository(ctrl)

				rpo.EXPECT().FindBookById(401).Return(&model.Book{
					Site:  "test-explore-existing",
					ID:    401,
					Error: errors.New("not found"),
				}, nil)

				c.EXPECT().Get(gomock.Any(), "https://test.com/book/401").Return("content", nil)

				p.EXPECT().ParseBook("content").Return(nil, parse.ErrParseBookFieldsNotFound)

				err := fmt.Errorf("parse html fail: %w", parse.ErrParseBookFieldsNotFound)
				rpo.EXPECT().SaveError(&model.Book{
					Site:   "test-explore-existing",
					ID:     401,
					Status: model.StatusError,
					Error:  err,
				}, err).Return(nil)

				rpo.EXPECT().FindBookById(402).Return(&model.Book{
					Site:  "test-explore-existing",
					ID:    402,
					Error: errors.New("not found"),
				}, nil)

				c.EXPECT().Get(gomock.Any(), "https://test.com/book/402").Return("content 402", nil)

				p.EXPECT().ParseBook("content 402").Return(parse.NewParsedBookFields(
					"title", "writer", "type", "date", "chapter",
				), nil)

				bk := &model.Book{
					Site:          "test-explore-existing",
					ID:            402,
					Title:         "title",
					Writer:        model.Writer{Name: "writer"},
					Type:          "type",
					Status:        model.StatusInProgress,
					UpdateDate:    "date",
					UpdateChapter: "chapter",
				}

				rpo.EXPECT().UpdateBook(bk).Return(nil)
				rpo.EXPECT().SaveWriter(&bk.Writer).Return(nil)
				rpo.EXPECT().SaveError(bk, nil)

				return ServiceImp{
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(1),
					client: c,
					parser: p,
					rpo:    rpo,
					conf: config.SiteConfig{
						URL:             config.URLConfig{Base: "https://test.com/book/%v"},
						MaxExploreError: 1,
					},
				}
			},
			summary: repo.Summary{
				LatestSuccessID: 400,
				MaxBookID:       402,
			},
			errorCount:       0,
			expectErrorCount: 0,
		},
		{
			name: "stop before reaching Max Book ID if it reach the Max fail count",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)

				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBookById(201).Return(nil, errors.New("db error"))
				rpo.EXPECT().FindBookById(202).Return(nil, errors.New("db error"))

				return ServiceImp{
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(1),
					client: c,
					rpo:    rpo,
					conf: config.SiteConfig{
						URL:             config.URLConfig{Base: "https://test.com/book/%v"},
						MaxExploreError: 1,
					},
				}
			},
			summary: repo.Summary{
				LatestSuccessID: 200,
				MaxBookID:       205,
			},
			errorCount:       0,
			expectErrorCount: 2,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := test.setupServ(ctrl)
			serv.exploreExisting(test.summary, &test.errorCount)

			assert.Equal(t, test.expectErrorCount, test.errorCount)
		})
	}
}

func TestServiceImp_exploreNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		setupServ        func(ctrl *gomock.Controller) ServiceImp
		summary          repo.Summary
		errorCount       int
		expectErrorCount int
	}{
		{
			name: "happy flow",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)

				p := mockparser.NewMockParser(ctrl)

				rpo := mockrepo.NewMockRepository(ctrl)

				err := fmt.Errorf("parse html fail: %w", parse.ErrParseBookFieldsNotFound)

				for i := 301; i < 303; i++ {
					c.EXPECT().Get(gomock.Any(), fmt.Sprintf("https://test.com/book/%v", i)).Return(fmt.Sprintf("content %v", i), nil)

					p.EXPECT().ParseBook(fmt.Sprintf("content %v", i)).Return(nil, parse.ErrParseBookFieldsNotFound)

					rpo.EXPECT().CreateBook(&model.Book{
						Site:     "test-explore-new",
						ID:       i,
						HashCode: model.GenerateHash(),
					}).Return(nil)
					rpo.EXPECT().SaveError(&model.Book{Site: "test-explore-new", ID: i, HashCode: model.GenerateHash(), Error: err}, err).Return(nil)
				}

				return ServiceImp{
					name:   "test-explore-new",
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(1),
					client: c,
					parser: p,
					rpo:    rpo,
					conf: config.SiteConfig{
						URL:             config.URLConfig{Base: "https://test.com/book/%v"},
						MaxExploreError: 1,
					},
				}
			},
			summary: repo.Summary{
				MaxBookID: 300,
			},
			errorCount:       0,
			expectErrorCount: 2,
		},
		{
			name: "clear error if explore book success",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)

				p := mockparser.NewMockParser(ctrl)

				rpo := mockrepo.NewMockRepository(ctrl)

				err := fmt.Errorf("parse html fail: %w", parse.ErrParseBookFieldsNotFound)

				c.EXPECT().Get(gomock.Any(), "https://test.com/book/301").Return("content 301", nil)

				p.EXPECT().ParseBook("content 301").Return(nil, parse.ErrParseBookFieldsNotFound)

				rpo.EXPECT().CreateBook(&model.Book{
					Site:     "test-explore-new",
					ID:       301,
					HashCode: model.GenerateHash(),
				}).Return(nil)
				rpo.EXPECT().SaveError(&model.Book{Site: "test-explore-new", ID: 301, HashCode: model.GenerateHash(), Error: err}, err).Return(nil)

				c.EXPECT().Get(gomock.Any(), "https://test.com/book/302").Return("content 302", nil)

				p.EXPECT().ParseBook("content 302").Return(parse.NewParsedBookFields(
					"title", "writer", "type", "date", "chapter",
				), nil)

				rpo.EXPECT().CreateBook(&model.Book{
					Site:     "test-explore-new",
					ID:       302,
					HashCode: model.GenerateHash(),
				}).Return(nil)
				bk := &model.Book{
					Site:          "test-explore-new",
					ID:            302,
					HashCode:      model.GenerateHash(),
					Title:         "title",
					Writer:        model.Writer{Name: "writer"},
					Status:        model.StatusInProgress,
					Type:          "type",
					UpdateDate:    "date",
					UpdateChapter: "chapter",
				}
				rpo.EXPECT().UpdateBook(bk).Return(nil)
				rpo.EXPECT().SaveWriter(&bk.Writer).Return(nil)
				rpo.EXPECT().SaveError(bk, nil).Return(nil)

				for i := 303; i <= 305; i++ {
					c.EXPECT().Get(gomock.Any(), fmt.Sprintf("https://test.com/book/%v", i)).Return(fmt.Sprintf("content %v", i), nil)

					p.EXPECT().ParseBook(fmt.Sprintf("content %v", i)).Return(nil, parse.ErrParseBookFieldsNotFound)

					rpo.EXPECT().CreateBook(&model.Book{
						Site:     "test-explore-new",
						ID:       i,
						HashCode: model.GenerateHash(),
					}).Return(nil)
					rpo.EXPECT().SaveError(&model.Book{Site: "test-explore-new", ID: i, HashCode: model.GenerateHash(), Error: err}, err).Return(nil)
				}

				return ServiceImp{
					name:   "test-explore-new",
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(1),
					client: c,
					parser: p,
					rpo:    rpo,
					conf: config.SiteConfig{
						URL:             config.URLConfig{Base: "https://test.com/book/%v"},
						MaxExploreError: 2,
					},
				}
			},
			summary: repo.Summary{
				MaxBookID: 300,
			},
			errorCount:       0,
			expectErrorCount: 3,
		},
		{
			name: "do nothing if it already reach max explore error at the beginning",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {

				return ServiceImp{
					name: "test-explore-new",
					conf: config.SiteConfig{
						MaxExploreError: 50,
					},
				}
			},
			summary: repo.Summary{
				MaxBookID: 300,
			},
			errorCount:       100,
			expectErrorCount: 100,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := test.setupServ(ctrl)
			serv.exploreNew(test.summary, &test.errorCount)

			assert.Equal(t, test.expectErrorCount, test.errorCount)
		})
	}
}

func TestServiceImp_Explore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setupServ      func(ctrl *gomock.Controller) ServiceImp
		expectedError  bool
		expectErrorStr string
	}{
		{
			name: "happy flow",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				c := mockclient.NewMockBookClient(ctrl)

				p := mockparser.NewMockParser(ctrl)

				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().Stats().Return(repo.Summary{
					LatestSuccessID: 400,
					MaxBookID:       402,
				})

				rpo.EXPECT().FindBookById(401).Return(&model.Book{
					Site:  "test-explore",
					ID:    401,
					Error: errors.New("not found"),
				}, nil)

				c.EXPECT().Get(gomock.Any(), "https://test.com/book/401").Return("content 401", nil)

				p.EXPECT().ParseBook("content 401").Return(parse.NewParsedBookFields(
					"title", "writer", "type", "date", "chapter",
				), nil)

				bk := &model.Book{
					Site:          "test-explore",
					ID:            401,
					Title:         "title",
					Writer:        model.Writer{Name: "writer"},
					Type:          "type",
					Status:        model.StatusInProgress,
					UpdateDate:    "date",
					UpdateChapter: "chapter",
				}

				rpo.EXPECT().UpdateBook(bk).Return(nil)
				rpo.EXPECT().SaveWriter(&bk.Writer).Return(nil)
				rpo.EXPECT().SaveError(bk, nil)

				err := fmt.Errorf("parse html fail: %w", parse.ErrParseBookFieldsNotFound)

				rpo.EXPECT().FindBookById(402).Return(&model.Book{
					Site:  "test-explore",
					ID:    402,
					Error: errors.New("not found"),
				}, nil)

				c.EXPECT().Get(gomock.Any(), "https://test.com/book/402").Return("content 402", nil)

				p.EXPECT().ParseBook("content 402").Return(nil, parse.ErrParseBookFieldsNotFound)

				rpo.EXPECT().SaveError(&model.Book{Site: "test-explore", ID: 402, Error: err}, err).Return(nil)

				for i := 403; i < 405; i++ {
					c.EXPECT().Get(gomock.Any(), fmt.Sprintf("https://test.com/book/%v", i)).Return(fmt.Sprintf("content %v", i), nil)

					p.EXPECT().ParseBook(fmt.Sprintf("content %v", i)).Return(nil, parse.ErrParseBookFieldsNotFound)

					rpo.EXPECT().CreateBook(&model.Book{
						Site:     "test-explore",
						ID:       i,
						HashCode: model.GenerateHash(),
					}).Return(nil)
					rpo.EXPECT().SaveError(&model.Book{Site: "test-explore", ID: i, HashCode: model.GenerateHash(), Error: err}, err).Return(nil)
				}

				return ServiceImp{
					name:   "test-explore",
					ctx:    context.Background(),
					sema:   semaphore.NewWeighted(1),
					client: c,
					parser: p,
					rpo:    rpo,
					conf: config.SiteConfig{
						URL:             config.URLConfig{Base: "https://test.com/book/%v"},
						MaxExploreError: 2,
					},
				}
			},
			expectedError:  false,
			expectErrorStr: "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			serv := test.setupServ(ctrl)
			err := serv.Explore()

			if test.expectedError {
				assert.Error(t, err)
			}

			if err != nil {
				assert.EqualError(t, err, test.expectErrorStr)
			}
		})
	}
}
