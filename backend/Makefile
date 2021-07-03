all: controller backend
	mkdir -p ./build
	mv ./src/controller ./build/
	mv ./src/backend ./build/
	cp ./src/public/* ./build -r

controller: ./src/main/controller.go
	cd ./src && go build ./main/controller.go

backend: ./src/main/controller.go
	cd ./src && go build ./main/backend.go

test: ./src/helper/ ./src/model
	cd ./src && go test github.com/htchan/BookSpider/helper -v
	cd ./src && go test github.com/htchan/BookSpider/model -v

clean:
	rm ./build/ -r