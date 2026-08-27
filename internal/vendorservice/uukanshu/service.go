package uukanshu

import (
	"github.com/htchan/BookSpider/internal/config/v2"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/htchan/BookSpider/internal/service"
	serviceV1 "github.com/htchan/BookSpider/internal/service/v1"
	vendor "github.com/htchan/BookSpider/internal/vendorservice"
	"golang.org/x/sync/semaphore"
)

type VendorService struct {
	*vendor.BaseClient
}

var _ vendor.VendorService = (*VendorService)(nil)

func NewService(rpo repo.Repository, sema *semaphore.Weighted, conf config.SiteConfig) service.Service {
	baseClient, vendorSema := vendor.NewBaseClientFromConfig(vendor.ClientConfig{
		Name:                           Host,
		DecodeMethod:                   conf.DecodeMethod,
		RequestTimeout:                 conf.RequestTimeout,
		RateLimitQueueSize:             conf.ClientConfig.RateLimit.QueueSize,
		RateLimitInterval:              conf.ClientConfig.RateLimit.Interval,
		CircuitBreakerFailureThreshold: conf.ClientConfig.CircuitBreaker.FailureThreshold,
		CircuitBreakerSuccessThreshold: conf.ClientConfig.CircuitBreaker.SuccessThreshold,
		CircuitBreakerRecoverDuration:  conf.ClientConfig.CircuitBreaker.RecoverDuration,
		CircuitBreakerOpenQueueRatio:   conf.ClientConfig.CircuitBreaker.OpenQueueRatio,
		RetryMaxRetries:                conf.ClientConfig.Retry.MaxRetries,
		RetryBaseInterval:              conf.ClientConfig.Retry.BaseInterval,
		RetryIntervalType:              conf.ClientConfig.Retry.IntervalType,
	})

	return serviceV1.NewService(Host, rpo, &VendorService{BaseClient: baseClient}, sema, conf, vendorSema)
}
