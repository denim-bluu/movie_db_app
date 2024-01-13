include .env

DOCKER_EXEC := docker exec -it
MIGRATE_CMD := migrate -path db/migrations -database ${DATABASE_URL}

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

.DEFAULT_GOAL := help
.PHONY: help
help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
	
.PHONY: wait
wait: ## Wait for 5 seconds
	@echo -n 'Waiting for 5 seconds...' && sleep 5

.PHONY: confirm
confirm: ## Confirm action
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

.PHONY: audit 
audit: ## audit: tidy dependencies and format, vet and test all code
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

.PHONY: vendor
vendor: ## vendor: tidy and vendor dependencies
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #
.PHONY: run/api
run/api: ## Run API server
	go run ./cmd/api

.PHONY: db/psql
db/psql: ## Access psql shell
	$(DOCKER_EXEC) ${POSTGRES_CONTAINER_NAME} psql -U ${POSTGRES_USER} -d ${POSTGRES_DB}

# ==================================================================================== #
# SETUP
# ==================================================================================== #
.PHONY: db/start
db/start: ## Start PostgreSQL container
	docker run --name ${DOCKER_IMAGE_NAME} -p ${POSTGRES_PORT}:${POSTGRES_PORT} -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:${POSTGRES_VERSION}

.PHONY: db/createdb
db/createdb: ## Create PostgreSQL database and enable citext extension if exists skip
	$(DOCKER_EXEC) ${POSTGRES_CONTAINER_NAME} createdb --username=${POSTGRES_USER} ${POSTGRES_DB}

.PHONY: db/addcitext
db/addcitext: ## Add citext extension to PostgreSQL database
	$(DOCKER_EXEC) ${POSTGRES_CONTAINER_NAME} psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} -c "CREATE EXTENSION IF NOT EXISTS citext;"

.PHONY: db/drop
db/drop: ## Drop PostgreSQL database
	$(DOCKER_EXEC) ${POSTGRES_CONTAINER_NAME} dropdb --username=${POSTGRES_USER} ${POSTGRES_DB}

.PHONY: docker/up
docker/up: ## Start Docker Compose services
	docker-compose up -d

.PHONY: docker/down
docker/down: ## down: Stop Docker Compose services
	docker-compose down

.PHONY: db/migration/up
db/migration/up: ## Apply database migrations
	$(MIGRATE_CMD) up

.PHONY: db/migration/up/force
db/migration/up/force: ## Apply database migrations (force) & with parameters
	$(MIGRATE_CMD) force ${VERSION}
	$(MIGRATE_CMD) up

.PHONY: db/migration/down
db/migration/down: ## Rollback database migrations
	$(MIGRATE_CMD) down

.PHONY: db/migration/down/force
db/migration/down/force: ## Rollback database migrations (force) & with parameters
	$(MIGRATE_CMD) force ${VERSION}
	$(MIGRATE_CMD) down

.PHONY: setup/init 
setup/init: confirm docker/up wait db/addcitext db/migration/up ## Start PostgreSQL container, create database and apply migrations

.PHONY: setup/teardown
setup/teardown: confirm db/migration/down db/drop docker/down ## Take down services and clean database

# ==================================================================================== #
# BUILD
# ==================================================================================== #

.PHONY: build/api
build/api:  ## build/api: build the cmd/api application
	@echo 'Building cmd/api...'
	go build -ldflags='-s -w' -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o=./bin/linux_amd64/api ./cmd/api