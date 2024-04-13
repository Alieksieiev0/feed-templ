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
	Id        string
	CreatedAt time.Time
}

type Post struct {
	Base
	Title     string
	Body      string
	OwnerName string
	OwnerId   string
}

type ResponseError struct {
	Error string `json:"error"`
}
