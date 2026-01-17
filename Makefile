# Makefile

.PHONY: run test migrate-up migrate-down docker-up docker-down

run:
	go run cmd/api/main.go

test:
	go test -v -race ./...

test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

migrate-up:
	migrate -path migrations -database "postgresql://user:pass@localhost:5432/invoices?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgresql://user:pass@localhost:5432/invoices?sslmode=disable" down

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

lint:
	golangci-lint run

build:
	go build -o bin/api cmd/api/main.go