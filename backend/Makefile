.PHONY: help run test build migrate-up migrate-down tidy

help:
	@echo "Targets:"
	@echo "  make run          - Levanta el server en :8080"
	@echo "  make test         - Corre tests"
	@echo "  make build        - Compila binario en ./bin/api"
	@echo "  make tidy         - go mod tidy"
	@echo "  make migrate-up   - Aplica migraciones (preparado, sin migrations aún)"

run:
	@if [ ! -f .env ]; then cp .env.example .env && echo "Creado .env desde .env.example"; fi
	go run ./cmd/api

test:
	go test ./... -v

build:
	go build -o ./bin/api ./cmd/api

tidy:
	go mod tidy

migrate-up:
	@echo "Sin migraciones aún (se agregan en Fase 1)"

migrate-down:
	@echo "Sin migraciones aún"
