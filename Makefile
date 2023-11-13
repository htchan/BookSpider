.PHONY: frontend backend controller build start

service ?= all

## help: show available command and description
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed  -e 's/^/ /'

## backend-*: redirect the command to makefile in backend directory
backend-%:
	make -C backend $*

## build service=<service>: build docker image of specified service (default all)
build:
	docker buildx bake backend

## build_flutter: build the image for compile flutter frontend
build_flutter:
	cd frontend ; docker build . -f Dockerfile.flutter -t flutter:stable

## frontend: compile flutter frontend
frontend:
	docker compose --profile frontend up

## backend: deploy backend container
api:
	docker compose up -d api

## batch: deploy batch container
worker:
	docker compose up -d --force-recreate worker

start:
	docker compose pull api worker
	docker compose up -d api worker
