FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o CRUD_API ./cmd/main.go

CMD ["./CRUD_API"]