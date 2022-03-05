.PHONY: console
all: controller backend
	mkdir -p ./build
	mv ./src/controller ./build/
	mv ./src/backend ./build/
	cp ./src/public/* ./build -r

controller:
	go build -o ../../bin github.com/htchan/BookSpider/cmd/controller

backend:
	go build github.com/htchan/BookSpider/cmd/backend -o ../../bin

# test: ./src/helper/ ./src/model
# 	cd ./src && go test github.com/htchan/BookSpider/helper -v
# 	cd ./src && go test github.com/htchan/BookSpider/model -v
# 	cd ./src && go test github.com/htchan/BookSpider/public -v

test: 
	# go clean --testcache
	go test github.com/htchan/BookSpider/internal/logging
	go test github.com/htchan/BookSpider/internal/database
	go test github.com/htchan/BookSpider/internal/database/sqlite
	go test github.com/htchan/BookSpider/internal/utils
	go test github.com/htchan/BookSpider/pkg/configs
	go test github.com/htchan/BookSpider/pkg/books
	go test github.com/htchan/BookSpider/pkg/flags
	go test github.com/htchan/BookSpider/pkg/sites
	go test github.com/htchan/BookSpider/cmd/backend
	go test github.com/htchan/BookSpider/cmd/controller

coverage: 
	# go clean --testcache
	go test -coverprofile profile.txt \
		github.com/htchan/BookSpider/internal/database \
		github.com/htchan/BookSpider/internal/database/sqlite \
		github.com/htchan/BookSpider/internal/utils \
		github.com/htchan/BookSpider/pkg/configs \
		github.com/htchan/BookSpider/pkg/books \
		github.com/htchan/BookSpider/pkg/flags

apply_db_schema: ../assets/schema/schema.sql
	rm -f ../assets/template.db
	cat ../assets/schema/schema.sql | sqlite3 ../assets/template.db

populate_test_data:
	cp ../assets/template.db ../assets/test-data/internal_database_sqlite.db
	cat ../assets/schema/internal_database_sqlite_schema.sql | sqlite3 ../assets/test-data/internal_database_sqlite.db

console:
	go build github.com/htchan/BookSpider/cmd/console -o ../../bin
	mv ../../bin/console .
	./console
	rm ./console

clean:
	rm ./build/ -r