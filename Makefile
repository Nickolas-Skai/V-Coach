## Filename Makefile
include .envrc

.PHONY: run/tests
run/tests: vet
	go test -v ./...

.PHONY: fmt
fmt: 
	go fmt ./...

.PHONY: vet
vet: fmt
	go vet ./...

.PHONY: run
run: vet
	go run ./cmd/web -addr=${ADDRESS} -dsn=${VCOACH_DB_DSN}

.PHONY: start
start:
	@echo "Starting sql server..."
	sudo service postgresql start

.PHONY: db/psql
db/psql:
	sudo -u vcoach psql --host=localhost --dbname=vcoach --username=vcoach

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database "postgres://vcoach:seo@localhost:5432/vcoach?sslmode=disable" up

# db/migrations/down-1: undo the last migration
.PHONY: db/migrations/down-1
db/migrations/down-1:
	@echo 'Undoing the last migration...'
	migrate -path ./migrations -database "postgres://vcoach:seo@localhost:5432/vcoach?sslmode=disable" down 1

## db/migrations/fix: fix a SQL migration
.PHONY: db/migrations/fix
db/migrations/fix:
	@echo 'Checking migration status...'
	@migrate -path ./migrations -database "postgres://vcoach:seo@localhost:5432/vcoach?sslmode=disable" version > /tmp/migrate_version 2>&1
	@cat /tmp/migrate_version
	@if grep -q "dirty" /tmp/migrate_version; then \
		version=$$(grep -o '[0-9]\+' /tmp/migrate_version | head -1); \
		echo "Found dirty migration at version $$version"; \
		echo "Forcing version $$version..."; \
		migrate -path ./migrations -database "postgres://vcoach:seo@localhost:5432/vcoach?sslmode=disable" force $$version; \
		echo "Running down migration..."; \
		migrate -path ./migrations -database "postgres://vcoach:seo@localhost:5432/vcoach?sslmode=disable" down 1; \
		echo "Running up migration..."; \
		migrate -path ./migrations -database "postgres://vcoach:seo@localhost:5432/vcoach?sslmode=disable" up; \
	else \
		echo "No dirty migration found"; \
	fi
	@rm -f /tmp/migrate_version

## db/migrations/force: force a migration to a specific version
	.PHONY: db/reset
db/reset:
	@echo "Dropping all tables and resetting migrations..."
	psql --host=localhost --dbname=vcoach --username=vcoach -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO public;"
	migrate -path ./migrations -database "postgres://vcoach:seo@localhost:5432/vcoach?sslmode=disable" force 0
	migrate -path ./migrations -database "postgres://vcoach:seo@localhost:5432/vcoach?sslmode=disable" up
