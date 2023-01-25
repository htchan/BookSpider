package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/htchan/BookSpider/internal/model"
)

var (
	ErrBookIsNotDownload = errors.New("book is not download")
	ErrBookFileNotFound  = errors.New("file not found")
)

func (serv *ServiceImp) BookFileLocation(bk *model.Book) string {
	filename := fmt.Sprintf("%d.txt", bk.ID)
	if bk.HashCode > 0 {
		filename = fmt.Sprintf("%d-v%s.txt", bk.ID, strconv.FormatInt(int64(bk.HashCode), 36))
	}

	return filepath.Join(serv.conf.Storage, filename)
}

func (serv *ServiceImp) Info(bk *model.Book) string {
	bytes, err := json.Marshal(bk)
	if err != nil {
		return fmt.Sprintf("%s-%v#%v", bk.Site, bk.ID, strconv.FormatInt(int64(bk.HashCode), 36))
	}
	return string(bytes)
}

func (serv *ServiceImp) BookContent(bk *model.Book) (string, error) {
	if !bk.IsDownloaded {
		return "", ErrBookIsNotDownload
	}

	location := serv.BookFileLocation(bk)
	if _, err := os.Stat(location); err != nil {
		return "", ErrBookFileNotFound
	}

	content, err := os.ReadFile(location)
	if err != nil {
		return "", fmt.Errorf("get book content error: %w", err)
	}

	return string(content), nil
}

func (serv *ServiceImp) Chapters(bk *model.Book) (model.Chapters, error) {
	content, err := serv.BookContent(bk)
	if err != nil {
		return nil, fmt.Errorf("load content failed: %w", err)
	}

	chapters, err := model.StringToChapters(content)
	if err != nil {
		return nil, fmt.Errorf("parse chapter failed: %w", err)
	}

	return chapters, nil
}
