package services

import (
	"log"
	
	"github.com/htchan/BookSpider/pkg/flags"
	"github.com/htchan/BookSpider/pkg/sites"
	"github.com/htchan/BookSpider/internal/utils"

	"encoding/json"
	"io/ioutil"
)

func Validate(siteMap map[string]sites.Site, flags flags.Flags) {
	exploreResult := make(map[string]float64)
	downloadResult := make(map[string]float64)
	for name, site := range siteMap {
		if *flags.Site != "" && name != *flags.Site { continue }
		log.Println(name + "\tvalidate explore")
		exploreResult[name] = site.Validate()
		log.Println(name + "\tvalidate download")
		downloadResult[name] = site.ValidateDownload()
	}
	b, err := json.Marshal(exploreResult)
	utils.CheckError(err)
	err = ioutil.WriteFile("./validate.json", b, 0644)
	utils.CheckError(err)
	b, err = json.Marshal(downloadResult)
	utils.CheckError(err)
	err = ioutil.WriteFile("./validate-download.json", b, 0644)
	utils.CheckError(err)

}