package models

import (
	"time"

	"github.com/google/uuid"
)

// Todo defines tasks need to be done
type Todo struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	DueDate     *time.Time `json:"dueDate"`
	Tasks       []Task     `json:"tasks"`
}
