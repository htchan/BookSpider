all: controller backend
	mkdir -p ./build
	mv ./src/controller ./build/
	mv ./src/backend ./build/
	cp ./src/public/* ./build -r

controller: ./src/main/controller.go
	cd ./src && go build ./main/controller.go

backend: ./src/main/controller.go
	cd ./src && go build ./main/backend.go

test: ./src/helper/helper ./src/model
	cd ./src/helper && go test
	cd ./src/model && go test

clean:
	rm ./build/ -r