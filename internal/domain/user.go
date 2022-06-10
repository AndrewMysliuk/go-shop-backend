package domain

import "time"

type User struct {
	Id        string    `json:"id"`
	Name      string    `json:"name" binding:"required,min=1"`
	Sername   string    `json:"sername" binding:"required,min=1"`
	Email     string    `json:"email" binding:"required,email"`
	Phone     string    `json:"phone" binding:"required"`
	Password  string    `json:"password" binding:"required,min=5"`
	CreatedAt time.Time `json:"created_at"`
}

type UserSignIn struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
