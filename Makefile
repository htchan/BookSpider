.PHONY: frontend backend controller build

build:
	DOCKER_BUILDKIT=1 docker-compose --profile all build ${service}

build_flutter:
	cd frontend ; docker build . -f Dockerfile.flutter -t flutter:stable

frontend:
	docker-compose --profile frontend up

backend:
	docker-compose --profile backend up -d

controller:
	command=${command} params="${params}" docker-compose --profile controller up --force-recreate

batch:
	docker-compose --profile batch up -d

test:
	docker-compose --profile test up
