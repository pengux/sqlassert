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

lint:
	@golangci-lint run

.PHONY: default docker-up docker-down test lint
