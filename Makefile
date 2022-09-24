.PHONY: frontend backend controller build

service ?= all

## help: show available command and description
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed  -e 's/^/ /'

## build service=<service>: build docker image of specified service (default all)
build:
	DOCKER_BUILDKIT=1 docker-compose --profile ${service} build

## build_flutter: build the image for compile flutter frontend
build_flutter:
	cd frontend ; docker build . -f Dockerfile.flutter -t flutter:stable

## frontend: compile flutter frontend
frontend:
	docker-compose --profile frontend up

## backend: deploy backend container
backend:
	docker-compose --profile backend up -d

## controller: deploy controller container
console:
	command=${command} params="${params}" docker-compose --profile console up --force-recreate

## batch: deploy batch container
batch:
	docker-compose --profile batch up -d --force-recreate

## test: deploy test container
test:
	docker-compose --profile test up
