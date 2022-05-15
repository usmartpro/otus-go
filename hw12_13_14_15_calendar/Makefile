BIN := "./bin/calendar"
BIN_SCHEDULER := "./bin/calendar_scheduler"
BIN_SENDER := "./bin/calendar_sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

DB_CONN := "postgresql://calendar:postgres@localhost:5432/calendar?sslmode=disable"

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

run: build
	$(BIN) -config ./configs/calendar_config.yaml

run-calendar:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	$(BIN) -config ./configs/calendar_config.yaml

run-scheduler:
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	$(BIN_SCHEDULER) -config ./configs/scheduler_config.yaml

run-sender:
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender
	$(BIN_SENDER) -config ./configs/sender_config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

migrate:
	goose --dir=migrations postgres ${DB_CONN} up

generate:
	go generate ./...

.PHONY: build run build-img run-img version test lint
