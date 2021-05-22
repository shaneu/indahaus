module github.com/shaneu/indahaus

go 1.16

require (
	github.com/99designs/gqlgen v0.13.0
	github.com/google/go-cmp v0.3.0
	github.com/google/uuid v1.2.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/labstack/echo/v4 v4.3.0
	github.com/mattn/go-sqlite3 v1.14.7
	github.com/pkg/errors v0.8.1
	github.com/spf13/viper v1.7.1
	// explicitly requiring v2.1.0, please see https://github.com/99designs/gqlgen/issues/1402 and https://stackoverflow.com/a/67187051/7571000
	github.com/vektah/gqlparser/v2 v2.1.0
)
