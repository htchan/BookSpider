package config

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Test_LoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		envMap         map[string]string
		configLocation string
		configContent  string
		expectedConf   *Config
		expectError    bool
	}{
		{
			name: "happy flow with default",
			envMap: map[string]string{
				"CONFIG_LOCATION":       "./test.yaml",
				"API_AVAILABLE_SITES":   "xbiquge,xqishu",
				"MAX_WORKING_THREADS":   "1000",
				"BATCH_AVAILABLE_SITES": "xbiquge,xqishu",
				"PSQL_HOST":             "host",
				"PSQL_PORT":             "12345",
				"PSQL_USER":             "user",
				"PSQL_PASSWORD":         "password",
				"PSQL_NAME":             "name",
			},
			configLocation: "./test.yaml",
			configContent: `
				xbiquge:
					decode_method: gbk
					max_threads: 1000
					request_timeout: 30s
					circuit_breaker:
						max_fail_count: 1000
						max_fail_multiplier: 1.5
						sleep_interval: 10s
					retry_map:
						default: 10
						unavailable: 100

					storage: /storage
					backup_directory: /backup

					urls: #desktop
						base: https://www.xbiquge.so/book/%%v/
						download: https://www.xbiquge.so/book/%%v/
						chapter_prefix: https://www.xbiquge.so
					max_explore_error: 100
					max_download_concurrency: 10
					update_date_layout: null
					goquery_selectors:
						title:
							selector: title
							unwanted_content:
								- a
								- b
						writer:
							selector: writer
							attr:
						book_type:
							selector: book-type
							attr:
						update_date:
							selector: update-date
							attr:
						update_chapter:
							selector: update-chapter
							attr:
						book_chapter_url: 
							selector: book-chapter
							attr: href
						book_chapter_title: 
							selector: book-chapter
							attr: 
						chapter_title: 
							selector: chapter-title
							attr:
						chapter_content:
							selector: chapter-content
							attr:
					availability:
						url: availability
						check_string: check

				xqishu:
					decode_method: utf8
					max_threads: 200
					request_timeout: 30s
					circuit_breaker:
						max_fail_count: 10
						max_fail_multiplier: 2
						sleep_interval: 5s
					retry_map:
						default: 10
						unavailable: 100

					storage: /storage
					backup_directory: /backup

					urls: #desktop
						base: http://www.aidusk.com/txt%%v/
						download: http://www.aidusk.com/t/%%v/
						chapter_prefix: http://www.aidusk.com
					max_explore_error: 100
					max_download_concurrency: 10
					update_date_layout: null
					goquery_selectors:
						title:
							selector: title
							attr:
						writer:
							selector: writer
							attr:
						book_type:
							selector: book-type
							attr:
						update_date:
							selector: update-date
							attr:
						update_chapter:
							selector: update-chapter
							attr:
						book_chapter_url: 
							selector: book-chapter
							attr: href
						book_chapter_title: 
							selector: book-chapter
							attr: 
						chapter_title: 
							selector: chapter-title
							attr:
						chapter_content:
							selector: chapter-content
							attr:
					availability:
						url: availability
						check_string: check`,
			expectedConf: &Config{
				APIConfig: APIConfig{
					APIRoutePrefix:     "/api/novel",
					LiteRoutePrefix:    "/lite/novel",
					AvailableSiteNames: []string{"xbiquge", "xqishu"},
				},
				BatchConfig: BatchConfig{
					MaxWorkingThreads:  1000,
					AvailableSiteNames: []string{"xbiquge", "xqishu"},
				},
				SiteConfigs: map[string]SiteConfig{
					"xbiquge": {
						DecodeMethod: "gbk",
						MaxThreads:   1000,
						CircuitBreakerConfig: CircuitBreakerClientConfig{
							MaxFailCount:      1000,
							MaxFailMultiplier: 1.5,
							SleepInterval:     10 * time.Second,
						},
						RequestTimeout: 30 * time.Second,
						RetryConfig: map[string]int{
							"unavailable": 100,
							"default":     10,
						},
						Storage:         "/storage",
						BackupDirectory: "/backup",
						URL: URLConfig{
							Base:          "https://www.xbiquge.so/book/%%v/",
							Download:      "https://www.xbiquge.so/book/%%v/",
							ChapterPrefix: "https://www.xbiquge.so",
						},
						MaxExploreError:        100,
						MaxDownloadConcurrency: 10,
						GoquerySelectorsConfig: GoquerySelectorsConfig{
							Title:            GoquerySelectorConfig{"title", "", []string{"a", "b"}},
							Writer:           GoquerySelectorConfig{"writer", "", nil},
							BookType:         GoquerySelectorConfig{"book-type", "", nil},
							LastUpdate:       GoquerySelectorConfig{"update-date", "", nil},
							LastChapter:      GoquerySelectorConfig{"update-chapter", "", nil},
							BookChapterURL:   GoquerySelectorConfig{"book-chapter", "href", nil},
							BookChapterTitle: GoquerySelectorConfig{"book-chapter", "", nil},
							ChapterTitle:     GoquerySelectorConfig{"chapter-title", "", nil},
							ChapterContent:   GoquerySelectorConfig{"chapter-content", "", nil},
						},
						AvailabilityConfig: AvailabilityConfig{
							URL:         "availability",
							CheckString: "check",
						},
					},
					"xqishu": {
						DecodeMethod: "utf8",
						MaxThreads:   200,
						CircuitBreakerConfig: CircuitBreakerClientConfig{
							MaxFailCount:      10,
							MaxFailMultiplier: 2,
							SleepInterval:     5 * time.Second,
						},
						RequestTimeout: 30 * time.Second,
						RetryConfig: map[string]int{
							"unavailable": 100,
							"default":     10,
						},
						Storage:         "/storage",
						BackupDirectory: "/backup",
						URL: URLConfig{
							Base:          "http://www.aidusk.com/txt%%v/",
							Download:      "http://www.aidusk.com/t/%%v/",
							ChapterPrefix: "http://www.aidusk.com",
						},
						MaxExploreError:        100,
						MaxDownloadConcurrency: 10,
						GoquerySelectorsConfig: GoquerySelectorsConfig{
							Title:            GoquerySelectorConfig{"title", "", nil},
							Writer:           GoquerySelectorConfig{"writer", "", nil},
							BookType:         GoquerySelectorConfig{"book-type", "", nil},
							LastUpdate:       GoquerySelectorConfig{"update-date", "", nil},
							LastChapter:      GoquerySelectorConfig{"update-chapter", "", nil},
							BookChapterURL:   GoquerySelectorConfig{"book-chapter", "href", nil},
							BookChapterTitle: GoquerySelectorConfig{"book-chapter", "", nil},
							ChapterTitle:     GoquerySelectorConfig{"chapter-title", "", nil},
							ChapterContent:   GoquerySelectorConfig{"chapter-content", "", nil},
						},
						AvailabilityConfig: AvailabilityConfig{
							URL:         "availability",
							CheckString: "check",
						},
					},
				},
				DatabaseConfig: DatabaseConfig{
					Host:     "host",
					Port:     "12345",
					User:     "user",
					Password: "password",
					Name:     "name",
				},
				ConfigLocation: "./test.yaml",
			},
			expectError: false,
		},
		{
			name: "happy flow without default",
			envMap: map[string]string{
				"CONFIG_LOCATION":       "./test.yaml",
				"API_AVAILABLE_SITES":   "xbiquge",
				"MAX_WORKING_THREADS":   "1000",
				"BATCH_AVAILABLE_SITES": "xbiquge",
				"API_ROUTE_PREFIX":      "/api-novel",
				"LITE_ROUTE_PREFIX":     "/lite-novel",
				"PSQL_HOST":             "host",
				"PSQL_PORT":             "12345",
				"PSQL_USER":             "user",
				"PSQL_PASSWORD":         "password",
				"PSQL_NAME":             "name",
			},
			configLocation: "./test.yaml",
			configContent: `
				xbiquge:
					decode_method: gbk
					max_threads: 1000
					request_timeout: 30s
					circuit_breaker:
						max_fail_count: 1000
						max_fail_multiplier: 1.5
						sleep_interval: 10s
					retry_map:
						default: 10
						unavailable: 100

					storage: /storage
					backup_directory: /backup

					urls: #desktop
						base: https://www.xbiquge.so/book/%%v/
						download: https://www.xbiquge.so/book/%%v/
						chapter_prefix: https://www.xbiquge.so
					max_explore_error: 100
					max_download_concurrency: 10
					update_date_layout: null
					availability:
						url: availability
						check_string: check`,
			expectedConf: &Config{
				APIConfig: APIConfig{
					APIRoutePrefix:     "/api-novel",
					LiteRoutePrefix:    "/lite-novel",
					AvailableSiteNames: []string{"xbiquge"},
				},
				BatchConfig: BatchConfig{
					MaxWorkingThreads:  1000,
					AvailableSiteNames: []string{"xbiquge"},
				},
				SiteConfigs: map[string]SiteConfig{
					"xbiquge": {
						DecodeMethod: "gbk",
						MaxThreads:   1000,
						CircuitBreakerConfig: CircuitBreakerClientConfig{
							MaxFailCount:      1000,
							MaxFailMultiplier: 1.5,
							SleepInterval:     10 * time.Second,
						},
						RequestTimeout: 30 * time.Second,
						RetryConfig: map[string]int{
							"unavailable": 100,
							"default":     10,
						},
						Storage:         "/storage",
						BackupDirectory: "/backup",
						URL: URLConfig{
							Base:          "https://www.xbiquge.so/book/%%v/",
							Download:      "https://www.xbiquge.so/book/%%v/",
							ChapterPrefix: "https://www.xbiquge.so",
						},
						MaxExploreError:        100,
						MaxDownloadConcurrency: 10,
						AvailabilityConfig: AvailabilityConfig{
							URL:         "availability",
							CheckString: "check",
						},
					},
				},
				DatabaseConfig: DatabaseConfig{
					Host:     "host",
					Port:     "12345",
					User:     "user",
					Password: "password",
					Name:     "name",
				},
				ConfigLocation: "./test.yaml",
			},
			expectError: false,
		},
		{
			name: "empty sites",
			envMap: map[string]string{
				"CONFIG_LOCATION":       "./test.yaml",
				"API_AVAILABLE_SITES":   "xbiquge",
				"MAX_WORKING_THREADS":   "1000",
				"BATCH_AVAILABLE_SITES": "xbiquge",
				"PSQL_HOST":             "host",
				"PSQL_PORT":             "12345",
				"PSQL_USER":             "user",
				"PSQL_PASSWORD":         "password",
				"PSQL_NAME":             "name",
			},
			configLocation: "./test.yaml",
			configContent:  ``,
			expectedConf: &Config{
				APIConfig: APIConfig{
					APIRoutePrefix:     "/api/novel",
					LiteRoutePrefix:    "/lite/novel",
					AvailableSiteNames: []string{"xbiquge"},
				},
				BatchConfig: BatchConfig{
					MaxWorkingThreads:  1000,
					AvailableSiteNames: []string{"xbiquge"},
				},
				DatabaseConfig: DatabaseConfig{
					Host:     "host",
					Port:     "12345",
					User:     "user",
					Password: "password",
					Name:     "name",
				},
				SiteConfigs:    nil,
				ConfigLocation: "./test.yaml",
			},
			expectError: false,
		},
		{
			name: "sites config file not exist",
			envMap: map[string]string{
				"CONFIG_LOCATION":       "./test-not-exist.yaml",
				"API_AVAILABLE_SITES":   "xbiquge",
				"MAX_WORKING_THREADS":   "1000",
				"BATCH_AVAILABLE_SITES": "xbiquge",
				"PSQL_HOST":             "host",
				"PSQL_PORT":             "12345",
				"PSQL_USER":             "user",
				"PSQL_PASSWORD":         "password",
				"PSQL_NAME":             "name",
			},
			configLocation: "./test.yaml",
			configContent:  ``,
			expectedConf:   nil,
			expectError:    true,
		},
		{
			name:           "fail to provide required env",
			envMap:         map[string]string{},
			configLocation: "./test.yaml",
			configContent:  ``,
			expectedConf:   nil,
			expectError:    true,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			// populate env
			for key, value := range test.envMap {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			configContent := strings.ReplaceAll(test.configContent, "\t", "  ")
			os.WriteFile(test.configLocation, []byte(configContent), 0644)
			defer os.Remove(test.configLocation)

			conf, err := LoadConfig()
			if (err != nil) != test.expectError {
				t.Error("error diff")
				t.Errorf("expected error exist: %v", test.expectError)
				t.Errorf("actual: %v", err)
			}

			if !cmp.Equal(test.expectedConf, conf) {
				t.Errorf("conf diff: %v", cmp.Diff(test.expectedConf, conf))
			}
		})
	}
}
