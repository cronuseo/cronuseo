FROM golang:1.19-alpine 

ENV GO111MODULE=on

ARG APP_ENV
ENV APP_ENV=$APP_ENV

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod tidy
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server
EXPOSE 3333
ENTRYPOINT ./server -config "./config/${APP_ENV}.yml"