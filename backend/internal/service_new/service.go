package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/htchan/BookSpider/internal/client"
	config "github.com/htchan/BookSpider/internal/config_new"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/parse"
	"github.com/htchan/BookSpider/internal/parse/goquery"
	"github.com/htchan/BookSpider/internal/repo"
	psql "github.com/htchan/BookSpider/internal/repo/psql"
	"golang.org/x/sync/semaphore"
)

type BookOperation func(*model.Book) error

type SiteOperation func() error

//go:generate mockgen -source=./$GOFILE -destination=../mock/$GOFILE -package=mock
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
	client client.Client
	parser parse.Parser
	conf   config.SiteConfig
	rpo    repo.Repostory
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
	weight *semaphore.Weighted,
	ctx *context.Context,
) (Service, error) {
	parser, err := goquery.LoadParser(&conf.GoquerySelectorsConfig)
	if err != nil {
		return nil, fmt.Errorf("load %v service failed: %w", name, err)
	}

	return &ServiceImp{
		name:   name,
		conf:   conf,
		client: client.NewClientV2(&conf, weight, ctx),
		parser: parser,
		rpo:    psql.NewRepo(name, db),
	}, nil
}
