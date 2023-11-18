FROM golang:1.19-alpine as build

ENV GO111MODULE=on

ARG APP_ENV
ENV APP_ENV=$APP_ENV

WORKDIR /app

COPY go.mod .
COPY go.sum .

COPY cmd cmd
COPY config config
COPY docs docs
COPY internal internal
COPY proto proto

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o check_server ./cmd/check_server

FROM alpine:latest

WORKDIR /app/
COPY --from=build /app/check_server .
COPY --from=build /app/config/*.yml ./config/
COPY --from=build /app/config/*.yml ./config/

# Create a non-root user
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid 10014 \
    "cronuseo"
# Use the above created unprivileged user
USER 10014

ENTRYPOINT ./check_server -config "./config/${APP_ENV}.yml"
