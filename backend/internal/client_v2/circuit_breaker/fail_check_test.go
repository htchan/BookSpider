package circuitbreaker

import (
	"testing"

	client "github.com/htchan/BookSpider/internal/client_v2"
	"github.com/stretchr/testify/assert"
)

func Test_newFailCheck(t *testing.T) {
	t.Parallel()

	statusCodeConf := CheckConfig{
		Type:  CheckTypeStatusCodes,
		Value: []int{502},
	}

	type args struct {
		body string
		err  error
	}

	tests := []struct {
		name string
		args args
		conf CheckConfig
		want bool
	}{
		{
			name: "status codes check/no error",
			args: args{body: "test", err: nil},
			conf: statusCodeConf,
			want: false,
		},
		{
			name: "status codes check/no match status codes list",
			args: args{body: "", err: client.StatusCodeError{StatusCode: 500}},
			conf: statusCodeConf,
			want: false,
		},
		{
			name: "status codes check/match status codes list",
			args: args{body: "", err: client.StatusCodeError{StatusCode: 502}},
			conf: statusCodeConf,
			want: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			check := newFailCheck(test.conf)
			got := check(test.args.body, test.args.err)
			assert.Equal(t, test.want, got)
		})
	}
}
