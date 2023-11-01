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
│   └── config.yaml.example - example config file
├── docs - documentation for swagger
├── internal
│   ├── common - common service
│   ├── film - film service
│   └── user - user service
├── pkg - common packages for project
├── repositories - external repositories like postgreDB, redis, etc.
├── Dockerfile - dockerfile for local project
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

## Swagger for REST API

http://localhost:8088/api/v1/film/swagger/index.html