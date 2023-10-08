package service

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	mockrepo "github.com/htchan/BookSpider/internal/mock/repo"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestServiceImp_Book(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupServ    func(ctrl *gomock.Controller) ServiceImp
		id           int
		hash         string
		wantBook     *model.Book
		wantError    bool
		wantErrorStr string
	}{
		{
			name: "find book with empty hash",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBookById(1).Return(&model.Book{
					ID: 1,
				}, nil)

				return ServiceImp{
					rpo: rpo,
				}
			},
			id:           1,
			hash:         "",
			wantBook:     &model.Book{ID: 1},
			wantError:    false,
			wantErrorStr: "",
		},
		{
			name: "find book with non empty hash",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBookByIdHash(1, 100).Return(&model.Book{
					ID:       1,
					HashCode: 100,
				}, nil)

				return ServiceImp{
					rpo: rpo,
				}
			},
			id:           1,
			hash:         "2s",
			wantBook:     &model.Book{ID: 1, HashCode: 100},
			wantError:    false,
			wantErrorStr: "",
		},
		{
			name: "getting error in find book",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBookById(1).Return(nil, errors.New("find book by id error"))

				return ServiceImp{
					rpo: rpo,
				}
			},
			id:           1,
			hash:         "",
			wantBook:     nil,
			wantError:    true,
			wantErrorStr: "find book by id error",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := test.setupServ(ctrl)
			bk, err := serv.Book(test.id, test.hash)

			if test.wantError {
				assert.Error(t, err)
			}

			if err != nil {
				assert.EqualError(t, err, test.wantErrorStr)
			}

			assert.Equal(t, test.wantBook, bk)
		})
	}
}

func TestServiceImp_QueryBooks(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		setupServ    func(ctrl *gomock.Controller) ServiceImp
		title        string
		writer       string
		limit        int
		offset       int
		wantBooks    []model.Book
		wantError    bool
		wantErrorStr string
	}{
		{
			name: "happy flow",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBooksByTitleWriter("title", "writer", 1, 1).Return([]model.Book{
					{ID: 1, Title: "title", Writer: model.Writer{Name: "somebody"}},
					{ID: 1, Title: "some text", Writer: model.Writer{Name: "writer"}},
				}, nil)

				return ServiceImp{rpo: rpo}
			},
			title:  "title",
			writer: "writer",
			limit:  1,
			offset: 1,
			wantBooks: []model.Book{
				{ID: 1, Title: "title", Writer: model.Writer{Name: "somebody"}},
				{ID: 1, Title: "some text", Writer: model.Writer{Name: "writer"}},
			},
			wantError:    false,
			wantErrorStr: "",
		},
		{
			name: "getting error",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBooksByTitleWriter("title", "writer", 1, 1).Return(nil, errors.New("some error"))

				return ServiceImp{rpo: rpo}
			},
			title:        "title",
			writer:       "writer",
			limit:        1,
			offset:       1,
			wantBooks:    nil,
			wantError:    true,
			wantErrorStr: "some error",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := test.setupServ(ctrl)
			bks, err := serv.QueryBooks(test.title, test.writer, test.limit, test.offset)

			if test.wantError {
				assert.Error(t, err)
			}

			if err != nil {
				assert.EqualError(t, err, test.wantErrorStr)
			}

			assert.Equal(t, test.wantBooks, bks)
		})
	}
}

func TestServiceImp_RandomBooks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupServ    func(ctrl *gomock.Controller) ServiceImp
		limit        int
		wantBooks    []model.Book
		wantError    bool
		wantErrorStr string
	}{
		{
			name: "happy flow",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBooksByRandom(10).Return([]model.Book{{ID: 1}}, nil)

				return ServiceImp{rpo: rpo}
			},
			limit:        10,
			wantBooks:    []model.Book{{ID: 1}},
			wantError:    false,
			wantErrorStr: "",
		},
		{
			name: "getting error",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().FindBooksByRandom(10).Return(nil, errors.New("some error"))

				return ServiceImp{rpo: rpo}
			},
			limit:        10,
			wantBooks:    nil,
			wantError:    true,
			wantErrorStr: "some error",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := test.setupServ(ctrl)
			bks, err := serv.RandomBooks(test.limit)

			if test.wantError {
				assert.Error(t, err)
			}

			if err != nil {
				assert.EqualError(t, err, test.wantErrorStr)
			}

			assert.Equal(t, test.wantBooks, bks)
		})
	}
}
