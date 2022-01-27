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

var dbString = "postgresql://gorm:password@localhost:5432/jobstore?sslmode=disable"

func Db() *sql.DB {
	db, err := sql.Open("pgx", dbString)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Printf("unable to reach database: %v", err)
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
	// Type        string `json:"type"`             //onsite, remote or hybrid
	// Category    string `json:"category"`         // industry the job is in
	// Location    string `json:"location"`
	//    Expires     time.Time `json:"expires"`
	Created_at string `json:"created_at"`
}

//return all jobs
func GetJobs(w http.ResponseWriter, _ *http.Request) {
	rows, err := db.Query("SELECT * FROM jobs limit 10")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	jobs := []Job{}
	for rows.Next() {
		job := Job{}
		if err := rows.Scan(&job.Id, &job.Title, &job.Description, &job.Email, &job.Created_at); err != nil {
			log.Printf("could not scan row: %v", err)
			http.Error(w, err.Error(), 500)
			return
		}

		jobs = append(jobs, job)
	}
	jobbytes, err := json.Marshal(jobs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jobbytes)

}

//return a job by its id
func GetOneJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, http.StatusText(400), 400)
	}

	row := db.QueryRow("SELECT * FROM jobs WHERE id = $1", id)

	jb := new(Job)
	err := row.Scan(&jb.Id, &jb.Title, &jb.Description, &jb.Email, &jb.Created_at)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, "%s, %s, %s\n", jb.Title, jb.Description, jb.Email)
}

//create a job listing
func CreateJobListing(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	description := r.FormValue("description")
	email := r.FormValue("email")

	if title == "" || description == "" || email == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	newJob := Job{
		Title:       title,
		Description: description,
		Email:       email,
	}
	result, err := db.Exec("INSERT INTO jobs (title, description, email) VALUES($1, $2, $3)", newJob.Title, newJob.Description, newJob.Email)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Println("inserted", rowsAffected, "rows")
}

//delete a job given its id
func DeleteAJobListing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, err := db.Exec("DELETE FROM jobs WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("could not get affected rows: %v", err)
		return
	}

	fmt.Println("deleted", rowsAffected, "row")
}

//edit a job by its id
func EditAJobListing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("edit job"))
}
