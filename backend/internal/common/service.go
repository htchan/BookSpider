package common

import (
	"database/sql"
	"slices"

	config "github.com/htchan/BookSpider/internal/config_new"
	repo "github.com/htchan/BookSpider/internal/repo/sqlc"
	"github.com/htchan/BookSpider/internal/service"
	"github.com/htchan/BookSpider/internal/vendorservice/baling"
	"github.com/htchan/BookSpider/internal/vendorservice/bestory"
	"github.com/htchan/BookSpider/internal/vendorservice/ck101"
	"github.com/htchan/BookSpider/internal/vendorservice/hjwzw"
	"github.com/htchan/BookSpider/internal/vendorservice/uukanshu"
	"github.com/htchan/BookSpider/internal/vendorservice/xbiquge"
	"github.com/htchan/BookSpider/internal/vendorservice/xqishu"
	"golang.org/x/sync/semaphore"
)

func LoadServices(vendors []string, db *sql.DB, conf *config.Config) map[string]service.Service {
	result := make(map[string]service.Service)

	publicSema := semaphore.NewWeighted(int64(conf.BatchConfig.MaxWorkingThreads))

	if slices.Contains(vendors, baling.Host) {
		rpo := repo.NewRepo(baling.Host, db)

		result[baling.Host] = baling.NewService(rpo, publicSema, conf.SiteConfigs[baling.Host])
	}

	if slices.Contains(vendors, bestory.Host) {
		rpo := repo.NewRepo(bestory.Host, db)

		result[bestory.Host] = bestory.NewService(rpo, publicSema, conf.SiteConfigs[bestory.Host])
	}

	if slices.Contains(vendors, ck101.Host) {
		rpo := repo.NewRepo(ck101.Host, db)

		result[ck101.Host] = ck101.NewService(rpo, publicSema, conf.SiteConfigs[ck101.Host])
	}

	if slices.Contains(vendors, hjwzw.Host) {
		rpo := repo.NewRepo(hjwzw.Host, db)

		result[hjwzw.Host] = hjwzw.NewService(rpo, publicSema, conf.SiteConfigs[hjwzw.Host])
	}

	if slices.Contains(vendors, xbiquge.Host) {
		rpo := repo.NewRepo(xbiquge.Host, db)

		result[xbiquge.Host] = xbiquge.NewService(rpo, publicSema, conf.SiteConfigs[xbiquge.Host])
	}

	if slices.Contains(vendors, xqishu.Host) {
		rpo := repo.NewRepo(xqishu.Host, db)

		result[xqishu.Host] = xqishu.NewService(rpo, publicSema, conf.SiteConfigs[xqishu.Host])
	}

	if slices.Contains(vendors, uukanshu.Host) {
		rpo := repo.NewRepo(uukanshu.Host, db)

		result[uukanshu.Host] = uukanshu.NewService(rpo, publicSema, conf.SiteConfigs[uukanshu.Host])
	}

	return result
}
