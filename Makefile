
migrateup:
	 migrate -path db/migration -database "postgresql://gorm:password@localhost:5432/jobstore?sslmode=disable" -verbose up
migratedown:
	 migrate -path db/migration -database "postgresql://gorm:password@localhost:5432/jobstore?sslmode=disable" -verbose down

