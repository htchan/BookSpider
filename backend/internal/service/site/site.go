package site

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/go-cmp/cmp"
	"github.com/htchan/BookSpider/internal/client"
	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	psql "github.com/htchan/BookSpider/internal/repo/psql"
	"golang.org/x/sync/semaphore"
)

type Site struct {
	Name   string
	rp     repo.Repostory
	BkConf *config.BookConfig
	StConf *config.SiteConfig
	Client *client.CircuitBreakerClient
}

func NewSite(name string, bkConf *config.BookConfig, stConf *config.SiteConfig, clientConf *config.CircuitBreakerClientConfig, commonWeighted *semaphore.Weighted, commonCtx *context.Context) (*Site, error) {
	db, err := psql.OpenDatabase(name)
	if err != nil {
		return nil, fmt.Errorf("new site error: %w", err)
	}
	c := client.NewClient(*clientConf, commonWeighted, commonCtx)
	return &Site{
		Name:   name,
		rp:     psql.NewRepo(name, db),
		BkConf: bkConf,
		StConf: stConf,
		Client: &c,
	}, nil
}
func MockSite(name string, rp repo.Repostory, bkConf *config.BookConfig, stConf *config.SiteConfig, c *client.CircuitBreakerClient) *Site {
	return &Site{
		Name:   name,
		rp:     rp,
		BkConf: bkConf,
		StConf: stConf,
		Client: c,
	}
}

func containsString(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}
func LoadSitesFromConfigDirectory(directory string, enabledSiteNames []string) (map[string]*Site, error) {
	siteConfigs, err := config.LoadSiteConfigs(directory)
	if err != nil {
		return nil, fmt.Errorf("load site config: %w", err)
	}
	// TODO: load book config
	bookConfigs, err := config.LoadBookConfigs(directory)
	if err != nil {
		return nil, fmt.Errorf("load book config: %w", err)
	}
	// TODO: load client config
	clientConfigs, err := config.LoadClientConfigs(directory)
	if err != nil {
		return nil, fmt.Errorf("load client config: %w", err)
	}

	sites := make(map[string]*Site)
	for siteName, stConf := range siteConfigs {
		if !containsString(enabledSiteNames, siteName) {
			continue
		}
		bkConf := bookConfigs[stConf.BookKey]
		clientConf := clientConfigs[stConf.ClientKey]
		sites[siteName], err = NewSite(siteName, bkConf, stConf, clientConf, nil, nil)
		if err != nil {
			log.Printf("fail to load site %v: %v", siteName, err)
		}
		err = sites[siteName].Migrate()
		if err != nil {
			log.Printf("fail to migrate for site: %v; err: %v", siteName, err)
		}
	}
	return sites, nil
}

func (st *Site) Info() repo.Summary {
	return st.rp.Stats()
}

func (st *Site) Stat() sql.DBStats {
	return st.rp.DBStats()
}

func (st *Site) Close() error {
	st.Client.Close()
	return st.rp.Close()
}

func (st Site) Equal(compare Site) bool {
	return st.Name == compare.Name &&
		cmp.Equal(&st.BkConf, &compare.BkConf) &&
		cmp.Equal(&st.StConf, &compare.StConf)
}

func (st *Site) CreateBook(bk *model.Book) error {
	return st.rp.CreateBook(bk)
}
func (st *Site) UpdateBook(bk *model.Book) error {
	return st.rp.UpdateBook(bk)
}

func (st *Site) SaveWriter(bk *model.Book) error {
	return st.rp.SaveWriter(&bk.Writer)
}
func (st *Site) SaveError(bk *model.Book) error {
	return st.rp.SaveError(bk, bk.Error)
}

func (st *Site) Migrate() error {
	return st.rp.Migrate()
}
