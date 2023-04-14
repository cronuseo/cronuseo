FROM golang:1.19-alpine as build

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor

ARG APP_ENV
ENV APP_ENV=$APP_ENV

WORKDIR /app

COPY go.mod .
COPY go.sum .

COPY cmd cmd
COPY config config
COPY docs docs
COPY internal internal
COPY vendor vendor

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o check_server ./cmd/check_server

FROM alpine:latest

WORKDIR /app/
COPY --from=build /app/check_server .
COPY --from=build /app/config/*.yml ./config/
COPY --from=build /app/config/*.yml ./config/

ENTRYPOINT ./check_server -config "./config/${APP_ENV}.yml"