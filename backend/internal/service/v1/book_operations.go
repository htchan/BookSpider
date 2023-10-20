package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/repo"
	serv "github.com/htchan/BookSpider/internal/service"
	vendor "github.com/htchan/BookSpider/internal/vendorservice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/semaphore"
)

func isNewBook(bk *model.Book, bkInfo *vendor.BookInfo) bool {
	return bk.Status != model.Error && (bk.Title != bkInfo.Title || bk.Writer.Name != bkInfo.Writer || bk.Type != bkInfo.Type)
}

func isBookUpdated(bk *model.Book, bkInfo *vendor.BookInfo) bool {
	return bk.UpdateDate != bkInfo.UpdateDate || bk.UpdateChapter != bkInfo.UpdateChapter
}

func (s *ServiceImpl) UpdateBook(ctx context.Context, bk *model.Book) error {
	body, err := s.cli.Get(ctx, s.urlBuilder.BookURL(strconv.FormatInt(int64(bk.ID), 10)))
	if err != nil {
		return fmt.Errorf("get book page failed: %w", err)
	}

	bkInfo, err := s.parser.ParseBook(body)
	if err != nil {
		return fmt.Errorf("parse book page failed: %w", err)
	}

	logger := zerolog.Ctx(ctx).With().
		Int("bk_id", bk.ID).
		Str("bk_hash_code", bk.FormatHashCode()).
		Str("bk_title", bkInfo.Title).
		Logger()

	if isNewBook(bk, bkInfo) {
		bk.Title, bk.Writer.Name, bk.Type = bkInfo.Title, bkInfo.Writer, bkInfo.Type
		bk.UpdateDate, bk.UpdateChapter = bkInfo.UpdateDate, bkInfo.UpdateChapter

		bk.HashCode = model.GenerateHash()
		bk.Status = model.InProgress
		bk.Error = nil

		saveWriterErr := s.rpo.SaveWriter(&bk.Writer)
		saveBkErr := s.rpo.CreateBook(bk)
		saveErrErr := s.rpo.SaveError(bk, bk.Error)
		if saveWriterErr != nil || saveBkErr != nil || saveErrErr != nil {
			return errors.Join(saveWriterErr, saveBkErr)
		}

		logger.Debug().Msg("new book found")
	} else if isBookUpdated(bk, bkInfo) {
		if bk.Status == model.Error {
			bk.Title, bk.Writer.Name, bk.Type = bkInfo.Title, bkInfo.Writer, bkInfo.Type
		}

		bk.UpdateDate, bk.UpdateChapter = bkInfo.UpdateDate, bkInfo.UpdateChapter

		bk.Status = model.InProgress
		bk.Error = nil

		if bk.Status == model.Error {
			bk.Status = model.InProgress
		}

		saveWriterErr := s.rpo.SaveWriter(&bk.Writer)
		saveBkErr := s.rpo.UpdateBook(bk)
		saveErrErr := s.rpo.SaveError(bk, bk.Error)
		if saveWriterErr != nil || saveBkErr != nil || saveErrErr != nil {
			return errors.Join(saveWriterErr, saveBkErr)
		}

		logger.Debug().Msg("book updated")
	} else {
		logger.Debug().Msg("book not updated")
	}

	return nil
}

func (s *ServiceImpl) Update(ctx context.Context) error {
	var wg sync.WaitGroup

	bkChan, err := s.rpo.FindBooksForUpdate()
	if err != nil {
		return fmt.Errorf("fail to load books from DB: %w", err)
	}

	for bk := range bkChan {
		bk := bk
		s.sema.Acquire(ctx, 1)
		wg.Add(1)

		go func(bk *model.Book) {
			defer wg.Done()
			defer s.sema.Release(1)

			err := s.UpdateBook(ctx, bk)
			if err != nil {
				zerolog.Ctx(ctx).Error().Err(err).
					Int("bk_id", bk.ID).
					Str("bk_hash_code", bk.FormatHashCode()).
					Msg("update book failed")
			}
		}(&bk)
	}

	wg.Wait()

	return nil
}

func (s *ServiceImpl) ExploreBook(ctx context.Context, bk *model.Book) error {
	if bk.Status != model.Error {
		return serv.ErrBookStatusNotError
	}

	//TODO: find a new method to check if we should create the book
	isNew := bk.Error == nil
	if isNew {
		s.rpo.CreateBook(bk)
	}

	logger := zerolog.Ctx(ctx).With().
		Int("bk_id", bk.ID).
		Str("bk_hash_code", bk.FormatHashCode()).
		Logger()

	err := s.UpdateBook(ctx, bk)
	if err != nil {
		logger.Error().Err(err).Msg("explore book fail")
		bk.Error = err
		saveErr := s.rpo.SaveError(bk, bk.Error)
		if saveErr != nil {
			return fmt.Errorf("save error fail: %w", saveErr)
		}

		return fmt.Errorf("explore book fail: %w", err)
	}

	return err
}

func (s *ServiceImpl) Explore(ctx context.Context) error {
	summary := s.rpo.Stats()
	errorCount := 0

	var wg sync.WaitGroup

	for i := summary.LatestSuccessID + 1; i <= summary.MaxBookID && errorCount < s.conf.MaxExploreError; i++ {
		i := i

		s.sema.Acquire(ctx, 1)
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			defer s.sema.Release(1)

			bk, err := s.rpo.FindBookById(id)
			if err != nil {
				errorCount += 1
				return
			}

			err = s.ExploreBook(ctx, bk)
			if err != nil {
				errorCount += 1
			} else {
				errorCount = 0
			}
		}(i)
	}

	wg.Wait()

	for i := summary.MaxBookID + 1; errorCount < s.conf.MaxExploreError; i++ {
		i := i

		s.sema.Acquire(ctx, 1)
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			defer s.sema.Release(1)

			bk := model.NewBook(s.name, i)
			err := s.ExploreBook(ctx, &bk)
			if err != nil {
				errorCount += 1
			} else {
				errorCount = 0
			}
		}(i)
	}

	wg.Wait()

	return nil
}

func (s *ServiceImpl) downloadChapter(ctx context.Context, ch *model.Chapter) error {
	body, err := s.cli.Get(ctx, ch.URL)
	if err != nil {
		ch.Error = err

		return fmt.Errorf("get chapter page failed: %w", err)
	}

	chapter, err := s.parser.ParseChapter(body)
	if err != nil {
		ch.Error = err

		return fmt.Errorf("parse chapter page failed: %w", err)
	}

	ch.Title, ch.Content = chapter.Title, chapter.Body

	return nil
}

func (s *ServiceImpl) DownloadBook(ctx context.Context, bk *model.Book) error {
	if bk.Status != model.End {
		return serv.ErrBookStatusNotEnd
	} else if bk.IsDownloaded {
		return serv.ErrBookAlreadyDownloaded
	}

	logger := zerolog.Ctx(ctx).With().Str("site", s.name).Int("bk_id", bk.ID).Str("bk_hash_code", bk.FormatHashCode()).Logger()

	logger.Info().Msg("get chapter list")

	body, err := s.cli.Get(ctx, s.urlBuilder.ChapterListURL(strconv.FormatInt(int64(bk.ID), 10)))
	if err != nil {
		return fmt.Errorf("get chapter list failed: %w", err)
	}

	chapterList, err := s.parser.ParseChapterList(body)
	if err != nil {
		return fmt.Errorf("parse chapter list failed: %w", err)
	}

	logger.Info().Msg("download chapters")
	chapters := make(model.Chapters, len(chapterList))
	var wg sync.WaitGroup
	failedChapterCount := 0

	for i := range chapters {
		i := i
		chapters[i] = model.NewChapter(i, (chapterList)[i].URL, (chapterList)[i].Title)
		wg.Add(1)
		s.sema.Acquire(ctx, 1)

		go func(ch *model.Chapter) {
			defer wg.Done()
			defer s.sema.Release(1)

			err := s.downloadChapter(ctx, ch)
			if err != nil {
				failedChapterCount += 1
				logger.Error().Err(err).
					Str("chapter_url", ch.URL).
					Str("chapter_title", ch.Title).
					Msg("download chapter failed")
			}
		}(&chapters[i])
	}

	wg.Wait()

	if failedChapterCount > 50 || failedChapterCount*10 > len(chapters) {
		return fmt.Errorf("Download chapters fail: %w (%v/%v)", serv.ErrTooManyFailedChapters, failedChapterCount, len(chapters))
	}

	logger.Info().Msg("save chapters")
	file, err := os.Create(s.bookFileLocation(bk))
	if err != nil {
		return fmt.Errorf("save book fail: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(bk.HeaderInfo())
	if err != nil {
		return fmt.Errorf("save book fail: %w", err)
	}

	for _, chapter := range chapters {
		_, err := file.WriteString(chapter.ContentString())
		if err != nil {
			return fmt.Errorf("save book fail: %w", err)
		}
	}

	logger.Info().Msg("update book is_downloaded")
	bk.IsDownloaded = true
	err = s.rpo.UpdateBook(bk)
	if err != nil {
		return fmt.Errorf("update book download fail: %w", err)
	}

	return nil
}

func (s *ServiceImpl) Download(ctx context.Context) error {
	se := semaphore.NewWeighted(int64(s.conf.MaxDownloadConcurrency))
	var wg sync.WaitGroup

	bkChan, err := s.rpo.FindBooksForDownload()
	if err != nil {
		return fmt.Errorf("fail to fetch books: %w", err)
	}

	for bk := range bkChan {
		bk := bk
		s.sema.Acquire(ctx, 1)
		se.Acquire(ctx, 1)
		wg.Add(1)

		go func(bk *model.Book) {
			defer wg.Done()
			defer se.Release(1)
			defer s.sema.Release(1)

			err := s.DownloadBook(ctx, bk)
			if err != nil {
				log.Error().Err(err).Str("site", s.name).Int("bk_id", bk.ID).Str("bk_hash_code", bk.FormatHashCode()).Msg("download book failed")
			}
		}(&bk)
	}

	wg.Wait()

	return nil
}

func isEnd(bk *model.Book) bool {
	//TODO: fetch all chapter
	//hint: use book.generateEmptyChapters
	//TODO: check last n chapter to see if they contains any end keywords
	//hint: use len(chapters) and the n should come from book config
	if bk.UpdateDate < strconv.Itoa(time.Now().Year()-1) {
		return true
	}

	chapter := strings.ReplaceAll(bk.UpdateChapter, " ", "")
	for _, keyword := range repo.ChapterEndKeywords {
		if strings.Contains(chapter, keyword) {
			return true
		}
	}
	return false
}

func (s *ServiceImpl) ValidateBookEnd(ctx context.Context, bk *model.Book) error {
	isUpdated, isBookEnded := false, isEnd(bk)
	if isBookEnded && bk.Status != model.End {
		bk.IsDownloaded = false
		bk.Status = model.End
		isUpdated = true
	} else if !isBookEnded && bk.Status != model.InProgress {
		bk.Status = model.InProgress
		isUpdated = true
	}

	if isUpdated {
		err := s.rpo.UpdateBook(bk)
		if err != nil {
			return fmt.Errorf("update book in DB fail: %w", err)
		}
	}

	return nil
}

func (s *ServiceImpl) ValidateEnd(ctx context.Context) error {
	return s.rpo.UpdateBooksStatus()
}