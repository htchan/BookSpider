package service

import (
	"context"
	"fmt"

	"github.com/htchan/BookSpider/internal/model"
	serv "github.com/htchan/BookSpider/internal/service"
	"github.com/rs/zerolog"
)

func (s *ServiceImpl) ProcessBook(ctx context.Context, bk *model.Book) error {
	ctx = zerolog.Ctx(ctx).With().Str("site", s.name).Logger().WithContext(ctx)

	updateErr := s.UpdateBook(ctx, bk, nil)
	if updateErr != nil {
		return fmt.Errorf("Update book error: %w", updateErr)
	}

	validateErr := s.ValidateBookEnd(ctx, bk)
	if validateErr != nil {
		return fmt.Errorf("validate book error: %w", validateErr)
	}

	bookStorageUpdated := s.checkBookStorage(bk, nil)
	if bookStorageUpdated {
		s.rpo.UpdateBook(bk)
	}

	downloadErr := s.DownloadBook(ctx, bk, nil)
	if downloadErr != nil {
		return fmt.Errorf("Process book error: %w", downloadErr)
	}

	bookStorageUpdated = s.checkBookStorage(bk, nil)
	if bookStorageUpdated {
		s.rpo.UpdateBook(bk)
	}

	return nil
}

func (s *ServiceImpl) Process(ctx context.Context) error {
	ctx = zerolog.Ctx(ctx).With().Str("site", s.name).Logger().WithContext(ctx)

	// zerolog.Ctx(ctx).Trace().Str("operation", "backup").Msg("start")
	// backupErr := s.Backup()
	// zerolog.Ctx(ctx).Trace().Str("operation", "backup").Msg("complete")
	// if backupErr != nil {
	// 	return fmt.Errorf("Backup fail: %w", backupErr)
	// }

	checkAvailabilityCtx := zerolog.Ctx(ctx).With().Str("operation", "check-availability").Logger().WithContext(ctx)
	zerolog.Ctx(checkAvailabilityCtx).Trace().Msg("start")
	checkAvailabilityErr := s.CheckAvailability(checkAvailabilityCtx)
	zerolog.Ctx(checkAvailabilityCtx).Trace().Msg("complete")
	if checkAvailabilityErr != nil {
		return fmt.Errorf("check availability fail: %w", checkAvailabilityErr)
	}

	updateCtx := zerolog.Ctx(ctx).With().Str("operation", "update").Logger().WithContext(ctx)
	zerolog.Ctx(updateCtx).Trace().Msg("start")
	updateStats := new(serv.UpdateStats)
	updateErr := s.Update(updateCtx, updateStats)
	zerolog.Ctx(updateCtx).Trace().
		Int64("total", updateStats.Total.Load()).
		Int64("fail", updateStats.Fail.Load()).
		Int64("unchanged", updateStats.Unchanged.Load()).
		Int64("new_chapter", updateStats.NewChapter.Load()).
		Int64("new_entity", updateStats.NewEntity.Load()).
		Int64("error_updated", updateStats.ErrorUpdated.Load()).
		Int64("in_progress_updated", updateStats.InProgressUpdated.Load()).
		Int64("end_updated", updateStats.EndUpdated.Load()).
		Int64("downloaded_updated", updateStats.DownloadedUpdated.Load()).
		Msg("complete")
	if updateErr != nil {
		return fmt.Errorf("Update fail: %w", updateErr)
	}
	exploreCtx := zerolog.Ctx(ctx).With().Str("operation", "explore").Logger().WithContext(ctx)
	zerolog.Ctx(exploreCtx).Trace().Msg("start")
	exploreStats := new(serv.UpdateStats)
	exploreErr := s.Explore(exploreCtx, exploreStats)
	zerolog.Ctx(exploreCtx).Trace().
		Int64("total", exploreStats.Total.Load()).
		Int64("fail", exploreStats.Fail.Load()).
		Int64("unchanged", exploreStats.Unchanged.Load()).
		Int64("new_chapter", exploreStats.NewChapter.Load()).
		Int64("new_entity", exploreStats.NewEntity.Load()).
		Int64("error_updated", exploreStats.ErrorUpdated.Load()).
		Int64("in_progress_updated", exploreStats.InProgressUpdated.Load()).
		Int64("end_updated", exploreStats.EndUpdated.Load()).
		Int64("downloaded_updated", exploreStats.DownloadedUpdated.Load()).
		Msg("complete")
	if exploreErr != nil {
		return fmt.Errorf("Explore fail: %w", exploreErr)
	}

	checkCtx := zerolog.Ctx(ctx).With().Str("operation", "check status").Logger().WithContext(ctx)
	zerolog.Ctx(checkCtx).Trace().Msg("start")
	checkErr := s.ValidateEnd(checkCtx)
	zerolog.Ctx(checkCtx).Trace().Msg("complete")
	if checkErr != nil {
		return fmt.Errorf("Update Status fail: %w", checkErr)
	}

	downloadCtx := zerolog.Ctx(ctx).With().Str("operation", "download").Logger().WithContext(ctx)
	zerolog.Ctx(downloadCtx).Trace().Msg("start")
	downloadStats := new(serv.DownloadStats)
	downloadErr := s.Download(downloadCtx, downloadStats)
	zerolog.Ctx(downloadCtx).Trace().
		Int64("total", downloadStats.Total.Load()).
		Int64("success", downloadStats.Success.Load()).
		Int64("no_chapter_error", downloadStats.NoChapter.Load()).
		Int64("too_many_failed_chapters", downloadStats.TooManyFailChapters.Load()).
		Int64("request_fail", downloadStats.RequestFail.Load()).
		Msg("complete")
	if downloadErr != nil {
		return fmt.Errorf("Download fail: %w", downloadErr)
	}

	patchStatusCtx := zerolog.Ctx(ctx).With().Str("operation", "patch-status").Logger().WithContext(ctx)
	zerolog.Ctx(patchStatusCtx).Trace().Msg("start")
	patchStorageStats := new(serv.PatchStorageStats)
	patchDownloadStatusErr := s.PatchDownloadStatus(patchStatusCtx, patchStorageStats)
	zerolog.Ctx(patchStatusCtx).Trace().
		Int64("file_exist", patchStorageStats.FileExist.Load()).
		Int64("file_missing", patchStorageStats.FileMissing.Load()).
		Msg("complete")
	if patchDownloadStatusErr != nil {
		return fmt.Errorf("patch status fail: %w", patchDownloadStatusErr)
	}

	patchMissingCtx := zerolog.Ctx(ctx).With().Str("operation", "patch-missing-records").Logger().WithContext(ctx)
	zerolog.Ctx(patchMissingCtx).Trace().Msg("start")
	patchMissingStats := new(serv.UpdateStats)
	patchMissingRecordsErr := s.PatchMissingRecords(patchMissingCtx, patchMissingStats)
	zerolog.Ctx(patchMissingCtx).Trace().
		Int64("total", patchMissingStats.Total.Load()).
		Int64("fail", patchMissingStats.Fail.Load()).
		Int64("unchanged", patchMissingStats.Unchanged.Load()).
		Int64("new_chapter", patchMissingStats.NewChapter.Load()).
		Int64("new_entity", patchMissingStats.NewEntity.Load()).
		Int64("error_updated", patchMissingStats.ErrorUpdated.Load()).
		Int64("in_progress_updated", patchMissingStats.InProgressUpdated.Load()).
		Int64("end_updated", patchMissingStats.EndUpdated.Load()).
		Int64("downloaded_updated", patchMissingStats.DownloadedUpdated.Load()).Msg("complete")
	if patchMissingRecordsErr != nil {
		return fmt.Errorf("patch status fail: %w", patchMissingRecordsErr)
	}

	return nil
}
