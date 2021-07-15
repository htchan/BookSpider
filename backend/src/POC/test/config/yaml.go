package main;

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"io/ioutil"
)

type SiteConfig struct {
	Name, Decode, configLocation, databaseLocation, downloadLocation string
}

type Config struct {
	Sites map[string]map[string]string
	Api []string
}

func main() () {
	s, err := ioutil.ReadFile("src/public/config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	//var mapConfig map[string]interface{}
	var mapConfig Config
	err = yaml.Unmarshal(s, &mapConfig)
	//log.Println(mapConfig)
	log.Println(mapConfig.Sites)
	log.Println(mapConfig.Sites["ck101"])
	log.Println(mapConfig.Api)
	/*
	for i, site := range map[string]interface{}(mapConfig["sites"]) {
		log.Print(i)
		log.Println(site)
	}
	*/
}
