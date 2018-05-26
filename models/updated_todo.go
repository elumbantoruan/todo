package models

import "time"

// UpdatedTodo .
type UpdatedTodo struct {
	Completed bool       `json:"completed"`
	DueDate   *time.Time `json:"dueDate"`
}
