GO_TEST_FILES := $(shell find . -name '*_test.go' | grep -v /vendor/)
include .test.env
export

default: test-integration

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

test: docker-down docker-up
	go test -v -race -cover -count=1 ./...
	docker-compose down

deps:
	@go mod download

lint:
	@golangci-lint run

coverage.out: ${GO_TEST_FILES} docker-down docker-up
	@go test -coverprofile=coverage.out ./...
	docker-compose down

show-coverage: coverage.out
	@go tool cover -html=coverage.out

.PHONY: default docker-up docker-down test deps lint show-coverage
