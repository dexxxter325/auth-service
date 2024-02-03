# Stage 1: Build
FROM golang:1.21-alpine3.19 AS builder

RUN go version

WORKDIR /CRUD_API

COPY . .

RUN go mod download && go get -u ./...
RUN go build -o ./bin/api ./cmd/main.go

# Stage 2: Final Image
FROM alpine:latest

WORKDIR /CRUD_API

COPY --from=builder /CRUD_API/.env .
COPY --from=builder /CRUD_API/bin/api .

EXPOSE 80

CMD ["./api"]


#выполняется бинарный запуск приложения



#Dockerfile-файл с настройками и зависимостями приложения.Докер-образ(изобр.(image))-это результат команды docker build(содержащий все необх.настройки)