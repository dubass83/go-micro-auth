.PHONY: *

ENV = dev
DB_URL = "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable"

docker_up:
	limactl start docker

docker_down: 
	limactl stop docker

postgres_up:
	docker start ${ENV}-postgres \
	|| docker run --name ${ENV}-postgres \
	-e POSTGRES_PASSWORD=postgres \
	-e POSTGRES_DB=simple_bank \
	-p 5432:5432 -d postgres

postgres_down:
	docker stop ${ENV}-postgres
	docker rm ${ENV}-postgres

sqlc:
	sqlc generate

db_docs:
	dbdocs build docs/db.dbml

db_schema:
	dbml2sql --postgres -o docs/schema.sql docs/db.dbml