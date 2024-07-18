postgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=test -d postgres:16-alpine

redis:
	docker run --name redis16 -p 5433:6379 -d redis

postgres-log:
	docker exec -it postgres16 psql -U root dev

create-db:
	docker exec -it postgres16 createdb --username=root --owner=root dev

drop-db:
	docker exec -it postgres16 dropdb dev

migrate-up:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/dev?sslmode=disable" -verbose up

new-migration:
	@[ "${name}" ] || ( echo ">> USAGE: make generate-migration name=CHANGE_ME"; exit 1 )
	migrate create -dir ./migrations/ -ext sql -seq ${name}
	rm ./migrations/*${name}.down.sql


.PHONY: all