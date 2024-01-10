include .env

DOCKER_EXEC := docker exec -it
DB_URL := "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"
MIGRATE_CMD := migrate -path db/migrations -database ${DB_URL}

## start: Start PostgreSQL container
postgres:
	docker run --name ${DOCKER_IMAGE_NAME} -p ${POSTGRES_PORT}:${POSTGRES_PORT} -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:${POSTGRES_VERSION}

## createdb: Create PostgreSQL database and enable citext extension if exists skip
createdb:
	$(DOCKER_EXEC) ${POSTGRES_CONTAINER_NAME} createdb --username=${POSTGRES_USER} ${POSTGRES_DB}

## addcitext: Add citext extension to PostgreSQL database
addcitext:
	$(DOCKER_EXEC) ${POSTGRES_CONTAINER_NAME} psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "CREATE EXTENSION IF NOT EXISTS citext;"

## psql: Access psql shell
psql:
	$(DOCKER_EXEC) ${POSTGRES_CONTAINER_NAME} psql -U ${POSTGRES_USER} -d ${POSTGRES_DB}

## dropdb: Drop PostgreSQL database
dropdb:
	$(DOCKER_EXEC) ${POSTGRES_CONTAINER_NAME} dropdb --username=${POSTGRES_USER} ${POSTGRES_DB}

## migrateup: Apply database migrations
migrateup:
	$(MIGRATE_CMD) up

## migrateup: Apply database migrations (force) & with parameters
migrateupf:
	$(MIGRATE_CMD) force ${VERSION}
	$(MIGRATE_CMD) up

## migratedown: Rollback database migrations
migratedown:
	$(MIGRATE_CMD) down

## migratedown: Rollback database migrations (force) & with parameters
migratedownf:
	$(MIGRATE_CMD) force ${VERSION}
	$(MIGRATE_CMD) down

## test: Run tests
test:
	go test -v -cover ./...

## up: Start Docker Compose services
up:
	docker-compose up -d

## down: Stop Docker Compose services
down:
	docker-compose down

wait:
	sleep 5

## init: Start PostgreSQL container, create database and apply migrations
init: up wait addcitext migrateup

## downclean: Take down services and clean database
downclean: migratedown dropdb down 

## help: Show this help.
help: Makefile
	@sed -n 's/^##//p' $<

.PHONY: postgres createdb dropdb migrateup migratedown test psql up down init downclean help