# Spendr - Backend

Spendr is a minimalist Expense Trakcer app, this repository holds Spendr code for its backend services.

## ERD
DB Diagram: https://dbdiagram.io/d/Spendr-6548afd77d8bbd6465901f12

## Development

- Copy `config.example.yml` into `config.yml` in `configs` folder then adjust your env:
```bash
make env
```

- Adjust db config from `congfigs/config.yml` into `Makefile` (if not correct)

- Use [go migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) to run db migrations:
```bash
migrate -path pkg/db/migration -database "postgresql://<user>:<password>@localhost:5432/spendr?sslmode=disable" -verbose up
```
- if you want to create migration you can use: `migrate create -ext sql -dir pkg/db/migration -seq <migration_name>`

- Use [air](https://github.com/cosmtrek/air) for live reload and start server:
```bash
make server
```

- Or you can start server without live reload:
```bash
go run cmd/api/main.go
```
