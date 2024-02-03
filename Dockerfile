FROM golang:1.21-alpine3.19 AS builder

RUN go version

COPY . /CRUD_API


WORKDIR /CRUD_API

RUN go mod download && go get -u ./...
RUN go build -o ./bin/api ./cmd/main.go

#lightweight docker container with binary(чтобы докер образ создавался бинарно и занимал меньше места)
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /CRUD_API/.env .
COPY --from=0 /CRUD_API/bin/api .

EXPOSE 80

CMD ["./api"]

#выполняется бинарный запуск приложения



#Dockerfile-файл с настройками и зависимостями приложения.Докер-образ(изобр.(image))-это результат команды docker build(содержащий все необх.настройки)