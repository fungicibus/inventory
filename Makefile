.PHONY: generate-api-v1
generate-api-v1:
	oapi-codegen -package="v1" -generate types -o internal/api/v1/openapi_types.gen.go modules/docs/api/inventory/v1/openapi.yaml
	oapi-codegen -package="v1" -generate chi-server -o internal/api/v1/openapi_server.gen.go modules/docs/api/inventory/v1/openapi.yaml
	oapi-codegen -package="v1" -generate spec -o internal/api/v1/specs.gen.go modules/docs/api/inventory/v1/openapi.yaml
	go mod tidy

include .env
DB_URL ?= $(PG_RW_DSN)
MIGRATION_PATH ?= ./cmd/migrations


.PHONY: migrations-status migrations-up migrations-down migrations-reset

migrations-status:
	goose -dir $(MIGRATION_PATH) postgres "$(DB_URL)" status

migrations-up:
	goose -dir $(MIGRATION_PATH) postgres "$(DB_URL)" up

migrations-down:
	goose -dir $(MIGRATION_PATH) postgres "$(DB_URL)" down

migrations-reset:
	goose -dir $(MIGRATION_PATH) postgres "$(DB_URL)" reset


.PHONY: compose-up compose-down

compose-up:
	docker compose -f docker/docker-compose.yml -p inventory up 

compose-down:
	docker compose -f docker/docker-compose.yml -p inventory down 
