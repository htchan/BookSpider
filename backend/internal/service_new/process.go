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
	log.Printf("[operation.%v.backup] start", serv.Name)
	backupErr := serv.Backup()
	log.Printf("[operation.%v.backup] complete", serv.Name)
	if backupErr != nil {
		return fmt.Errorf("Backup fail: %w", backupErr)
	}

	log.Printf("[operation.%v.check-availability] start", serv.Name)
	checkAvailabilityErr := serv.CheckAvailability()
	log.Printf("[operation.%v.check-availability] complete", serv.Name)
	if checkAvailabilityErr != nil {
		return fmt.Errorf("Check Availability fail: %w", checkAvailabilityErr)
	}

	log.Printf("[operation.%v.update] start", serv.Name)
	updateErr := serv.Update()
	log.Printf("[operation.%v.update] complete", serv.Name)
	if updateErr != nil {
		return fmt.Errorf("Update fail: %w", updateErr)
	}

	log.Printf("[operation.%v.explore] start", serv.Name)
	exploreErr := serv.Explore()
	log.Printf("[operation.%v.explore] complete", serv.Name)
	if exploreErr != nil {
		return fmt.Errorf("Explore fail: %w", exploreErr)
	}

	log.Printf("[operation.%v.update-status] start", serv.Name)
	checkErr := serv.ValidateEnd()
	log.Printf("[operation.%v.update-status] complete", serv.Name)
	if checkErr != nil {
		return fmt.Errorf("Update Status fail: %w", checkErr)
	}

	log.Printf("[operation.%v.download] start", serv.Name)
	downloadErr := serv.Download()
	log.Printf("[operation.%v.download] complete", serv.Name)
	if downloadErr != nil {
		return fmt.Errorf("Download fail: %w", downloadErr)
	}

	log.Printf("[operation.%v.patch-status] start", serv.Name)
	patchDownloadStatusErr := serv.PatchDownloadStatus()
	log.Printf("[operation.%v.patch-status] complete", serv.Name)
	if patchDownloadStatusErr != nil {
		return fmt.Errorf("Patch Status fail: %w", patchDownloadStatusErr)
	}

	log.Printf("[operation.%v.patch-missing-records] start", serv.Name)
	patchMissingRecordsErr := serv.PatchMissingRecords()
	log.Printf("[operation.%v.patch-missing-records] complete", serv.Name)
	if patchMissingRecordsErr != nil {
		return fmt.Errorf("Patch Status fail: %w", patchMissingRecordsErr)
	}

	return nil
}
