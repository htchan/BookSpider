package config

import (
	"time"

	client "github.com/htchan/BookSpider/internal/client/v2"
)

type SiteConfig struct {
	DecodeMethod   client.DecodeMethod `yaml:"decode_method" validate:"oneof=gbk big5 utf8"`
	MaxThreads     int                 `yaml:"max_threads" validate:"min=1"`
	ClientConfig   ClientConfig        `yaml:"client" validate:"dive"`
	RequestTimeout time.Duration       `yaml:"request_timeout" validate:"min=1s"`

	Storage         string `yaml:"storage" validate:"dir"`
	BackupDirectory string `yaml:"backup_directory" validate:"min=1"`

	URL                    URLConfig              `yaml:"urls"`
	MaxExploreError        int                    `yaml:"max_explore_error" validate:"min=1"`
	MaxDownloadConcurrency int                    `yaml:"max_download_concurrency" validate:"min=1"`
	GoquerySelectorsConfig GoquerySelectorsConfig `yaml:"goquery_selectors"`
	AvailabilityConfig     AvailabilityConfig     `yaml:"availability"`
}

type ClientConfig struct {
	RateLimit      RateLimitConfig      `yaml:"rate_limit" validate:"dive"`
	CircuitBreaker CircuitBreakerConfig `yaml:"circuit_breaker" validate:"dive"`
	Retry          RetryConfig          `yaml:"retry" validate:"dive"`
}

type RateLimitConfig struct {
	QueueSize int           `yaml:"queue_size" validate:"min=1"`
	Interval  time.Duration `yaml:"interval" validate:"min=100ms"`
}

type CircuitBreakerConfig struct {
	FailureThreshold int           `yaml:"failure_threshold" validate:"min=1"`
	SuccessThreshold int           `yaml:"success_threshold" validate:"min=1"`
	RecoverDuration  time.Duration `yaml:"recover_duration" validate:"min=1s"`
}

type RetryConfig struct {
	MaxRetries    int           `yaml:"max_retries" validate:"min=0"`
	BaseInterval  time.Duration `yaml:"base_interval" validate:"min=100ms"`
	IntervalType  string        `yaml:"interval_type" validate:"oneof=static linear exponential"`
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
