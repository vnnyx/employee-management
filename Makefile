DB_USER ?= postgres
DB_PASS ?= postgres
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_NAME ?= employee_management
DB_SSLMODE ?= disable
DB_URL = postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

migrate-up:
	migrate -database "$(DB_URL)" -path database/migrations up

migrate-down:
	migrate -database "$(DB_URL)" -path database/migrations down

migrate-create:
	migrate create -ext sql -dir database/migrations -seq $(name)

mock-repo:
	mockgen -source internal/$(domain)/repository.go -destination internal/$(domain)/mock/repository_mock.go -package=mocks -typed

mock-usecase:
	mockgen -source internal/$(domain)/usecase.go -destination internal/$(domain)/mock/usecase_mock.go -package=mocks -typed