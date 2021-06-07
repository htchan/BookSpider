pwd:=$(shell pwd)
user:=$(shell whoami)
frontend_src_volume = $(pwd)/bin/frontend
frontend_dst_volume = novel_frontend_volume

database_volume = $(pwd)/bin/database
log_volume = $(pwd)/bin/log
book_volume = /mnt/addition/download/Books
backup_volume = $(pwd)/bin/backup

.PHONY: frontend backend controller build

build:
	docker build -f ./backend/Dockerfile.backend -t novel_backend ./backend
	docker build -f ./backend/Dockerfile.controller -t novel_controller ./backend
	docker build -f ./backend/Dockerfile.backup -t novel_backup ./backend/src/operation
	docker image prune -f

frontend:
	docker run -v ${frontend_dst_volume}:/frontend \
		--name novel_frontend busybox true
	docker cp ${frontend_src_volume}/. novel_frontend:/frontend/
	docker rm novel_frontend

backend:
	docker run --name novel_backend_container -d \
		--network=router \
		-v ${database_volume}:/database \
		-v ${log_volume}:/log \
		-v ${book_volume}:/books \
		novel_backend sh -c "./backend > /log/backend.log"

controller:
	docker run --rm --name novel_controller_container \
		-v ${database_volume}:/database \
		-v ${log_volume}:/log \
		-v ${book_volume}:/books \
		novel_controller sh -c "./controller ${command} >> /log/controller.log"

backup:
	docker run --rm --name novel_backup_container \
		-v ${database_volume}:/database \
		-v ${backup_volume}:/backup \
		novel_backup python ./backup.py
	sudo chown -R ${user} ${backup_volume}
