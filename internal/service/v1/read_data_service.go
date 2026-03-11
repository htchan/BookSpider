package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/rs/zerolog"
)

type readDataServiceImpl struct {
	rpo         repo.Repository
	storagePath string
}

var _ ReadDataService = (*readDataServiceImpl)(nil)

func NewReadDataService(rpo repo.Repository, storagePath string) ReadDataService {
	return &readDataServiceImpl{
		rpo:         rpo,
		storagePath: storagePath,
	}
}

func (s *readDataServiceImpl) bookFileLocation(bk *model.Book) string {
	filename := fmt.Sprintf("%d.txt", bk.ID)
	if bk.HashCode > 0 {
		filename = fmt.Sprintf("%d-v%s.txt", bk.ID, bk.FormatHashCode())
	}

	return filepath.Join(s.storagePath, bk.Site, filename)
}

func (s *readDataServiceImpl) Book(ctx context.Context, site, id, hash string) (*model.Book, error) {
	bkID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, ErrInvalidBookID
	}

	if hash == "" {
		return s.rpo.FindBookById(ctx, site, int(bkID))
	}

	hashcode, err := strconv.ParseInt(hash, 36, 64)
	if err != nil {
		return nil, ErrInvalidHashCode
	}

	return s.rpo.FindBookByIdHash(ctx, site, int(bkID), int(hashcode))
}

func (s *readDataServiceImpl) BookContent(ctx context.Context, bk *model.Book) (string, error) {
	if !bk.IsDownloaded {
		return "", ErrBookNotDownload
	}

	location := s.bookFileLocation(bk)
	if _, err := os.Stat(location); err != nil {
		return "", ErrBookFileNotFound
	}

	content, err := os.ReadFile(location)
	if err != nil {
		return "", fmt.Errorf("read file fail: %w", err)
	}

	return string(content), nil
}

func (s *readDataServiceImpl) BookChapters(ctx context.Context, bk *model.Book) (model.Chapters, error) {
	content, err := s.BookContent(ctx, bk)
	if err != nil {
		return nil, fmt.Errorf("load content failed: %w", err)
	}

	chapters, err := model.StringToChapters(content)
	if err != nil {
		return nil, fmt.Errorf("parse chapter failed: %w", err)
	}

	return chapters, nil

}

func (s *readDataServiceImpl) BookGroup(ctx context.Context, site, id, hash string) (*model.Book, *model.BookGroup, error) {
	var (
		bkIndex  = -1
		group    model.BookGroup
		hashcode int64
		err      error
	)

	bkID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, nil, ErrInvalidBookID
	}

	if hash == "" {
		group, err = s.rpo.FindBookGroupByID(ctx, site, int(bkID))
		if err != nil {
			return nil, nil, err
		}
	} else {
		hashcode, err = strconv.ParseInt(hash, 36, 64)
		if err != nil {
			return nil, nil, ErrInvalidHashCode
		}

		group, err = s.rpo.FindBookGroupByIDHash(ctx, site, int(bkID), int(hashcode))
		if err != nil {
			return nil, nil, err
		}
	}

	groupIDs := make([]string, 0, len(group))

	for i, bk := range group {
		groupIDs = append(groupIDs, bk.String())

		if group[i].Site == site && group[i].ID == int(bkID) && group[i].HashCode == int(hashcode) {
			bkIndex = i
			break
		}
	}

	if bkIndex < 0 {
		err := errors.New("books not found")
		zerolog.Ctx(ctx).
			Error().
			Err(err).
			Str("site", site).
			Str("id", id).
			Int64("hashcode", hashcode).
			Strs("book group id", groupIDs).
			Msg("find book group failed")
		return nil, nil, err
	}

	bk := group[bkIndex]
	if bkIndex+1 >= len(group) {
		group = group[:bkIndex]
	} else {
		group = append(group[:bkIndex], group[bkIndex+1:]...)
	}

	return &bk, &group, nil
}

func (s *readDataServiceImpl) SearchBooks(ctx context.Context, title, writer string, limit, offset int) ([]model.Book, error) {
	return s.rpo.FindBooksByTitleWriter(ctx, title, writer, limit, offset)
}

func (s *readDataServiceImpl) RandomBooks(ctx context.Context, limit int) ([]model.Book, error) {
	return s.rpo.FindBooksByRandom(ctx, limit)
}

func (s *readDataServiceImpl) Stats(ctx context.Context, site string) repo.Summary {
	return s.rpo.Stats(ctx, site)
}

func (s *readDataServiceImpl) DBStats(ctx context.Context) sql.DBStats {
	return s.rpo.DBStats(ctx)
}
