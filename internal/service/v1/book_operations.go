package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/htchan/BookSpider/internal/model"
	serv "github.com/htchan/BookSpider/internal/service"
	vendor "github.com/htchan/BookSpider/internal/vendorservice"
	"github.com/rs/zerolog"
	"golang.org/x/sync/semaphore"
)

func isNewBook(bk *model.Book, bkInfo *vendor.BookInfo) bool {
	return bk.Status != model.StatusError && (bk.Title != bkInfo.Title || bk.Writer.Name != bkInfo.Writer || bk.Type != bkInfo.Type)
}

func isBookUpdated(bk *model.Book, bkInfo *vendor.BookInfo) bool {
	return bk.UpdateDate != bkInfo.UpdateDate || bk.UpdateChapter != bkInfo.UpdateChapter
}

func (s *ServiceImpl) UpdateBook(ctx context.Context, bk *model.Book, stats *serv.UpdateStats) error {
	if stats == nil {
		stats = new(serv.UpdateStats)
	}

	body, err := s.cli.Get(ctx, s.vendorService.BookURL(strconv.FormatInt(int64(bk.ID), 10)))
	if err != nil {
		stats.Fail.Add(1)
		return fmt.Errorf("get book page failed: %w", err)
	}

	bkInfo, err := s.vendorService.ParseBook(body)
	if err != nil {
		stats.Fail.Add(1)
		return fmt.Errorf("parse book page failed: %w", err)
	}

	logger := zerolog.Ctx(ctx).With().Str("bk_title", bkInfo.Title).Logger()

	if isNewBook(bk, bkInfo) {
		logger.Debug().
			Interface("existing_book", bk).
			Interface("new_book", bkInfo).
			Msg("Found new book")

		stats.NewEntity.Add(1)
		switch bk.Status {
		case model.StatusError:
			stats.ErrorUpdated.Add(1)
		case model.StatusInProgress:
			stats.InProgressUpdated.Add(1)
		case model.StatusEnd:
			if bk.IsDownloaded {
				stats.DownloadedUpdated.Add(1)
			} else {
				stats.EndUpdated.Add(1)
			}
		}

		bk.Title, bk.Writer.Name, bk.Type = bkInfo.Title, bkInfo.Writer, bkInfo.Type
		bk.UpdateDate, bk.UpdateChapter = bkInfo.UpdateDate, bkInfo.UpdateChapter

		bk.HashCode = model.GenerateHash()
		bk.Status = model.StatusInProgress
		bk.Error = nil

		saveWriterErr := s.rpo.SaveWriter(&bk.Writer)
		saveBkErr := s.rpo.CreateBook(bk)
		saveErrErr := s.rpo.SaveError(bk, bk.Error)
		if saveWriterErr != nil || saveBkErr != nil || saveErrErr != nil {
			return errors.Join(saveWriterErr, saveBkErr)
		}
	} else if isBookUpdated(bk, bkInfo) {
		logger.Debug().
			Str("old_updated_data", bk.UpdateDate).Str("new_updated_data", bkInfo.UpdateDate).
			Str("old_updated_chapter", bk.UpdateChapter).Str("new_updated_chapter", bkInfo.UpdateChapter).
			Msg("Found updated book")

		stats.NewChapter.Add(1)
		switch bk.Status {
		case model.StatusError:
			stats.ErrorUpdated.Add(1)
		case model.StatusInProgress:
			stats.InProgressUpdated.Add(1)
		case model.StatusEnd:
			if bk.IsDownloaded {
				stats.EndUpdated.Add(1)
			} else {
				stats.EndUpdated.Add(1)
			}
		}

		if bk.Status == model.StatusError {
			bk.Title, bk.Writer.Name, bk.Type = bkInfo.Title, bkInfo.Writer, bkInfo.Type
		}

		bk.UpdateDate, bk.UpdateChapter = bkInfo.UpdateDate, bkInfo.UpdateChapter

		bk.Status = model.StatusInProgress
		bk.Error = nil

		if bk.Status == model.StatusError {
			bk.Status = model.StatusInProgress
		}

		saveWriterErr := s.rpo.SaveWriter(&bk.Writer)
		saveBkErr := s.rpo.UpdateBook(bk)
		saveErrErr := s.rpo.SaveError(bk, bk.Error)
		if saveWriterErr != nil || saveBkErr != nil || saveErrErr != nil {
			return errors.Join(saveWriterErr, saveBkErr)
		}
	} else {
		logger.Debug().Msg("book not updated")
		stats.Unchanged.Add(1)
	}

	return nil
}

func (s *ServiceImpl) Update(ctx context.Context, stats *serv.UpdateStats) error {
	var wg sync.WaitGroup
	if stats == nil {
		stats = new(serv.UpdateStats)
	}

	bkChan, err := s.rpo.FindBooksForUpdate()
	if err != nil {
		return fmt.Errorf("fail to load books from DB: %w", err)
	}

	for bk := range bkChan {
		bk := bk
		s.sema.Acquire(ctx, 1)
		wg.Add(1)
		stats.Total.Add(1)

		go func(bk *model.Book) {
			defer wg.Done()
			defer s.sema.Release(1)

			logger := zerolog.Ctx(ctx).With().
				Int("bk_id", bk.ID).
				Str("bk_hash_code", bk.FormatHashCode()).
				Str("worker_id", uuid.New().String()).
				Logger()
			err := s.UpdateBook(logger.WithContext(ctx), bk, stats)
			if err != nil {
				logger.Error().Err(err).
					Msg("update book failed")
			}
		}(&bk)

		// give chance to others service running at the same time
		time.Sleep(time.Millisecond)
	}

	wg.Wait()

	return nil
}

func (s *ServiceImpl) ExploreBook(ctx context.Context, bk *model.Book, stats *serv.UpdateStats) error {
	if bk.Status != model.StatusError {
		return serv.ErrBookStatusNotError
	}

	//TODO: find a new method to check if we should create the book
	isNew := bk.Error == nil
	if isNew {
		s.rpo.CreateBook(bk)
	}

	err := s.UpdateBook(ctx, bk, stats)
	if err != nil {
		bk.Error = err
		saveErr := s.rpo.SaveError(bk, bk.Error)
		if saveErr != nil {
			return fmt.Errorf("explore book fail: %w; save error fail: %w", err, saveErr)
		}

		return fmt.Errorf("explore book fail: %w", err)
	}

	return err
}

func (s *ServiceImpl) Explore(ctx context.Context, stats *serv.UpdateStats) error {
	summary := s.rpo.Stats()
	var errorCount atomic.Int64

	if stats == nil {
		stats = new(serv.UpdateStats)
	}

	var wg sync.WaitGroup

	for i := summary.LatestSuccessID + 1; i <= summary.MaxBookID && int(errorCount.Load()) < s.conf.MaxExploreError; i++ {
		i := i

		s.sema.Acquire(ctx, 1)
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			defer s.sema.Release(1)

			bk, err := s.rpo.FindBookById(id)
			if err != nil {
				errorCount.Add(1)
				return
			}

			logger := zerolog.Ctx(ctx).With().
				Str("worker_id", uuid.New().String()).
				Int("bk_id", bk.ID).
				Str("bk_hash_code", bk.FormatHashCode()).
				Logger()
			err = s.ExploreBook(logger.WithContext(ctx), bk, stats)
			if err != nil {
				logger.Error().Err(err).
					Msg("explore book failed")
				errorCount.Add(1)
			} else {
				errorCount.Store(0)
			}
		}(i)

		// give chance to others service running at the same time
		time.Sleep(time.Millisecond)
	}

	wg.Wait()

	for i := summary.MaxBookID + 1; int(errorCount.Load()) < s.conf.MaxExploreError; i++ {
		i := i

		s.sema.Acquire(ctx, 1)
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			defer s.sema.Release(1)

			bk := model.NewBook(s.name, i)
			logger := zerolog.Ctx(ctx).With().
				Str("worker_id", uuid.New().String()).
				Int("bk_id", bk.ID).
				Str("bk_hash_code", bk.FormatHashCode()).
				Logger()

			err := s.ExploreBook(logger.WithContext(ctx), &bk, stats)
			if err != nil {
				logger.Error().Err(err).
					Msg("explore book failed")
				errorCount.Add(1)
			} else {
				errorCount.Store(0)
			}
		}(i)

		// give chance to others service running at the same time
		time.Sleep(time.Millisecond)
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

	chapter, err := s.vendorService.ParseChapter(body)
	if err != nil {
		ch.Error = err

		return fmt.Errorf("parse chapter page failed: %w", err)
	}

	ch.Title, ch.Content = chapter.Title, chapter.Body

	ch.OptimizeContent()

	return nil
}

func (s *ServiceImpl) DownloadBook(ctx context.Context, bk *model.Book, stats *serv.DownloadStats) error {
	if stats == nil {
		stats = new(serv.DownloadStats)
	}

	if bk.Status != model.StatusEnd {
		return serv.ErrBookStatusNotEnd
	} else if bk.IsDownloaded {
		return serv.ErrBookAlreadyDownloaded
	}

	logger := zerolog.Ctx(ctx)

	logger.Info().Msg("get chapter list")

	body, err := s.cli.Get(ctx, s.vendorService.ChapterListURL(strconv.FormatInt(int64(bk.ID), 10)))
	if err != nil {
		stats.RequestFail.Add(1)

		return fmt.Errorf("get chapter list failed: %w", err)
	}

	chapterList, err := s.vendorService.ParseChapterList(strconv.Itoa(bk.ID), body)
	if err != nil {
		if errors.Is(err, vendor.ErrChapterListEmpty) {
			stats.NoChapter.Add(1)
		} else {
			stats.RequestFail.Add(1)
		}

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

			chapterLogger := logger.With().
				Str("chapter_worker_id", uuid.New().String()).
				Str("chapter_url", ch.URL).
				Logger()
			err := s.downloadChapter(chapterLogger.WithContext(ctx), ch)
			if err != nil {
				failedChapterCount += 1
				chapterLogger.Error().Err(err).
					Str("chapter_title", ch.Title).
					Msg("download chapter failed")
			}
		}(&chapters[i])
	}

	wg.Wait()

	if failedChapterCount > 50 || failedChapterCount*10 > len(chapters) {
		stats.TooManyFailChapters.Add(1)

		return fmt.Errorf("Download chapters fail: %w (%v/%v)", serv.ErrTooManyFailedChapters, failedChapterCount, len(chapters))
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
	err = s.rpo.UpdateBook(bk)
	if err != nil {
		return fmt.Errorf("update book is_downloaded fail: %w", err)
	}

	stats.Success.Add(1)

	return nil
}

func (s *ServiceImpl) Download(ctx context.Context, stats *serv.DownloadStats) error {
	se := semaphore.NewWeighted(int64(s.conf.MaxDownloadConcurrency))
	var wg sync.WaitGroup

	if stats == nil {
		stats = new(serv.DownloadStats)
	}

	bkChan, err := s.rpo.FindBooksForDownload()
	if err != nil {
		return fmt.Errorf("fail to fetch books: %w", err)
	}

	for bk := range bkChan {
		bk := bk
		s.sema.Acquire(ctx, 1)
		se.Acquire(ctx, 1)
		wg.Add(1)

		stats.Total.Add(1)

		go func(bk *model.Book) {
			defer wg.Done()
			defer se.Release(1)
			defer s.sema.Release(1)

			logger := zerolog.Ctx(ctx).With().
				Str("worker_id", uuid.New().String()).
				Int("bk_id", bk.ID).
				Str("bk_hash_code", bk.FormatHashCode()).
				Logger()

			err := s.DownloadBook(logger.WithContext(ctx), bk, stats)
			if err != nil {
				logger.Error().Err(err).Msg("download book failed")
			}
		}(&bk)

		// give chance to others service running at the same time
		time.Sleep(time.Millisecond)
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
	for _, keyword := range model.ChapterEndKeywords {
		if strings.Contains(chapter, keyword) {
			return true
		}
	}
	return false
}

func (s *ServiceImpl) ValidateBookEnd(ctx context.Context, bk *model.Book) error {
	isUpdated, isBookEnded := false, isEnd(bk)
	if isBookEnded && bk.Status != model.StatusEnd {
		bk.IsDownloaded = false
		bk.Status = model.StatusEnd
		isUpdated = true
	} else if !isBookEnded && bk.Status != model.StatusInProgress {
		bk.Status = model.StatusInProgress
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
