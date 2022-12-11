FROM golang:alpine AS build
RUN apk update && \
    apk add curl \
            git \
            bash \
            make \
            ca-certificates && \
    rm -rf /var/cache/apk/*

WORKDIR /app

COPY go.* ./
RUN go mod download
RUN go mod verify

COPY . .
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates bash
WORKDIR /app/
COPY --from=build /app/server .
COPY --from=build /app/cmd/server/entrypoint.sh .
COPY --from=build /app/config/*.yml ./config/
RUN ls -la
RUN chmod +x entrypoint.sh
ENTRYPOINT ["./entrypoint.sh"]