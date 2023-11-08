package config

import (
	"flag"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
	circuitbreaker "github.com/htchan/BookSpider/internal/client/v2/circuit_breaker"
	"github.com/htchan/BookSpider/internal/client/v2/retry"
	"github.com/htchan/BookSpider/internal/client/v2/simple"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "check for memory leaks")
	flag.Parse()

	if *leak {
		goleak.VerifyTestMain(m)
	} else {
		os.Exit(m.Run())
	}
}

func Test_validate_Config(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		conf  Config
		valid bool
	}{
		{
			name: "valid conf",
			conf: Config{
				APIConfig: APIConfig{
					APIRoutePrefix:     "/data",
					LiteRoutePrefix:    "/data",
					AvailableSiteNames: []string{"data"},
				},
				BatchConfig: BatchConfig{
					MaxWorkingThreads:  1,
					AvailableSiteNames: []string{"data"},
				},
				SiteConfigs: map[string]SiteConfig{},
				DatabaseConfig: DatabaseConfig{
					Host:     "host",
					Port:     "port",
					User:     "user",
					Password: "pwd",
					Name:     "name",
				},
				ConfigDirectory: ".",
			},
			valid: true,
		},
		{
			name: "invalid APIConfig",
			conf: Config{
				APIConfig: APIConfig{
					APIRoutePrefix:     "/",
					LiteRoutePrefix:    "/data",
					AvailableSiteNames: []string{"data"},
				},
				BatchConfig: BatchConfig{
					MaxWorkingThreads:  1,
					AvailableSiteNames: []string{"data"},
				},
				SiteConfigs: map[string]SiteConfig{},
				DatabaseConfig: DatabaseConfig{
					Host:     "host",
					Port:     "port",
					User:     "user",
					Password: "pwd",
					Name:     "name",
				},
				ConfigDirectory: ".",
			},
			valid: false,
		},
		{
			name: "invalid BatchConfig",
			conf: Config{
				APIConfig: APIConfig{
					APIRoutePrefix:     "/data",
					LiteRoutePrefix:    "/data",
					AvailableSiteNames: []string{"data"},
				},
				BatchConfig: BatchConfig{
					MaxWorkingThreads:  0,
					AvailableSiteNames: []string{"data"},
				},
				SiteConfigs: map[string]SiteConfig{},
				DatabaseConfig: DatabaseConfig{
					Host:     "host",
					Port:     "port",
					User:     "user",
					Password: "pwd",
					Name:     "name",
				},
				ConfigDirectory: ".",
			},
			valid: false,
		},
		{
			name: "invalid SiteConfig",
			conf: Config{
				APIConfig: APIConfig{
					APIRoutePrefix:     "/data",
					LiteRoutePrefix:    "/data",
					AvailableSiteNames: []string{"data"},
				},
				BatchConfig: BatchConfig{
					MaxWorkingThreads:  1,
					AvailableSiteNames: []string{"data"},
				},
				SiteConfigs: map[string]SiteConfig{"data": {}},
				DatabaseConfig: DatabaseConfig{
					Host:     "host",
					Port:     "port",
					User:     "user",
					Password: "pwd",
					Name:     "name",
				},
				ConfigDirectory: ".",
			},
			valid: false,
		},
		{
			name: "invalid DatabaseConfig",
			conf: Config{
				APIConfig: APIConfig{
					APIRoutePrefix:     "/data",
					LiteRoutePrefix:    "/data",
					AvailableSiteNames: []string{"data"},
				},
				BatchConfig: BatchConfig{
					MaxWorkingThreads:  1,
					AvailableSiteNames: []string{"data"},
				},
				SiteConfigs: map[string]SiteConfig{},
				DatabaseConfig: DatabaseConfig{
					Host:     "",
					Port:     "port",
					User:     "user",
					Password: "pwd",
					Name:     "name",
				},
				ConfigDirectory: ".",
			},
			valid: false,
		},
		{
			name: "invalid ConfigDirectory",
			conf: Config{
				APIConfig: APIConfig{
					APIRoutePrefix:     "/data",
					LiteRoutePrefix:    "/data",
					AvailableSiteNames: []string{"data"},
				},
				BatchConfig: BatchConfig{
					MaxWorkingThreads:  1,
					AvailableSiteNames: []string{"data"},
				},
				SiteConfigs: map[string]SiteConfig{},
				DatabaseConfig: DatabaseConfig{
					Host:     "host",
					Port:     "port",
					User:     "user",
					Password: "pwd",
					Name:     "name",
				},
				ConfigDirectory: "./not-exist/",
			},
			valid: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := validator.New().Struct(test.conf)
			if !assert.Equal(t, test.valid, err == nil) {
				t.Errorf("getting error: %v", err)
			}
		})
	}
}

func Test_validate_APIConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		conf  APIConfig
		valid bool
	}{
		{
			name: "valid conf",
			conf: APIConfig{
				APIRoutePrefix:     "/data",
				LiteRoutePrefix:    "/data",
				AvailableSiteNames: []string{"data"},
			},
			valid: true,
		},
		{
			name: "invalid APIRoutePrefix - non / prefix",
			conf: APIConfig{
				APIRoutePrefix:     "data",
				LiteRoutePrefix:    "/data",
				AvailableSiteNames: []string{"data"},
			},
			valid: false,
		},
		{
			name: "invalid APIRoutePrefix - / suffix",
			conf: APIConfig{
				APIRoutePrefix:     "/data/",
				LiteRoutePrefix:    "/data",
				AvailableSiteNames: []string{"data"},
			},
			valid: false,
		},
		{
			name: "invalid LiteRoutePrefix - non / prefix",
			conf: APIConfig{
				APIRoutePrefix:     "/data",
				LiteRoutePrefix:    "data",
				AvailableSiteNames: []string{"data"},
			},
			valid: false,
		},
		{
			name: "invalid LiteRoutePrefix - / suffix",
			conf: APIConfig{
				APIRoutePrefix:     "/data",
				LiteRoutePrefix:    "/data/",
				AvailableSiteNames: []string{"data"},
			},
			valid: false,
		},
		{
			name: "invalid AvailableSiteNames - empty",
			conf: APIConfig{
				APIRoutePrefix:     "/data",
				LiteRoutePrefix:    "/data",
				AvailableSiteNames: []string{},
			},
			valid: false,
		},
		{
			name: "invalid AvailableSiteNames - empty content",
			conf: APIConfig{
				APIRoutePrefix:     "/data",
				LiteRoutePrefix:    "/data",
				AvailableSiteNames: []string{""},
			},
			valid: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := validator.New().Struct(test.conf)
			if !assert.Equal(t, test.valid, err == nil) {
				t.Errorf("getting error: %v", err)
			}
		})
	}
}

func Test_validate_BatchConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		conf  BatchConfig
		valid bool
	}{
		{
			name: "valid conf",
			conf: BatchConfig{
				MaxWorkingThreads:  1,
				AvailableSiteNames: []string{"data"},
			},
			valid: true,
		},
		{
			name: "invalid MaxWorkingThreads - zero",
			conf: BatchConfig{
				MaxWorkingThreads:  0,
				AvailableSiteNames: []string{"data"},
			},
			valid: false,
		},
		{
			name: "invalid AvailableSiteNames - empty",
			conf: BatchConfig{
				MaxWorkingThreads:  1,
				AvailableSiteNames: []string{},
			},
			valid: false,
		},
		{
			name: "invalid AvailableSiteNames - empty content",
			conf: BatchConfig{
				MaxWorkingThreads:  1,
				AvailableSiteNames: []string{""},
			},
			valid: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := validator.New().Struct(test.conf)
			if !assert.Equal(t, test.valid, err == nil) {
				t.Errorf("getting error: %v", err)
			}
		})
	}
}

func Test_validate_DatabaseConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		conf  DatabaseConfig
		valid bool
	}{
		{
			name: "valid conf",
			conf: DatabaseConfig{
				Host:     "host",
				Port:     "port",
				User:     "user",
				Password: "pwd",
				Name:     "name",
			},
			valid: true,
		},
		{
			name: "invalid Host - empty",
			conf: DatabaseConfig{
				Host:     "",
				Port:     "port",
				User:     "user",
				Password: "pwd",
				Name:     "name",
			},
			valid: false,
		},
		{
			name: "invalid Port - empty",
			conf: DatabaseConfig{
				Host:     "host",
				Port:     "",
				User:     "user",
				Password: "pwd",
				Name:     "name",
			},
			valid: false,
		},
		{
			name: "invalid User - empty",
			conf: DatabaseConfig{
				Host:     "host",
				Port:     "port",
				User:     "",
				Password: "pwd",
				Name:     "name",
			},
			valid: false,
		},
		{
			name: "invalid Password - empty",
			conf: DatabaseConfig{
				Host:     "host",
				Port:     "port",
				User:     "user",
				Password: "",
				Name:     "name",
			},
			valid: false,
		},
		{
			name: "invalid Name - empty",
			conf: DatabaseConfig{
				Host:     "host",
				Port:     "port",
				User:     "user",
				Password: "pwd",
				Name:     "",
			},
			valid: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := validator.New().Struct(test.conf)
			if !assert.Equal(t, test.valid, err == nil) {
				t.Errorf("getting error: %v", err)
			}
		})
	}
}

func Test_LoadConfig(t *testing.T) {
	siteSelectorData1 := `xbiquge_selector: &xbiquge_selector
	decode_method: gbk
	urls: #desktop
		base: https://www.xbiquge.so/book/%%v/
		download: https://www.xbiquge.so/book/%%v/
		chapter_prefix: https://www.xbiquge.so
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
`
	siteSelectorData2 := `xqishu_selector: &xqishu_selector
	decode_method: utf8
	urls: #desktop
		base: http://www.aidusk.com/txt%%v/
		download: http://www.aidusk.com/t/%%v/
		chapter_prefix: http://www.aidusk.com
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
		check_string: check
`
	siteClientData1 := `default_circuit_breaker_config: &default_circuit_breaker_config
  open_threshold: 1000
  acquire_timeout: 500ms
  max_concurrency_threads: 1000
  recover_threads: [1, 2, 5, 10, 50, 100, 500]
  open_duration: 30s
  recover_duration: 10s
  check_configs:
  - type: status-codes
    value: [502]
default_retry_config: &default_retry_config
	max_retry_weight: 1000
	retry_conditions:
	- type: status-code
		value: [500, 502]
		weight: 10
		pause_interval: 1s
		pause_interval_type: exponential
	- type: body-contains
		value: []
		weight: 100
		pause_interval: 1s
		pause_interval_type: linear
	- type: error
		weight: 100
		pause_interval: 1s
		pause_interval_type: linear
xbiquge_client: &xbiquge_client
  # retry client
  retry: *default_retry_config
  # circuit breaker client
  circuit_breaker: *default_circuit_breaker_config
  # simple client
  simple:
    request_timeout: 30s
    decode_method: gbk
`
	siteClientData2 := `xqishu_client: &xqishu_client
  # retry client
  retry: *default_retry_config
  # circuit breaker client
  circuit_breaker: 
    <<: *default_circuit_breaker_config
    open_threshold: 10
    max_concurrency_threads: 200
  # simple client
  simple:
    request_timeout: 30s
    decode_method: utf8
`

	tests := []struct {
		name                string
		envMap              map[string]string
		stubConfFileFunc    func()
		cleanupConfFileFunc func()
		expectedConf        *Config
		expectError         bool
	}{
		{
			name: "happy flow with default",
			envMap: map[string]string{
				"CONFIG_DIRECTORY":               ".",
				"NOVEL_SPIDER_API_ROUTE_PREFIX":  "/api-novel",
				"NOVEL_SPIDER_LITE_ROUTE_PREFIX": "/lite-novel",
				"API_AVAILABLE_SITES":            "xbiquge,xqishu",
				"MAX_WORKING_THREADS":            "1000",
				"BATCH_AVAILABLE_SITES":          "xbiquge,xqishu",
				"PSQL_HOST":                      "host",
				"PSQL_PORT":                      "12345",
				"PSQL_USER":                      "user",
				"PSQL_PASSWORD":                  "password",
				"PSQL_NAME":                      "name",
			},
			stubConfFileFunc: func() {
				os.Mkdir("./selectors", os.ModePerm)
				os.WriteFile("./selectors/xbiquge.yaml", []byte(strings.ReplaceAll(siteSelectorData1, "\t", "  ")), 0644)
				os.WriteFile("./selectors/xqishu.yaml", []byte(strings.ReplaceAll(siteSelectorData2, "\t", "  ")), 0644)

				os.Mkdir("./clients", os.ModePerm)
				os.WriteFile("./clients/xbiquge.yaml", []byte(strings.ReplaceAll(siteClientData1, "\t", "  ")), 0644)
				os.WriteFile("./clients/xqishu.yaml", []byte(strings.ReplaceAll(siteClientData2, "\t", "  ")), 0644)

				confData := `sites:
	xbiquge:
		<<: *xbiquge_selector
		client: *xbiquge_client
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

		max_explore_error: 100
		max_download_concurrency: 10
		update_date_layout: null

	xqishu:
		<<: *xqishu_selector
		client: *xqishu_client
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

		max_explore_error: 100
		max_download_concurrency: 10
		update_date_layout: null
`
				os.WriteFile(`./main.yaml`, []byte(strings.ReplaceAll(confData, "\t", "  ")), 0644)
			},
			cleanupConfFileFunc: func() {
				os.RemoveAll("./clients")
				os.RemoveAll("./selectors")
				os.Remove("./main.yaml")
			},
			expectedConf: &Config{
				APIConfig: APIConfig{
					APIRoutePrefix:     "/api-novel",
					LiteRoutePrefix:    "/lite-novel",
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
						ClientConfig: ClientConfig{
							Simple: simple.SimpleClientConfig{
								RequestTimeout: 30 * time.Second,
								DecodeMethod:   "gbk",
							},
							Retry: retry.RetryClientConfig{
								RetryConditions: []retry.RetryCondition{
									{
										Type:              "status-code",
										Value:             []any{500, 502},
										Weight:            10,
										PauseInterval:     time.Second,
										PauseIntervalType: "exponential",
									},
									{
										Type:              "body-contains",
										Value:             []any{},
										Weight:            100,
										PauseInterval:     time.Second,
										PauseIntervalType: "linear",
									},
									{
										Type:              "error",
										Weight:            100,
										PauseInterval:     time.Second,
										PauseIntervalType: "linear",
									},
								},
								MaxRetryWeight: 1000,
							},
							CircuitBreaker: circuitbreaker.CircuitBreakerClientConfig{
								OpenThreshold:         1000,
								AcquireTimeout:        500 * time.Millisecond,
								MaxConcurrencyThreads: 1000,
								RecoverThreads:        []int64{1, 2, 5, 10, 50, 100, 500},
								OpenDuration:          30 * time.Second,
								RecoverDuration:       10 * time.Second,
								CheckConfigs: []circuitbreaker.CheckConfig{
									{Type: "status-codes", Value: []any{502}},
								},
							},
						},
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
						ClientConfig: ClientConfig{
							Simple: simple.SimpleClientConfig{
								RequestTimeout: 30 * time.Second,
								DecodeMethod:   "utf8",
							},
							Retry: retry.RetryClientConfig{
								RetryConditions: []retry.RetryCondition{
									{
										Type:              "status-code",
										Value:             []any{500, 502},
										Weight:            10,
										PauseInterval:     time.Second,
										PauseIntervalType: "exponential",
									},
									{
										Type:              "body-contains",
										Value:             []any{},
										Weight:            100,
										PauseInterval:     time.Second,
										PauseIntervalType: "linear",
									},
									{
										Type:              "error",
										Weight:            100,
										PauseInterval:     time.Second,
										PauseIntervalType: "linear",
									},
								},
								MaxRetryWeight: 1000,
							},
							CircuitBreaker: circuitbreaker.CircuitBreakerClientConfig{
								OpenThreshold:         10,
								AcquireTimeout:        500 * time.Millisecond,
								MaxConcurrencyThreads: 200,
								RecoverThreads:        []int64{1, 2, 5, 10, 50, 100, 500},
								OpenDuration:          30 * time.Second,
								RecoverDuration:       10 * time.Second,
								CheckConfigs: []circuitbreaker.CheckConfig{
									{Type: "status-codes", Value: []any{502}},
								},
							},
						},
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
				ConfigDirectory: ".",
			},
			expectError: false,
		},
		{
			name: "happy flow without default",
			envMap: map[string]string{
				"CONFIG_DIRECTORY":               ".",
				"API_AVAILABLE_SITES":            "xbiquge",
				"MAX_WORKING_THREADS":            "1000",
				"BATCH_AVAILABLE_SITES":          "xbiquge",
				"NOVEL_SPIDER_API_ROUTE_PREFIX":  "/api-novel",
				"NOVEL_SPIDER_LITE_ROUTE_PREFIX": "/lite-novel",
				"PSQL_HOST":                      "host",
				"PSQL_PORT":                      "12345",
				"PSQL_USER":                      "user",
				"PSQL_PASSWORD":                  "password",
				"PSQL_NAME":                      "name",
			},
			stubConfFileFunc: func() {
				confData := `sites:
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
			check_string: check
`
				os.WriteFile(`./main.yaml`, []byte(strings.ReplaceAll(confData, "\t", "  ")), 0644)
			},
			cleanupConfFileFunc: func() {
				os.Remove("./xbiquge.yaml")
				os.Remove("./main.yaml")
			},
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
				ConfigDirectory: ".",
			},
			expectError: false,
		},
		{
			name: "empty sites",
			envMap: map[string]string{
				"CONFIG_DIRECTORY":               ".",
				"API_AVAILABLE_SITES":            "xbiquge",
				"NOVEL_SPIDER_API_ROUTE_PREFIX":  "/api-novel",
				"NOVEL_SPIDER_LITE_ROUTE_PREFIX": "/lite-novel",
				"MAX_WORKING_THREADS":            "1000",
				"BATCH_AVAILABLE_SITES":          "xbiquge",
				"PSQL_HOST":                      "host",
				"PSQL_PORT":                      "12345",
				"PSQL_USER":                      "user",
				"PSQL_PASSWORD":                  "password",
				"PSQL_NAME":                      "name",
			},
			stubConfFileFunc: func() {
				os.WriteFile("./main.yaml", []byte(""), 0644)
			},
			cleanupConfFileFunc: func() {
				os.Remove("./main.yaml")
			},
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
				DatabaseConfig: DatabaseConfig{
					Host:     "host",
					Port:     "12345",
					User:     "user",
					Password: "password",
					Name:     "name",
				},
				SiteConfigs:     nil,
				ConfigDirectory: ".",
			},
			expectError: false,
		},
		{
			name: "sites config file not exist",
			envMap: map[string]string{
				"CONFIG_DIRECTORY":      ".",
				"API_AVAILABLE_SITES":   "xbiquge",
				"MAX_WORKING_THREADS":   "1000",
				"BATCH_AVAILABLE_SITES": "xbiquge",
				"PSQL_HOST":             "host",
				"PSQL_PORT":             "12345",
				"PSQL_USER":             "user",
				"PSQL_PASSWORD":         "password",
				"PSQL_NAME":             "name",
			},
			stubConfFileFunc: func() {
				os.WriteFile("./main.yaml", []byte(""), 0644)
			},
			cleanupConfFileFunc: func() {
				os.Remove("./main.yaml")
			},
			expectedConf: nil,
			expectError:  true,
		},
		{
			name: "sites config file not exist",
			envMap: map[string]string{
				"CONFIG_DIRECTORY":      "./not-exist/",
				"API_AVAILABLE_SITES":   "xbiquge",
				"MAX_WORKING_THREADS":   "1000",
				"BATCH_AVAILABLE_SITES": "xbiquge",
				"PSQL_HOST":             "host",
				"PSQL_PORT":             "12345",
				"PSQL_USER":             "user",
				"PSQL_PASSWORD":         "password",
				"PSQL_NAME":             "name",
			},
			stubConfFileFunc: func() {
				os.WriteFile("./main.yaml", []byte(""), 0644)
			},
			cleanupConfFileFunc: func() {
				os.Remove("./main.yaml")
			},
			expectedConf: nil,
			expectError:  true,
		},
		{
			name:   "fail to provide required env",
			envMap: map[string]string{},
			stubConfFileFunc: func() {
				os.WriteFile("./main.yaml", []byte(""), 0644)
			},
			cleanupConfFileFunc: func() {
				os.Remove("./main.yaml")
			},
			expectedConf: nil,
			expectError:  true,
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

			test.stubConfFileFunc()
			defer test.cleanupConfFileFunc()

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
