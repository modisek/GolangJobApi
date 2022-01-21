package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func Db() *sql.DB {
	db, err := sql.Open("pgx", "postgresql://gorm:password@localhost:5432/jobstore?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("unable to reach database: %v", err)
	}
	fmt.Println("database is reachable")
	// Maximum Idle Connections
	db.SetMaxIdleConns(5)
	// Maximum Open Connections
	db.SetMaxOpenConns(10)
	// Idle Connection Timeout
	db.SetConnMaxIdleTime(1 * time.Second)
	// Connection Lifetime
	db.SetConnMaxLifetime(30 * time.Second)
	return db
}

var db = Db()

type Job struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Email       string `json:"email"`
	Created_at  string `json:"created_at"`
}

func GetJobs(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM jobs limit 10")
	if err != nil {
		log.Fatalf("could not execute query: %v", err)
	}

	jobs := []Job{}

	for rows.Next() {
		job := Job{}

		if err := rows.Scan(&job.Id, &job.Title, &job.Description, &job.Email, &job.Created_at); err != nil {
			log.Fatalf("could not scan row: %v", err)
		}
		jobs = append(jobs, job)
	}
	jobbytes, err := json.Marshal(jobs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jobbytes)
	// for _, jbs := range jobs {
	// 	fmt.Fprintf(w, "%s, %s, %s\n ", jbs.Title, jbs.Description, jbs.Email)
	// }
}

func GetOneJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	row := db.QueryRow("SELECT * FROM jobs WHERE id = $1", id)

	jb := new(Job)
	err := row.Scan(&jb.Id, &jb.Title, &jb.Description, &jb.Email, &jb.Created_at)
	if err != nil {
		log.Fatalf("error %v", err)
	}

	fmt.Fprintf(w, "%s, %s, %s\n", jb.Title, jb.Description, jb.Email)
}

func CreateJobListing(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	description := r.FormValue("description")
	email := r.FormValue("email")

	newJob := Job{
		Title:       title,
		Description: description,
		Email:       email,
	}
	result, err := db.Exec("INSERT INTO jobs (title, description, email) VALUES($1, $2, $3)", newJob.Title, newJob.Description, newJob.Email)
	if err != nil {
		log.Fatalf("could not insert row: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("could not get affected rows: %v", err)
	}
	fmt.Println("inserted", rowsAffected, "rows")
}

func DeleteAJobListing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, err := db.Exec("DELETE FROM jobs WHERE id = $1", id)
	if err != nil {
		log.Fatalf("could not delete row: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("could not get affected rows: %v", err)
	}

	fmt.Println("deleted", rowsAffected, "row")
}

func EditAJobListing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("edit job"))
}
