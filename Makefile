include .test.env
export

.PHONY: default
default: test-integration

.PHONY: docker-up
docker-up:
	docker-compose up -d

.PHONY: docker-down
docker-down:
	docker-compose down

.PHOHY: test-integration
test: docker-down docker-up
	go test -v -race -cover -count=1 ./...
	docker-compose down
