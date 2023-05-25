FROM golang:1.19-alpine as build

# Create a non-root user
RUN useradd -m -u 10001 myuser
USER myuser

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
COPY policy policy
COPY internal internal
COPY vendor vendor

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM alpine:latest

WORKDIR /app/
COPY --from=build /app/server .
COPY --from=build /app/config/*.yml ./config/
COPY --from=build /app/config/*.yml ./config/
COPY --from=build /app/policy/*.rego ./policy/

ENTRYPOINT ./server -config "./config/${APP_ENV}.yml"
