package config

import (
	"time"

	circuitbreaker "github.com/htchan/BookSpider/internal/client/v2/circuit_breaker"
	"github.com/htchan/BookSpider/internal/client/v2/retry"
	"github.com/htchan/BookSpider/internal/client/v2/simple"
)

type SiteConfig struct {
	DecodeMethod         string                     `yaml:"decode_method" validate:"oneof=gbk big5 utf8"`
	MaxThreads           int                        `yaml:"max_threads" validate:"min=1"`
	CircuitBreakerConfig CircuitBreakerClientConfig `yaml:"circuit_breaker"`
	ClientConfig         ClientConfig               `yaml:"client" validate:"dive"`
	RequestTimeout       time.Duration              `yaml:"request_timeout" validate:"min=1s"`
	RetryConfig          map[string]int             `yaml:"retry_map" validate:"min=1,dive,min=1"`

	Storage         string `yaml:"storage" validate:"dir"`
	BackupDirectory string `yaml:"backup_directory" validate:"min=1"`

	URL                    URLConfig              `yaml:"urls"`
	MaxExploreError        int                    `yaml:"max_explore_error" validate:"min=1"`
	MaxDownloadConcurrency int                    `yaml:"max_download_concurrency" validate:"min=1"`
	GoquerySelectorsConfig GoquerySelectorsConfig `yaml:"goquery_selectors"`
	AvailabilityConfig     AvailabilityConfig     `yaml:"availability"`
	// UpdateDateLayour string    `yaml:"update_date_layout"`
}

type ClientConfig struct {
	Simple         simple.SimpleClientConfig                 `yaml:"simple" validate:"dive"`
	Retry          retry.RetryClientConfig                   `yaml:"retry" validate:"dive"`
	CircuitBreaker circuitbreaker.CircuitBreakerClientConfig `yaml:"circuit_breaker" validate:"dive"`
}

type CircuitBreakerClientConfig struct {
	MaxFailCount      int           `yaml:"max_fail_count" validate:"min=1"`
	MaxFailMultiplier float64       `yaml:"max_fail_multiplier" validate:"min=1"`
	SleepInterval     time.Duration `yaml:"sleep_interval" validate:"min=1s"`
}

type URLConfig struct {
	Base          string `yaml:"base" validate:"startswith=http://|startswith=https://"`
	Download      string `yaml:"download" validate:"startswith=http://|startswith=https://"`
	ChapterPrefix string `yaml:"chapter_prefix" validate:"startswith=http://|startswith=https://"`
}

type AvailabilityConfig struct {
	URL         string `yaml:"url" validate:"url"`
	CheckString string `yaml:"check_string" validate:"min=1"`
}

type GoquerySelectorsConfig struct {
	Title            GoquerySelectorConfig `yaml:"title"`
	Writer           GoquerySelectorConfig `yaml:"writer"`
	BookType         GoquerySelectorConfig `yaml:"book_type"`
	LastUpdate       GoquerySelectorConfig `yaml:"update_date"`
	LastChapter      GoquerySelectorConfig `yaml:"update_chapter"`
	BookChapterURL   GoquerySelectorConfig `yaml:"book_chapter_url"`
	BookChapterTitle GoquerySelectorConfig `yaml:"book_chapter_title"`
	ChapterTitle     GoquerySelectorConfig `yaml:"chapter_title"`
	ChapterContent   GoquerySelectorConfig `yaml:"chapter_content"`
}

type GoquerySelectorConfig struct {
	Selector        string   `yaml:"selector" validate:"min=1"`
	Attr            string   `yaml:"attr"`
	UnwantedContent []string `yaml:"unwanted_content" validate:"dive,min=1"`
}
