package bestory

import (
	config "github.com/htchan/BookSpider/internal/config_new"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/htchan/BookSpider/internal/service"
	vendor "github.com/htchan/BookSpider/internal/vendorservice"
	"golang.org/x/sync/semaphore"
)

type VendorService struct {
}

var _ vendor.VendorService = (*VendorService)(nil)

func NewService(rpo repo.Repository, sema *semaphore.Weighted, conf config.SiteConfig) service.Service {
	panic("bestory is not available")
	// return serviceV1.NewService(Host, rpo, &VendorService{}, sema, conf)
}
