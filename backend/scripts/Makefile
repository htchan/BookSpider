.PHONY: console

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
	go clean --testcache
	# go test github.com/htchan/BookSpider/internal/logging -coverprofile ./profile.out
	# go test github.com/htchan/BookSpider/internal/database -coverprofile ./profile.out
	# go test github.com/htchan/BookSpider/internal/database/sqlite -coverprofile ./profile.out
	# go test github.com/htchan/BookSpider/internal/utils -coverprofile ./profile.out
	# go test github.com/htchan/BookSpider/pkg/configs -coverprofile ./profile.out
	# go test github.com/htchan/BookSpider/pkg/books -coverprofile ./profile.out
	# go test github.com/htchan/BookSpider/pkg/flags -coverprofile ./profile.out
	# go test github.com/htchan/BookSpider/pkg/sites -coverprofile ./profile.out
	go test github.com/htchan/BookSpider/cmd/backend -coverprofile ./profile.out
	# go test github.com/htchan/BookSpider/cmd/controller -coverprofile ./profile.out
	go tool cover -html=profile.out -o coverage.html

coverage: 
	# go clean --testcache
	go test -coverprofile profile.txt \
		github.com/htchan/BookSpider/internal/database \
		github.com/htchan/BookSpider/internal/database/sqlite \
		github.com/htchan/BookSpider/internal/utils \
		github.com/htchan/BookSpider/pkg/configs \
		github.com/htchan/BookSpider/pkg/books \
		github.com/htchan/BookSpider/pkg/flags

## apply_db_schema: construct ./assets/template.db by sql writern in ./assets/schema/schema.sql
apply_db_schema: ../assets/schema/schema.sql
	rm -f ../assets/template.db
	cat ../assets/schema/schema.sql | sqlite3 ../assets/template.db

## populate_test_data: construct ./assets/test-data/internal_database_sqlite.db by sql writern in
populate_test_data:
	cp ../assets/template.db ../assets/test-data/internal_database_sqlite.db
	cat ../assets/schema/internal_database_sqlite_schema.sql | sqlite3 ../assets/test-data/internal_database_sqlite.db

## console: compile and run console
console:
	go build github.com/htchan/BookSpider/cmd/console -o ../../bin
	mv ../../bin/console .
	./console
	rm ./console

## clean: clean
clean:
	rm ./build/ -r