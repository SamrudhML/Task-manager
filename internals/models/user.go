package models

import "time"

type User struct {
	ID           string    `json:"id" bson:"id"`
	Username     string    `json:"username" bson:"username"`
	PasswordHash string    `json:"-" bson:"password_hash"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}
