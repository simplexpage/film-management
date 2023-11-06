# Description: Makefile for film-management
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

.PHONY: cover
cover:
	go test ./... -short -count=1 -race -coverprofile=coverage.out .. && go tool cover -html coverage.out && rm coverage.out

race:
	go test -v -race -count=1 ./...

test:
	go test -v -count=1 ./...
