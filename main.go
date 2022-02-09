package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	database "newwebapp/db"
	"newwebapp/models"

	"github.com/gorilla/mux"
)

type httpServerHelper struct {
	cancelFunc context.CancelFunc
}

func newRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to the job api"))
	})
	r.HandleFunc("/jobs", createJob).Methods("POST")
	r.HandleFunc("/jobs", listJobs).Methods("GET")
	r.HandleFunc("/jobs/{id}", getJobById).Methods("GET")
	r.HandleFunc("/jobs/{id}", editJob).Methods("PUT")
	r.HandleFunc("/jobs/{id}", deleteJob).Methods("DELETE")

	return r
}

func main() {

	r := newRouter()

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe():%v ", err)
		}
	}()
	log.Printf("Server started on %s", srv.Addr)

	<-done
	log.Println("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %+v", err)
	}
	log.Println("Server exited gracefully")

}

var dbString = "postgresql://gorm:password@localhost:5432/jobstore?sslmode=disable"
var db, _ = database.NewJobStore("pgx", dbString)

func listJobs(w http.ResponseWriter, r *http.Request) {

	jobs, err := db.GetJobs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jobbytes, err := json.Marshal(jobs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jobbytes)

	// for _, i := range jobs {
	// 	fmt.Fprintf(w, "%s %s %s ", i.Title, i.Description, i.Email)
	// }
}

func getJobById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	validId, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
	}
	job, err := db.GetOneJob(validId)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jobbytes, err := json.Marshal(job)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jobbytes)

}

func createJob(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	description := r.FormValue("description")
	email := r.FormValue("email")

	if title == "" || description == "" || email == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	newJob := models.Job{
		Title:       title,
		Description: description,
		Email:       email,
	}

	if err := db.CreateJobListing(&newJob); err != nil {

		w.WriteHeader(http.StatusInternalServerError)

	}
	log.Println([]byte("inserted one row"))
}

func deleteJob(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	validId, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
	}
	if err := db.DeleteAJobListing(validId); err != nil {

		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write([]byte("deleted one row"))

}
func editJob(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	description := r.FormValue("description")
	email := r.FormValue("email")

	vars := mux.Vars(r)
	id := vars["id"]

	validId, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		return
	}

	if title == "" || description == "" || email == "" || id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	newJob := models.Job{
		Title:       title,
		Description: description,
		Email:       email,
	}
	if err := db.EditAJobListing(&newJob, validId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write([]byte("edited "))
}
