-include .env
export

.PHONY: docs

docs:
	swag init -g ./cmd/api/main.go -o ./docs

run:
	go run ./cmd/api

seed:
	go run ./cmd/seed

lint:
	golangci-lint run --timeout=120s

migrate-up:
	migrate -path ./migrations -database $(DB_URL) up

migrate-down:
	migrate -path ./migrations -database $(DB_URL) down 1

migrate-down-all:
	migrate -path ./migrations -database $(DB_URL) down -all

migrate-drop:
	migrate -path ./migrations -database $(DB_URL) drop -f
