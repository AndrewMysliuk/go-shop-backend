postgres:
	docker run --name db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:12-alpine

createdb:
	docker exec -it db createdb --username=root --owner=root shop-db

dropdb:
	docker exec -it db dropdb shop-db

migrateup:
	migrate -path ./schema -database "postgresql://root:root@localhost:5432/shop-db?sslmode=disable" -verbose up

migratedown:
	migrate -path ./schema -database "postgresql://root:root@localhost:5432/shop-db?sslmode=disable" -verbose down

swag:
	swag init -g cmd/main.go

test:
	go test -v -cover ./...

.PHONY:postgres createdb dropdb migrateup migratedown swag test