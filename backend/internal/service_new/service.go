package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	client "github.com/htchan/BookSpider/internal/client_v2"
	circuitbreaker "github.com/htchan/BookSpider/internal/client_v2/circuit_breaker"
	"github.com/htchan/BookSpider/internal/client_v2/retry"
	"github.com/htchan/BookSpider/internal/client_v2/simple"
	config "github.com/htchan/BookSpider/internal/config_new"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/parse"
	"github.com/htchan/BookSpider/internal/parse/goquery"
	"github.com/htchan/BookSpider/internal/repo"
	sqlc "github.com/htchan/BookSpider/internal/repo/sqlc"
	"golang.org/x/sync/semaphore"
)

type BookOperation func(*model.Book) error

type SiteOperation func() error

//go:generate mockgen -destination=../mock/service/v2/service.go -package=mockservice . Service
type Service interface {
	Name() string
	Backup() error
	ValidateEnd() error
	Download() error
	Explore() error
	PatchDownloadStatus() error
	PatchMissingRecords() error
	Process() error
	Update() error
	CheckAvailability() error

	DownloadBook(*model.Book) error
	ExploreBook(*model.Book) error
	ProcessBook(*model.Book) error
	UpdateBook(*model.Book) error
	ValidateBookEnd(*model.Book) error

	BookInfo(*model.Book) string
	BookContent(*model.Book) (string, error)

	Book(id int, hash string) (*model.Book, error)
	BookGroup(id int, hash string) (*model.Book, *model.BookGroup, error)
	QueryBooks(title, writer string, limit, offset int) ([]model.Book, error)
	RandomBooks(limit int) ([]model.Book, error)

	Stats() repo.Summary
	DBStats() sql.DBStats
}

var (
	ErrInvalidSite = errors.New("invalid site")
)

type ServiceImp struct {
	name   string
	ctx    context.Context
	client client.BookClient
	sema   *semaphore.Weighted
	parser parse.Parser
	conf   config.SiteConfig
	rpo    repo.Repository
}

var _ Service = (*ServiceImp)(nil)

func (serv ServiceImp) Name() string {
	return serv.name
}
func (serv *ServiceImp) Stats() repo.Summary {
	return serv.rpo.Stats()
}
func (serv *ServiceImp) DBStats() sql.DBStats {
	return serv.rpo.DBStats()
}

func LoadService(
	name string,
	conf config.SiteConfig,
	db *sql.DB,
	ctx context.Context,
	sema *semaphore.Weighted,
) (Service, error) {
	parser, err := goquery.LoadParser(&conf.GoquerySelectorsConfig)
	if err != nil {
		return nil, fmt.Errorf("load %v service failed: %w", name, err)
	}

	return &ServiceImp{
		name: name,
		conf: conf,
		ctx:  ctx,
		client: retry.NewClient(
			&conf.ClientConfig.Retry,
			circuitbreaker.NewClient(
				&conf.ClientConfig.CircuitBreaker,
				simple.NewClient(&conf.ClientConfig.Simple),
			),
		),
		sema:   sema,
		parser: parser,
		rpo:    sqlc.NewRepo(name, db),
	}, nil
}
