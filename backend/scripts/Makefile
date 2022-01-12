all: controller backend
	mkdir -p ./build
	mv ./src/controller ./build/
	mv ./src/backend ./build/
	cp ./src/public/* ./build -r

controller: ./src/main/controller.go
	cd ./src && go build ./main/controller.go

backend: ./src/main/controller.go
	cd ./src && go build ./main/backend.go

# test: ./src/helper/ ./src/model
# 	cd ./src && go test github.com/htchan/BookSpider/helper -v
# 	cd ./src && go test github.com/htchan/BookSpider/model -v
# 	cd ./src && go test github.com/htchan/BookSpider/public -v

test: 
	go test github.com/htchan/BookSpider/internal/database
	go test github.com/htchan/BookSpider/internal/database/sqlite
	go test github.com/htchan/BookSpider/internal/utils
	go test github.com/htchan/BookSpider/pkg/configs

apply_db_schema: ../assets/schema/schema.sql
	rm -f ../assets/template.db
	cat ../assets/schema/schema.sql | sqlite3 ../assets/template.db

populate_test_data:
	cp ../assets/template.db ../assets/test-data/internal_database_sqlite.db
	cat ../assets/schema/internal_database_sqlite_schema.sql | sqlite3 ../assets/test-data/internal_database_sqlite.db

clean:
	rm ./build/ -r