.PHONY: frontend backend controller build

build:
	docker build -f ./frontend/Dockerfile.flutter -t flutter ./frontend
	docker-compose --profile all build ${service}

frontend:
	docker-compose --profile frontend up

backend:
	docker-compose --profile backend up -d

controller:
	command=${command} params="${params}" docker-compose --profile controller up

test:
	docker-compose --profile test up
