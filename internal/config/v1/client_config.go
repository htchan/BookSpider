package config

import "time"

type ClientConfig struct {
	DecodeMethod string           `yaml:"decode_method" json:"decode_method"`
	Retry        RetryConfig      `yaml:"retry" json:"retry"`
	Pool         ClientPoolConfig `yaml:"pool" json:"pool"`
}

type RetryConfig struct {
	MaxRetryCount       int           `yaml:"max_retry_count" json:"max_retry_count"`
	LinearRetryInterval time.Duration `yaml:"linear_retry_interval" json:"linear_retry_interval"`
}

type ClientPoolConfig struct {
	RefreshInterval            time.Duration `yaml:"refresh_interval" json:"refresh_interval"`
	DropClientFailureThreshold int           `yaml:"drop_client_failure_threshold" json:"drop_client_failure_threshold"`
	FailureCooldownInterval    time.Duration `yaml:"failure_cooldown_interval" json:"failure_cooldown_interval"`
	SuccessCooldownInterval    time.Duration `yaml:"success_cooldown_interval" json:"success_cooldown_interval"`
	ClientTimeout              time.Duration `yaml:"client_timeout" json:"client_timeout"`
	Socks5ProxySourceURL       string        `yaml:"socks5_proxy_source_url" json:"socks5_proxy_source_url"`
	Socks4ProxySourceURL       string        `yaml:"socks4_proxy_source_url" json:"socks4_proxy_source_url"`
	HTTPProxySourceURL         string        `yaml:"http_proxy_source_url" json:"http_proxy_source_url"`
}
