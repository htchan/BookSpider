package config

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

var (
	standardURLConf = URLConfig{
		Base:          "https://test.com",
		Download:      "http://test.com",
		ChapterPrefix: "https://test.com",
	}
	standardGoquerySelectorsConf = GoquerySelectorsConfig{
		Title:            GoquerySelectorConfig{Selector: "data"},
		Writer:           GoquerySelectorConfig{Selector: "data", Attr: "data"},
		BookType:         GoquerySelectorConfig{Selector: "data", UnwantedContent: []string{"a"}},
		LastUpdate:       GoquerySelectorConfig{Selector: "data", Attr: "data", UnwantedContent: []string{"a"}},
		LastChapter:      GoquerySelectorConfig{Selector: "data"},
		BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
		BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
		ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
		ChapterContent:   GoquerySelectorConfig{Selector: "data"},
	}
	standardAvailabilityConf = AvailabilityConfig{
		URL:         "http://test.com",
		CheckString: "sth",
	}
	standardClientConf = ClientConfig{
		RateLimit: RateLimitConfig{
			QueueSize: 10,
			Interval:  time.Second,
		},
		CircuitBreaker: CircuitBreakerConfig{
			FailureThreshold: 3,
			SuccessThreshold: 1,
			RecoverDuration:  5 * time.Second,
		},
		Retry: RetryConfig{
			MaxRetries:   3,
			BaseInterval: time.Second,
			IntervalType: "exponential",
		},
	}
)

func Test_validate_SiteConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		conf  SiteConfig
		valid bool
	}{
		{
			name: "valid conf",
			conf: SiteConfig{
				DecodeMethod:   "gbk",
				MaxThreads:     1,
				ClientConfig:   standardClientConf,
				RequestTimeout: 1 * time.Second,

				Storage:         ".",
				BackupDirectory: ".",

				URL:                    standardURLConf,
				MaxExploreError:        1,
				MaxDownloadConcurrency: 1,
				GoquerySelectorsConfig: standardGoquerySelectorsConf,
				AvailabilityConfig:     standardAvailabilityConf,
			},
			valid: true,
		},
		{
			name: "invalid DecodeMethod",
			conf: SiteConfig{
				DecodeMethod:   "unknown",
				MaxThreads:     1,
				ClientConfig:   standardClientConf,
				RequestTimeout: 1 * time.Second,

				Storage:         ".",
				BackupDirectory: ".",

				URL:                    standardURLConf,
				MaxExploreError:        1,
				MaxDownloadConcurrency: 1,
				GoquerySelectorsConfig: standardGoquerySelectorsConf,
				AvailabilityConfig:     standardAvailabilityConf,
			},
			valid: false,
		},
		{
			name: "invalid MaxThreads",
			conf: SiteConfig{
				DecodeMethod:   "gbk",
				MaxThreads:     0,
				ClientConfig:   standardClientConf,
				RequestTimeout: 1 * time.Second,

				Storage:         ".",
				BackupDirectory: ".",

				URL:                    standardURLConf,
				MaxExploreError:        1,
				MaxDownloadConcurrency: 1,
				GoquerySelectorsConfig: standardGoquerySelectorsConf,
				AvailabilityConfig:     standardAvailabilityConf,
			},
			valid: false,
		},
		{
			name: "invalid RequestTimeout",
			conf: SiteConfig{
				DecodeMethod: "gbk",
				MaxThreads:   1,
				ClientConfig: standardClientConf,

				Storage:         ".",
				BackupDirectory: ".",

				URL:                    standardURLConf,
				MaxExploreError:        1,
				MaxDownloadConcurrency: 1,
				GoquerySelectorsConfig: standardGoquerySelectorsConf,
				AvailabilityConfig:     standardAvailabilityConf,
			},
			valid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := validator.New().Struct(test.conf)
			if !assert.Equal(t, test.valid, err == nil) {
				t.Errorf("getting error: %v", err)
			}
		})
	}
}

func Test_validate_ClientConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		conf  ClientConfig
		valid bool
	}{
		{
			name:  "valid conf",
			conf:  standardClientConf,
			valid: true,
		},
		{
			name: "invalid rate limit queue size",
			conf: ClientConfig{
				RateLimit: RateLimitConfig{
					QueueSize: 0,
					Interval:  time.Second,
				},
				CircuitBreaker: standardClientConf.CircuitBreaker,
				Retry:          standardClientConf.Retry,
			},
			valid: false,
		},
		{
			name: "invalid circuit breaker failure threshold",
			conf: ClientConfig{
				RateLimit: standardClientConf.RateLimit,
				CircuitBreaker: CircuitBreakerConfig{
					FailureThreshold: 0,
					SuccessThreshold: 1,
					RecoverDuration:  5 * time.Second,
				},
				Retry: standardClientConf.Retry,
			},
			valid: false,
		},
		{
			name: "invalid retry interval type",
			conf: ClientConfig{
				RateLimit:      standardClientConf.RateLimit,
				CircuitBreaker: standardClientConf.CircuitBreaker,
				Retry: RetryConfig{
					MaxRetries:   3,
					BaseInterval: time.Second,
					IntervalType: "invalid",
				},
			},
			valid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := validator.New().Struct(test.conf)
			if !assert.Equal(t, test.valid, err == nil) {
				t.Errorf("getting error: %v", err)
			}
		})
	}
}

func Test_validate_URLConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		conf  URLConfig
		valid bool
	}{
		{
			name:  "valid conf",
			conf:  standardURLConf,
			valid: true,
		},
		{
			name: "invalid base - empty",
			conf: URLConfig{
				Base:          "",
				Download:      "http://test.com",
				ChapterPrefix: "https://test.com",
			},
			valid: false,
		},
		{
			name: "invalid base - no http prefix",
			conf: URLConfig{
				Base:          "test.com",
				Download:      "http://test.com",
				ChapterPrefix: "https://test.com",
			},
			valid: false,
		},
		{
			name: "invalid download - empty",
			conf: URLConfig{
				Base:          "http://test.com",
				Download:      "",
				ChapterPrefix: "https://test.com",
			},
			valid: false,
		},
		{
			name: "invalid chapter prefix - empty",
			conf: URLConfig{
				Base:          "http://test.com",
				Download:      "http://test.com",
				ChapterPrefix: "",
			},
			valid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := validator.New().Struct(test.conf)
			if !assert.Equal(t, test.valid, err == nil) {
				t.Errorf("getting error: %v", err)
			}
		})
	}
}

func Test_validate_AvailablityConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		conf  AvailabilityConfig
		valid bool
	}{
		{
			name:  "valid conf",
			conf:  standardAvailabilityConf,
			valid: true,
		},
		{
			name: "invalid URL - empty",
			conf: AvailabilityConfig{
				URL:         "",
				CheckString: "sth",
			},
			valid: false,
		},
		{
			name: "invalid CheckString - empty",
			conf: AvailabilityConfig{
				URL:         "http://test.com",
				CheckString: "",
			},
			valid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := validator.New().Struct(test.conf)
			if !assert.Equal(t, test.valid, err == nil) {
				t.Errorf("getting error: %v", err)
			}
		})
	}
}

func Test_validate_GoquerySelectorsConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		conf  GoquerySelectorsConfig
		valid bool
	}{
		{
			name:  "valid conf",
			conf:  standardGoquerySelectorsConf,
			valid: true,
		},
		{
			name: "invalid Title",
			conf: GoquerySelectorsConfig{
				Title:            GoquerySelectorConfig{Selector: ""},
				Writer:           GoquerySelectorConfig{Selector: "data"},
				BookType:         GoquerySelectorConfig{Selector: "data"},
				LastUpdate:       GoquerySelectorConfig{Selector: "data"},
				LastChapter:      GoquerySelectorConfig{Selector: "data"},
				BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
				BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
				ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
				ChapterContent:   GoquerySelectorConfig{Selector: "data"},
			},
			valid: false,
		},
		{
			name: "invalid Writer",
			conf: GoquerySelectorsConfig{
				Title:            GoquerySelectorConfig{Selector: "data"},
				Writer:           GoquerySelectorConfig{Selector: ""},
				BookType:         GoquerySelectorConfig{Selector: "data"},
				LastUpdate:       GoquerySelectorConfig{Selector: "data"},
				LastChapter:      GoquerySelectorConfig{Selector: "data"},
				BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
				BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
				ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
				ChapterContent:   GoquerySelectorConfig{Selector: "data"},
			},
			valid: false,
		},
		{
			name: "invalid BookType",
			conf: GoquerySelectorsConfig{
				Title:            GoquerySelectorConfig{Selector: "data"},
				Writer:           GoquerySelectorConfig{Selector: "data"},
				BookType:         GoquerySelectorConfig{Selector: ""},
				LastUpdate:       GoquerySelectorConfig{Selector: "data"},
				LastChapter:      GoquerySelectorConfig{Selector: "data"},
				BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
				BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
				ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
				ChapterContent:   GoquerySelectorConfig{Selector: "data"},
			},
			valid: false,
		},
		{
			name: "invalid LastUpdate",
			conf: GoquerySelectorsConfig{
				Title:            GoquerySelectorConfig{Selector: "data"},
				Writer:           GoquerySelectorConfig{Selector: "data"},
				BookType:         GoquerySelectorConfig{Selector: "data"},
				LastUpdate:       GoquerySelectorConfig{Selector: ""},
				LastChapter:      GoquerySelectorConfig{Selector: "data"},
				BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
				BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
				ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
				ChapterContent:   GoquerySelectorConfig{Selector: "data"},
			},
			valid: false,
		},
		{
			name: "invalid LastChapter",
			conf: GoquerySelectorsConfig{
				Title:            GoquerySelectorConfig{Selector: "data"},
				Writer:           GoquerySelectorConfig{Selector: "data"},
				BookType:         GoquerySelectorConfig{Selector: "data"},
				LastUpdate:       GoquerySelectorConfig{Selector: "data"},
				LastChapter:      GoquerySelectorConfig{Selector: ""},
				BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
				BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
				ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
				ChapterContent:   GoquerySelectorConfig{Selector: "data"},
			},
			valid: false,
		},
		{
			name: "invalid BookChapterURL",
			conf: GoquerySelectorsConfig{
				Title:            GoquerySelectorConfig{Selector: "data"},
				Writer:           GoquerySelectorConfig{Selector: "data"},
				BookType:         GoquerySelectorConfig{Selector: "data"},
				LastUpdate:       GoquerySelectorConfig{Selector: "data"},
				LastChapter:      GoquerySelectorConfig{Selector: "data"},
				BookChapterURL:   GoquerySelectorConfig{Selector: ""},
				BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
				ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
				ChapterContent:   GoquerySelectorConfig{Selector: "data"},
			},
			valid: false,
		},
		{
			name: "invalid BookChapterTitle",
			conf: GoquerySelectorsConfig{
				Title:            GoquerySelectorConfig{Selector: "data"},
				Writer:           GoquerySelectorConfig{Selector: "data"},
				BookType:         GoquerySelectorConfig{Selector: "data"},
				LastUpdate:       GoquerySelectorConfig{Selector: "data"},
				LastChapter:      GoquerySelectorConfig{Selector: "data"},
				BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
				BookChapterTitle: GoquerySelectorConfig{Selector: ""},
				ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
				ChapterContent:   GoquerySelectorConfig{Selector: "data"},
			},
			valid: false,
		},
		{
			name: "invalid ChapterTitle",
			conf: GoquerySelectorsConfig{
				Title:            GoquerySelectorConfig{Selector: "data"},
				Writer:           GoquerySelectorConfig{Selector: "data"},
				BookType:         GoquerySelectorConfig{Selector: "data"},
				LastUpdate:       GoquerySelectorConfig{Selector: "data"},
				LastChapter:      GoquerySelectorConfig{Selector: "data"},
				BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
				BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
				ChapterTitle:     GoquerySelectorConfig{Selector: ""},
				ChapterContent:   GoquerySelectorConfig{Selector: "data"},
			},
			valid: false,
		},
		{
			name: "invalid ChapterContent",
			conf: GoquerySelectorsConfig{
				Title:            GoquerySelectorConfig{Selector: "data"},
				Writer:           GoquerySelectorConfig{Selector: "data"},
				BookType:         GoquerySelectorConfig{Selector: "data"},
				LastUpdate:       GoquerySelectorConfig{Selector: "data"},
				LastChapter:      GoquerySelectorConfig{Selector: "data"},
				BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
				BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
				ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
				ChapterContent:   GoquerySelectorConfig{Selector: ""},
			},
			valid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := validator.New().Struct(test.conf)
			if !assert.Equal(t, test.valid, err == nil) {
				t.Errorf("getting error: %v", err)
			}
		})
	}
}

func Test_validate_GoquerySelectorConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		conf  GoquerySelectorConfig
		valid bool
	}{
		{
			name: "valid conf",
			conf: GoquerySelectorConfig{
				Selector:        "data",
				Attr:            "data",
				UnwantedContent: []string{"data"},
			},
			valid: true,
		},
		{
			name: "invalid selector - empty",
			conf: GoquerySelectorConfig{
				Selector:        "",
				Attr:            "data",
				UnwantedContent: []string{"data"},
			},
			valid: false,
		},
		{
			name: "valid Attr - empty",
			conf: GoquerySelectorConfig{
				Selector:        "data",
				Attr:            "",
				UnwantedContent: []string{"data"},
			},
			valid: true,
		},
		{
			name: "valid UnwantedContent - empty ",
			conf: GoquerySelectorConfig{
				Selector:        "data",
				Attr:            "data",
				UnwantedContent: []string{},
			},
			valid: true,
		},
		{
			name: "invalid UnwantedContent - empty content",
			conf: GoquerySelectorConfig{
				Selector:        "data",
				Attr:            "data",
				UnwantedContent: []string{""},
			},
			valid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := validator.New().Struct(test.conf)
			if !assert.Equal(t, test.valid, err == nil) {
				t.Errorf("getting error: %v", err)
			}
		})
	}
}
