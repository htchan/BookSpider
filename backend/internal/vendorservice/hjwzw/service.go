package hjwzw

import (
	config "github.com/htchan/BookSpider/internal/config_new"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/htchan/BookSpider/internal/service"
	serviceV1 "github.com/htchan/BookSpider/internal/service/v1"
	"golang.org/x/sync/semaphore"
)

func NewService(rpo repo.Repository, sema *semaphore.Weighted, conf config.SiteConfig) service.Service {
	return serviceV1.NewService(Host, rpo, NewParser(), NewURLBuilder(), sema, conf)
}
