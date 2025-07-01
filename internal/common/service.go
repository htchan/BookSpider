package common

import (
	"database/sql"
	"slices"

	"github.com/htchan/BookSpider/internal/config/v2"
	repo "github.com/htchan/BookSpider/internal/repo/sqlc"
	"github.com/htchan/BookSpider/internal/service"
	service_v1 "github.com/htchan/BookSpider/internal/service/v1"
	"github.com/htchan/BookSpider/internal/vendorservice/baling"
	"github.com/htchan/BookSpider/internal/vendorservice/bestory"
	"github.com/htchan/BookSpider/internal/vendorservice/ck101"
	"github.com/htchan/BookSpider/internal/vendorservice/hjwzw"
	"github.com/htchan/BookSpider/internal/vendorservice/uukanshu"
	"github.com/htchan/BookSpider/internal/vendorservice/xbiquge"
	"github.com/htchan/BookSpider/internal/vendorservice/xqishu"
	"golang.org/x/sync/semaphore"
)

func LoadServices(vendors []string, db *sql.DB, siteConf map[string]config.SiteConfig, maxThreads int64) map[string]service.Service {
	result := make(map[string]service.Service)

	publicSema := semaphore.NewWeighted(maxThreads)
	rpo := repo.NewRepo(db)

	if slices.Contains(vendors, baling.Host) {
		result[baling.Host] = baling.NewService(rpo, publicSema, siteConf[baling.Host])
	}

	if slices.Contains(vendors, bestory.Host) {
		result[bestory.Host] = bestory.NewService(rpo, publicSema, siteConf[bestory.Host])
	}

	if slices.Contains(vendors, ck101.Host) {
		result[ck101.Host] = ck101.NewService(rpo, publicSema, siteConf[ck101.Host])
	}

	if slices.Contains(vendors, hjwzw.Host) {
		result[hjwzw.Host] = hjwzw.NewService(rpo, publicSema, siteConf[hjwzw.Host])
	}

	if slices.Contains(vendors, xbiquge.Host) {
		result[xbiquge.Host] = xbiquge.NewService(rpo, publicSema, siteConf[xbiquge.Host])
	}

	if slices.Contains(vendors, xqishu.Host) {
		result[xqishu.Host] = xqishu.NewService(rpo, publicSema, siteConf[xqishu.Host])
	}

	if slices.Contains(vendors, uukanshu.Host) {
		result[uukanshu.Host] = uukanshu.NewService(rpo, publicSema, siteConf[uukanshu.Host])
	}

	return result
}

func LoadReadDataService(db *sql.DB, siteConf map[string]config.SiteConfig) service.ReadDataService {
	rpo := repo.NewRepo(db)

	return service_v1.NewReadDataService(rpo, siteConf)
}
