FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download && go get -u ./...
RUN go build -o CRUD_API ./cmd/main.go

RUN go get -u -d -v github.com/golang-migrate/migrate/cmd/migrate
RUN go install -v github.com/golang-migrate/migrate/cmd/migrate

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate.linux-amd64 /usr/local/bin/migrate

CMD ["./CRUD_API"]