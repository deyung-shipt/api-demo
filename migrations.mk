MIGRATIONS_DIR = ./migrations
PSQL_URL ?= postgresql://pguser:pgpass@localhost:5432/tempest-template?sslmode=disable

.PHONY: migrate migration

# Migrate the database up.
migrate:
	migrate -path "$(MIGRATIONS_DIR)" -database "$(PSQL_URL)" up

# Create a new migration.
#
# Migration name is set via the NAME variable.
migration:
	@[ "$(NAME)" ] || (echo "NAME must be given"; exit 1)
	migrate create -dir "$(MIGRATIONS_DIR)" -ext sql "$(NAME)"
