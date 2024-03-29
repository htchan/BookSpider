package service

import (
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/htchan/BookSpider/internal/config/v2"
	mockclient "github.com/htchan/BookSpider/internal/mock/client/v2"
)

func TestServiceImp_CheckAvailability(t *testing.T) {
	t.Parallel()

	conf := config.SiteConfig{
		AvailabilityConfig: config.AvailabilityConfig{
			URL:         "test",
			CheckString: "123",
		},
	}

	tests := []struct {
		name        string
		setupServ   func(*gomock.Controller) ServiceImp
		expectError error
	}{
		{
			name: "client response contains check string",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				client := mockclient.NewMockBookClient(ctrl)
				client.EXPECT().Get(gomock.Any(), "test").Return("abc123def", nil)

				return ServiceImp{
					conf:   conf,
					client: client,
				}
			},
			expectError: nil,
		},
		{
			name: "client response not contains check string",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				client := mockclient.NewMockBookClient(ctrl)
				client.EXPECT().Get(gomock.Any(), "test").Return("abcdef", nil)

				return ServiceImp{
					conf:   conf,
					client: client,
				}
			},
			expectError: ErrInvalidSite,
		},
		{
			name: "client return error",
			setupServ: func(ctrl *gomock.Controller) ServiceImp {
				client := mockclient.NewMockBookClient(ctrl)
				client.EXPECT().Get(gomock.Any(), "test").Return("", http.ErrNotSupported)

				return ServiceImp{
					conf:   conf,
					client: client,
				}
			},
			expectError: http.ErrNotSupported,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := test.setupServ(ctrl)
			var op SiteOperation = serv.CheckAvailability

			err := op()
			if !errors.Is(err, test.expectError) {
				t.Errorf("error diff:\n%v\n%v", err, test.expectError)
			}
		})
	}

}
