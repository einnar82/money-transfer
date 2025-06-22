ifneq (,$(wildcard .env))
    include .env
    export
endif

DATABASE_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

migrate:
	docker compose run --rm migrate -path=/migrations -database="$(DATABASE_URL)" up

migrate-rollback:
	docker compose run --rm migrate -path=/migrations -database="$(DATABASE_URL)" down

migrate-force:
	docker compose run --rm migrate -path=/migrations -database="$(DATABASE_URL)" force $(version)

migrate-status:
	docker compose run --rm migrate -path=/migrations -database="$(DATABASE_URL)" version

migrations:
	@read -p "Migration name: " name; \
	docker run --rm -v $$PWD/migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq $$name

run-functional-tests:
	@echo "🧹 Cleaning up..."
	docker compose down -v
	@echo "🚀 Starting test_db..."
	docker compose up -d test_db
	@sleep 3

	@echo "🧨 Dropping existing schema..."
	docker compose run --rm migrate \
	  -path=/migrations \
	  -database "postgres://postgres:postgres@test_db:5432/testdb?sslmode=disable" \
	  drop -f

	@echo "🛠  Running migrations..."
	docker compose run --rm migrate \
	  -path=/migrations \
	  -database "postgres://postgres:postgres@test_db:5432/testdb?sslmode=disable" \
	  up

	@echo "✅ Running tests..."
	go test ./tests -v
