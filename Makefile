.PHONY: frontend backend controller build

build:
	docker-compose --profile all build ${service}

frontend:
	docker-compose --profile frontend up

backend:
	docker-compose --profile backend up

controller:
	command=${command} params=${params} docker-compose --profile controller up

test:
	docker-compose --profile test up