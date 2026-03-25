-include .env
export

.PHONY: docs

docs:
	swag init -g ./cmd/api/main.go -o ./docs

run:
	go run ./cmd/api

seed:
	go run ./cmd/seed

migrate-up:
	migrate -path ./migrations -database $(DB_URL) up

migrate-down:
	migrate -path ./migrations -database $(DB_URL) down 1
