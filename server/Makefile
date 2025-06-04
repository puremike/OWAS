include .env
export $(shell sed 's/=.*//' .env)

.PHONY: mup
mup:
	@echo "Running DB migration on: $(DB_ADDR)"
	migrate -database "$(DB_ADDR)" -path migrate/migrations up


.PHONY: mdown
mdown:
	@echo "Running DB migration on: $(DB_ADDR)"
	migrate -database "$(DB_ADDR)" -path migrate/migrations down

.PHONY: si
si:
	@echo "Generating Swagger docs..."
	swag init --dir cmd/api/ --parseDependency --parseInternal --parseDepth 3
	@echo "Swagger docs generated successfully."

.PHONY: sf
sf:
	@echo "Formatting Swagger docs..."
	swag fmt

.PHONY: dkup
dkup:
	@echo "Starting Docker containers..."
	docker compose up -d

.PHONY: dkdown
dkdown:
	@echo "Stopping Docker containers..."
	docker compose down