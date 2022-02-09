package db

import (
	"database/sql"
	"log"
	models "newwebapp/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var jb = &models.Job{
	Title:       "sre engineer",
	Description: "5 yr experience sre engineer wanted",
	Email:       "hr@handoz.com",
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Default().Fatalf("an error '%s' was not expected when opening a stubdatabase connection", err)

	}
	return db, mock
}

func TestGetAllJobs(t *testing.T) {
	db, mock := NewMock()
	job := &jobStore{db}
	defer func() {
		job.Close()
	}()

	query := "SELECT title, description , email FROM jobs"

	rows := sqlmock.NewRows([]string{"title", "description", "email"}).
		AddRow(jb.Title, jb.Description, jb.Email)

	mock.ExpectQuery(query).WillReturnRows(rows)

	jobs, _ := job.GetJobs()
	assert.NotEmpty(t, jobs)
	assert.Len(t, jobs, 1)
	//assert.Len(t, jobs, 1)
}
