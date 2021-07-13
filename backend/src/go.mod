module github.com/htchan/BookSpider

go 1.13

replace github.com/htchan/BookSpider/helper v0.0.0 => ./helper

replace github.com/htchan/BookSpider/model v0.0.0 => ./model

require (
	github.com/htchan/BookSpider/helper v0.0.0 // indirect
	github.com/htchan/BookSpider/model v0.0.0 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/text v0.3.6 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
