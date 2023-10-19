package service

import (
	"context"
	"fmt"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/rs/zerolog"
)

func (s *ServiceImpl) ProcessBook(ctx context.Context, bk *model.Book) error {
	ctx = zerolog.Ctx(ctx).With().Str("site", s.name).Logger().WithContext(ctx)

	updateErr := s.UpdateBook(ctx, bk)
	if updateErr != nil {
		return fmt.Errorf("Update book error: %w", updateErr)
	}

	validateErr := s.ValidateBookEnd(ctx, bk)
	if validateErr != nil {
		return fmt.Errorf("validate book error: %w", validateErr)
	}

	bookStorageUpdated := s.checkBookStorage(bk)
	if bookStorageUpdated {
		s.rpo.UpdateBook(bk)
	}

	downloadErr := s.DownloadBook(ctx, bk)
	if downloadErr != nil {
		return fmt.Errorf("Process book error: %w", downloadErr)
	}

	bookStorageUpdated = s.checkBookStorage(bk)
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
	updateErr := s.Update(updateCtx)
	zerolog.Ctx(updateCtx).Trace().Msg("complete")
	if updateErr != nil {
		return fmt.Errorf("Update fail: %w", updateErr)
	}
	exploreCtx := zerolog.Ctx(ctx).With().Str("operation", "explore").Logger().WithContext(ctx)
	zerolog.Ctx(exploreCtx).Trace().Msg("start")
	exploreErr := s.Explore(exploreCtx)
	zerolog.Ctx(exploreCtx).Trace().Msg("complete")
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
	downloadErr := s.Download(downloadCtx)
	zerolog.Ctx(downloadCtx).Trace().Msg("complete")
	if downloadErr != nil {
		return fmt.Errorf("Download fail: %w", downloadErr)
	}

	patchStatusCtx := zerolog.Ctx(ctx).With().Str("operation", "patch-status").Logger().WithContext(ctx)
	zerolog.Ctx(patchStatusCtx).Trace().Msg("start")
	patchDownloadStatusErr := s.PatchDownloadStatus(patchStatusCtx)
	zerolog.Ctx(patchStatusCtx).Trace().Msg("complete")
	if patchDownloadStatusErr != nil {
		return fmt.Errorf("patch status fail: %w", patchDownloadStatusErr)
	}

	patchMissingCtx := zerolog.Ctx(ctx).With().Str("operation", "patch-missing-records").Logger().WithContext(ctx)
	zerolog.Ctx(patchMissingCtx).Trace().Msg("start")
	patchMissingRecordsErr := s.PatchMissingRecords(patchMissingCtx)
	zerolog.Ctx(patchMissingCtx).Trace().Msg("complete")
	if patchMissingRecordsErr != nil {
		return fmt.Errorf("patch status fail: %w", patchMissingRecordsErr)
	}

	return nil
}
