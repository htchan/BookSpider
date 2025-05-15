package service

import (
	"strconv"
	"testing"
	"time"

	repomock "github.com/htchan/BookSpider/internal/mock/repo"
	"github.com/htchan/BookSpider/internal/model"
	serv "github.com/htchan/BookSpider/internal/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_isEnd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		bk   *model.Book
		want bool
	}{
		{
			name: "book is end due to update date",
			bk:   &model.Book{UpdateDate: strconv.Itoa(time.Now().Year() - 3)},
			want: true,
		},
		{
			name: "book is end due to update chapter",
			bk:   &model.Book{UpdateDate: strconv.Itoa(time.Now().Year() - 3), UpdateChapter: "番外"},
			want: true,
		},
		{
			name: "book is not end",
			bk:   &model.Book{UpdateDate: strconv.Itoa(time.Now().Year()), UpdateChapter: ""},
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

func TestServiceImpl_ValidateBookEnd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		getServ   func(ctrl *gomock.Controller) *ServiceImpl
		bk        *model.Book
		wantBk    *model.Book
		wantError error
	}{
		{
			name: "book is end, but its status is not",
			getServ: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := repomock.NewMockRepository(ctrl)
				rpo.EXPECT().UpdateBook(gomock.Any(), &model.Book{
					UpdateDate: strconv.Itoa(time.Now().Year() - 3),
					Status:     model.StatusEnd,
				}).Return(nil)

				return &ServiceImpl{rpo: rpo}
			},
			bk:        &model.Book{UpdateDate: strconv.Itoa(time.Now().Year() - 3), Status: model.StatusInProgress},
			wantBk:    &model.Book{UpdateDate: strconv.Itoa(time.Now().Year() - 3), Status: model.StatusEnd},
			wantError: nil,
		},
		{
			name: "book is end, and its status is also end",
			getServ: func(ctrl *gomock.Controller) *ServiceImpl {
				return &ServiceImpl{}
			},
			bk:        &model.Book{UpdateDate: strconv.Itoa(time.Now().Year() - 3), Status: model.StatusEnd},
			wantBk:    &model.Book{UpdateDate: strconv.Itoa(time.Now().Year() - 3), Status: model.StatusEnd},
			wantError: nil,
		},
		{
			name: "book is not end and its status is also not end",
			getServ: func(ctrl *gomock.Controller) *ServiceImpl {
				return &ServiceImpl{}
			},
			bk:        &model.Book{UpdateDate: strconv.Itoa(time.Now().Year()), Status: model.StatusInProgress},
			wantBk:    &model.Book{UpdateDate: strconv.Itoa(time.Now().Year()), Status: model.StatusInProgress},
			wantError: nil,
		},
		{
			name: "book is not end, but its status is end",
			getServ: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := repomock.NewMockRepository(ctrl)
				rpo.EXPECT().UpdateBook(gomock.Any(), &model.Book{
					UpdateDate: strconv.Itoa(time.Now().Year()),
					Status:     model.StatusInProgress,
				}).Return(nil)

				return &ServiceImpl{rpo: rpo}
			},
			bk:        &model.Book{UpdateDate: strconv.Itoa(time.Now().Year()), Status: model.StatusEnd},
			wantBk:    &model.Book{UpdateDate: strconv.Itoa(time.Now().Year()), Status: model.StatusInProgress},
			wantError: nil,
		},
		{
			name: "update book return error",
			getServ: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := repomock.NewMockRepository(ctrl)
				rpo.EXPECT().UpdateBook(gomock.Any(), &model.Book{
					UpdateDate: strconv.Itoa(time.Now().Year()),
					Status:     model.StatusInProgress,
				}).Return(serv.ErrUnavailable)

				return &ServiceImpl{rpo: rpo}
			},
			bk:        &model.Book{UpdateDate: strconv.Itoa(time.Now().Year()), Status: model.StatusEnd},
			wantBk:    &model.Book{UpdateDate: strconv.Itoa(time.Now().Year()), Status: model.StatusInProgress},
			wantError: serv.ErrUnavailable,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			err := test.getServ(ctrl).ValidateBookEnd(t.Context(), test.bk)
			assert.Equal(t, test.wantBk, test.bk)
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}

func TestServiceImpl_ValidateEnd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		getServ   func(ctrl *gomock.Controller) *ServiceImpl
		wantError error
	}{
		{
			name: "call rpo function",
			getServ: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := repomock.NewMockRepository(ctrl)
				rpo.EXPECT().UpdateBooksStatus(gomock.Any()).Return(serv.ErrUnavailable)

				return &ServiceImpl{rpo: rpo}
			},
			wantError: serv.ErrUnavailable,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			err := test.getServ(ctrl).ValidateEnd(t.Context())
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}
