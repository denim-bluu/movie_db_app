include .env

postgres:
	docker run --name ${DOCKER_IMAGE_NAME} -p ${POSTGRES_PORT}:${POSTGRES_PORT} -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:${POSTGRES_VERSION}
createdb:
	docker exec -it ${POSTGRES_CONTAINER_NAME} createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} ${POSTGRES_DB}
psql:
	docker exec -it ${POSTGRES_CONTAINER_NAME} psql -U ${POSTGRES_USER} -d ${POSTGRES_DB}
dropdb:
	docker exec -it ${POSTGRES_CONTAINER_NAME} dropdb ${POSTGRES_DB}
migrateup:
	migrate -path db/migrations -database "${DATABASE_URL}" -verbose up
migratedown:
	migrate -path db/migrations -database "${DATABASE_URL}" -verbose down
test:
	go test -v -cover ./...
docker_compose_up:
	docker-compose up -d
docker_compose_down:
	docker-compose down

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test psql docker_compose_up docker_compose_down 