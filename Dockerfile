FROM golang:latest

RUN go version


COPY . /CRUD_API
WORKDIR /CRUD_API

RUN go mod download && go get -u ./...
RUN go build -o CRUD_API ./cmd/main.go

CMD ["./CRUD_API"]



#Dockerfile-файл с настройками и зависимостями приложения.Докер-образ(изобр.(image))-это результат команды docker build(содержащий все необх.настройки)