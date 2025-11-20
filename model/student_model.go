package model

import "time"

// Student struct
type Student struct {
	ID         int64     `json:"id" example:"1"`
	Name       string    `json:"name" example:"John Doe" binding:"required"`
	Email      string    `json:"email" example:"john.doe@example.com" binding:"required,email"`
	Department string    `json:"department" example:"Computer Science"`
	CreatedAt  time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

// Students array of Student type
type Students []Student
