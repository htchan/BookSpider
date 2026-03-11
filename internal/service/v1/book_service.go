package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"github.com/htchan/BookSpider/internal/client/v1"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	"github.com/rs/zerolog"
)

type bookServiceImpl struct {
	clis        map[string]client.Client
	rpo         repo.Repository
	storagePath string
}

func NewBookService(clis map[string]client.Client, rpo repo.Repository, storagePath string) BookService {
	return &bookServiceImpl{
		clis:        clis,
		rpo:         rpo,
		storagePath: storagePath,
	}
}

func (s *bookServiceImpl) SupportBook(bk *model.Book) bool {
	_, ok := s.clis[bk.Site]
	return ok
}

func (s *bookServiceImpl) bookFileLocation(bk *model.Book) string {
	filename := fmt.Sprintf("%d.txt", bk.ID)
	if bk.HashCode > 0 {
		filename = fmt.Sprintf("%d-v%s.txt", bk.ID, bk.FormatHashCode())
	}

	return filepath.Join(s.storagePath, bk.Site, filename)
}

func (s *bookServiceImpl) bookClient(bk *model.Book) client.Client {
	return s.clis[bk.Site]
}

func isNewBook(bk *model.Book, bkInfo *client.BookInfo) bool {
	return (bk.Status == model.StatusError && bk.Error.Err == nil) ||
		(bk.Status != model.StatusError && (bk.Title != bkInfo.Title || bk.Writer.Name != bkInfo.Author || bk.Type != bkInfo.Type))
}

func isBookUpdated(bk *model.Book, bkInfo *client.BookInfo) bool {
	return bk.UpdateDate != bkInfo.UpdateDate.String() || bk.UpdateChapter != bkInfo.UpdateChapter
}

func (s *bookServiceImpl) UpdateBook(ctx context.Context, bk *model.Book) error {
	bkInfo, err := s.bookClient(bk).GetBookInfo(ctx, strconv.Itoa(bk.ID))
	if err != nil {
		if bk.Status == model.StatusError && bk.Error.Err == nil {
			bk.HashCode = model.GenerateHash()
			bk.Error.Err = err

			saveBkErr := s.rpo.CreateBook(ctx, bk)
			saveErrErr := s.rpo.SaveError(ctx, bk, bk.Error.Err)
			return errors.Join(err, saveBkErr, saveErrErr)
		} else {
			return err
		}
	}

	if isNewBook(bk, bkInfo) {
		bk.Title, bk.Writer.Name, bk.Type = bkInfo.Title, bkInfo.Author, bkInfo.Type
		bk.UpdateDate, bk.UpdateChapter = bkInfo.UpdateDate.UTC().String(), bkInfo.UpdateChapter

		bk.HashCode = model.GenerateHash()
		bk.Status = model.StatusInProgress
		if bk.IsEnd() {
			bk.Status = model.StatusEnd
		}
		bk.Error.Err = nil

		saveWriterErr := s.rpo.SaveWriter(ctx, &bk.Writer)
		saveBkErr := s.rpo.CreateBook(ctx, bk)
		saveErrErr := s.rpo.SaveError(ctx, bk, bk.Error.Err)
		if saveWriterErr != nil || saveBkErr != nil || saveErrErr != nil {
			return errors.Join(saveWriterErr, saveBkErr, saveErrErr)
		}
	} else if isBookUpdated(bk, bkInfo) {
		if bk.Status == model.StatusError {
			bk.Title, bk.Writer.Name, bk.Type = bkInfo.Title, bkInfo.Author, bkInfo.Type
		}

		bk.UpdateDate, bk.UpdateChapter = bkInfo.UpdateDate.String(), bkInfo.UpdateChapter

		bk.Status = model.StatusInProgress
		if bk.IsEnd() {
			bk.Status = model.StatusEnd
		}
		bk.Error.Err = nil

		saveWriterErr := s.rpo.SaveWriter(ctx, &bk.Writer)
		saveBkErr := s.rpo.UpdateBook(ctx, bk)
		saveErrErr := s.rpo.SaveError(ctx, bk, bk.Error.Err)
		if saveWriterErr != nil || saveBkErr != nil || saveErrErr != nil {
			return errors.Join(saveWriterErr, saveBkErr, saveErrErr)
		}
	}

	return nil
}

func (s *bookServiceImpl) downloadChapter(ctx context.Context, bk *model.Book, ch *model.Chapter) error {
	chapter, err := s.bookClient(bk).GetChapterContent(ctx, client.ChapterEntry{URL: ch.URL, Title: ch.Title})
	if err != nil {
		ch.Error = err

		return fmt.Errorf("get chapter page failed: %w", err)
	}

	ch.Title, ch.Content = chapter.Title, chapter.Body
	ch.OptimizeContent()

	return nil
}

func (s *bookServiceImpl) DownloadBook(ctx context.Context, bk *model.Book) error {
	if bk.Status != model.StatusEnd {
		return ErrBookStatusNotEnd
	} else if bk.IsDownloaded {
		return ErrBookAlreadyDownloaded
	}

	logger := zerolog.Ctx(ctx)

	logger.Info().Msg("get chapter list")

	chapterList, err := s.bookClient(bk).GetBookChapterList(ctx, strconv.Itoa(bk.ID))
	if err != nil {
		return fmt.Errorf("get chapter list failed: %w", err)
	}

	logger.Info().Msg("download chapters")
	chapters := make(model.Chapters, len(chapterList))
	failedChapterCount := 0

	for i := range chapters {
		chapters[i] = model.NewChapter(i, (chapterList)[i].URL, (chapterList)[i].Title)

		func(ch *model.Chapter) {
			chapterLogger := logger.With().
				Str("chapter_worker_id", uuid.New().String()).
				Str("chapter_url", ch.URL).
				Logger()
			err := s.downloadChapter(chapterLogger.WithContext(ctx), bk, ch)
			if err != nil {
				failedChapterCount += 1
				chapterLogger.Error().Err(err).
					Str("chapter_title", ch.Title).
					Msg("download chapter failed")
			}
		}(&chapters[i])
	}

	if failedChapterCount > 50 || failedChapterCount*10 > len(chapters) {

		return fmt.Errorf("Download chapters fail: %w (%v/%v)", ErrTooManyFailedChapters, failedChapterCount, len(chapters))
	}

	logger.Info().Msg("save chapters")
	file, err := os.Create(s.bookFileLocation(bk))
	if err != nil {
		return fmt.Errorf("create file to save chapters fail: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(bk.HeaderInfo())
	if err != nil {
		return fmt.Errorf("write book header in save chapter fail: %w", err)
	}

	for _, chapter := range chapters {
		_, err := file.WriteString(chapter.ContentString())
		if err != nil {
			return fmt.Errorf("write chapter %s in save chapters fail: %w", chapter.URL, err)
		}
	}

	logger.Info().Msg("update book is_downloaded")
	bk.IsDownloaded = true
	err = s.rpo.UpdateBook(ctx, bk)
	if err != nil {
		return fmt.Errorf("update book is_downloaded fail: %w", err)
	}

	return nil
}

func (s *bookServiceImpl) ProcessBook(ctx context.Context, bk *model.Book) error {
	updateErr := s.UpdateBook(ctx, bk)
	if updateErr != nil {
		return fmt.Errorf("update book fail: %w", updateErr)
	}

	if bk.Status == model.StatusEnd && !bk.IsDownloaded {
		downloadErr := s.DownloadBook(ctx, bk)
		if downloadErr != nil {
			return fmt.Errorf("download book fail: %w", downloadErr)
		}
	}

	return nil
}
