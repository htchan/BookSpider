package main;

import (
	"fmt"
	"flag"
)

func testA(i int) {
	flag.Parse()
	fmt.Println(i)
}

func testB(s string) {
	flag.Parse()
	fmt.Println(s)
}

func main() {
	id := flag.Int("id", -1, "the target book id")
	site := flag.String("site", "empty", "the target book id")
	flag.Parse()
	flag.Visit(func(f *flag.Flag) {
		fmt.Println(f.Name)
		if f.Name == "id" {
			testA(*id)
		} else if f.Name == "site" {
			testB(*site)
		}
	})
}