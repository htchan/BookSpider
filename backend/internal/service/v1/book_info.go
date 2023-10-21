package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/htchan/BookSpider/internal/model"
	serv "github.com/htchan/BookSpider/internal/service"
	"github.com/rs/zerolog"
)

func (s *ServiceImpl) BookInfo(ctx context.Context, bk *model.Book) string {
	bytes, err := json.Marshal(bk)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to marshal book info")

		return fmt.Sprintf("%s-%v#%v", bk.Site, bk.ID, bk.FormatHashCode())
	}

	return string(bytes)
}

func (s *ServiceImpl) BookContent(ctx context.Context, bk *model.Book) (string, error) {
	if !bk.IsDownloaded {
		return "", serv.ErrBookNotDownload
	}

	location := s.bookFileLocation(bk)
	if _, err := os.Stat(location); err != nil {
		return "", serv.ErrBookFileNotFound
	}

	content, err := os.ReadFile(location)
	if err != nil {
		return "", fmt.Errorf("read file fail: %w", err)
	}

	return string(content), nil
}

func (s *ServiceImpl) BookChapters(ctx context.Context, bk *model.Book) (model.Chapters, error) {
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

func (s *ServiceImpl) Book(ctx context.Context, id, hash string) (*model.Book, error) {
	bkID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, serv.ErrInvalidBookID
	}

	if hash == "" {
		return s.rpo.FindBookById(int(bkID))
	}

	hashcode, err := strconv.ParseInt(hash, 36, 64)
	if err != nil {
		return nil, serv.ErrInvalidHashCode
	}

	return s.rpo.FindBookByIdHash(int(bkID), int(hashcode))
}

func (s *ServiceImpl) BookGroup(
	ctx context.Context, id, hash string,
) (*model.Book, *model.BookGroup, error) {
	var (
		bkIndex  = -1
		group    model.BookGroup
		hashcode int64
		err      error
	)

	bkID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, nil, serv.ErrInvalidBookID
	}

	if hash == "" {
		group, err = s.rpo.FindBookGroupByID(int(bkID))
		if err != nil {
			return nil, nil, err
		}
	} else {
		hashcode, err = strconv.ParseInt(hash, 36, 64)
		if err != nil {
			return nil, nil, serv.ErrInvalidHashCode
		}

		group, err = s.rpo.FindBookGroupByIDHash(int(bkID), int(hashcode))
		if err != nil {
			return nil, nil, err
		}
	}

	groupIDs := make([]string, 0, len(group))

	for i, bk := range group {
		groupIDs = append(groupIDs, bk.String())

		if group[i].Site == s.name && group[i].ID == int(bkID) && group[i].HashCode == int(hashcode) {
			bkIndex = i
			break
		}
	}

	if bkIndex < 0 {
		err := errors.New("books not found")
		zerolog.Ctx(ctx).
			Error().
			Err(err).
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

func (s *ServiceImpl) QueryBooks(
	ctx context.Context, title, writer string, limit, offset int,
) ([]model.Book, error) {
	return s.rpo.FindBooksByTitleWriter(title, writer, limit, offset)
}

func (s *ServiceImpl) RandomBooks(ctx context.Context, limit int) ([]model.Book, error) {
	return s.rpo.FindBooksByRandom(limit)
}
