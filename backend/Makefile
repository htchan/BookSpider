all: controller backend
	mkdir -p ./build
	mv ./controller ./build/
	mv ./backend ./build/
	cp ./src/public/* ./build -r

controller:
	go build ./src/main/controller.go

backend:
	go build ./src/main/backend.go

test:
	cd ./src/helper && go test
	cd ./src/model && go test

clean:
	rm ./build/ -r