package service

import (
	"fmt"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/rs/zerolog/log"
)

func (serv *ServiceImp) ProcessBook(bk *model.Book) error {
	updateErr := serv.UpdateBook(bk)
	if updateErr != nil {
		return fmt.Errorf("Update book error: %w", updateErr)
	}

	validateErr := serv.ValidateBookEnd(bk)
	if validateErr != nil {
		return fmt.Errorf("validate book error: %w", validateErr)
	}

	bookStorageUpdated := serv.checkBookStorage(bk)
	if bookStorageUpdated {
		serv.rpo.UpdateBook(bk)
	}

	downloadErr := serv.DownloadBook(bk)
	if downloadErr != nil {
		return fmt.Errorf("Process book error: %w", downloadErr)
	}

	bookStorageUpdated = serv.checkBookStorage(bk)
	if bookStorageUpdated {
		serv.rpo.UpdateBook(bk)
	}

	return nil
}

func (serv *ServiceImp) Process() error {
	log.Trace().Str("operation", "backup").Str("site", serv.name).Msg("start")
	backupErr := serv.Backup()
	log.Trace().Str("operation", "backup").Str("site", serv.name).Msg("complete")
	if backupErr != nil {
		return fmt.Errorf("Backup fail: %w", backupErr)
	}

	log.Trace().Str("operation", "check-availability").Str("site", serv.name).Msg("start")
	checkAvailabilityErr := serv.CheckAvailability()
	log.Trace().Str("operation", "check-availability").Str("site", serv.name).Msg("complete")
	if checkAvailabilityErr != nil {
		return fmt.Errorf("check availability fail: %w", checkAvailabilityErr)
	}

	log.Trace().Str("operation", "update").Str("site", serv.name).Msg("start")
	updateErr := serv.Update()
	log.Trace().Str("operation", "update").Str("site", serv.name).Msg("complete")
	if updateErr != nil {
		return fmt.Errorf("Update fail: %w", updateErr)
	}

	log.Trace().Str("operation", "explore").Str("site", serv.name).Msg("start")
	exploreErr := serv.Explore()
	log.Trace().Str("operation", "explore").Str("site", serv.name).Msg("complete")
	if exploreErr != nil {
		return fmt.Errorf("Explore fail: %w", exploreErr)
	}

	log.Trace().Str("operation", "update-status").Str("site", serv.name).Msg("start")
	checkErr := serv.ValidateEnd()
	log.Trace().Str("operation", "update-status").Str("site", serv.name).Msg("complete")
	if checkErr != nil {
		return fmt.Errorf("Update Status fail: %w", checkErr)
	}

	log.Trace().Str("operation", "download").Str("site", serv.name).Msg("start")
	downloadErr := serv.Download()
	log.Trace().Str("operation", "download").Str("site", serv.name).Msg("complete")
	if downloadErr != nil {
		return fmt.Errorf("Download fail: %w", downloadErr)
	}

	log.Trace().Str("operation", "patch-status").Str("site", serv.name).Msg("start")
	patchDownloadStatusErr := serv.PatchDownloadStatus()
	log.Trace().Str("operation", "patch-status").Str("site", serv.name).Msg("complete")
	if patchDownloadStatusErr != nil {
		return fmt.Errorf("patch status fail: %w", patchDownloadStatusErr)
	}

	log.Trace().Str("operation", "patch-missing-records").Str("site", serv.name).Msg("start")
	patchMissingRecordsErr := serv.PatchMissingRecords()
	log.Trace().Str("operation", "patch-missing-records").Str("site", serv.name).Msg("complete")
	if patchMissingRecordsErr != nil {
		return fmt.Errorf("patch status fail: %w", patchMissingRecordsErr)
	}

	return nil
}
