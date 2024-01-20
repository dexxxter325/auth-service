build: #собираем образ компиляции
	go build -o ./.bin/api cmd/main.go  #main.go-файл,который будет преобразован в бинарный

run: build #компилируем
	./.bin/api

test:
	go test ./...

migrate:
	migrate -path ./schema -database "postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable" up

migrate-down:
	migrate -path ./schema -database "postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable" down

swag:
	swag init -g cmd/main.go