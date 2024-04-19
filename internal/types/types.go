package types

import (
	"time"
)

type UserToken struct {
	User  User  `json:"user"`
	Token Token `json:"token"`
}

type UserBase struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type User struct {
	UserBase
	Password    string     `json:"password"`
	Posts       []Post     `json:"posts"`
	Subscribers []UserBase `json:"subscribers"`
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
