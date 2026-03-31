package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/htchan/goclient"
	client "github.com/htchan/BookSpider/internal/client/v2"
	"github.com/htchan/BookSpider/internal/config/v2"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	serv "github.com/htchan/BookSpider/internal/service"
	vendor "github.com/htchan/BookSpider/internal/vendorservice"
	circuitbreaker "github.com/htchan/goclient/middlewares/circuit_breaker"
	ratelimit "github.com/htchan/goclient/middlewares/rate_limit"
	"github.com/htchan/goclient/middlewares/retry"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/semaphore"
)

type ServiceImpl struct {
	name          string
	cli           client.BookClient
	rpo           repo.Repository
	vendorService vendor.VendorService

	conf config.SiteConfig
	sema *semaphore.Weighted
}

var _ serv.Service = (*ServiceImpl)(nil)

func isServerError(req *http.Request, resp *http.Response, err error) bool {
	if err != nil {
		return true
	}
	return resp != nil && resp.StatusCode >= 500
}

func retryOnError(req *http.Request, resp *http.Response, err error) bool {
	if err != nil {
		return true
	}
	if resp != nil && resp.StatusCode >= 500 {
		return true
	}
	return false
}

func newRetryIntervalCalculator(intervalType string, baseInterval time.Duration) retry.IntervalCalculator {
	switch intervalType {
	case "linear":
		return retry.LinearRetryInterval(baseInterval)
	case "exponential":
		return retry.ExponentialRetryInterval(baseInterval)
	default:
		return retry.StaticRetryInterval(baseInterval)
	}
}

func NewService(
	name string, rpo repo.Repository,
	vendorService vendor.VendorService,
	sema *semaphore.Weighted, conf config.SiteConfig,
) *ServiceImpl {
	queue := ratelimit.NewQueue(conf.ClientConfig.RateLimit.QueueSize)

	breaker := circuitbreaker.NewCircuitBreaker(
		conf.ClientConfig.CircuitBreaker.FailureThreshold,
		conf.ClientConfig.CircuitBreaker.SuccessThreshold,
		conf.ClientConfig.CircuitBreaker.RecoverDuration,
		isServerError,
		circuitbreaker.WithOnStateChange(func(from, to circuitbreaker.State) {
			if to == circuitbreaker.StateOpen {
				queue.Resize(conf.ClientConfig.RateLimit.QueueSize / 2)
			} else if to == circuitbreaker.StateClosed {
				queue.Resize(conf.ClientConfig.RateLimit.QueueSize)
			}
		}),
	)

	cli := goclient.NewClient(
		goclient.WithMiddlewares(
			retry.NewRetryMiddleware(
				conf.ClientConfig.Retry.MaxRetries,
				retryOnError,
				newRetryIntervalCalculator(
					conf.ClientConfig.Retry.IntervalType,
					conf.ClientConfig.Retry.BaseInterval,
				),
			),
			circuitbreaker.NewCircuitBreakerMiddleware(breaker),
			ratelimit.NewRateLimitMiddleware(queue, conf.ClientConfig.RateLimit.Interval),
		),
	)

	return &ServiceImpl{
		name: name,
		cli: client.NewClient(
			cli,
			conf.DecodeMethod,
		),
		rpo:           rpo,
		vendorService: vendorService,

		sema: sema,
		conf: conf,
	}
}

func (s *ServiceImpl) Name() string {
	return s.name
}

func (s *ServiceImpl) bookFileLocation(bk *model.Book) string {
	filename := fmt.Sprintf("%d.txt", bk.ID)
	if bk.HashCode > 0 {
		filename = fmt.Sprintf("%d-v%s.txt", bk.ID, bk.FormatHashCode())
	}

	return filepath.Join(s.conf.Storage, filename)
}

func (s *ServiceImpl) checkBookStorage(bk *model.Book, stats *serv.PatchStorageStats) bool {
	isDownloadUpdated, fileExist := false, true
	if stats == nil {
		stats = new(serv.PatchStorageStats)
	}

	if _, err := os.Stat(s.bookFileLocation(bk)); err != nil {
		fileExist = false
	}

	if fileExist && !bk.IsDownloaded {
		log.Info().Str("book", bk.String()).Msg("file exist for not downloaded book")
		bk.IsDownloaded = true
		isDownloadUpdated = true
		stats.FileExist.Add(1)
	} else if !fileExist && bk.IsDownloaded {
		log.Info().Str("book", bk.String()).Msg("file not exist for downloaded book")
		bk.IsDownloaded = false
		isDownloadUpdated = true
		stats.FileMissing.Add(1)
	}

	return isDownloadUpdated
}

func (s *ServiceImpl) PatchDownloadStatus(ctx context.Context, stats *serv.PatchStorageStats) error {
	if stats == nil {
		stats = new(serv.PatchStorageStats)
	}

	bks, err := s.rpo.FindAllBooks(ctx, s.name)
	if err != nil {
		return fmt.Errorf("patch download status fail: %w", err)
	}

	var wg sync.WaitGroup
	zerolog.Ctx(ctx).Info().Str("site", s.name).Msg("update books is_downloaded by storage")

	for bk := range bks {
		s.sema.Acquire(ctx, 1)
		wg.Add(1)

		go func(bk *model.Book) {
			defer wg.Done()
			defer s.sema.Release(1)

			needUpdate := s.checkBookStorage(bk, stats)
			if needUpdate {
				err := s.rpo.UpdateBook(ctx, bk)
				if err != nil {
					zerolog.Ctx(ctx).Error().Err(err).
						Str("site", s.name).
						Int("bk_id", bk.ID).
						Str("bk_hash_code", bk.FormatHashCode()).
						Msg("update book is_downloaded fail")
				}
			}
		}(&bk)
	}

	wg.Wait()

	return nil
}

func (s *ServiceImpl) PatchMissingRecords(ctx context.Context, stats *serv.UpdateStats) error {
	zerolog.Ctx(ctx).Info().Msg("patch missing records")

	if stats == nil {
		stats = new(serv.UpdateStats)
	}

	var wg sync.WaitGroup
	allBkIDs, err := s.rpo.FindAllBookIDs(ctx, s.name)
	if err != nil {
		return fmt.Errorf("find all book ids fail: %w", err)
	}

	missingIDs := s.vendorService.FindMissingIds(allBkIDs)
	for _, bookID := range missingIDs {
		s.sema.Acquire(ctx, 1)
		wg.Add(1)
		stats.Total.Add(1)

		go func(id int) {
			defer s.sema.Release(1)
			defer wg.Done()

			zerolog.Ctx(ctx).Error().Err(err).Int("id", id).Msg("book not exist in database")
			bk := model.NewBook(s.name, id)
			s.ExploreBook(ctx, &bk, stats)
		}(bookID)
	}
	wg.Wait()

	return nil
}

func (s *ServiceImpl) CheckAvailability(ctx context.Context) error {
	body, err := s.cli.Get(ctx, s.vendorService.AvailabilityURL())
	if err != nil {
		return fmt.Errorf("get availability page failed: %w", err)
	}

	if !s.vendorService.IsAvailable(body) {
		return serv.ErrUnavailable
	}

	return nil
}
