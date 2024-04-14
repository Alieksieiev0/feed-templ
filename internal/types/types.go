package types

import "time"

type UserToken struct {
	User  User  `json:"user"`
	Token Token `json:"token"`
}

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Value string `json:"value"`
}

type Base struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type Post struct {
	Base
	Title     string `json:"title"`
	Body      string `json:"body"`
	OwnerName string `json:"owner_name"`
	OwnerId   string `json:"owner_id"`
}

type ResponseError struct {
	Error string `json:"error"`
}
