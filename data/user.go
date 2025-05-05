package data

import "time"

type User struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Created  time.Time `json:"created"`
}
