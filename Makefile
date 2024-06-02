ENV_FILE := .env
include $(ENV_FILE)
export

MIGRATION_UP_PATH=./migrations/up
MIGRATION_DOWN_PATH=./migrations/down
SEED_SCRIPT=./seeds/seed_data.sql
DATABASE_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: migrate-up migrate-down

build:
	docker-compose up --build  

run: 
	go run main.go   

migrate-up:
	migrate -database "$(DATABASE_URL)" -path $(MIGRATION_UP_PATH) up

migrate-down:
	migrate -database "$(DATABASE_URL)" -path $(MIGRATION_DOWN_PATH) down

seed:
	psql $(DATABASE_URL) -f $(SEED_SCRIPT)

tidy:
	go mod tidy

test:
	go test ./...

mockgen:
	mockgen -source=repository/contribution_repository.go -destination=mocks/mock_contribution_repository.go -package=mocks /
	mockgen -source=repository/db_transaction_repository.go -destination=mocks/mock_db_transaction_repository.go -package=mocks
	mockgen -source=repository/referral_link_repository.go -destination=mocks/mock_referral_link_repository.go -package=mocks
	mockgen -source=repository/role_repository.go -destination=mocks/mock_role_repository.go -package=mocks /
	mockgen -source=repository/user_repository.go -destination=mocks/mock_user_repository.go -package=mocks
