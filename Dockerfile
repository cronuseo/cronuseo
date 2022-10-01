FROM golang:1.19.1-alpine3.16

COPY . /api
WORKDIR /api

CMD cd src; go mod download; go run main.go

EXPOSE 8080