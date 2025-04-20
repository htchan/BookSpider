package xbiquge

import (
	"flag"
	"log"
	"os"
	"testing"

	"go.uber.org/goleak"
)

var (
	testBookBytes        []byte
	testChapterBytes     []byte
	testChapterListBytes []byte
)

func TestMain(m *testing.M) {
	var err error
	testBookBytes, err = os.ReadFile("../test_resources/xbiquge_book.html")
	if err != nil {
		log.Fatalf("could not read book string")
	}

	testChapterBytes, err = os.ReadFile("../test_resources/xbiquge_chapter.html")
	if err != nil {
		log.Fatalf("could not read book string")
	}

	testChapterListBytes, err = os.ReadFile("../test_resources/xbiquge_chapter_list.html")
	if err != nil {
		log.Fatalf("could not read book string")
	}

	leak := flag.Bool("leak", false, "check for memory leaks")
	flag.Parse()

	if *leak {
		goleak.VerifyTestMain(m)
	} else {
		os.Exit(m.Run())
	}
}
