.PHONY: dev test test-cover build lint migrate-up

dev:
	go run ./cmd/api/...

test:
	go test ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

build:
	go build -o bin/server ./cmd/api/...

lint:
	golangci-lint run ./...

# Migrations live in migrations/ in this repo (not a separate lms-iam-db) —
# see README's "Migration scope" for why.
migrate-up:
	docker compose -f deploy/compose.yml up
