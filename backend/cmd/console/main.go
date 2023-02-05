package main

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/htchan/BookSpider/internal/arguement"
	config "github.com/htchan/BookSpider/internal/config_new"
	"github.com/htchan/BookSpider/internal/model"
	repo "github.com/htchan/BookSpider/internal/repo/psql"
	service_new "github.com/htchan/BookSpider/internal/service_new"
	"golang.org/x/sync/semaphore"
)

func OperateAllSites(services map[string]service_new.Service, operation string) error {
	// loop all sites by calling process
	var wg sync.WaitGroup
	for _, serv := range services {
		serv := serv
		wg.Add(1)
		go func(serv service_new.Service) {
			defer wg.Done()
			err := OperateSite(serv, operation)

			if err != nil {
				log.Printf("[%v] Operate fail: %v\n", serv.Name(), err)
			}
		}(serv)
	}
	wg.Wait()
	return nil
}

func OperateSite(serv service_new.Service, operation string) error {
	log.Println(operation)
	if serv == nil {
		return errors.New("Site not found")
	}
	var err error
	switch operation {
	case "backup":
		err = serv.Backup()
	case "update-status":
		err = serv.ValidateEnd()
	case "download":
		err = serv.Download()
	case "explore":
		err = serv.Explore()
	case "patch-missing-records":
		err = serv.PatchMissingRecords()
	case "patch-download-status":
		err = serv.PatchDownloadStatus()

	case "update":
		err = serv.Update()
	default:
		err = errors.New("operation not found")
	}

	return err
}

func OperateBook(serv service_new.Service, bk *model.Book, operation string) error {
	if serv == nil {
		return errors.New("Site not found")
	}
	var err error

	switch operation {
	case "download":
		err = serv.DownloadBook(bk)
	case "update":
		err = serv.UpdateBook(bk)
	case "explore":
		err = serv.ExploreBook(bk)
	case "validate":
		err = serv.ValidateBookEnd(bk)
	default:
		err = errors.New("operation not found")
	}

	if err != nil {
		return err
	}

	return nil
}

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("load backend config: %v", err)
		return
	}

	validErr := conf.Validate()
	if validErr != nil {
		log.Fatalf("validate config fail: %v", validErr)
	}

	ctx := context.Background()
	services := make(map[string]service_new.Service)
	for _, siteName := range conf.BatchConfig.AvailableSiteNames {
		db, err := repo.OpenDatabase(siteName)
		if err != nil {
			log.Fatalf("load db Fail. site: %v; err: %v", siteName, err)
		}

		sema := semaphore.NewWeighted(int64(conf.SiteConfigs[siteName].MaxThreads))

		serv, err := service_new.LoadService(
			siteName, conf.SiteConfigs[siteName], db, sema, &ctx,
		)
		if err != nil {
			log.Fatalf("load service fail. site: %v, err: %v", siteName, err)
		}

		services[siteName] = serv
	}

	// load arguements
	args := arguement.LoadArgs()
	if args.IsAllSite() {
		err = OperateAllSites(services, *args.Operation)
	} else if args.IsSite() {
		err = OperateSite(args.GetSite(services), *args.Operation)
	} else if args.IsBook() {
		err = OperateBook(args.GetSite(services), args.GetBook(services), *args.Operation)
	} else {
		err = errors.New("invalid arguements")
	}

	if err != nil {
		operation := "process"
		if *args.Operation != "" {
			operation = *args.Operation
		}
		log.Fatalf("%v failed: %v", operation, err)
	}
}
