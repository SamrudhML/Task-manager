package models

import "time"

type Task struct {
	ID          string    `json:"id" bson:"id"`
	UserID      string    `json:"-" bson:"user_id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description,omitempty" bson:"description,omitempty"`
	Status      string    `json:"status" bson:"status"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
}
