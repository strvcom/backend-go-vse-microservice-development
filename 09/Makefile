MIGRATIONS_DIR := ./repository/sql/migrations

start-local-database: DATABASE_NAME=postgres
start-local-database:
	@docker run \
			-d \
			--name postgres \
			-p 5432:5432 \
			-e POSTGRES_USER=root \
			-e POSTGRES_PASSWORD=1234 \
			-e POSTGRES_DB=$(DATABASE_NAME) \
			postgres:15 \
			>/dev/null
	$(info Database $(DATABASE_NAME) started on port :5432)

stop-local-database:
	@docker stop postgres && docker rm postgres >/dev/null
	$(info Database container stopped and removed)

migrate-up: DATABASE_NAME=postgres
migrate-up:
	@migrate -path $(MIGRATIONS_DIR) -database 'postgres://root:1234@localhost:5432/$(DATABASE_NAME)?sslmode=disable' up

migrate-down: DATABASE_NAME=postgres
migrate-down:
	@migrate -path $(MIGRATIONS_DIR) -database 'postgres://root:1234@localhost:5432/$(DATABASE_NAME)?sslmode=disable' down

# Use this to force a migration to a specific version.
# This is useful when a migration fails and you need to rollback to a specific version.
migrate-force: DATABASE_NAME=postgres
migrate-force: MIGRATION_VERSION=0
migrate-force:
	@migrate -path $(MIGRATIONS_DIR) -database 'postgres://root:1234@localhost:5432/$(DATABASE_NAME)?sslmode=disable' force $(MIGRATION_VERSION)

migration: MIGRATION_NAME=unnamed
migration:
	@migrate create -seq -digits 3 -ext sql -dir $(MIGRATIONS_DIR) $(MIGRATION_NAME)
