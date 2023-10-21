package simple

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	client "github.com/htchan/BookSpider/internal/client/v2"
)

func Test_validate_SimpleConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		conf  SimpleClientConfig
		valid bool
	}{
		{
			name: "valie config",
			conf: SimpleClientConfig{
				RequestTimeout: 1 * time.Second,
				DecodeMethod:   client.DecodeMethodGBK,
			},
			valid: true,
		},
		{
			name: "invalid request timeout",
			conf: SimpleClientConfig{
				RequestTimeout: 999 * time.Millisecond,
				DecodeMethod:   client.DecodeMethodBig5,
			},
			valid: false,
		},
		{
			name: "invalid decode method",
			conf: SimpleClientConfig{
				RequestTimeout: 1 * time.Second,
				DecodeMethod:   "invalid",
			},
			valid: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := validator.New().Struct(test.conf)
			if test.valid && err != nil {
				t.Errorf("unexpected error: %s", err)
			}
		})
	}
}
