package hjwzw

import (
	client "github.com/htchan/BookSpider/internal/client_v2"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/htchan/BookSpider/internal/service"
	serviceV1 "github.com/htchan/BookSpider/internal/service/v1"
)

func NewService(cli client.BookClient, rpo repo.Repository) service.Service {
	return serviceV1.NewService(Host, cli, rpo, NewParser(), NewURLBuilder())
}
