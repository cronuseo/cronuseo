FROM golang:1.19.1-alpine3.16

COPY . /api
WORKDIR /api
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/server/main.go
EXPOSE 8080
RUN go mod download
RUN cd cmd/server && go build -o ./../../
ENV GO111MODULE=off
RUN pwd && ls
CMD [ "./server" ]