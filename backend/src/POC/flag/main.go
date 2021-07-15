package main;

import (
	"fmt"
	"flag"
)

func testA(i int) {
	flag.Parse()
	log.Println(i)
}

func testB(s string) {
	flag.Parse()
	log.Println(s)
}

func main() {
	id := flag.Int("id", -1, "the target book id")
	site := flag.String("site", "empty", "the target book id")
	flag.Parse()
	flag.Visit(func(f *flag.Flag) {
		log.Println(f.Name)
		if f.Name == "id" {
			testA(*id)
		} else if f.Name == "site" {
			testB(*site)
		}
	})
}