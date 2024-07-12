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
	go test -cover ./... -v

mockgen:
	mockgen -package=mocks -destination=mocks/mock_contribution_repository.go referral-system/repository IContributionRepository
	mockgen -package=mocks -destination=mocks/db_transaction_repository.go referral-system/repository IDBTransactionRepository
	mockgen -package=mocks -destination=mocks/referral_link_repository.go referral-system/repository IReferralLinkRepository
	mockgen -package=mocks -destination=mocks/role_repository.go referral-system/repository IRoleRepository
	mockgen -package=mocks -destination=mocks/user_repository.go referral-system/repository IUserRepository

clear-mock:
	rm -rf mocks

swag:
	swag init