package models

import (
	"github.com/google/uuid"
)

// Task represent an activity to be completed
type Task struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
}
