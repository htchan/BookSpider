package service

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	config "github.com/htchan/BookSpider/internal/config_new"
	"github.com/htchan/BookSpider/internal/mock"
)

func TestServiceImp_Backup(t *testing.T) {
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
			name: "calls rpo.Backup",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				rpo := mock.NewMockRepostory(ctrl)
				rpo.EXPECT().Backup("some dir")

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
			var op SiteOperation = serv.Backup

			err := op()
			if !errors.Is(err, test.expectError) {
				t.Errorf("error diff:\n%v\n%v", err, test.expectError)
			}
		})
	}

}
