include .env  # Load environment variables

.PHONY: migrate-up migrate-down  # Declare these as phony targets

migrate-up:
	@echo "Running database migrations..."
	migrate -path migrations -database "$(DATABASE_URL)" -verbose up

migrate-down:
	@echo "Rolling back database migrations..."
	migrate -path migrations -database "$(DATABASE_URL)" -verbose down
