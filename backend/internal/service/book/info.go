package book

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/model"
)

func BookFileLocation(bk *model.Book, stConf *config.SiteConfig) string {
	filename := fmt.Sprintf("%d.txt", bk.ID)
	if bk.HashCode > 0 {
		filename = fmt.Sprintf("%d-v%s.txt", bk.ID, strconv.FormatInt(int64(bk.HashCode), 36))
	}
	return filepath.Join(stConf.Storage, filename)
}

func Info(bk *model.Book) string {
	bytes, err := json.Marshal(bk)
	if err != nil {
		return fmt.Sprintf("%s-%v#%v", bk.Site, bk.ID, strconv.FormatInt(int64(bk.HashCode), 36))
	}
	return string(bytes)
}

func Content(bk *model.Book, stConf *config.SiteConfig) (string, error) {
	if !bk.IsDownloaded {
		return "", errors.New("book is not download")
	}

	location := BookFileLocation(bk, stConf)
	if _, err := os.Stat(location); err != nil {
		return "", errors.New("file not found")
	}

	content, err := os.ReadFile(location)
	if err != nil {
		return "", fmt.Errorf("get book content error: %w", err)
	}

	return string(content), nil
}

func Chapters(bk *model.Book, stConf *config.SiteConfig) (model.Chapters, error) {
	content, err := Content(bk, stConf)
	if err != nil {
		return nil, fmt.Errorf("parse chapter failed: %w", err)
	}

	chapters, err := model.StringToChapters(content)
	if err != nil {
		return nil, fmt.Errorf("parse chapter failed: %w", err)
	}
	return chapters, nil
}
