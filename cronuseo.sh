#!/bin/sh

cd src
golangci-lint run
if [ $# -eq 1 ]; then
    export PATH=$(go env GOPATH)/bin:$PATH
    swag init -g server.go  
fi
go run server.go
