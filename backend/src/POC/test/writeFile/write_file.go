package main

import (
	"path/filepath"
	"fmt"
	"os"
	"log"
)

func main() () {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fmt.Println(dir + "/test.txt")
	/*
    err := ioutil.WriteFile(dir + "/test.txt", []byte("Hi\n"), 0644)
    if err != nil {
        log.Fatal(err)
	}
	*/
    file, err := os.OpenFile(dir + "/test.txt", os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0664)
    if err != nil {
        log.Fatal(err)
	}
	_, err = file.WriteString("hii")
    if err != nil {
        log.Fatal(err)
	}

    file.Close()

}