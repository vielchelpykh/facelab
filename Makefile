include .env
export

export PROJECT_ROOT = $(shell pwd)

postgres-up:
	@docker compose up -d postgres

postgres-down:
	@docker compose down postgres

port-forwarder-up:
	@docker compose up -d port-forwarder

port-forwarder-down:
	@docker compose down port-forwarder

migrate-create:
	if [ -z "$(seq)" ]; then \
		echo "You must write the parametre 'seq'. For example: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm migrate \
	create \
	-ext sql \
	-dir /migrations \
	-seq "$(seq)"

migrate-action:
	if [ -z "$(action)" ]; then \
		echo "You must write the parametre 'action'. For example: make migrate-create action=up"; \
		exit 1; \
	docker compose run --rm migrate \
	-path /migrations \
	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}/${POSTGRES_PORT}?sslmode=disable \
	"$(action)"

migrate-up:
	make migrate-action action=up

migrate-down:
	make migrate-action action=down

go-run:
	export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/app/main.go

fix-perms:
	sudo chown -R $(shell whoami):$(shell whoami) ${PROJECT_ROOT}/out/
	sudo chmod -R 755 ${PROJECT_ROOT}/out/