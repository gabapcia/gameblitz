test:			## Run unit tests
	@go test ./...

api-docs:		## Generate API Docs
	./scripts/api-docs

db-queries:		## Generate SQLC code.
	./scripts/pg-management --gen-queries

db-migration:	## Run database migrations. 
	./scripts/pg-management --create-migration

db-migrate:		## Apply the migrations to the database
	./scripts/pg-management --apply-migration
