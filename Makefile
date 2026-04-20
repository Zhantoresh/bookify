.PHONY: help build run test test-cover test-unit test-integration migrate-up migrate-down docker-up docker-down docker-logs docker-rebuild swagger lint fmt vet mod-tidy clean all

help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application locally"
	@echo "  make test         - Run tests"
	@echo "  make test-cover   - Run tests with coverage"
	@echo "  make migrate-up   - Apply database migrations"
	@echo "  make migrate-down - Rollback database migrations"
	@echo "  make docker-up    - Start all services with Docker Compose"
	@echo "  make docker-down  - Stop all services"
	@echo "  make clean        - Clean build artifacts"

build:
	GOCACHE=/tmp/go-build go build -o bin/api ./cmd/api

run:
	GOCACHE=/tmp/go-build go run ./cmd/api

test:
	GOCACHE=/tmp/go-build go test -v ./...

test-cover:
	GOCACHE=/tmp/go-build go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-unit:
	GOCACHE=/tmp/go-build go test -v ./tests/unit/...

test-integration:
	GOCACHE=/tmp/go-build go test -v ./tests/integration/...

migrate-up:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f api

docker-rebuild:
	docker compose up -d --build

swagger:
	swag init -g cmd/api/main.go -o docs

lint:
	golangci-lint run ./...

fmt:
	gofmt -w cmd internal pkg tests

vet:
	GOCACHE=/tmp/go-build go vet ./...

mod-tidy:
	GOCACHE=/tmp/go-build go mod tidy

clean:
	rm -rf bin
	rm -f coverage.out coverage.html

all: fmt vet test build
