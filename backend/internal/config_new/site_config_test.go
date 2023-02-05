package config

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
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
				DecodeMethod: "gbk",
				MaxThreads:   1,
				CircuitBreakerConfig: CircuitBreakerClientConfig{
					MaxFailCount:      1,
					MaxFailMultiplier: 1,
					SleepInterval:     1 * time.Second,
				},
				RequestTimeout: 1 * time.Second,
				RetryConfig:    map[string]int{"default": 1},

				Storage:         ".",
				BackupDirectory: ".",

				URL: URLConfig{
					Base:          "https://test.com",
					Download:      "http://test.com",
					ChapterPrefix: "https://test.com",
				},
				MaxExploreError:        1,
				MaxDownloadConcurrency: 1,
				GoquerySelectorsConfig: GoquerySelectorsConfig{
					Title:            GoquerySelectorConfig{Selector: "data"},
					Writer:           GoquerySelectorConfig{Selector: "data", Attr: "data"},
					BookType:         GoquerySelectorConfig{Selector: "data", UnwantedContent: []string{"a"}},
					LastUpdate:       GoquerySelectorConfig{Selector: "data", Attr: "data", UnwantedContent: []string{"a"}},
					LastChapter:      GoquerySelectorConfig{Selector: "data"},
					BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
					BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
					ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
					ChapterContent:   GoquerySelectorConfig{Selector: "data"},
				},
				AvailabilityConfig: AvailabilityConfig{
					URL:         "http://test.com",
					CheckString: "sth",
				},
			},
			valid: true,
		},
		{
			name: "invalid DecodeMethod",
			conf: SiteConfig{
				DecodeMethod: "unknown",
				MaxThreads:   1,
				CircuitBreakerConfig: CircuitBreakerClientConfig{
					MaxFailCount:      1,
					MaxFailMultiplier: 1,
					SleepInterval:     1 * time.Second,
				},
				RequestTimeout: 1 * time.Second,
				RetryConfig:    map[string]int{"default": 1},

				Storage:         ".",
				BackupDirectory: ".",

				URL: URLConfig{
					Base:          "https://test.com",
					Download:      "http://test.com",
					ChapterPrefix: "https://test.com",
				},
				MaxExploreError:        1,
				MaxDownloadConcurrency: 1,
				GoquerySelectorsConfig: GoquerySelectorsConfig{
					Title:            GoquerySelectorConfig{Selector: "data"},
					Writer:           GoquerySelectorConfig{Selector: "data", Attr: "data"},
					BookType:         GoquerySelectorConfig{Selector: "data", UnwantedContent: []string{"a"}},
					LastUpdate:       GoquerySelectorConfig{Selector: "data", Attr: "data", UnwantedContent: []string{"a"}},
					LastChapter:      GoquerySelectorConfig{Selector: "data"},
					BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
					BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
					ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
					ChapterContent:   GoquerySelectorConfig{Selector: "data"},
				},
				AvailabilityConfig: AvailabilityConfig{
					URL:         "http://test.com",
					CheckString: "sth",
				},
			},
			valid: false,
		},
		{
			name: "invalid MaxThreads",
			conf: SiteConfig{
				DecodeMethod: "gbk",
				MaxThreads:   0,
				CircuitBreakerConfig: CircuitBreakerClientConfig{
					MaxFailCount:      1,
					MaxFailMultiplier: 1,
					SleepInterval:     1 * time.Second,
				},
				RequestTimeout: 1 * time.Second,
				RetryConfig:    map[string]int{"default": 1},

				Storage:         ".",
				BackupDirectory: ".",

				URL: URLConfig{
					Base:          "https://test.com",
					Download:      "http://test.com",
					ChapterPrefix: "https://test.com",
				},
				MaxExploreError:        1,
				MaxDownloadConcurrency: 1,
				GoquerySelectorsConfig: GoquerySelectorsConfig{
					Title:            GoquerySelectorConfig{Selector: "data"},
					Writer:           GoquerySelectorConfig{Selector: "data", Attr: "data"},
					BookType:         GoquerySelectorConfig{Selector: "data", UnwantedContent: []string{"a"}},
					LastUpdate:       GoquerySelectorConfig{Selector: "data", Attr: "data", UnwantedContent: []string{"a"}},
					LastChapter:      GoquerySelectorConfig{Selector: "data"},
					BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
					BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
					ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
					ChapterContent:   GoquerySelectorConfig{Selector: "data"},
				},
				AvailabilityConfig: AvailabilityConfig{
					URL:         "http://test.com",
					CheckString: "sth",
				},
			},
			valid: false,
		},
		{
			name: "invalid CircuitBreakerConfig",
			conf: SiteConfig{
				DecodeMethod: "gbk",
				MaxThreads:   1,
				CircuitBreakerConfig: CircuitBreakerClientConfig{
					MaxFailCount:      0,
					MaxFailMultiplier: 1,
					SleepInterval:     1 * time.Second,
				},
				RequestTimeout: 1 * time.Second,
				RetryConfig:    map[string]int{"default": 1},

				Storage:         ".",
				BackupDirectory: ".",

				URL: URLConfig{
					Base:          "https://test.com",
					Download:      "http://test.com",
					ChapterPrefix: "https://test.com",
				},
				MaxExploreError:        1,
				MaxDownloadConcurrency: 1,
				GoquerySelectorsConfig: GoquerySelectorsConfig{
					Title:            GoquerySelectorConfig{Selector: "data"},
					Writer:           GoquerySelectorConfig{Selector: "data", Attr: "data"},
					BookType:         GoquerySelectorConfig{Selector: "data", UnwantedContent: []string{"a"}},
					LastUpdate:       GoquerySelectorConfig{Selector: "data", Attr: "data", UnwantedContent: []string{"a"}},
					LastChapter:      GoquerySelectorConfig{Selector: "data"},
					BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
					BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
					ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
					ChapterContent:   GoquerySelectorConfig{Selector: "data"},
				},
				AvailabilityConfig: AvailabilityConfig{
					URL:         "http://test.com",
					CheckString: "sth",
				},
			},
			valid: false,
		},
		{
			name: "invalid RequestConfig - empty",
			conf: SiteConfig{
				DecodeMethod: "gbk",
				MaxThreads:   1,
				CircuitBreakerConfig: CircuitBreakerClientConfig{
					MaxFailCount:      1,
					MaxFailMultiplier: 1,
					SleepInterval:     1 * time.Second,
				},
				RequestTimeout: 1 * time.Second,
				RetryConfig:    map[string]int{},

				Storage:         ".",
				BackupDirectory: ".",

				URL: URLConfig{
					Base:          "https://test.com",
					Download:      "http://test.com",
					ChapterPrefix: "https://test.com",
				},
				MaxExploreError:        1,
				MaxDownloadConcurrency: 1,
				GoquerySelectorsConfig: GoquerySelectorsConfig{
					Title:            GoquerySelectorConfig{Selector: "data"},
					Writer:           GoquerySelectorConfig{Selector: "data", Attr: "data"},
					BookType:         GoquerySelectorConfig{Selector: "data", UnwantedContent: []string{"a"}},
					LastUpdate:       GoquerySelectorConfig{Selector: "data", Attr: "data", UnwantedContent: []string{"a"}},
					LastChapter:      GoquerySelectorConfig{Selector: "data"},
					BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
					BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
					ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
					ChapterContent:   GoquerySelectorConfig{Selector: "data"},
				},
				AvailabilityConfig: AvailabilityConfig{
					URL:         "http://test.com",
					CheckString: "sth",
				},
			},
			valid: false,
		},
		{
			name: "invalid RequestConfig - value is 0",
			conf: SiteConfig{
				DecodeMethod: "gbk",
				MaxThreads:   1,
				CircuitBreakerConfig: CircuitBreakerClientConfig{
					MaxFailCount:      1,
					MaxFailMultiplier: 1,
					SleepInterval:     1 * time.Second,
				},
				RequestTimeout: 1 * time.Second,
				RetryConfig:    map[string]int{"default": 0},

				Storage:         ".",
				BackupDirectory: ".",

				URL: URLConfig{
					Base:          "https://test.com",
					Download:      "http://test.com",
					ChapterPrefix: "https://test.com",
				},
				MaxExploreError:        1,
				MaxDownloadConcurrency: 1,
				GoquerySelectorsConfig: GoquerySelectorsConfig{
					Title:            GoquerySelectorConfig{Selector: "data"},
					Writer:           GoquerySelectorConfig{Selector: "data", Attr: "data"},
					BookType:         GoquerySelectorConfig{Selector: "data", UnwantedContent: []string{"a"}},
					LastUpdate:       GoquerySelectorConfig{Selector: "data", Attr: "data", UnwantedContent: []string{"a"}},
					LastChapter:      GoquerySelectorConfig{Selector: "data"},
					BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
					BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
					ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
					ChapterContent:   GoquerySelectorConfig{Selector: "data"},
				},
				AvailabilityConfig: AvailabilityConfig{
					URL:         "http://test.com",
					CheckString: "sth",
				},
			},
			valid: false,
		},
		{
			name: "invalid storage",
			conf: SiteConfig{
				DecodeMethod: "gbk",
				MaxThreads:   1,
				CircuitBreakerConfig: CircuitBreakerClientConfig{
					MaxFailCount:      1,
					MaxFailMultiplier: 1,
					SleepInterval:     1 * time.Second,
				},
				RequestTimeout: 1 * time.Second,
				RetryConfig:    map[string]int{"default": 1},

				Storage:         "./not-exist-dir",
				BackupDirectory: ".",

				URL: URLConfig{
					Base:          "https://test.com",
					Download:      "http://test.com",
					ChapterPrefix: "https://test.com",
				},
				MaxExploreError:        1,
				MaxDownloadConcurrency: 1,
				GoquerySelectorsConfig: GoquerySelectorsConfig{
					Title:            GoquerySelectorConfig{Selector: "data"},
					Writer:           GoquerySelectorConfig{Selector: "data", Attr: "data"},
					BookType:         GoquerySelectorConfig{Selector: "data", UnwantedContent: []string{"a"}},
					LastUpdate:       GoquerySelectorConfig{Selector: "data", Attr: "data", UnwantedContent: []string{"a"}},
					LastChapter:      GoquerySelectorConfig{Selector: "data"},
					BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
					BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
					ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
					ChapterContent:   GoquerySelectorConfig{Selector: "data"},
				},
				AvailabilityConfig: AvailabilityConfig{
					URL:         "http://test.com",
					CheckString: "sth",
				},
			},
			valid: false,
		},
		{
			name: "invalid BackupDirectory - empty",
			conf: SiteConfig{
				DecodeMethod: "gbk",
				MaxThreads:   1,
				CircuitBreakerConfig: CircuitBreakerClientConfig{
					MaxFailCount:      1,
					MaxFailMultiplier: 1,
					SleepInterval:     1 * time.Second,
				},
				RequestTimeout: 1 * time.Second,
				RetryConfig:    map[string]int{"default": 1},

				Storage:         ".",
				BackupDirectory: "",

				URL: URLConfig{
					Base:          "https://test.com",
					Download:      "http://test.com",
					ChapterPrefix: "https://test.com",
				},
				MaxExploreError:        1,
				MaxDownloadConcurrency: 1,
				GoquerySelectorsConfig: GoquerySelectorsConfig{
					Title:            GoquerySelectorConfig{Selector: "data"},
					Writer:           GoquerySelectorConfig{Selector: "data", Attr: "data"},
					BookType:         GoquerySelectorConfig{Selector: "data", UnwantedContent: []string{"a"}},
					LastUpdate:       GoquerySelectorConfig{Selector: "data", Attr: "data", UnwantedContent: []string{"a"}},
					LastChapter:      GoquerySelectorConfig{Selector: "data"},
					BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
					BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
					ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
					ChapterContent:   GoquerySelectorConfig{Selector: "data"},
				},
				AvailabilityConfig: AvailabilityConfig{
					URL:         "http://test.com",
					CheckString: "sth",
				},
			},
			valid: false,
		},
		{
			name: "invalid URL",
			conf: SiteConfig{
				DecodeMethod: "gbk",
				MaxThreads:   1,
				CircuitBreakerConfig: CircuitBreakerClientConfig{
					MaxFailCount:      1,
					MaxFailMultiplier: 1,
					SleepInterval:     1 * time.Second,
				},
				RequestTimeout: 1 * time.Second,
				RetryConfig:    map[string]int{"default": 1},

				Storage:         ".",
				BackupDirectory: ".",

				URL: URLConfig{
					Base:          "test.com",
					Download:      "http://test.com",
					ChapterPrefix: "https://test.com",
				},
				MaxExploreError:        1,
				MaxDownloadConcurrency: 1,
				GoquerySelectorsConfig: GoquerySelectorsConfig{
					Title:            GoquerySelectorConfig{Selector: "data"},
					Writer:           GoquerySelectorConfig{Selector: "data", Attr: "data"},
					BookType:         GoquerySelectorConfig{Selector: "data", UnwantedContent: []string{"a"}},
					LastUpdate:       GoquerySelectorConfig{Selector: "data", Attr: "data", UnwantedContent: []string{"a"}},
					LastChapter:      GoquerySelectorConfig{Selector: "data"},
					BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
					BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
					ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
					ChapterContent:   GoquerySelectorConfig{Selector: "data"},
				},
				AvailabilityConfig: AvailabilityConfig{
					URL:         "http://test.com",
					CheckString: "sth",
				},
			},
			valid: false,
		},
		{
			name: "invalid MaxExploreError",
			conf: SiteConfig{
				DecodeMethod: "gbk",
				MaxThreads:   1,
				CircuitBreakerConfig: CircuitBreakerClientConfig{
					MaxFailCount:      1,
					MaxFailMultiplier: 1,
					SleepInterval:     1 * time.Second,
				},
				RequestTimeout: 1 * time.Second,
				RetryConfig:    map[string]int{"default": 1},

				Storage:         ".",
				BackupDirectory: ".",

				URL: URLConfig{
					Base:          "https://test.com",
					Download:      "http://test.com",
					ChapterPrefix: "https://test.com",
				},
				MaxExploreError:        0,
				MaxDownloadConcurrency: 1,
				GoquerySelectorsConfig: GoquerySelectorsConfig{
					Title:            GoquerySelectorConfig{Selector: "data"},
					Writer:           GoquerySelectorConfig{Selector: "data", Attr: "data"},
					BookType:         GoquerySelectorConfig{Selector: "data", UnwantedContent: []string{"a"}},
					LastUpdate:       GoquerySelectorConfig{Selector: "data", Attr: "data", UnwantedContent: []string{"a"}},
					LastChapter:      GoquerySelectorConfig{Selector: "data"},
					BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
					BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
					ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
					ChapterContent:   GoquerySelectorConfig{Selector: "data"},
				},
				AvailabilityConfig: AvailabilityConfig{
					URL:         "http://test.com",
					CheckString: "sth",
				},
			},
			valid: false,
		},
		{
			name: "invalid MaxDownloadConcurrency",
			conf: SiteConfig{
				DecodeMethod: "gbk",
				MaxThreads:   1,
				CircuitBreakerConfig: CircuitBreakerClientConfig{
					MaxFailCount:      1,
					MaxFailMultiplier: 1,
					SleepInterval:     1 * time.Second,
				},
				RequestTimeout: 1 * time.Second,
				RetryConfig:    map[string]int{"default": 1},

				Storage:         ".",
				BackupDirectory: ".",

				URL: URLConfig{
					Base:          "https://test.com",
					Download:      "http://test.com",
					ChapterPrefix: "https://test.com",
				},
				MaxExploreError:        1,
				MaxDownloadConcurrency: 0,
				GoquerySelectorsConfig: GoquerySelectorsConfig{
					Title:            GoquerySelectorConfig{Selector: "data"},
					Writer:           GoquerySelectorConfig{Selector: "data", Attr: "data"},
					BookType:         GoquerySelectorConfig{Selector: "data", UnwantedContent: []string{"a"}},
					LastUpdate:       GoquerySelectorConfig{Selector: "data", Attr: "data", UnwantedContent: []string{"a"}},
					LastChapter:      GoquerySelectorConfig{Selector: "data"},
					BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
					BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
					ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
					ChapterContent:   GoquerySelectorConfig{Selector: "data"},
				},
				AvailabilityConfig: AvailabilityConfig{
					URL:         "http://test.com",
					CheckString: "sth",
				},
			},
			valid: false,
		},
		{
			name: "invalid GoquerySelectorsConfig",
			conf: SiteConfig{
				DecodeMethod: "gbk",
				MaxThreads:   1,
				CircuitBreakerConfig: CircuitBreakerClientConfig{
					MaxFailCount:      1,
					MaxFailMultiplier: 1,
					SleepInterval:     1 * time.Second,
				},
				RequestTimeout: 1 * time.Second,
				RetryConfig:    map[string]int{"default": 1},

				Storage:         ".",
				BackupDirectory: ".",

				URL: URLConfig{
					Base:          "https://test.com",
					Download:      "http://test.com",
					ChapterPrefix: "https://test.com",
				},
				MaxExploreError:        1,
				MaxDownloadConcurrency: 1,
				GoquerySelectorsConfig: GoquerySelectorsConfig{
					Title:            GoquerySelectorConfig{Selector: ""},
					Writer:           GoquerySelectorConfig{Selector: "data", Attr: "data"},
					BookType:         GoquerySelectorConfig{Selector: "data", UnwantedContent: []string{"a"}},
					LastUpdate:       GoquerySelectorConfig{Selector: "data", Attr: "data", UnwantedContent: []string{"a"}},
					LastChapter:      GoquerySelectorConfig{Selector: "data"},
					BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
					BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
					ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
					ChapterContent:   GoquerySelectorConfig{Selector: "data"},
				},
				AvailabilityConfig: AvailabilityConfig{
					URL:         "http://test.com",
					CheckString: "sth",
				},
			},
			valid: false,
		},
		{
			name: "invalid AvailabilityConfig",
			conf: SiteConfig{
				DecodeMethod: "gbk",
				MaxThreads:   1,
				CircuitBreakerConfig: CircuitBreakerClientConfig{
					MaxFailCount:      1,
					MaxFailMultiplier: 1,
					SleepInterval:     1 * time.Second,
				},
				RequestTimeout: 1 * time.Second,
				RetryConfig:    map[string]int{"default": 1},

				Storage:         ".",
				BackupDirectory: ".",

				URL: URLConfig{
					Base:          "https://test.com",
					Download:      "http://test.com",
					ChapterPrefix: "https://test.com",
				},
				MaxExploreError:        1,
				MaxDownloadConcurrency: 1,
				GoquerySelectorsConfig: GoquerySelectorsConfig{
					Title:            GoquerySelectorConfig{Selector: "data"},
					Writer:           GoquerySelectorConfig{Selector: "data", Attr: "data"},
					BookType:         GoquerySelectorConfig{Selector: "data", UnwantedContent: []string{"a"}},
					LastUpdate:       GoquerySelectorConfig{Selector: "data", Attr: "data", UnwantedContent: []string{"a"}},
					LastChapter:      GoquerySelectorConfig{Selector: "data"},
					BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
					BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
					ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
					ChapterContent:   GoquerySelectorConfig{Selector: "data"},
				},
				AvailabilityConfig: AvailabilityConfig{
					URL:         "http://test.com",
					CheckString: "",
				},
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

func Test_validate_CircuitBreakerClientConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		conf  CircuitBreakerClientConfig
		valid bool
	}{
		{
			name: "valid conf",
			conf: CircuitBreakerClientConfig{
				MaxFailCount:      1,
				MaxFailMultiplier: 1,
				SleepInterval:     1 * time.Second,
			},
			valid: true,
		},
		{
			name: "invalid MaxFailCount - zero",
			conf: CircuitBreakerClientConfig{
				MaxFailCount:      0,
				MaxFailMultiplier: 1,
				SleepInterval:     1 * time.Second,
			},
			valid: false,
		},
		{
			name: "invalid MaxFailMultiplier - zero",
			conf: CircuitBreakerClientConfig{
				MaxFailCount:      1,
				MaxFailMultiplier: 0,
				SleepInterval:     1 * time.Second,
			},
			valid: false,
		},
		{
			name: "invalid SleepInterval - smaller than 1 second",
			conf: CircuitBreakerClientConfig{
				MaxFailCount:      1,
				MaxFailMultiplier: 1,
				SleepInterval:     1 * time.Millisecond,
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

func Test_validate_URLConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		conf  URLConfig
		valid bool
	}{
		{
			name: "valid conf",
			conf: URLConfig{
				Base:          "http://test.com",
				Download:      "http://test.com",
				ChapterPrefix: "http://test.com",
			},
			valid: true,
		},
		{
			name: "invalid Base - not url",
			conf: URLConfig{
				Base:          "test.com",
				Download:      "http://test.com",
				ChapterPrefix: "http://test.com",
			},
			valid: false,
		},
		{
			name: "invalid Download - not url",
			conf: URLConfig{
				Base:          "http://test.com",
				Download:      "test.com",
				ChapterPrefix: "http://test.com",
			},
			valid: false,
		},
		{
			name: "invalid ChapterPrefix - not url",
			conf: URLConfig{
				Base:          "http://test.com",
				Download:      "http://test.com",
				ChapterPrefix: "test.com",
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

func Test_validate_AvailablityConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		conf  AvailabilityConfig
		valid bool
	}{
		{
			name: "valid conf",
			conf: AvailabilityConfig{
				URL:         "http://test.com",
				CheckString: "data",
			},
			valid: true,
		},
		{
			name: "invalid URL - not url",
			conf: AvailabilityConfig{
				URL:         "test.com",
				CheckString: "data",
			},
			valid: false,
		},
		{
			name: "invald CheckString - empty",
			conf: AvailabilityConfig{
				URL:         "http://test.com",
				CheckString: "",
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

func Test_validate_GoquerySelectorsConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		conf  GoquerySelectorsConfig
		valid bool
	}{
		{
			name: "valid conf",
			conf: GoquerySelectorsConfig{
				Title:            GoquerySelectorConfig{Selector: "data"},
				Writer:           GoquerySelectorConfig{Selector: "data"},
				BookType:         GoquerySelectorConfig{Selector: "data"},
				LastUpdate:       GoquerySelectorConfig{Selector: "data"},
				LastChapter:      GoquerySelectorConfig{Selector: "data"},
				BookChapterURL:   GoquerySelectorConfig{Selector: "data"},
				BookChapterTitle: GoquerySelectorConfig{Selector: "data"},
				ChapterTitle:     GoquerySelectorConfig{Selector: "data"},
				ChapterContent:   GoquerySelectorConfig{Selector: "data"},
			},
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
