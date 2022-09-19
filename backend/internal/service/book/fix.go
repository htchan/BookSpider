package book

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/model"
)

func BookFileLocation(bk *model.Book, stConf config.SiteConfig) string {
	filename := fmt.Sprintf("%d.txt", bk.ID)
	if bk.HashCode > 0 {
		filename = fmt.Sprintf("%d-%s.txt", bk.ID, strconv.FormatInt(int64(bk.HashCode), 36))
	}
	return filepath.Join(stConf.Storage, filename)
}

func checkStorage(bk *model.Book, stConf config.SiteConfig) bool {
	isUpdated, fileExist := false, true

	if _, err := os.Stat(BookFileLocation(bk, stConf)); err != nil {
		fileExist = false
	}

	if fileExist && !bk.IsDownloaded {
		// TODO: log book status
		bk.IsDownloaded = true
		isUpdated = true
	} else if !fileExist && bk.IsDownloaded {
		bk.IsDownloaded = false
		isUpdated = true
	}
	return isUpdated
}

func Fix(bk *model.Book, stConf config.SiteConfig) (bool, error) {
	return checkStorage(bk, stConf), nil
}
