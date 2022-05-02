package main

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	database "newwebapp/db"
	"newwebapp/models"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/crypto/bcrypt"

	_ "newwebapp/docs"
)

var sessionManager *scs.SessionManager
var userSess = uuid.New().String()

func newRouter() *mux.Router {

	r := mux.NewRouter()

	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	r.Use(sessionManager.LoadAndSave)

	apiv1 := r.PathPrefix("/api").Subrouter()

	apiv1.HandleFunc("/jobs", protected(createJob)).Methods("POST")
	apiv1.HandleFunc("/jobs", listJobs).Methods("GET")
	apiv1.HandleFunc("/jobs/{id}", getJobById).Methods("GET")
	apiv1.HandleFunc("/jobs/{id}", protected(editJob)).Methods("PUT")
	apiv1.HandleFunc("/jobs/{id}", protected(deleteJob)).Methods("DELETE")

	apiv1.HandleFunc("/signup", signup).Methods("POST")
	apiv1.HandleFunc("/login", login).Methods("POST")
	apiv1.HandleFunc("/logout", logout).Methods("GET")

	r.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

	return r
}

// @title User API documentation
// @version 1.0.0
// @host localhost:8000
// @BasePath /api/jobs
func main() {
	fmt.Println(userSess)

	gob.Register(models.User{})
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

// var dbString = os.Getenv("DB_CONN")
var db, _ = database.NewJobStore("pgx", dbString)

//listJobs godoc
// @Summary Get all jobs
// @Description get all jobs
// @Tags Users
// @Success 200 {array} model.Job
// @Failure 404 {object} object
// @Router / [get]

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

func signup(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), 5)

	user := models.User{
		Username: username,
		Password: string(hashedPass),
	}

	if err := db.CreateUser(&user); err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("created a user")

}
func login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	ok, err := db.GetUser(username)
	if err != nil {
		log.Println(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(ok.Password), []byte(password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	sessionManager.Put(r.Context(), userSess, ok)
	log.Println("logged in")

}

func logout(w http.ResponseWriter, r *http.Request) {
	if err := sessionManager.Destroy(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//createJob godoc
// @Summary Create new job  based on paramters
// @Description Create new job
// @Tags  Job
// @Accept json
// @Param job body model.Job true "Job Data"
// @Success 200 {object} object
// @Failure 400,500 {object} object
// @Router / [post]

func createJob(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	description := r.FormValue("description")
	email := r.FormValue("email")
	typeOfJob := r.FormValue("type")
	category := r.FormValue("category")
	location := r.FormValue("location")

	if title == "" || description == "" || email == "" || typeOfJob == "" || category == "" || location == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	newJob := models.Job{
		Title:       title,
		Description: description,
		Email:       email,
		Type:        typeOfJob,
		Category:    category,
		Location:    location,
	}

	if err := db.CreateJobListing(&newJob); err != nil {

		w.WriteHeader(http.StatusInternalServerError)

	}
	log.Println("inserted one row")
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
	typeOfJob := r.FormValue("type")
	category := r.FormValue("category")
	location := r.FormValue("location")

	vars := mux.Vars(r)
	id := vars["id"]
	validId, _ := strconv.Atoi(id)

	if title == "" || description == "" || email == "" || typeOfJob == "" || category == "" || location == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	newJob := models.Job{
		Title:       title,
		Description: description,
		Email:       email,
		Type:        typeOfJob,
		Category:    category,
		Location:    location,
	}

	if err := db.EditAJobListing(&newJob, validId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write([]byte("edited "))
}

func protected(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !sessionManager.Exists(r.Context(), userSess) {
			http.Error(w, "unauthorised please login", http.StatusUnauthorized)
			return
		}

		user := sessionManager.Get(r.Context(), userSess).(models.User)
		ctx := context.WithValue(r.Context(), userSess, user)
		h(w, r.WithContext(ctx))
	}
}
