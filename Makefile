pwd:=$(shell pwd)
frontend_src_volume = $(pwd)/bin/frontend
frontend_dst_volume = novel_frontend_volume

database_volume = $(pwd)/bin/database
log_volume = $(pwd)/bin/log
book_volume = /mnt/addition/download/Books
backup_volume = $(pwd)/bin/backup

.PHONY: frontend backend controller

frontend:
	docker run -v ${frontend_dst_volume}:/frontend \
		--name novel_frontend busybox true
	docker cp ${frontend_src_volume}/. novel_frontend:/frontend/
	docker rm novel_frontend

backend:
	docker build -f ./backend/Dockerfile.backend -t novel_backend ./backend
	docker image prune -f
	docker run --name novel_backend_container -d \
		--network=router \
		-v ${database_volume}:/database \
		-v ${log_volume}:/log \
		-v ${book_volume}:/books \
		novel_backend ./backend > ./backend.log

controller:
	docker build -f ./backend/Dockerfile.controller -t novel_controller ./backend
	docker run --name novel_controller_container -d \
		-v ${database_volume}:/database \
		-v ${log_volume}:/log \
		-v ${book_volume}:/books \
		novel_controller ./controller ${command} > /log/controller.log

backup:
	echo backup start
	docker build -f ./backend/Dockerfile.backup -t novel_backup ./backend/src/operation
	docker image prune -f
	docker run --rm --name novel_backup_container -d \
		-v ${database_volume}:/database \
		-v ${backup_volume}:/backup \
		novel_backup ./backup.py