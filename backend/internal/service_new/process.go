package service

import (
	"fmt"
	"log"

	"github.com/htchan/BookSpider/internal/model"
)

func (serv *ServiceImp) ProcessBook(bk *model.Book) error {
	updateErr := serv.UpdateBook(bk)
	if updateErr != nil {
		return fmt.Errorf("Update book error: %w", updateErr)
	}

	validateErr := serv.ValidateBookEnd(bk)
	if validateErr != nil {
		return fmt.Errorf("Validate book error: %w", validateErr)
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
	log.Printf("[operation.%v.backup] start", serv.name)
	backupErr := serv.Backup()
	log.Printf("[operation.%v.backup] complete", serv.name)
	if backupErr != nil {
		return fmt.Errorf("Backup fail: %w", backupErr)
	}

	log.Printf("[operation.%v.check-availability] start", serv.name)
	checkAvailabilityErr := serv.CheckAvailability()
	log.Printf("[operation.%v.check-availability] complete", serv.name)
	if checkAvailabilityErr != nil {
		return fmt.Errorf("Check Availability fail: %w", checkAvailabilityErr)
	}

	log.Printf("[operation.%v.update] start", serv.name)
	updateErr := serv.Update()
	log.Printf("[operation.%v.update] complete", serv.name)
	if updateErr != nil {
		return fmt.Errorf("Update fail: %w", updateErr)
	}

	log.Printf("[operation.%v.explore] start", serv.name)
	exploreErr := serv.Explore()
	log.Printf("[operation.%v.explore] complete", serv.name)
	if exploreErr != nil {
		return fmt.Errorf("Explore fail: %w", exploreErr)
	}

	log.Printf("[operation.%v.update-status] start", serv.name)
	checkErr := serv.ValidateEnd()
	log.Printf("[operation.%v.update-status] complete", serv.name)
	if checkErr != nil {
		return fmt.Errorf("Update Status fail: %w", checkErr)
	}

	log.Printf("[operation.%v.download] start", serv.name)
	downloadErr := serv.Download()
	log.Printf("[operation.%v.download] complete", serv.name)
	if downloadErr != nil {
		return fmt.Errorf("Download fail: %w", downloadErr)
	}

	log.Printf("[operation.%v.patch-status] start", serv.name)
	patchDownloadStatusErr := serv.PatchDownloadStatus()
	log.Printf("[operation.%v.patch-status] complete", serv.name)
	if patchDownloadStatusErr != nil {
		return fmt.Errorf("Patch Status fail: %w", patchDownloadStatusErr)
	}

	log.Printf("[operation.%v.patch-missing-records] start", serv.name)
	patchMissingRecordsErr := serv.PatchMissingRecords()
	log.Printf("[operation.%v.patch-missing-records] complete", serv.name)
	if patchMissingRecordsErr != nil {
		return fmt.Errorf("Patch Status fail: %w", patchMissingRecordsErr)
	}

	return nil
}
