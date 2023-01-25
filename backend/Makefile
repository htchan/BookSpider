.PHONY: console

table = test
target = ./...

## help: show available command and description
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed  -e 's/^/ /'

## setup: run go mod tidy
setup:
	go mod tidy

## controller: build executable of controller, backend, batch
build:
	go build -o ../../bin github.com/htchan/BookSpider/cmd/controller
	go build -o ../../bin github.com/htchan/BookSpider/cmd/backend
	go build -o ../../bin github.com/htchan/BookSpider/cmd/batch

## test: test packages and show coverage
test:
	# go clean --testcache
	go test ${target}

## benchmark: benchmark packages and show coverage
benchmark:
	# go clean --testcache
	go test -bench="Book*"
	# ../internal/client/... \
	# ../internal/decoder/... \
	# ../pkg/config/...
	# -coverprofile ./profile.out
	# go test ../internal/client/... -coverprofile ./profile.out
	go tool cover -html=profile.out -o coverage.html

## coverage: check coverage of backend
coverage: 
	# go clean --testcache
	go test $(target) -coverprofile profile.txt
	go tool cover -html=profile.txt -o coverage.html
	rm profile.txt
	# google-chrome ./coverage.html &

## open-coverage: open coverage file in chrome is it exist
open-coverage: ./coverage.html
	google-chrome ./coverage.html &

create_database:
	PGPASSWORD=books psql -h localhost -p 5432 -U books -c "create database ${table}"

## mockgen: generate mock code to internal/mock package
generate:
	go generate ./...

## clean: clean
clean:
	rm ./build/ -r