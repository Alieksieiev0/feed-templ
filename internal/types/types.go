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

type Post struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	OwnerName string    `json:"owner_name"`
	OwnerId   string    `json:"owner_id"`
}

type ResponseError struct {
	Error string `json:"error"`
}

type Notification struct {
	Id       string `json:"id"`
	NotifyId string `json:"notify_id"`
	FromId   string `json:"from_id"`
	FromName string `json:"from_name"`
	TargetId string `json:"target_id"`
	Type     string `json:"type"`
	Status   string `json:"status"`
}
