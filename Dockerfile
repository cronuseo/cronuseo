FROM golang:1.19.1-alpine3.16

COPY . /api
WORKDIR /api

CMD go mod download; go run cmd/server/main.go     

EXPOSE 8080