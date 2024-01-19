build:
	go build -o ./.bin/api cmd/main.go

run: build
	./.bin/api

test:
	go test -v ./...

migrate:
	migrate -path ./schema -database "postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable" up

swag:
	swag init -g cmd/main.go