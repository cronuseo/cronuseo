FROM golang:1.19-alpine as build

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

ENV GO111MODULE=on

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

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM alpine:latest

WORKDIR /app/
COPY --from=build /app/server .
COPY --from=build /app/config/*.yml ./config/
COPY --from=build /app/config/*.yml ./config/
COPY --from=build /app/policy/*.rego ./policy/

ENTRYPOINT ./server -config "./config/${APP_ENV}.yml"
