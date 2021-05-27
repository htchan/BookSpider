package main

import (
	"fmt"
	"encoding/json"
	"log"
	"io/ioutil"
)

type Config struct {
	Sites map[string]map[string]string
	Api []string
}

func main() () {
	s, err := ioutil.ReadFile("src/public/config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	//var mapConfig map[string]interface{}
	var mapConfig Config
	err = json.Unmarshal(s, &mapConfig)
	//fmt.Println(mapConfig)
	fmt.Println(mapConfig.Sites)
	fmt.Println(mapConfig.Sites["ck101"])
	fmt.Println(mapConfig.Api)
}