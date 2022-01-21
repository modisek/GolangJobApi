package main

import (
	"fmt"
	"log"
	"net/http"
	"newwebapp/models"

	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to the job api"))
	})

	r.HandleFunc("/jobs", models.CreateJobListing).Methods("POST")
	r.HandleFunc("/jobs", models.GetJobs).Methods("GET")
	r.HandleFunc("/jobs/{id}", models.GetOneJob).Methods("GET")
	r.HandleFunc("/jobs/{id}", models.EditAJobListing).Methods("PUT")
	r.HandleFunc("/jobs/{id}", models.DeleteAJobListing).Methods("DELETE")

	return r
}

func main() {
	r := newRouter()

	fmt.Println("Server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
