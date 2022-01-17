migrateup:
	migrate -path db/migration -database "postgresql://postgres:admin@localhost:5432/test?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://postgres@admin:5432/test?sslmode=disable" -verbose down