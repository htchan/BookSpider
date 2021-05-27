all: controller backend
	mkdir -p ./build
	mv ./controller ./build/
	mv ./backend ./build/
	cp ./src/public/* ./build -r

controller: ./src/main/controller.go
	go build ./src/main/controller.go

backend: ./src/main/controller.go
	go build ./src/main/backend.go

test: ./src/helper/helper ./src/model
	cd ./src/helper && go test
	cd ./src/model && go test

clean:
	rm ./build/ -r