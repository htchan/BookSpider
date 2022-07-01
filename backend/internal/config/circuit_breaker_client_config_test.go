package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadClientConfigs(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() {
		os.Remove("client_test")
	})

	tests := []struct {
		name             string
		yamlClientString string
		want             map[string]*CircuitBreakerClientConfig
		wantErr          bool
	}{
		{
			name: "load circuit breaker client success",
			yamlClientString: `testClientKey:
        maxFailCount: 1
        maxFailMultiplier: 0.5
        circuitBreakingSleep: 2
        intervalSleep: 3
        timeout: 4
        retry503: 5
        retryErr: 6
        decoder:
          method: big5`,
			want: map[string]*CircuitBreakerClientConfig{
				"testClientKey": &CircuitBreakerClientConfig{
					MaxFailCount:         1,
					MaxFailMultiplier:    0.5,
					CircuitBreakingSleep: 2,
					IntervalSleep:        3,
					Timeout:              4,
					Retry503:             5,
					RetryErr:             6,
					DecoderConfig:        DecoderConfig{Method: "big5"},
				},
			},
			wantErr: false,
		},
	}

	for i, test := range tests {
		test := test
		dirName := fmt.Sprintf("client_test/%d", i)
		os.MkdirAll(dirName, 0750)
		os.WriteFile(fmt.Sprintf("./%v/client_configs.yaml", dirName), []byte(test.yamlClientString), 0644)
		t.Run(test.name, func(t *testing.T) {
			t.Cleanup(func() {
				os.Remove(fmt.Sprintf("./%v/client_configs.yaml", dirName))
				os.Remove(dirName)
			})

			got, err := LoadClientConfigs(dirName)
			if (err != nil) != test.wantErr {
				t.Errorf("LoadClientConfigs() return error %v, wantErr %v", err, test.wantErr)
			}

			if !cmp.Equal(got, test.want) {
				t.Error(cmp.Diff(got, test.want))
			}
		})
	}
}
