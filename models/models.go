package models

import "time"

type JobInterface interface {
	Close()
	GetJobs() ([]*Job, error)
	GetOneJob(id int) (*Job, error)
	CreateJobListing(job *Job) error
	EditAJobListing(job *Job, id int) error
	DeleteAJobListing(id int) error
}

type Job struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Email       string    `json:"email"`
	Type        string    `json:"type"`     //onsite, remote or hybrid
	Category    string    `json:"category"` // industry the job is in
	Location    string    `json:"location"`
	Expires     time.Time `json:"expires"`
	Created_at  string    `json:"created_at"`
}
