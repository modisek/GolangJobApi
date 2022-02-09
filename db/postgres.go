package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	models "newwebapp/models"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type jobStore struct {
	db *sql.DB
}

//var dbString = "postgresql://gorm:password@localhost:5432/jobstore?sslmode=disable"

func NewJobStore(dialect, dsn string) (*jobStore, error) {
	db, err := sql.Open(dialect, dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Printf("unable to reach database: %v", err)
	}
	log.Println("database is reachable")
	// Maximum Idle Connections
	db.SetMaxIdleConns(5)
	// Maximum Open Connections
	db.SetMaxOpenConns(10)
	// Idle Connection Timeout
	db.SetConnMaxIdleTime(1 * time.Second)
	// Connection Lifetime
	db.SetConnMaxLifetime(30 * time.Second)
	return &jobStore{db}, nil
}

func (st *jobStore) Close() {
	st.db.Close()
}

//return all jobs
func (st *jobStore) GetJobs() ([]*models.Job, error) {

	jobs := make([]*models.Job, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := st.db.QueryContext(ctx, "SELECT * FROM jobs limit 10")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		job := new(models.Job)
		if err := rows.Scan(&job.Id, &job.Title, &job.Description, &job.Email, &job.Created_at); err != nil {

			return nil, err
		}

		jobs = append(jobs, job)
	}

	return jobs, nil

}

// //return a job by its id
func (st *jobStore) GetOneJob(id int) (*models.Job, error) {
	job := new(models.Job)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := st.db.QueryRowContext(ctx, "SELECT * FROM jobs WHERE id = $1", id)

	err := row.Scan(&job.Id, &job.Title, &job.Description, &job.Email, &job.Created_at)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}
	return job, nil

}

//
// //create a job listing
func (st *jobStore) CreateJobListing(job *models.Job) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := st.db.ExecContext(ctx, "INSERT INTO jobs (title, description, email) VALUES($1, $2, $3)", job.Title, job.Description, job.Email)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return err
	}
	return err

}

//
// //delete a job given its id
func (st *jobStore) DeleteAJobListing(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := st.db.ExecContext(ctx, "DELETE FROM jobs WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return err
	}

	return err
}

//edit a job by its id
func (st *jobStore) EditAJobListing(job *models.Job, id int) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := st.db.ExecContext(ctx, "UPDATE jobs SET title=$1, description=$2, email=$3 WHERE id = $4 ", job.Title, job.Description, job.Email, id)
	if err != nil {
		return err
	}
	return err
}
