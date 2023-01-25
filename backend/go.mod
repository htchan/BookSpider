module github.com/htchan/BookSpider

go 1.17

require (
	github.com/PuerkitoBio/goquery v1.8.0
	github.com/caarlos0/env/v6 v6.10.1
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/cors v1.2.1
	github.com/golang-migrate/migrate/v4 v4.15.2
	github.com/golang/mock v1.6.0
	github.com/google/go-cmp v0.5.8
	github.com/htchan/ApiParser v0.0.4
	github.com/lib/pq v1.10.6
	github.com/stretchr/testify v1.7.0
	golang.org/x/sync v0.1.0
	golang.org/x/text v0.6.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/net v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/htchan/ApiParser => /home/htchan/Project/ApiParser
