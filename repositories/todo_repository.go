package repositories

import (
	"time"

	"github.com/elumbantoruan/todo/models"
	"github.com/google/uuid"
)

// TodoRepository is an interface for repository
type TodoRepository interface {
	AddTodo(todo models.Todo) error
	AddTask(todoID uuid.UUID, task models.Task) error
	GetTodo() ([]models.Todo, error)
	GetTodoByID(todoID uuid.UUID) (*models.Todo, error)
	UpdateTodo(todoID uuid.UUID, completed bool, dueDate *time.Time) error
	UpdateTask(todoID uuid.UUID, taskID uuid.UUID, completed bool) error
	DeleteTask(todoID, taskID uuid.UUID) error
	DeleteTodo(todoID uuid.UUID) error
}
