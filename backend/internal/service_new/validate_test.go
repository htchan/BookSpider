package service

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	config "github.com/htchan/BookSpider/internal/config_new"
	"github.com/htchan/BookSpider/internal/mock"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/stretchr/testify/assert"
)

func Test_isEnd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		chapter string
		want    bool
	}{
		{
			name:    "chapter of ended book",
			chapter: "last chatper （完）",
			want:    true,
		},
		{
			name:    "chapter of not ended book",
			chapter: "chapter in middle",
			want:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := isEnd(test.chapter)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestServiceImp_ValidateBookEnd(t *testing.T) {
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
			name: "update in progress book to end",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mock.NewMockRepostory(ctrl)
				rpo.EXPECT().UpdateBook(&model.Book{
					ID: 1, Status: model.End, UpdateChapter: "結尾", IsDownloaded: false,
				})

				return ServiceImp{rpo: rpo}
			},
			bk:           &model.Book{ID: 1, Status: model.InProgress, UpdateChapter: "結尾", IsDownloaded: true},
			wantBook:     &model.Book{ID: 1, Status: model.End, UpdateChapter: "結尾", IsDownloaded: false},
			wantError:    false,
			wantErrorStr: "",
		},
		{
			name: "update end book to in progress",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mock.NewMockRepostory(ctrl)
				rpo.EXPECT().UpdateBook(&model.Book{
					ID: 1, Status: model.InProgress, UpdateChapter: "中間", IsDownloaded: true,
				})

				return ServiceImp{rpo: rpo}
			},
			bk:           &model.Book{ID: 1, UpdateChapter: "中間", Status: model.End, IsDownloaded: true},
			wantBook:     &model.Book{ID: 1, Status: model.InProgress, UpdateChapter: "中間", IsDownloaded: true},
			wantError:    false,
			wantErrorStr: "",
		},
		{
			name: "do nothing on end book with end chapter",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				return ServiceImp{}
			},
			bk:           &model.Book{ID: 1, UpdateChapter: "結尾", Status: model.End, IsDownloaded: false},
			wantBook:     &model.Book{ID: 1, Status: model.End, UpdateChapter: "結尾", IsDownloaded: false},
			wantError:    false,
			wantErrorStr: "",
		},
		{
			name: "do nothing on non end book with non end chapter",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				return ServiceImp{}
			},
			bk:           &model.Book{ID: 1, UpdateChapter: "中間", Status: model.InProgress, IsDownloaded: true},
			wantBook:     &model.Book{ID: 1, Status: model.InProgress, UpdateChapter: "中間", IsDownloaded: true},
			wantError:    false,
			wantErrorStr: "",
		},
		{
			name: "update book fail",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mock.NewMockRepostory(ctrl)
				rpo.EXPECT().UpdateBook(&model.Book{
					ID: 1, Status: model.InProgress, UpdateChapter: "中間", IsDownloaded: true,
				}).Return(errors.New("some error"))

				return ServiceImp{rpo: rpo}
			},
			bk:           &model.Book{ID: 1, UpdateChapter: "中間", Status: model.End, IsDownloaded: true},
			wantBook:     &model.Book{ID: 1, Status: model.InProgress, UpdateChapter: "中間", IsDownloaded: true},
			wantError:    true,
			wantErrorStr: "update book in DB fail: some error",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := test.setupServ(ctrl)
			err := serv.ValidateBookEnd(test.bk)

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

func TestServiceImp_ValidateEnd(t *testing.T) {
	t.Parallel()

	conf := config.SiteConfig{
		BackupDirectory: "some dir",
	}

	tests := []struct {
		name        string
		setupServ   func(*gomock.Controller) ServiceImp
		expectError error
	}{
		{
			name: "calls rpo.UpdateBooksStatus",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mock.NewMockRepostory(ctrl)
				rpo.EXPECT().UpdateBooksStatus()

				return ServiceImp{
					conf: conf,
					rpo:  rpo,
				}
			},
			expectError: nil,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := test.setupServ(ctrl)
			var op SiteOperation = serv.ValidateEnd

			err := op()
			if !errors.Is(err, test.expectError) {
				t.Errorf("error diff:\n%v\n%v", err, test.expectError)
			}
		})
	}

}
