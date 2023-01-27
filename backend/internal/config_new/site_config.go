package config

import "time"

type SiteConfig struct {
	DecodeMethod         string                     `yaml:"decode_method"`
	MaxThreads           int                        `yaml:"max_threads"`
	CircuitBreakerConfig CircuitBreakerClientConfig `yaml:"circuit_breaker"`
	RequestTimeout       time.Duration              `yaml:"request_timeout"`
	RetryConfig          map[string]int             `yaml:"retry_map"`

	Storage         string `yaml:"storage"`
	BackupDirectory string `yaml:"backup_directory"`

	URL                    URLConfig              `yaml:"urls"`
	MaxExploreError        int                    `yaml:"max_explore_error"`
	MaxDownloadConcurrency int                    `yaml:"max_download_concurrency"`
	GoquerySelectorsConfig GoquerySelectorsConfig `yaml:"goquery_selectors"`
	AvailabilityConfig     AvailabilityConfig     `yaml:"availability"`
	// UpdateDateLayour string    `yaml:"update_date_layout"`
}

type CircuitBreakerClientConfig struct {
	MaxFailCount      int           `yaml:"max_fail_count"`
	MaxFailMultiplier float64       `yaml:"max_fail_multiplier"`
	SleepInterval     time.Duration `yaml:"sleep_interval"`
}

type URLConfig struct {
	Base          string `yaml:"base"`
	Download      string `yaml:"download"`
	ChapterPrefix string `yaml:"chapter_prefix"`
}

type AvailabilityConfig struct {
	URL         string `yaml:"url"`
	CheckString string `yaml:"check_string"`
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
	Selector        string   `yaml:"selector"`
	Attr            string   `yaml:"attr"`
	UnwantedContent []string `yaml:"unwanted_content"`
}
