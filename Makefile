.PHONY: frontend


all: controller backend #frontend
	mkdir -p ./build
	mv ./controller ./build/
	mv ./backend ./build/
	#rm ./build/frontend -rf
	#mv ./src/frontend/build/ ./build/frontend
	cp ./src/public/* ./build -r

run:
	cd ./bin && nohup ./backend &
	cd ./bin/frontend && nohup python3 -m http.server 8427 > ../frontend.log &

test:
	cd ./src/helper && go test
	cd ./src/model && go test

controller:
	go build ./src/main/controller.go

backend:
	go build ./src/main/backend.go

frontend:
	cd ./src/frontend && npm run-script build

clean:
	rm ./build/ -r
