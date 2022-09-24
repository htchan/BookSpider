package main

import (
	"errors"
	"log"
	"os"
	"sync"

	"github.com/htchan/BookSpider/internal/arguement"
	"github.com/htchan/BookSpider/internal/config"
	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/service/book"
	"github.com/htchan/BookSpider/internal/service/site"
)

func OperateAllSites(sites map[string]*site.Site, operation string) error {
	// loop all sites by calling process
	var wg sync.WaitGroup
	for _, st := range sites {
		st := st
		go func(st *site.Site) {
			defer wg.Done()
			err := OperateSite(st, operation)

			if err != nil {
				log.Printf("Process fail: %v\n", err)
			}
		}(st)
	}
	wg.Wait()
	return nil
}

func OperateSite(st *site.Site, operation string) error {
	var err error
	switch operation {
	case "backup":
		err = site.Backup(st)
	case "check":
		err = site.Check(st)
	case "download":
		err = site.Download(st)
	case "exolore":
		err = site.Explore(st)
	case "fix":
		err = site.Fix(st)
	case "update":
		err = site.Update(st)
	case "validate":
		err = site.Validate(st)
	default:
		err = errors.New("operation not found")
	}

	return err
}

func OperateBook(st *site.Site, bk *model.Book, operation string) error {
	var (
		isUpdated bool
		err       error
	)

	switch operation {
	case "download":
		isUpdated, err = book.Download(bk, st.BkConf, st.StConf, st.Client)
	case "fix":
		isUpdated, err = book.Fix(bk, st.StConf)
	case "update":
		isUpdated, err = book.Update(bk, st.BkConf, st.StConf, st.Client)
	case "validate":
		isUpdated, err = book.Validate(bk)
	default:
		err = errors.New("operation not found")
	}

	if err != nil {
		return err
	}

	if isUpdated {
		st.SaveWriter(bk)
		st.SaveError(bk)
		st.UpdateBook(bk)
	}
	return nil
}

func main() {
	configLocation := os.Getenv("ASSETS_LOCATION") + "/config"

	// load backend config
	batchConfig, err := config.LoadBatchConfig(configLocation)
	if err != nil {
		log.Fatalf("load backend config: %v", err)
	}

	// load sites from config
	sites, err := site.LoadSitesFromConfigDirectory(configLocation, batchConfig.EnabledSites)
	if err != nil {
		log.Fatal(err)
	}

	// load arguements
	args := arguement.LoadArgs()
	if args.IsAllSite() {
		err = OperateAllSites(sites, *args.Operation)
	} else if args.IsSite() {
		err = OperateSite(args.GetSite(sites), *args.Operation)
	} else if args.IsBook() {
		err = OperateBook(args.GetSite(sites), args.GetBook(sites), *args.Operation)
	} else {
		err = errors.New("invalid arguements")
	}

	if err != nil {
		log.Fatal(err)
	}
}
