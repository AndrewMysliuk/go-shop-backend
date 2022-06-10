# CRUD API для магазина

### Стэк
- go 1.17
- postgres 

### Запуск
```go run cmd/main.go```

Для postgres можно использовать Docker

```go
docker run --name db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:12-alpine
```

```go
docker exec -it db createdb --username=root --owner=root shop-db
```

```go
migrate -path ./schema -database "postgresql://root:root@localhost:5432/shop-db?sslmode=disable" -verbose up
```

<!-- ### Swagger UI
```http://localhost:8080/swagger/index.html#/``` -->