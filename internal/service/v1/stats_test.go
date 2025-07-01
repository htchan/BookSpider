package service

import (
	"database/sql"
	"testing"

	mockrepo "github.com/htchan/BookSpider/internal/mock/repo"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestServiceImpl_Stats(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		getService func(*gomock.Controller) *ServiceImpl
		want       repo.Summary
	}{
		{
			name: "happy flow",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().Stats(gomock.Any(), "test").Return(repo.Summary{})

				return &ServiceImpl{name: "test", rpo: rpo}
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

			got := svc.Stats(t.Context())
			assert.Equal(t, test.want, got)
		})
	}
}

func TestServiceImpl_DBStats(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		getService func(*gomock.Controller) *ServiceImpl
		want       sql.DBStats
	}{
		{
			name: "happy flow",
			getService: func(ctrl *gomock.Controller) *ServiceImpl {
				rpo := mockrepo.NewMockRepository(ctrl)
				rpo.EXPECT().DBStats(gomock.Any()).Return(sql.DBStats{})

				return &ServiceImpl{rpo: rpo}
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
