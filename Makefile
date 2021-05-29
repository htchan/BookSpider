pwd:=$(shell pwd)
frontend_src_volume = $(pwd)/bin/frontend
frontend_dst_volume = novel_frontend_volume

database_volume = $(pwd)/bin/database
log_volume = $(pwd)/bin/logs
book_volume = /mnt/additional/download/Books

.PHONY: frontend backend controller

frontend:
	docker run -v ${frontend_src_volume}:/source \
		-v ${frontend_dst_volume}:/frontend \
		alpine cp /source/* /frontend -run

backend:
	docker build -f ./backend/Dockerfile.backend -t NovelBackend
	docker run --name novel_backend -d \
		-v database_volume:/database \
		-v log_volume:/log \
		-v book_volume:/book \
		NovelBackend ./backend > /log/backend.log

controller:
	docker build -f ./backend/Dockerfile.controller -t NovelController
	docker run --name novel_backend -d \
		-v database_volume:/database \
		-v log_volume:/log \
		-v book_volume:/book \
		NovelController ./controller ${command} > /log/controller.log