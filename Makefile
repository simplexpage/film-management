# Description: Makefile for film-management
start_local:
	docker-compose up

stop_local:
	docker-compose down

console_local:
	docker exec -ti go_jwt_auth bash

migrate_local:
	docker-compose run --rm go_film_management bash -c "./cmd/tmp/main -migrate-postgres-database"

format_local:
	docker-compose run --rm go_jwt_auth bash -c "gofmt -s -w ."

all: lint test

# Init swagger
init_swagger:
	docker-compose run --rm go_film_management bash -c "swag init -g internal/common/transport/http/http.go"

# Run lint
lint:
	golangci-lint run

.PHONY: cover
cover:
	go test ./... -short -count=1 -race -coverprofile=coverage.out .. && go tool cover -html coverage.out && rm coverage.out

race:
	go test -v -race -count=1 ./...

test:
	go test -v -count=1 ./...
