package service

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/htchan/BookSpider/internal/client/v1"
	mockclient "github.com/htchan/BookSpider/internal/mock/client/v1"
	mockrepo "github.com/htchan/BookSpider/internal/mock/repo"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewBookService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		clis        map[string]client.Client
		rpo         repo.Repository
		storagePath string
		want        BookService
	}{
		{
			name: "happy flow",
			clis: map[string]client.Client{
				"test": nil,
			},
			rpo:         nil,
			storagePath: "test",
			want: &bookServiceImpl{
				storagePath: "test",
				clis:        map[string]client.Client{"test": nil},
				rpo:         nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := NewBookService(test.clis, test.rpo, test.storagePath)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestBookServiceImpl_SupportBook(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		clis map[string]client.Client
		bk   *model.Book
		want bool
	}{
		{
			name: "happy flow/support book",
			clis: map[string]client.Client{
				"test": nil,
			},
			bk: &model.Book{
				Site: "test",
			},
			want: true,
		},
		{
			name: "happy flow/not support book",
			clis: map[string]client.Client{
				"test": nil,
			},
			bk: &model.Book{
				Site: "unknown",
			},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			serv := &bookServiceImpl{clis: test.clis}
			got := serv.SupportBook(test.bk)
			assert.Equal(t, test.want, got)
		})
	}
}

func Test_bookServiceImpl_bookFileLocation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		serv     *bookServiceImpl
		bk       *model.Book
		wantPath string
	}{
		{
			name: "happy flow with hash code = 0",
			serv: &bookServiceImpl{
				storagePath: "/data/books",
			},
			bk: &model.Book{
				ID:       123,
				Site:     "testsite",
				HashCode: 0,
			},
			wantPath: "/data/books/testsite/123.txt",
		},
		{
			name: "happy flow with hash code > 0",
			serv: &bookServiceImpl{
				storagePath: "/data/books",
			},
			bk: &model.Book{
				ID:       123,
				Site:     "testsite",
				HashCode: 999,
			},
			wantPath: "/data/books/testsite/123-vrr.txt",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			gotPath := test.serv.bookFileLocation(test.bk)
			assert.Equal(t, test.wantPath, gotPath)
		})
	}
}

func Test_bookServiceImpl_bookClient(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli1 := mockclient.NewMockClient(nil)
	cli2 := mockclient.NewMockClient(nil)

	tests := []struct {
		name string
		serv *bookServiceImpl
		bk   *model.Book
		want client.Client
	}{
		{
			name: "happy flow for site1",
			serv: &bookServiceImpl{
				clis: map[string]client.Client{
					"site1": cli1,
					"site2": cli2,
				},
			},
			bk: &model.Book{
				Site: "site1",
			},
			want: cli1,
		},
		{
			name: "happy flow for site2",
			serv: &bookServiceImpl{
				clis: map[string]client.Client{
					"site1": cli1,
					"site2": cli2,
				},
			},
			bk: &model.Book{
				Site: "site2",
			},
			want: cli2,
		},
		{
			name: " return nil if site not defined",
			serv: &bookServiceImpl{
				clis: map[string]client.Client{
					"site1": cli1,
					"site2": cli2,
				},
			},
			bk: &model.Book{
				Site: "site3",
			},
			want: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := test.serv.bookClient(test.bk)
			assert.Equal(t, test.want, got)
		})
	}
}

func Test_isNewBook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		bk     *model.Book
		bkInfo *client.BookInfo
		want   bool
	}{
		{
			name: "book with error status and nil error is new",
			bk: &model.Book{
				Status: model.StatusError,
				Error:  model.Error{Err: nil},
			},
			bkInfo: nil,
			want:   true,
		},
		{
			name: "book with error status and non-nil error is not new",
			bk: &model.Book{
				Status: model.StatusError,
				Error:  model.Error{Err: assert.AnError},
			},
			bkInfo: nil,
			want:   false,
		},
		{
			name: "book with status other than error and different title is new",
			bk: &model.Book{
				Status: model.StatusEnd,
				Title:  "old title",
			},
			bkInfo: &client.BookInfo{
				Title: "new title",
			},
			want: true,
		},
		{
			name: "book with status other than error and different author is new",
			bk: &model.Book{
				Status: model.StatusEnd,
				Writer: model.Writer{
					Name: "old author",
				},
			},
			bkInfo: &client.BookInfo{
				Author: "new author",
			},
			want: true,
		},
		{
			name: "book with status other than error and different book type is new",
			bk: &model.Book{
				Status: model.StatusEnd,
				Type:   "old type",
			},
			bkInfo: &client.BookInfo{
				Type: "new type",
			},
			want: true,
		},
		{
			name: "book with status other than error and same title, author and type is not new",
			bk: &model.Book{
				Status: model.StatusEnd,
				Title:  "same title",
				Writer: model.Writer{
					Name: "same author",
				},
				Type: "same type",
			},
			bkInfo: &client.BookInfo{
				Title:  "same title",
				Author: "same author",
				Type:   "same type",
			},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := isNewBook(test.bk, test.bkInfo)
			assert.Equal(t, test.want, got)
		})
	}
}

func Test_isBookUpdated(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		bk     *model.Book
		bkInfo *client.BookInfo
		want   bool
	}{
		{
			name: "book with different update date is updated",
			bk: &model.Book{
				UpdateDate: "2023-01-01 00:00:00 +0000 UTC",
			},
			bkInfo: &client.BookInfo{
				UpdateDate: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			want: true,
		},
		{
			name: "book with different update chapter is updated",
			bk: &model.Book{
				UpdateChapter: "Chapter 1",
			},
			bkInfo: &client.BookInfo{
				UpdateChapter: "Chapter 2",
			},
			want: true,
		},
		{
			name: "book with same update date and update chapter is not updated",
			bk: &model.Book{
				UpdateDate:    "2023-01-01 00:00:00 +0000 UTC",
				UpdateChapter: "Chapter 1",
			},
			bkInfo: &client.BookInfo{
				UpdateDate:    time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdateChapter: "Chapter 1",
			},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := isBookUpdated(test.bk, test.bkInfo)
			assert.Equal(t, test.want, got)
		})
	}
}

func Test_bookServiceImpl_UpdateBook(t *testing.T) {
	t.Parallel()
	currentYear := time.Now().Year()

	removeBookHaah := func(_ context.Context, bk *model.Book) error {
		bk.HashCode = 0
		return nil
	}

	inProgressBkInfo := &client.BookInfo{
		Title:         "test title",
		Author:        "writer",
		Type:          "book type",
		UpdateDate:    time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdateChapter: "Chapter 1",
	}
	inProgressBkWithHash := &model.Book{
		Site:          "test",
		ID:            1,
		HashCode:      int(time.Now().Unix()),
		Title:         "test title",
		Writer:        model.Writer{Name: "writer"},
		Type:          "book type",
		UpdateDate:    time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC).String(),
		UpdateChapter: "Chapter 1",
		Status:        model.StatusInProgress,
	}
	inProgressBkWithoutHash := &model.Book{
		Site:          "test",
		ID:            1,
		Title:         "test title",
		Writer:        model.Writer{Name: "writer"},
		Type:          "book type",
		UpdateDate:    time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC).String(),
		UpdateChapter: "Chapter 1",
		Status:        model.StatusInProgress,
	}

	endBkInfo := &client.BookInfo{
		Title:         "test title",
		Author:        "writer",
		Type:          "book type",
		UpdateDate:    time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdateChapter: "Chapter 結局",
	}
	endBkWithHash := &model.Book{
		Site:          "test",
		ID:            1,
		HashCode:      int(time.Now().Unix()),
		Title:         "test title",
		Writer:        model.Writer{Name: "writer"},
		Type:          "book type",
		UpdateDate:    time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC).String(),
		UpdateChapter: "Chapter 結局",
		Status:        model.StatusEnd,
	}
	endBkWithoutHash := &model.Book{
		Site:          "test",
		ID:            1,
		Title:         "test title",
		Writer:        model.Writer{Name: "writer"},
		Type:          "book type",
		UpdateDate:    time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC).String(),
		UpdateChapter: "Chapter 結局",
		Status:        model.StatusEnd,
	}

	tests := []struct {
		name    string
		bk      *model.Book
		serv    func(ctrl *gomock.Controller) BookService
		wantBk  *model.Book
		wantErr error
	}{
		{
			name: "happy flow/new book",
			bk:   &model.Book{Site: "test", ID: 1},
			serv: func(ctrl *gomock.Controller) BookService {
				cli := mockclient.NewMockClient(ctrl)
				cli.EXPECT().
					GetBookInfo(gomock.Any(), "1").
					Return(inProgressBkInfo, nil)

				repo := mockrepo.NewMockRepository(ctrl)
				repo.EXPECT().
					CreateBook(gomock.Any(), inProgressBkWithHash).DoAndReturn(removeBookHaah)
				repo.EXPECT().
					SaveWriter(gomock.Any(), &model.Writer{Name: "writer"}).
					Return(nil)
				repo.EXPECT().
					SaveError(gomock.Any(), inProgressBkWithoutHash, nil).
					Return(nil)

				return &bookServiceImpl{
					clis: map[string]client.Client{
						"test": cli,
					},
					rpo: repo,
				}
			},
			wantBk: inProgressBkWithoutHash,
		},
		{
			name: "happy flow/new book/chapter contains end keyword",
			bk:   &model.Book{Site: "test", ID: 1},
			serv: func(ctrl *gomock.Controller) BookService {
				cli := mockclient.NewMockClient(ctrl)
				cli.EXPECT().
					GetBookInfo(gomock.Any(), "1").
					Return(endBkInfo, nil)

				repo := mockrepo.NewMockRepository(ctrl)
				repo.EXPECT().
					CreateBook(gomock.Any(), endBkWithHash).DoAndReturn(removeBookHaah)
				repo.EXPECT().
					SaveWriter(gomock.Any(), &model.Writer{Name: "writer"}).
					Return(nil)
				repo.EXPECT().
					SaveError(gomock.Any(), endBkWithoutHash, nil).
					Return(nil)

				return &bookServiceImpl{
					clis: map[string]client.Client{
						"test": cli,
					},
					rpo: repo,
				}
			},
			wantBk: &model.Book{
				Site: "test", ID: 1, HashCode: 0,
				Writer:        model.Writer{Name: "writer"},
				Title:         "test title",
				Type:          "book type",
				UpdateDate:    time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC).String(),
				UpdateChapter: "Chapter 結局",
				Status:        model.StatusEnd,
			},
		},
		{
			name: "error flow/new book/db raise error",
			bk:   &model.Book{Site: "test", ID: 1},
			serv: func(ctrl *gomock.Controller) BookService {
				cli := mockclient.NewMockClient(ctrl)
				cli.EXPECT().
					GetBookInfo(gomock.Any(), "1").
					Return(inProgressBkInfo, nil)

				repo := mockrepo.NewMockRepository(ctrl)
				repo.EXPECT().
					CreateBook(gomock.Any(), inProgressBkWithHash).DoAndReturn(removeBookHaah)
				repo.EXPECT().
					SaveWriter(gomock.Any(), &model.Writer{Name: "writer"}).
					Return(assert.AnError)
				repo.EXPECT().
					SaveError(gomock.Any(), inProgressBkWithoutHash, nil).
					Return(assert.AnError)

				return &bookServiceImpl{
					clis: map[string]client.Client{
						"test": cli,
					},
					rpo: repo,
				}
			},
			wantBk: &model.Book{
				Site: "test", ID: 1, HashCode: 0,
				Writer:        model.Writer{Name: "writer"},
				Title:         "test title",
				Type:          "book type",
				UpdateDate:    time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC).String(),
				UpdateChapter: "Chapter 1",
				Status:        model.StatusInProgress,
			},
			wantErr: assert.AnError,
		},
		{
			name: "error flow/new book/client return error",
			bk:   &model.Book{Site: "test", ID: 1},
			serv: func(ctrl *gomock.Controller) BookService {
				cli := mockclient.NewMockClient(ctrl)
				cli.EXPECT().
					GetBookInfo(gomock.Any(), "1").
					Return(nil, assert.AnError)

				repo := mockrepo.NewMockRepository(ctrl)
				repo.EXPECT().
					CreateBook(gomock.Any(), &model.Book{
						ID: 1, Site: "test", HashCode: int(time.Now().Unix()), Status: model.StatusError, Error: model.Error{Err: assert.AnError},
					}).
					DoAndReturn(removeBookHaah)
				repo.EXPECT().
					SaveError(gomock.Any(), &model.Book{
						ID: 1, Site: "test", Status: model.StatusError, Error: model.Error{Err: assert.AnError},
					}, assert.AnError).
					Return(nil)

				return &bookServiceImpl{
					clis: map[string]client.Client{
						"test": cli,
					},
					rpo: repo,
				}
			},
			wantBk: &model.Book{
				ID:    1,
				Site:  "test",
				Error: model.Error{Err: assert.AnError},
			},
			wantErr: assert.AnError,
		},
		{
			name: "happy flow/existing book/in progress book updated",
			bk: &model.Book{
				Site: "test", ID: 1,
				Title: "test title", Type: "book type",
				Writer:        model.Writer{Name: "writer"},
				UpdateDate:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).String(),
				UpdateChapter: "Chapter 0",
				Status:        model.StatusInProgress,
			},
			serv: func(ctrl *gomock.Controller) BookService {
				cli := mockclient.NewMockClient(ctrl)
				cli.EXPECT().
					GetBookInfo(gomock.Any(), "1").
					Return(inProgressBkInfo, nil)

				repo := mockrepo.NewMockRepository(ctrl)
				repo.EXPECT().
					UpdateBook(gomock.Any(), inProgressBkWithoutHash).DoAndReturn(removeBookHaah)
				repo.EXPECT().
					SaveWriter(gomock.Any(), &model.Writer{Name: "writer"}).
					Return(nil)
				repo.EXPECT().
					SaveError(gomock.Any(), inProgressBkWithoutHash, nil).
					Return(nil)

				return &bookServiceImpl{
					clis: map[string]client.Client{
						"test": cli,
					},
					rpo: repo,
				}
			},
			wantBk: &model.Book{
				Site: "test", ID: 1, HashCode: 0,
				Writer:        model.Writer{Name: "writer"},
				Title:         "test title",
				Type:          "book type",
				UpdateDate:    time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC).String(),
				UpdateChapter: "Chapter 1",
				Status:        model.StatusInProgress,
			},
		},
		{
			name: "happy flow/existing book/error book updated",
			bk: &model.Book{
				Site: "test", ID: 1,
				Status: model.StatusError,
				Error:  model.Error{Err: assert.AnError},
			},
			serv: func(ctrl *gomock.Controller) BookService {
				cli := mockclient.NewMockClient(ctrl)
				cli.EXPECT().
					GetBookInfo(gomock.Any(), "1").
					Return(inProgressBkInfo, nil)

				repo := mockrepo.NewMockRepository(ctrl)
				repo.EXPECT().
					UpdateBook(gomock.Any(), inProgressBkWithoutHash).DoAndReturn(removeBookHaah)
				repo.EXPECT().
					SaveWriter(gomock.Any(), &model.Writer{Name: "writer"}).
					Return(nil)
				repo.EXPECT().
					SaveError(gomock.Any(), inProgressBkWithoutHash, nil).
					Return(nil)

				return &bookServiceImpl{
					clis: map[string]client.Client{
						"test": cli,
					},
					rpo: repo,
				}
			},
			wantBk: &model.Book{
				Site: "test", ID: 1, HashCode: 0,
				Writer:        model.Writer{Name: "writer"},
				Title:         "test title",
				Type:          "book type",
				UpdateDate:    time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC).String(),
				UpdateChapter: "Chapter 1",
				Status:        model.StatusInProgress,
			},
		},
		{
			name: "happy flow/existing book/in progress book updated/chapter contains end keyword",
			bk: &model.Book{
				Site: "test", ID: 1,
				Title: "test title", Type: "book type",
				Writer:        model.Writer{Name: "writer"},
				UpdateDate:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).String(),
				UpdateChapter: "Chapter 結局",
				Status:        model.StatusInProgress,
			},
			serv: func(ctrl *gomock.Controller) BookService {
				cli := mockclient.NewMockClient(ctrl)
				cli.EXPECT().
					GetBookInfo(gomock.Any(), "1").
					Return(endBkInfo, nil)

				repo := mockrepo.NewMockRepository(ctrl)
				repo.EXPECT().
					UpdateBook(gomock.Any(), endBkWithoutHash).DoAndReturn(removeBookHaah)
				repo.EXPECT().
					SaveWriter(gomock.Any(), &model.Writer{Name: "writer"}).
					Return(nil)
				repo.EXPECT().
					SaveError(gomock.Any(), endBkWithoutHash, nil).
					Return(nil)

				return &bookServiceImpl{
					clis: map[string]client.Client{
						"test": cli,
					},
					rpo: repo,
				}
			},
			wantBk: &model.Book{
				Site: "test", ID: 1, HashCode: 0,
				Writer:        model.Writer{Name: "writer"},
				Title:         "test title",
				Type:          "book type",
				UpdateDate:    time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC).String(),
				UpdateChapter: "Chapter 結局",
				Status:        model.StatusEnd,
			},
		},
		{
			name: "happy flow/existing book/error book updated/chapter contains end keyword",
			bk: &model.Book{
				Site: "test", ID: 1,
				Status: model.StatusError,
				Error:  model.Error{Err: assert.AnError},
			},
			serv: func(ctrl *gomock.Controller) BookService {
				cli := mockclient.NewMockClient(ctrl)
				cli.EXPECT().
					GetBookInfo(gomock.Any(), "1").
					Return(endBkInfo, nil)

				repo := mockrepo.NewMockRepository(ctrl)
				repo.EXPECT().
					UpdateBook(gomock.Any(), endBkWithoutHash).DoAndReturn(removeBookHaah)
				repo.EXPECT().
					SaveWriter(gomock.Any(), &model.Writer{Name: "writer"}).
					Return(nil)
				repo.EXPECT().
					SaveError(gomock.Any(), endBkWithoutHash, nil).
					Return(nil)

				return &bookServiceImpl{
					clis: map[string]client.Client{
						"test": cli,
					},
					rpo: repo,
				}
			},
			wantBk: &model.Book{
				Site: "test", ID: 1, HashCode: 0,
				Writer:        model.Writer{Name: "writer"},
				Title:         "test title",
				Type:          "book type",
				UpdateDate:    time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC).String(),
				UpdateChapter: "Chapter 結局",
				Status:        model.StatusEnd,
			},
		},
		{
			name: "happy flow/existing book/not updated",
			bk: &model.Book{
				Site: "test", ID: 1,
				Title: "test title", Type: "book type",
				Writer:        model.Writer{Name: "writer"},
				UpdateDate:    time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC).String(),
				UpdateChapter: "Chapter 1",
				Status:        model.StatusInProgress,
			},
			serv: func(ctrl *gomock.Controller) BookService {
				cli := mockclient.NewMockClient(ctrl)
				cli.EXPECT().
					GetBookInfo(gomock.Any(), "1").
					Return(inProgressBkInfo, nil)

				repo := mockrepo.NewMockRepository(ctrl)
				return &bookServiceImpl{
					clis: map[string]client.Client{
						"test": cli,
					},
					rpo: repo,
				}
			},
			wantBk: &model.Book{
				Site: "test", ID: 1, HashCode: 0,
				Writer:        model.Writer{Name: "writer"},
				Title:         "test title",
				Type:          "book type",
				UpdateDate:    time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC).String(),
				UpdateChapter: "Chapter 1",
				Status:        model.StatusInProgress,
			},
		},
		{
			name: "error flow/existing book/db raise error",
			bk: &model.Book{
				Site: "test", ID: 1,
				Title: "test title", Type: "book type",
				Writer:        model.Writer{Name: "writer"},
				UpdateDate:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).String(),
				UpdateChapter: "Chapter 0",
				Status:        model.StatusInProgress,
			},
			serv: func(ctrl *gomock.Controller) BookService {
				cli := mockclient.NewMockClient(ctrl)
				cli.EXPECT().
					GetBookInfo(gomock.Any(), "1").
					Return(inProgressBkInfo, nil)

				repo := mockrepo.NewMockRepository(ctrl)
				repo.EXPECT().
					UpdateBook(gomock.Any(), inProgressBkWithoutHash).DoAndReturn(removeBookHaah)
				repo.EXPECT().
					SaveWriter(gomock.Any(), &model.Writer{Name: "writer"}).
					Return(assert.AnError)
				repo.EXPECT().
					SaveError(gomock.Any(), inProgressBkWithoutHash, nil).
					Return(assert.AnError)

				return &bookServiceImpl{
					clis: map[string]client.Client{
						"test": cli,
					},
					rpo: repo,
				}
			},
			wantBk: &model.Book{
				Site: "test", ID: 1, HashCode: 0,
				Writer:        model.Writer{Name: "writer"},
				Title:         "test title",
				Type:          "book type",
				UpdateDate:    time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC).String(),
				UpdateChapter: "Chapter 1",
				Status:        model.StatusInProgress,
			},
			wantErr: assert.AnError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			err := test.serv(ctrl).UpdateBook(context.Background(), test.bk)
			assert.Equal(t, test.wantBk, test.bk)
			assert.ErrorIs(t, err, test.wantErr)
		})
	}
}

func Test_bookServiceImpl_downloadChapter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		bk          *model.Book
		ch          *model.Chapter
		mockClient  func(ctrl *gomock.Controller) client.Client
		wantChapter *model.Chapter
		wantErr     error
	}{
		{
			name: "happy flow",
			bk: &model.Book{
				ID:   1,
				Site: "test",
			},
			ch: &model.Chapter{
				Title: "Chapter 1",
				URL:   "http://example.com/ch1",
			},
			mockClient: func(ctrl *gomock.Controller) client.Client {
				cli := mockclient.NewMockClient(ctrl)
				cli.EXPECT().
					GetChapterContent(gomock.Any(), client.ChapterEntry{
						Title: "Chapter 1",
						URL:   "http://example.com/ch1",
					}).
					Return(&client.ChapterContent{
						Title: "Actual Chapter 1",
						Body:  "line 1.\n \nline 2.\n\n\nline 3.",
					}, nil)

				return cli
			},
			wantChapter: &model.Chapter{
				Title:   "Actual Chapter 1",
				URL:     "http://example.com/ch1",
				Content: "line 1.\n\nline 2.\n\nline 3.",
			},
			wantErr: nil,
		},
		{
			name: "client returns error",
			bk: &model.Book{
				ID:   1,
				Site: "test",
			},
			ch: &model.Chapter{
				Title: "Chapter 1",
				URL:   "http://example.com/ch1",
			},
			mockClient: func(ctrl *gomock.Controller) client.Client {
				cli := mockclient.NewMockClient(ctrl)
				cli.EXPECT().
					GetChapterContent(gomock.Any(), client.ChapterEntry{
						Title: "Chapter 1",
						URL:   "http://example.com/ch1",
					}).
					Return(nil, assert.AnError)

				return cli
			},
			wantChapter: &model.Chapter{
				Title: "Chapter 1",
				URL:   "http://example.com/ch1",
				Error: assert.AnError,
			},
			wantErr: assert.AnError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := &bookServiceImpl{
				clis: map[string]client.Client{
					"test": test.mockClient(ctrl),
				},
			}

			err := serv.downloadChapter(context.Background(), test.bk, test.ch)
			assert.Equal(t, test.wantChapter, test.ch)
			assert.ErrorIs(t, err, test.wantErr)
		})
	}
}

func Test_bookServiceImpl_DownloadBook(t *testing.T) {
	t.Parallel()

	assert.NoError(t, os.Mkdir("./testdata", os.ModePerm))
	assert.NoError(t, os.Mkdir("./testdata/test", os.ModePerm))

	t.Cleanup(func() {
		os.RemoveAll("./testdata")
	})

	tests := []struct {
		name            string
		bk              *model.Book
		serv            func(ctrl *gomock.Controller) BookService
		wantBookPath    string
		wantBookContent string
		wantErr         error
	}{
		{
			name: "happy flow",
			bk: &model.Book{
				ID:       1,
				Site:     "test",
				Title:    "Test Book",
				Writer:   model.Writer{Name: "Test Author"},
				Status:   model.StatusEnd,
				HashCode: 0,
			},
			serv: func(ctrl *gomock.Controller) BookService {
				cli := mockclient.NewMockClient(ctrl)
				chapterEntryList := client.ChapterEntryList{
					{Title: "Chapter 1", URL: "http://example.com/ch1"},
					{Title: "Chapter 2", URL: "http://example.com/ch2"},
				}
				cli.EXPECT().
					GetBookChapterList(gomock.Any(), "1").
					Return(chapterEntryList, nil)

				for _, chapter := range chapterEntryList {
					cli.EXPECT().
						GetChapterContent(gomock.Any(), chapter).
						Return(&client.ChapterContent{
							Title: chapter.Title,
							Body:  "Content of " + chapter.Title + ".",
						}, nil)
				}

				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().
					UpdateBook(gomock.Any(), &model.Book{
						ID:           1,
						Site:         "test",
						Title:        "Test Book",
						Writer:       model.Writer{Name: "Test Author"},
						Status:       model.StatusEnd,
						IsDownloaded: true,
						HashCode:     0,
					}).
					Return(nil)

				return &bookServiceImpl{
					clis: map[string]client.Client{
						"test": cli,
					},
					rpo:         rpo,
					storagePath: "./testdata",
				}
			},
			wantBookPath:    "./testdata/test/1.txt",
			wantBookContent: "Test Book\nTest Author\n--------------------\n\nChapter 1\n--------------------\nContent of Chapter 1.\n--------------------\nChapter 2\n--------------------\nContent of Chapter 2.\n--------------------\n",
			wantErr:         nil,
		},
		{
			name: "book status is not ended",
			bk: &model.Book{
				ID:     1,
				Site:   "test",
				Status: model.StatusInProgress,
			},
			serv: func(ctrl *gomock.Controller) BookService {
				return &bookServiceImpl{}
			},
			wantErr: ErrBookStatusNotEnd,
		},
		{
			name: "book is already downloaded",
			bk: &model.Book{
				ID:           1,
				Site:         "test",
				Status:       model.StatusEnd,
				IsDownloaded: true,
			},
			serv: func(ctrl *gomock.Controller) BookService {
				return &bookServiceImpl{}
			},
			wantErr: ErrBookAlreadyDownloaded,
		},
		{
			name: "book client return error in get book clapter list",
			bk: &model.Book{
				ID:     1,
				Site:   "test",
				Status: model.StatusEnd,
			},
			serv: func(ctrl *gomock.Controller) BookService {
				cli := mockclient.NewMockClient(ctrl)
				cli.EXPECT().
					GetBookChapterList(gomock.Any(), "1").
					Return(nil, assert.AnError)

				return &bookServiceImpl{
					clis: map[string]client.Client{
						"test": cli,
					},
				}
			},
			wantErr: assert.AnError,
		},
		{
			name: "book client return error in download chapter less than 10% threshold",
			bk: &model.Book{
				ID:     2,
				Site:   "test",
				Status: model.StatusEnd,
			},
			serv: func(ctrl *gomock.Controller) BookService {
				cli := mockclient.NewMockClient(ctrl)
				chapterEntryList := client.ChapterEntryList{
					{Title: "Chapter 1", URL: "http://example.com/ch1"},
					{Title: "Chapter 2", URL: "http://example.com/ch2"},
					{Title: "Chapter 3", URL: "http://example.com/ch3"},
					{Title: "Chapter 4", URL: "http://example.com/ch4"},
					{Title: "Chapter 5", URL: "http://example.com/ch5"},
					{Title: "Chapter 6", URL: "http://example.com/ch6"},
					{Title: "Chapter 7", URL: "http://example.com/ch7"},
					{Title: "Chapter 8", URL: "http://example.com/ch8"},
					{Title: "Chapter 9", URL: "http://example.com/ch9"},
					{Title: "Chapter 10", URL: "http://example.com/ch10"},
				}
				cli.EXPECT().
					GetBookChapterList(gomock.Any(), "2").
					Return(chapterEntryList, nil)

				for i, chapter := range chapterEntryList {
					if i == 0 {
						cli.EXPECT().
							GetChapterContent(gomock.Any(), chapter).
							Return(nil, assert.AnError)
					} else {
						cli.EXPECT().
							GetChapterContent(gomock.Any(), chapter).
							Return(&client.ChapterContent{
								Title: chapter.Title,
								Body:  "Content of " + chapter.Title + ".",
							}, nil)
					}
				}

				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().
					UpdateBook(gomock.Any(), &model.Book{
						ID:           2,
						Site:         "test",
						Status:       model.StatusEnd,
						IsDownloaded: true,
						HashCode:     0,
					}).
					Return(nil)
				return &bookServiceImpl{
					clis: map[string]client.Client{
						"test": cli,
					},
					rpo:         rpo,
					storagePath: "./testdata",
				}
			},
			wantBookPath:    "./testdata/test/2.txt",
			wantBookContent: "\n\n--------------------\n\nChapter 1\n--------------------\n\n--------------------\nChapter 2\n--------------------\nContent of Chapter 2.\n--------------------\nChapter 3\n--------------------\nContent of Chapter 3.\n--------------------\nChapter 4\n--------------------\nContent of Chapter 4.\n--------------------\nChapter 5\n--------------------\nContent of Chapter 5.\n--------------------\nChapter 6\n--------------------\nContent of Chapter 6.\n--------------------\nChapter 7\n--------------------\nContent of Chapter 7.\n--------------------\nChapter 8\n--------------------\nContent of Chapter 8.\n--------------------\nChapter 9\n--------------------\nContent of Chapter 9.\n--------------------\nChapter 10\n--------------------\nContent of Chapter 10.\n--------------------\n",
			wantErr:         nil,
		},
		{
			name: "too many failed chapters",
			bk: &model.Book{
				ID:     3,
				Site:   "test",
				Status: model.StatusEnd,
			},
			serv: func(ctrl *gomock.Controller) BookService {
				cli := mockclient.NewMockClient(ctrl)
				chapterEntryList := client.ChapterEntryList{
					{Title: "Chapter 1", URL: "http://example.com/ch1"},
					{Title: "Chapter 2", URL: "http://example.com/ch2"},
				}
				cli.EXPECT().
					GetBookChapterList(gomock.Any(), "3").
					Return(chapterEntryList, nil)

				for _, chapter := range chapterEntryList {
					cli.EXPECT().
						GetChapterContent(gomock.Any(), chapter).
						Return(nil, assert.AnError)
				}

				return &bookServiceImpl{
					clis: map[string]client.Client{
						"test": cli,
					},
					storagePath: "./testdata",
				}
			},
			wantErr: ErrTooManyFailedChapters,
		},
		{
			name: "error in update book record",
			bk: &model.Book{
				ID:     4,
				Site:   "test",
				Status: model.StatusEnd,
			},
			serv: func(ctrl *gomock.Controller) BookService {
				cli := mockclient.NewMockClient(ctrl)
				chapterEntryList := client.ChapterEntryList{
					{Title: "Chapter 1", URL: "http://example.com/ch1"},
					{Title: "Chapter 2", URL: "http://example.com/ch2"},
				}
				cli.EXPECT().
					GetBookChapterList(gomock.Any(), "4").
					Return(chapterEntryList, nil)

				for _, chapter := range chapterEntryList {
					cli.EXPECT().
						GetChapterContent(gomock.Any(), chapter).
						Return(&client.ChapterContent{
							Title: chapter.Title,
							Body:  "Content of " + chapter.Title + ".",
						}, nil)
				}

				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().
					UpdateBook(gomock.Any(), &model.Book{
						ID:           4,
						Site:         "test",
						Status:       model.StatusEnd,
						IsDownloaded: true,
						HashCode:     0,
					}).
					Return(assert.AnError)

				return &bookServiceImpl{
					clis: map[string]client.Client{
						"test": cli,
					},
					rpo:         rpo,
					storagePath: "./testdata",
				}
			},
			wantBookPath:    "./testdata/test/4.txt",
			wantBookContent: "\n\n--------------------\n\nChapter 1\n--------------------\nContent of Chapter 1.\n--------------------\nChapter 2\n--------------------\nContent of Chapter 2.\n--------------------\n",
			wantErr:         assert.AnError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			err := test.serv(ctrl).DownloadBook(context.Background(), test.bk)
			assert.ErrorIs(t, err, test.wantErr)

			if test.wantBookPath != "" {
				content, err := os.ReadFile(test.wantBookPath)
				assert.NoError(t, err)
				assert.Equal(t, test.wantBookContent, string(content))
			}
		})
	}
}

func Test_bookServiceImpl_ProcessBook(t *testing.T) {
	t.Parallel()
	assert.NoError(t, os.Mkdir("./process_test_data", os.ModePerm))
	assert.NoError(t, os.Mkdir("./process_test_data/test", os.ModePerm))

	t.Cleanup(func() {
		os.RemoveAll("./process_test_data")
	})

	tests := []struct {
		name            string
		bk              *model.Book
		serv            func(*gomock.Controller) BookService
		wantErr         error
		wantBookPath    string
		wantBookContent string
	}{
		{
			name: "happy flow/update and download success",
			bk: &model.Book{
				ID:     1,
				Site:   "test",
				Status: model.StatusInProgress,
			},
			serv: func(ctrl *gomock.Controller) BookService {
				cli := mockclient.NewMockClient(ctrl)
				cli.EXPECT().
					GetBookInfo(gomock.Any(), "1").
					Return(&client.BookInfo{
						Title:         "test title",
						Author:        "writer",
						Type:          "book type",
						UpdateDate:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdateChapter: "Chapter 結局",
					}, nil)

				chapterEntryList := client.ChapterEntryList{
					{Title: "Chapter 1", URL: "http://example.com/ch1"},
					{Title: "Chapter 2", URL: "http://example.com/ch2"},
				}
				cli.EXPECT().
					GetBookChapterList(gomock.Any(), "1").
					Return(chapterEntryList, nil)

				for _, chapter := range chapterEntryList {
					cli.EXPECT().
						GetChapterContent(gomock.Any(), chapter).
						Return(&client.ChapterContent{
							Title: chapter.Title,
							Body:  "Content of " + chapter.Title + ".",
						}, nil)
				}

				repo := mockrepo.NewMockRepository(ctrl)
				repo.EXPECT().
					CreateBook(gomock.Any(), &model.Book{
						Site:          "test",
						ID:            1,
						HashCode:      int(time.Now().Unix()),
						Title:         "test title",
						Writer:        model.Writer{Name: "writer"},
						Type:          "book type",
						UpdateDate:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).String(),
						UpdateChapter: "Chapter 結局",
						Status:        model.StatusEnd,
					}).DoAndReturn(func(_ context.Context, bk *model.Book) error {
					bk.HashCode = 0
					return nil
				})
				repo.EXPECT().
					SaveWriter(gomock.Any(), &model.Writer{Name: "writer"}).
					Return(nil)
				repo.EXPECT().
					SaveError(gomock.Any(), &model.Book{
						Site:          "test",
						ID:            1,
						Title:         "test title",
						Writer:        model.Writer{Name: "writer"},
						Type:          "book type",
						UpdateDate:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).String(),
						UpdateChapter: "Chapter 結局",
						Status:        model.StatusEnd,
					}, nil).
					Return(nil)
				repo.EXPECT().
					UpdateBook(gomock.Any(), &model.Book{
						Site:          "test",
						ID:            1,
						Title:         "test title",
						Writer:        model.Writer{Name: "writer"},
						Type:          "book type",
						UpdateDate:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).String(),
						UpdateChapter: "Chapter 結局",
						Status:        model.StatusEnd,
						IsDownloaded:  true,
					}).
					Return(nil)

				return &bookServiceImpl{
					rpo: repo,
					clis: map[string]client.Client{
						"test": cli,
					},
					storagePath: "./process_test_data",
				}
			},
			wantErr:         nil,
			wantBookPath:    "./process_test_data/test/1.txt",
			wantBookContent: "test title\nwriter\n--------------------\n\nChapter 1\n--------------------\nContent of Chapter 1.\n--------------------\nChapter 2\n--------------------\nContent of Chapter 2.\n--------------------\n",
		},
		{
			name: "happy flow/update success/book is not ended",
			bk: &model.Book{
				ID:     2,
				Site:   "test",
				Status: model.StatusEnd,
			},
			serv: func(ctrl *gomock.Controller) BookService {
				now := time.Now().UTC()
				cli := mockclient.NewMockClient(ctrl)
				cli.EXPECT().
					GetBookInfo(gomock.Any(), "2").
					Return(&client.BookInfo{
						Title:         "test title",
						Author:        "writer",
						Type:          "book type",
						UpdateDate:    now,
						UpdateChapter: "Chapter 1",
					}, nil)

				repo := mockrepo.NewMockRepository(ctrl)
				repo.EXPECT().
					CreateBook(gomock.Any(), &model.Book{
						Site:          "test",
						ID:            2,
						HashCode:      int(time.Now().Unix()),
						Title:         "test title",
						Writer:        model.Writer{Name: "writer"},
						Type:          "book type",
						UpdateDate:    now.String(),
						UpdateChapter: "Chapter 1",
						Status:        model.StatusInProgress,
					}).DoAndReturn(func(_ context.Context, bk *model.Book) error {
					bk.HashCode = 0
					return nil
				})
				repo.EXPECT().
					SaveWriter(gomock.Any(), &model.Writer{Name: "writer"}).
					Return(nil)
				repo.EXPECT().
					SaveError(gomock.Any(), &model.Book{
						Site:          "test",
						ID:            2,
						Title:         "test title",
						Writer:        model.Writer{Name: "writer"},
						Type:          "book type",
						UpdateDate:    now.String(),
						UpdateChapter: "Chapter 1",
						Status:        model.StatusInProgress,
					}, nil).
					Return(nil)

				return &bookServiceImpl{
					rpo: repo,
					clis: map[string]client.Client{
						"test": cli,
					},
				}
			},
			wantErr: nil,
		},
		{
			name: "error flow/update book return error",
			bk: &model.Book{
				ID:   3,
				Site: "test",
			},
			serv: func(ctrl *gomock.Controller) BookService {
				cli := mockclient.NewMockClient(ctrl)
				cli.EXPECT().
					GetBookInfo(gomock.Any(), "3").
					Return(nil, assert.AnError)

				repo := mockrepo.NewMockRepository(ctrl)
				repo.EXPECT().
					CreateBook(gomock.Any(), &model.Book{
						ID: 3, Site: "test", HashCode: int(time.Now().Unix()), Status: model.StatusError, Error: model.Error{Err: assert.AnError},
					}).
					DoAndReturn(func(_ context.Context, bk *model.Book) error {
						bk.HashCode = 0
						return nil
					})
				repo.EXPECT().
					SaveError(gomock.Any(), &model.Book{
						ID: 3, Site: "test", Status: model.StatusError, Error: model.Error{Err: assert.AnError},
					}, assert.AnError).
					Return(nil)

				return &bookServiceImpl{
					clis: map[string]client.Client{
						"test": cli,
					},
					rpo:         repo,
					storagePath: "./process_test_data",
				}
			},
			wantErr: assert.AnError,
		},
		{
			name: "error flow/download book return error",
			bk: &model.Book{
				ID:     4,
				Site:   "test",
				Status: model.StatusInProgress,
			},
			serv: func(ctrl *gomock.Controller) BookService {
				cli := mockclient.NewMockClient(ctrl)
				cli.EXPECT().
					GetBookInfo(gomock.Any(), "4").
					Return(&client.BookInfo{
						Title:         "test title",
						Author:        "writer",
						Type:          "book type",
						UpdateDate:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdateChapter: "Chapter 結局",
					}, nil)

				cli.EXPECT().
					GetBookChapterList(gomock.Any(), "4").
					Return(nil, assert.AnError)

				repo := mockrepo.NewMockRepository(ctrl)
				repo.EXPECT().
					CreateBook(gomock.Any(), &model.Book{
						Site:          "test",
						ID:            4,
						HashCode:      int(time.Now().Unix()),
						Title:         "test title",
						Writer:        model.Writer{Name: "writer"},
						Type:          "book type",
						UpdateDate:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).String(),
						UpdateChapter: "Chapter 結局",
						Status:        model.StatusEnd,
					}).DoAndReturn(func(_ context.Context, bk *model.Book) error {
					bk.HashCode = 0
					return nil
				})
				repo.EXPECT().
					SaveWriter(gomock.Any(), &model.Writer{Name: "writer"}).
					Return(nil)
				repo.EXPECT().
					SaveError(gomock.Any(), &model.Book{
						Site:          "test",
						ID:            4,
						Title:         "test title",
						Writer:        model.Writer{Name: "writer"},
						Type:          "book type",
						UpdateDate:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).String(),
						UpdateChapter: "Chapter 結局",
						Status:        model.StatusEnd,
					}, nil).
					Return(nil)

				return &bookServiceImpl{
					rpo: repo,
					clis: map[string]client.Client{
						"test": cli,
					},
					storagePath: "./process_test_data",
				}
			},
			wantErr: assert.AnError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			err := test.serv(ctrl).ProcessBook(context.Background(), test.bk)
			assert.ErrorIs(t, err, test.wantErr)

			if test.wantBookPath != "" {
				content, err := os.ReadFile(test.wantBookPath)
				assert.NoError(t, err)
				assert.Equal(t, test.wantBookContent, string(content))
			}
		})
	}
}
