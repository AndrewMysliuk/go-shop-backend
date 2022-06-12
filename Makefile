postgres:
	docker run --name db -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=qwerty1234 -d postgres:12-alpine

createdb:
	docker exec -it db createdb --username=postgres --owner=postgres shop-db

dropdb:
	docker exec -it db dropdb shop-db

migrateup:
	migrate -path ./schema -database "postgresql://postgres:qwerty1234@localhost:5432/shop-db?sslmode=disable" -verbose up

migratedown:
	migrate -path ./schema -database "postgresql://postgres:qwerty1234@localhost:5432/shop-db?sslmode=disable" -verbose down

docker-compose:
	docker-compose up --build

swag:
	swag init -g cmd/main.go

test:
	go test -v -cover ./...

.PHONY:postgres createdb dropdb migrateup migratedown docker-compose swag test