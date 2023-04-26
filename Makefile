MODULE = $(shell go list -m)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || echo "1.0.0")
PACKAGES := $(shell go list ./... | grep -v /vendor/)
LDFLAGS := -ldflags "-X main.Version=${VERSION}"

CONFIG_FILE ?= ./config/local.yml
APP_DSN ?= $(shell sed -n 's/^dsn:[[:space:]]*"\(.*\)"/\1/p' $(CONFIG_FILE))

.PHONY: default
default: help

.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: ## run the c6o server
	make -j 2 run-mgt-server run-check-server

.PHONY: run-mgt-server
run-mgt-server: ## run the c6o mgt server
	go run ${LDFLAGS} cmd/server/main.go -config ${CONFIG_FILE}

.PHONY: run-check-server
run-check-server: ## run the c6o check server
	go run ${LDFLAGS} cmd/check_server/main.go -config ${CONFIG_FILE}

.PHONY: build
build:  ## build the c6o server binary
	make -j 2 build-mgt-server build-check-server

.PHONY: build-mgt-server
build-mgt-server:  ## build the c6o mgt server binary
	CGO_ENABLED=0 go build ${LDFLAGS} -a -o server $(MODULE)/cmd/server

.PHONY: build-check-server
build-check-server:  ## build the c6o check server binary
	CGO_ENABLED=0 go build ${LDFLAGS} -a -o server $(MODULE)/cmd/check_server

.PHONY: build-docker
build-docker: ## build the servers as a docker image
	make -j 2  build-mgt-server-docker build-check-server-docker

.PHONY: build-mgt-server-docker
build-mgt-server-docker: ## build the mgt server as a docker image
	docker build -f cmd/server/Dockerfile -t server .

.PHONY: build-check-server-docker
build-check-server-docker: ## build the check server as a docker image
	docker build -f cmd/check_server/Dockerfile -t check_server .

.PHONY: version
version: ## display the version of the API server
	@echo $(VERSION)

.PHONY: lint4

lint: ## run golint on all Go package
	@golint $(PACKAGES)

.PHONY: setup-db
setup-db: ## strat the databases
	docker compose -f docker-compose-db.yml up -d
	
.PHONY: start
start: ## start backend, frontend, keto and redis
	docker compose up

.PHONY: test
test: ## run unit tests
	@echo "mode: count" > coverage-all.out
	@$(foreach pkg,$(PACKAGES), \
		go test -p=1 -cover -covermode=count -coverprofile=coverage.out ${pkg}; \
		tail -n +2 coverage.out >> coverage-all.out;)