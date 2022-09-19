package site

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/htchan/BookSpider/internal/config"
)

func Test_NewSite(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		stName     string
		bkConf     config.BookConfig
		stConf     config.SiteConfig
		clientConf config.CircuitBreakerClientConfig
		expect     Site
		expectErr  bool
	}{
		{
			name:       "works",
			stName:     "test",
			bkConf:     config.BookConfig{},
			stConf:     config.SiteConfig{},
			clientConf: config.CircuitBreakerClientConfig{},
			expect: Site{
				Name: "test",
			},
			expectErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result, err := NewSite(test.stName, test.bkConf, test.stConf, test.clientConf, nil, nil)
			if (err != nil) != test.expectErr {
				t.Errorf("got error: %v; want err: %v", err, test.expectErr)
			}

			if !cmp.Equal(*result, test.expect) {
				t.Errorf("site diff: %v", cmp.Diff(result, test.expect))
			}
		})
	}
}
