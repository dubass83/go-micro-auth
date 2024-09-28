.PHONY: *

ENV = devel
DB_URL = "postgresql://postgres:postgres@localhost:5432/go_micro?sslmode=disable"

docker_up:
	limactl start docker

docker_down: 
	limactl stop docker

postgres_up:
	docker start ${ENV}-postgres \
	|| docker run --name ${ENV}-postgres \
	-e POSTGRES_PASSWORD=postgres \
	-e POSTGRES_DB=go_micro \
	-p 5432:5432 -d postgres

postgres_down:
	docker stop ${ENV}-postgres
	docker rm ${ENV}-postgres

new_migration:
	migrate create -ext sql -dir data/migration -seq ${name}

migrate_up:
	migrate -path data/migration -database ${DB_URL} -verbose up

migrate_up1:
	migrate -path data/migration -database ${DB_URL} -verbose up 1

migrate_down:
	migrate -path data/migration -database ${DB_URL} -verbose down

migrate_down1:
	migrate -path data/migration -database ${DB_URL} -verbose down 1

sqlc:
	sqlc generate

db_docs:
	dbdocs build docs/db.dbml

db_schema:
	dbml2sql --postgres -o docs/schema.sql docs/db.dbml

test:
	go test -v -cover -count=1 -short ./...

server:
	go run main.go