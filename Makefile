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
