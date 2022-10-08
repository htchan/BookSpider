package book

import (
	"fmt"

	"github.com/htchan/BookSpider/internal/client"
	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/model"
)

func Process(bk *model.Book, bkConf *config.BookConfig, stConf *config.SiteConfig, c *client.CircuitBreakerClient) (bool, error) {
	isUpdated, err := Update(bk, bkConf, stConf, c)
	if err != nil {
		return isUpdated, fmt.Errorf("Process book error: %w", err)
	}

	isValidateUpdated, err := Validate(bk)
	isUpdated = isUpdated || isValidateUpdated
	if err != nil {
		return isUpdated, fmt.Errorf("Process book error: %w", err)
	}

	isFixUpdated, err := Fix(bk, stConf)
	isUpdated = isUpdated || isFixUpdated
	if err != nil {
		return isUpdated, fmt.Errorf("Process book error: %w", err)
	}

	fmt.Println(bk.Status, bk.IsDownloaded)
	isDownloadpdated, err := Download(bk, bkConf, stConf, c)
	isUpdated = isUpdated || isDownloadpdated
	if err != nil {
		return isUpdated, fmt.Errorf("Process book error: %w", err)
	}

	isFixUpdated, err = Fix(bk, stConf)
	isUpdated = isUpdated || isFixUpdated
	if err != nil {
		return isUpdated, fmt.Errorf("Process book error: %w", err)
	}

	return isUpdated, nil
}
