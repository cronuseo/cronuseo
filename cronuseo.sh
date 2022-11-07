#!/bin/sh
cd src
golangci-lint run
export PATH=$(go env GOPATH)/bin:$PATH
swag init -g server.go  
go run server.go
