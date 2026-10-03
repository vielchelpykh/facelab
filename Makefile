include .env
export

export PROJECT_ROOT := $(shell pwd)

# ============================================================
# DOCKER COMPOSE
# ============================================================

up:
	@docker compose up -d --build

down:
	@docker compose down

restart:
	@docker compose down
	@docker compose up -d --build

ps:
	@docker compose ps

logs:
	@docker compose logs -f

logs-go:
	@docker compose logs -f go-backend

logs-python:
	@docker compose logs -f python-service

logs-postgres:
	@docker compose logs -f postgres

# ============================================================
# POSTGRES
# ============================================================

postgres-up:
	@docker compose up -d postgres

postgres-down:
	@docker compose down postgres

port-forwarder-up:
	@docker compose up -d port-forwarder

port-forwarder-down:
	@docker compose down port-forwarder

# ============================================================
# МИГРАЦИИ
# ============================================================

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Пример: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Пример: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm migrate \
		-path /migrations \
		-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable" \
		"$(action)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-version:
	@docker compose run --rm migrate \
		-path /migrations \
		-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable" \
		version