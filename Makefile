# Description: Makefile for film-management

# Local development
start_local:
	docker-compose up

stop_local:
	docker-compose down

console_local:
	docker exec -ti go_film_management bash

migrate_local:
	docker-compose run --rm go_film_management bash -c "./cmd/tmp/main -migrate-postgres-database"

seed_data_local:
	docker-compose run --rm go_film_management bash -c "./cmd/tmp/main -seed-postgres-database"

format_local:
	docker-compose run --rm go_film_management bash -c "gofmt -s -w ."

# Init swagger
init_swagger:
	docker-compose run --rm go_film_management bash -c "swag init -g internal/common/transport/http/http.go"

# Run test
.PHONY: cover
cover:
	go test ./... -short -count=1 -race -coverprofile=coverage.out .. && go tool cover -html coverage.out && rm coverage.out

race:
	go test -v -race -count=1 ./...

test:
	go test -v -count=1 ./...


# Production

# Defining Variables
IMAGE_NAME = my_app
TAG = latest

# Path to Dockerfile for production
DOCKERFILE = Dockerfile_prod

# Target for building Docker image
build_prod:
	docker build -f $(DOCKERFILE) -t $(IMAGE_NAME):$(TAG) .

# Target for running a Docker container
run_prod:
	docker run -d -p 8088:8080 -p 8078:8081 $(IMAGE_NAME):$(TAG)