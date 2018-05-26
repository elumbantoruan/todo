package repositories

import (
	"fmt"
	"time"

	"github.com/elumbantoruan/todo/models"
	"github.com/google/uuid"
)

// MockTodoRepository is a mock implementation for TodoRepository
type MockTodoRepository struct{}

var list []models.Todo
var keys = make(map[string]interface{})

// AddTodo adds new todo
func (m MockTodoRepository) AddTodo(todo models.Todo) error {
	if _, ok := keys[todo.ID.String()]; ok {
		// dups
		return fmt.Errorf("duplicate todoID")
	}
	keys[todo.ID.String()] = nil
	list = append(list, todo)
	return nil
}

// AddTask adds new task to existing todo
func (m MockTodoRepository) AddTask(todoID uuid.UUID, task models.Task) error {
	for i := 0; i < len(list); i++ {
		if list[i].ID == todoID {
			for j := 0; j < len(list[i].Tasks); j++ {
				if list[i].Tasks[j].ID == task.ID {
					return fmt.Errorf("duplicate taskId")
				}
			}
			list[i].Tasks = append(list[i].Tasks, task)
			break
		}
		break
	}
	return nil
}

// GetTodo return list of todo
func (m MockTodoRepository) GetTodo() ([]models.Todo, error) {
	return list, nil
}

// GetTodoByID return specific todo
func (m MockTodoRepository) GetTodoByID(todoID uuid.UUID) (*models.Todo, error) {
	for _, t := range list {
		if t.ID == todoID {
			return &t, nil
		}
	}
	return nil, nil
}

// UpdateTodo updates todo
func (m MockTodoRepository) UpdateTodo(todoID uuid.UUID, completed bool, dueDate *time.Time) error {
	for i := 0; i < len(list); i++ {
		if list[i].ID == todoID {
			list[i].Completed = completed
			list[i].DueDate = dueDate
		}
		break
	}
	return nil
}

// UpdateTask updates task
func (m MockTodoRepository) UpdateTask(todoID uuid.UUID, taskID uuid.UUID, completed bool) error {
	for i := 0; i < len(list); i++ {
		if list[i].ID == todoID {
			for j := 0; j < len(list[i].Tasks); j++ {
				if list[i].Tasks[j].ID == taskID {
					list[i].Tasks[j].Completed = completed
					break
				}
			}
		}
		break
	}
	return nil
}

// DeleteTask deletes task
func (m MockTodoRepository) DeleteTask(todoID uuid.UUID, taskID uuid.UUID) error {
	for i := 0; i < len(list); i++ {
		if list[i].ID == todoID {
			for j := 0; j < len(list[i].Tasks); j++ {
				list[i].Tasks = append(list[i].Tasks[:j], list[i].Tasks[j+1:]...)
				break
			}
			break
		}
	}
	return nil
}

// DeleteTodo deletes todo
func (m MockTodoRepository) DeleteTodo(todoID uuid.UUID) error {
	for i := 0; i < len(list); i++ {
		if list[i].ID == todoID {
			list = append(list[:i], list[i+1:]...)
			break
		}
	}
	return nil
}

// Clear clears out the slice
func (m MockTodoRepository) Clear() {
	list = nil
	for k := range keys {
		delete(keys, k)
	}
}
