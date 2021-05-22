shell := /bin/bash

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: run
run:
	go run cmd/api/main.go

.PHONY: migrate
migrate:
	go run cmd/admin/main.go migrate

.PHONY: resetdb
resetdb:
	rm test.db*
	go run cmd/admin/main.go migrate

.PHONY: test
test:
	go test -v ./...

.PHONY: generate
generate:
	# see https://github.com/99designs/gqlgen/issues/1483 and https://github.com/golang/go/issues/44129
	go get github.com/99designs/gqlgen/cmd@v0.13.0 && go run github.com/99designs/gqlgen generate