
migrateup:
	 migrate -path db/migration -database "postgresql://gorm:password@localhost:5432/job_portal?sslmode=disable" -verbose up
migratedown:
	 migrate -path db/migration -database "postgresql://gorm:password@localhost:5432/job_portal?sslmode=disable" -verbose down

