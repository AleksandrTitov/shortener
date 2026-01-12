# db migrations

### Установка библиотеки

```shell
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest 
```

```shell
go get -u github.com/golang-migrate/migrate/v4 
```

### Создание миграции

```shell
migrate create -ext sql -dir ./migrations -seq <имя миграции>
```

### Применение

```shell
migrate -database "postgres://postgres:postgres@localhost:5432/app?sslmode=disable" -path ./migrations up
```

_Пример_
```shell
migrate -database "postgres://postgres:postgres@localhost:5432/app?sslmode=disable" -path ./migrations up
```

### Откат

```shell
migrate -database "<dsn>" -path ./migrations down
```

_Пример_
```shell
migrate -database "postgres://postgres:postgres@localhost:5432/app?sslmode=disable" -path ./migrations down
```