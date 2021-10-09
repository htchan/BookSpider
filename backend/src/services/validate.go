package services

import (
	"log"
	
	"github.com/htchan/BookSpider/models"
	"github.com/htchan/BookSpider/helper"

	"encoding/json"
	"io/ioutil"
)

func Validate(sites map[string]models.Site, flags models.Flags) {
	exploreResult := make(map[string]float64)
	downloadResult := make(map[string]float64)
	for name, site := range sites {
		if *flags.Site != "" && name != *flags.Site { continue }
		log.Println(name + "\tvalidate explore")
		exploreResult[name] = site.Validate()
		log.Println(name + "\tvalidate download")
		downloadResult[name] = site.ValidateDownload()
	}
	b, err := json.Marshal(exploreResult)
	helper.CheckError(err)
	err = ioutil.WriteFile("./validate.json", b, 0644)
	helper.CheckError(err)
	b, err = json.Marshal(downloadResult)
	helper.CheckError(err)
	err = ioutil.WriteFile("./validate-download.json", b, 0644)
	helper.CheckError(err)

}