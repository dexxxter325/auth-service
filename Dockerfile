FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o CRUD_API ./cmd/main.go

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate.linux-amd64 /usr/local/bin/migrate

CMD ["./CRUD_API"]