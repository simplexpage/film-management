# film-management

## Project structure

```
├── cmd
│   ├── server - temporary folder for local development
│   │   ├── commands
│   │   │   └── migrate - command for apply migrations
│   │   └── main.go - entry point for launch project
│   └── tmp - temporary folder for local development
├── config - config files
│   ├── ssl - ssl certificates
│   └── config.yaml.example - example config file
├── docs - documentation for swagger
├── internal
│   ├── common - common service
│   ├── film - film service
│   │    ├── domain - film domain
│   │    ├── endpoints - film endpoints
│   │    └── transport - film transport
│   │       └── http - http transport
│   └── user - user service
│       ├── domain - user domain
│       ├── endpoints - user endpoints
│       └── transport - user transport
│           └── http - http transport
├── pkg - common packages for project
├── repositories - external repositories like postgreDB, redis, etc.
│   └── storage - storage repositories
├── Dockerfile - dockerfile for local project
├── Dockerfile_prod - dockerfile for production project
├── Makefile - makefile for local project
└── .env.example - example .env file
```

## Prepare for a local start

### 1. Install docker

To run the project you have to install **docker**.

You can read about installation here https://docs.docker.com/install/, just choose your OS.

For UNIX users - nothing else.

For WINDOWS users - you have to install MAKE by your own.

### 2. Create a new config file

Copy example config file from /config/config.yaml.example in /config/config.yaml

### 3. Create a new .env file

Copy example .env file from /.env.example in /.env

### 4. Run a project

Run `make start_local` to start REST API. All containers will start automatically.

### 5. Run migrations

Run `make migrate_local` in another console to apply Postgres db migrations for film-management service.

### 5. Run seed test data

Run `make seed_data_local` in another console to apply Postgres db seed test data for film-management service.

### 7. Stop a project
Run `make stop_local` to stop REST API.

## Run tests

Run `make test` to run tests inside GO container.

## Swagger for REST API

http://localhost:8088/api/v1/film/swagger/index.html

## Database

### Postgresql schema
https://dbdiagram.io/d/film-management-6545f3d87d8bbd646577e9de

## Prometheus metrics

http://localhost:8078/metrics


## Docker for production

### 1. Build docker image

Run `make build_prod` to build docker image.
This command will create a Docker image named my_app:latest based on the Dockerfile_prod.

### 2. Run docker image

Run `make run_prod` to run docker image.

This command starts a Docker container and forwards ports 8080 and 8081, allowing you to access your application on ports 8088 and 8078 on your local machine.

Note: This example uses ports 8088 and 8078 for forwarding, but you can use other ports at your discretion.

Your application will now be running inside a Docker container and accessible at http://localhost:8088 and http://localhost:8078.



## Adding and Updating a Movie

* **Genres**: When adding or updating a movie, if the database does not have a corresponding genre, an appropriate error will be returned.
* **Director**: When adding or updating a movie, if the database does not have an existing director, a new director will be created.
* **Actors**: When adding or updating a movie, if the database does not have an existing actor, a new actor will be created.