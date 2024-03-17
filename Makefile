LOCAL_MIGRATION_DIR = ./migrations
LOCAL_MIGRATION_DSN = "host=localhost port=54321 dbname=order-info-service user=postgres password=qwerty sslmode=disable"

.PHONY: local-migration-status
local-migration-status:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

.PHONY: local-migration-up
local-migration-up:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

.PHONY: local-migration-down
local-migration-down:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

.PHONY: server-run
server-run:
	CONFIG_PATH="config/config.yml" PG_DSN='postgresql://postgres:qwerty@localhost:54321/order-info-service' go run cmd/order-info-service/main.go