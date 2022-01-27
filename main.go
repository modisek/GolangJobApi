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

	staticFileDir := http.Dir("./assets")
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDir))
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")
	return r
}

func main() {
	port := ":8000"
	r := newRouter()

	fmt.Printf("Server started on http://localhost%v", port)
	log.Fatal(http.ListenAndServe(port, r))
}
