package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadBackendConfigs(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() {
		os.Remove("backend_test")
	})

	tests := []struct {
		name              string
		yamlSiteString    string
		yamlBookString    string
		yamlClientString  string
		yamlBackendString string
		want              BackendConfig
		wantErr           bool
	}{
		{
			name: "load site config success",
			yamlSiteString: `test:
        sourceKey: testSourceKey
        clientKey: testClientKey
        storage: /book
        backupDirectory: /backup
        concurrency:
          fetchThreads: 1
          downloadThreads: 2
        database:
          engine: engine
          conn: conn
          timeout: 3
          idleConns: 4
          openConns: 5
        maxExploreError: 6
        availability:
          url: url
          check: check`,
			yamlBookString: `testSourceKey:
        url:
          base: base url
          download: download url
          chapterPrefix: chapter prefix url
        maxChapterError: 1
        updateDateLayout: layout`,
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
			yamlBackendString: `
      enabledRoutes:
        a: false
        b: true
      enabledSites:
        - c
        - d`,
			want: BackendConfig{
				EnabledRoutes:    map[string]bool{"a": false, "b": true},
				EnabledSiteNames: []string{"c", "d"},
				SiteConfigs: map[string]*SiteConfig{
					"test": &SiteConfig{
						SourceKey:         "testSourceKey",
						ClientKey:         "testClientKey",
						Storage:           "/book",
						BackupDirectory:   "/backup",
						ConcurrencyConfig: ConcurrencyConfig{FetchThreads: 1, DownloadThreads: 2},
						DatabaseConfig: DatabaseConfig{
							Engine:    "engine",
							Conn:      "conn",
							Timeout:   3,
							IdleConns: 4,
							OpenConns: 5,
						},
						MaxExploreError:    6,
						AvailabilityConfig: AvailabilityConfig{URL: "url", Check: "check"},
						BookConfig: &BookConfig{
							URLConfig:        URLConfig{Base: "base url", Download: "download url", ChapterPrefix: "chapter prefix url"},
							MaxChaptersError: 1,
							UpdateDateLayout: "layout",
							Storage:          "/book",
							SourceKey:        "testSourceKey",
						},
						CircuitBreakerClientConfig: &CircuitBreakerClientConfig{
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
				},
			},
			wantErr: false,
		},
	}

	for i, test := range tests {
		test := test
		dirName := fmt.Sprintf("backend_test/%d", i)
		os.MkdirAll(dirName, 0750)
		os.WriteFile(fmt.Sprintf("./%v/site_configs.yaml", dirName), []byte(test.yamlSiteString), 0644)
		os.WriteFile(fmt.Sprintf("./%v/book_configs.yaml", dirName), []byte(test.yamlBookString), 0644)
		os.WriteFile(fmt.Sprintf("./%v/client_configs.yaml", dirName), []byte(test.yamlClientString), 0644)
		os.WriteFile(fmt.Sprintf("./%v/backend_config.yaml", dirName), []byte(test.yamlBackendString), 0644)
		t.Run(test.name, func(t *testing.T) {
			t.Cleanup(func() {
				os.Remove(fmt.Sprintf("./%v/site_configs.yaml", dirName))
				os.Remove(fmt.Sprintf("./%v/book_configs.yaml", dirName))
				os.Remove(fmt.Sprintf("./%v/client_configs.yaml", dirName))
				os.Remove(fmt.Sprintf("./%v/backend_config.yaml", dirName))
				os.Remove(dirName)
			})

			got, err := LoadBackendConfig(dirName)
			if (err != nil) != test.wantErr {
				t.Errorf("LoadBackendConfig() return error %v, wantErr %v", err, test.wantErr)
			}

			if !cmp.Equal(got, test.want) {
				t.Error(cmp.Diff(got, test.want))
				t.Error(got)
				t.Error(test.want)
			}
		})
	}
}
