package models

type User struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type UserInterface interface {
	GetUser(password, username string) (*User, error)
	CreateUser(password, username string) error
}
