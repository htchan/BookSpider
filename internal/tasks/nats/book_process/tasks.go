package bookprocess

import (
	"context"
	"fmt"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/service/v1"
	"github.com/nats-io/nats.go"
)

type BookProcessTasks []*BookProcessTask

func NewTaskSet(nc *nats.Conn, service service.BookService, availableSites []string) BookProcessTasks {
	updateTasks := make(BookProcessTasks, 0, len(availableSites))
	for _, site := range availableSites {
		updateTasks = append(updateTasks, NewTask(site, nc, service))
	}

	return updateTasks
}

func (tasks BookProcessTasks) Publish(ctx context.Context, bk *model.Book) ([]string, error) {
	supportedTasks := make([]string, 0, len(tasks))

	// publish for supported service
	for _, t := range tasks {
		fmt.Println(t.Site)
		if t.Service.SupportBook(bk) && t.Site == bk.Site {
			supportedTasks = append(supportedTasks, t.Site)

			err := t.Publish(ctx, bk)
			if err != nil {
				return supportedTasks, err
			}
		}
	}

	if len(supportedTasks) == 0 {
		return nil, ErrNotSupportedBook
	}

	return supportedTasks, nil
}
