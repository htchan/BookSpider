package service

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/htchan/BookSpider/internal/config/v2"
	mockrepo "github.com/htchan/BookSpider/internal/mock/repo"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/stretchr/testify/assert"
)

func Test_isEnd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		bk   *model.Book
		want bool
	}{
		{
			name: "book with ended chapter",
			bk:   &model.Book{UpdateChapter: "last chatper （完）", UpdateDate: strconv.Itoa(time.Now().Year())},
			want: true,
		},
		{
			name: "book with 1 yr ago update date",
			bk:   &model.Book{UpdateChapter: "chapter in middle", UpdateDate: strconv.Itoa(time.Now().Year() - 1)},
			want: false,
		},
		{
			name: "book with 2 yr ago update date",
			bk:   &model.Book{UpdateChapter: "chapter in middle", UpdateDate: strconv.Itoa(time.Now().Year() - 2)},
			want: true,
		},
		{
			name: "book with not ended chapter",
			bk:   &model.Book{UpdateChapter: "chapter in middle", UpdateDate: strconv.Itoa(time.Now().Year())},
			want: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := isEnd(test.bk)
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
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().UpdateBook(&model.Book{
					ID: 1, Status: model.StatusEnd, UpdateChapter: "結尾", IsDownloaded: false, UpdateDate: strconv.Itoa(time.Now().Year() - 2),
				})

				return ServiceImp{rpo: rpo}
			},
			bk:           &model.Book{ID: 1, Status: model.StatusInProgress, UpdateChapter: "結尾", IsDownloaded: true, UpdateDate: strconv.Itoa(time.Now().Year() - 2)},
			wantBook:     &model.Book{ID: 1, Status: model.StatusEnd, UpdateChapter: "結尾", IsDownloaded: false, UpdateDate: strconv.Itoa(time.Now().Year() - 2)},
			wantError:    false,
			wantErrorStr: "",
		},
		{
			name: "update end book to in progress",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().UpdateBook(&model.Book{
					ID: 1, Status: model.StatusInProgress, UpdateChapter: "中間", IsDownloaded: true, UpdateDate: strconv.Itoa(time.Now().Year()),
				})

				return ServiceImp{rpo: rpo}
			},
			bk:           &model.Book{ID: 1, UpdateChapter: "中間", Status: model.StatusEnd, IsDownloaded: true, UpdateDate: strconv.Itoa(time.Now().Year())},
			wantBook:     &model.Book{ID: 1, Status: model.StatusInProgress, UpdateChapter: "中間", IsDownloaded: true, UpdateDate: strconv.Itoa(time.Now().Year())},
			wantError:    false,
			wantErrorStr: "",
		},
		{
			name: "do nothing on end book with end chapter",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				return ServiceImp{}
			},
			bk:           &model.Book{ID: 1, UpdateChapter: "結尾", Status: model.StatusEnd, IsDownloaded: false},
			wantBook:     &model.Book{ID: 1, Status: model.StatusEnd, UpdateChapter: "結尾", IsDownloaded: false},
			wantError:    false,
			wantErrorStr: "",
		},
		{
			name: "do nothing on non end book with non end chapter",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				return ServiceImp{}
			},
			bk:           &model.Book{ID: 1, UpdateChapter: "中間", Status: model.StatusInProgress, IsDownloaded: true, UpdateDate: strconv.Itoa(time.Now().Year())},
			wantBook:     &model.Book{ID: 1, Status: model.StatusInProgress, UpdateChapter: "中間", IsDownloaded: true, UpdateDate: strconv.Itoa(time.Now().Year())},
			wantError:    false,
			wantErrorStr: "",
		},
		{
			name: "update book fail",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().UpdateBook(&model.Book{
					ID: 1, Status: model.StatusInProgress, UpdateChapter: "中間", IsDownloaded: true, UpdateDate: strconv.Itoa(time.Now().Year()),
				}).Return(errors.New("some error"))

				return ServiceImp{rpo: rpo}
			},
			bk:           &model.Book{ID: 1, UpdateChapter: "中間", Status: model.StatusEnd, IsDownloaded: true, UpdateDate: strconv.Itoa(time.Now().Year())},
			wantBook:     &model.Book{ID: 1, Status: model.StatusInProgress, UpdateChapter: "中間", IsDownloaded: true, UpdateDate: strconv.Itoa(time.Now().Year())},
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
				rpo := mockrepo.NewMockRepository(ctrl)
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
