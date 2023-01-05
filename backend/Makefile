.PHONY: console

table = test

## help: show available command and description
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed  -e 's/^/ /'

## all: build all executable and put them to bin
all: controller backend batch

## controller: build executable of controller, backend, batch
build:
	go build -o ../../bin github.com/htchan/BookSpider/cmd/controller
	go build -o ../../bin github.com/htchan/BookSpider/cmd/backend
	go build -o ../../bin github.com/htchan/BookSpider/cmd/batch

## test: test packages and show coverage
test:
	# go clean --testcache
	go test \
	../internal/client/... \
	../internal/decoder/... \
	../pkg/book/... \
	../pkg/config/... \
	-coverprofile ./profile.out
	# go test ../internal/client/... -coverprofile ./profile.out
	go tool cover -html=profile.out -o coverage.html

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
	go test $(PKG) -coverprofile profile.txt ; go tool cover -html=profile.txt -o coverage.html
	rm profile.txt
	google-chrome ./coverage.html &

create_database:
	PGPASSWORD=books psql -h localhost -p 5432 -U books -c "create database ${table}"

## mockgen: generate mock code to internal/mock package
mockgen:
	mockgen -source=./internal/client/client.go -destination=./internal/mock/client.go -package=mock
	mockgen -source=./internal/repo/repository.go -destination=./internal/mock/repository.go -package=mock
	mockgen -source=./internal/service_new/service.go -destination=./internal/mock/service.go          -package=mock

## clean: clean
clean:
	rm ./build/ -r