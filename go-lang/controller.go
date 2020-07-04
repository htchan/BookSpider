package main

import (
	"golang.org/x/text/encoding/traditionalchinese"
	//"bufio"
	"os"
	"strings"
	"fmt"
	"./model"
)


func help() {
	fmt.Println("Command: ")
    fmt.Println("help" + strings.Repeat(" ", 14) + "show the functin list avaliable")
    fmt.Println("download" + strings.Repeat(" ", 10) + "download books")
    fmt.Println("update" + strings.Repeat(" ", 12) + "update books information")
    fmt.Println("explore" + strings.Repeat(" ", 11) + "explore new books in internet")
    fmt.Println("check" + strings.Repeat(" ", 13) + "check recorded books finished")
    fmt.Println("error" + strings.Repeat(" ", 13) + "update all website may have error")
    fmt.Println("backup" + strings.Repeat(" ", 12) + "backup the current database by the current date and time")
    fmt.Println("regular" + strings.Repeat(" ", 11) + "do the default operation (explore->update->download->check)")
    fmt.Println("\n")
}

func download(sites map[string]model.Site) {
	for name, site := range sites {
		fmt.Println(name + "\tupdate");
		site.Download();
	}
}
func update(sites map[string]model.Site) {
	for name, site := range sites {
		fmt.Println(name + "\tupdate");
		site.Update();
	}
}
func explore(sites map[string]model.Site, maxError int) {
	for name, site := range sites {
		fmt.Println(name + "\tupdate");
		site.Explore(maxError);
	}
}
func updateError(sites map[string]model.Site) {
	for name, site := range sites {
		fmt.Println(name + "\tupdate");
		site.UpdateError();
	}
}
func info(sites map[string]model.Site) {
	for name, site := range sites {
		fmt.Println(name + "\tupdate");
		fmt.Println(strings.Repeat("- ", 20));
		site.Info();
		fmt.Println(strings.Repeat("- ", 20));
	}
}
func check(sites map[string]model.Site) {
	for name, site := range sites {
		fmt.Println(name + "\tupdate");
		site.Check();
	}
}
func backup(sites map[string]model.Site) {
	for name, site := range sites {
		fmt.Println(name + "\tupdate");
		site.Backup();
	}
}

func main() {
	fmt.Println("test (v0.0.0) - - - - - - - - - -");
	if (len(os.Args) < 2) {
		help();
		fmt.Println("No arguements");
		return;
	}

	big5Decoder := traditionalchinese.Big5.NewDecoder()

	sites := make(map[string]model.Site);
	sites["ck101"] = model.NewSite("ck101", big5Decoder, "./book-config/ck101-desktop.json", "./database/ck101.db", "./");
	
	switch operation := strings.ToUpper(os.Args[1]); operation {
	case "UPDATE":
		update(sites);
	case "EXPLORE":
		explore(sites, 1000);
	case "DOWNLOAD":
		download(sites);
	case "ERROR":
		updateError(sites);
	case "INFO":
		info(sites);
	case "CHECK":
		check(sites);
	case "BACKUP":
		backup(sites);
	default:
		help();
	}
	/*
	book := model.NewBook("ck101", 1, big5Decoder,
		"https://www.ck101.org/book/1.html", "", "",
		"<h1>.*?<a.*?>(.*?)</a>.*?</h1>", "", "", "", "", "", "");
	fmt.Println(book.Site)
	book.Update();
	fmt.Println(book.Title);
	*/
	/*
	site := model.NewSite("ck101", big5Decoder, "./book-config/ck101-desktop.json", "./database/ck101.db", "./");
	fmt.Println(site.MetaBaseUrl);
	*/
	
	//book := site.Book(-1);
	//book.Update();
	/*
	fmt.Println(book.Title);
	fmt.Println(book.Writer)
	fmt.Println(book.Type)
	fmt.Println(book.LastUpdate)
	fmt.Println(book.LastChapter)
	*/
	//site.Update()
	/*
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	*/
}