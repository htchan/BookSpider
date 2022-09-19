.PHONY: console

table = test

## help: show available command and description
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed  -e 's/^/ /'

## all: build all executable and put them to bin
all: controller backend batch

## controller: build executable of controller
controller:
	go build -o ../../bin github.com/htchan/BookSpider/cmd/controller

## backend: build executable of backend
backend:
	go build -o ../../bin github.com/htchan/BookSpider/cmd/backend

## batch: build executable of batch
batch:
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

coverage: 
	# go clean --testcache
	go test $(PKG) -coverprofile profile.txt ; go tool cover -html=profile.txt -o coverage.html
	rm profile.txt
	google-chrome ./coverage.html &

create_database:
	PGPASSWORD=books psql -h localhost -p 5432 -U books -c "create database ${table}"

populate_env:
	export $(cat ../assets/.env.test | xargs)

## apply_db_schema: construct ./assets/template.db by sql writern in ./assets/schema/schema.sql
apply_db_schema: ../assets/schema/schema.sql
	# rm -f ../assets/template.db
	# cat ../assets/schema/schema.sql | sqlite3 ../assets/template.db
	PGPASSWORD=books psql -h localhost -p 5432 -U books -d ${table} -c "drop table if exists books, writers, errors"
	PGPASSWORD=books psql -h localhost -p 5432 -U books -d ${table} -c "drop index if exists books_index, books_title, books_writer, books_status, writers_id, writers_name, errors_index"
	PGPASSWORD=books psql -h localhost -p 5432 -U books -d ${table} -f ../assets/schema/schema.sql

## populate_test_data: construct ./assets/test-data/internal_database_sqlite.db by sql writern in
populate_test_data:
	cp ../assets/template.db ../assets/test-data/internal_database_sqlite.db
	# cat ../assets/schema/internal_database_sqlite_schema.sql | sqlite3 ../assets/test-data/internal_database_sqlite.db

## console: compile and run console
console:
	go build github.com/htchan/BookSpider/cmd/console -o ../../bin
	mv ../../bin/console .
	./console
	rm ./console

## clean: clean
clean:
	rm ./build/ -r