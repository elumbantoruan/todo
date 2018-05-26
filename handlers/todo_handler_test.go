package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/elumbantoruan/todo/models"
	"github.com/elumbantoruan/todo/repositories"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestTodoHandler_HandleAddTodo(t *testing.T) {
	mockRepo := repositories.MockTodoRepository{}
	mockRepo.Clear()

	url := "/v1/todo"
	todo := newTodo()
	bytes, _ := json.Marshal(todo)
	request, _ := http.NewRequest("POST", url, strings.NewReader(string(bytes)))
	responseRecorder := httptest.NewRecorder()

	h := NewTodoHandler(mockRepo)
	h.HandleAddTodo(responseRecorder, request)

	assert.Equal(t, http.StatusCreated, responseRecorder.Code)
}

func TestTodoHandler_HandleAddTodo_Duplicate(t *testing.T) {
	mockRepo := repositories.MockTodoRepository{}
	mockRepo.Clear()

	todoID := uuid.New()
	newTodo := newTodoID(todoID)
	mockRepo.AddTodo(newTodo)

	url := "/v1/todo"
	// todo2 uses the same todoID as previous todo instance
	// so this will test for http.StatusConflict
	todo2 := newTodoID(todoID)
	bytes, _ := json.Marshal(todo2)
	request, _ := http.NewRequest("POST", url, strings.NewReader(string(bytes)))
	responseRecorder := httptest.NewRecorder()

	h := NewTodoHandler(mockRepo)
	h.HandleAddTodo(responseRecorder, request)

	assert.Equal(t, http.StatusConflict, responseRecorder.Code)
}

func TestTodoHandler_HandleAddTask(t *testing.T) {
	mockRepo := repositories.MockTodoRepository{}
	mockRepo.Clear()

	todoID := uuid.New()
	todo := newTodoID(todoID)
	mockRepo.AddTodo(todo)

	val, _ := mockRepo.GetTodoByID(todoID)
	assert.Equal(t, 1, len(val.Tasks))

	task := newTask()

	url := fmt.Sprintf("/v1/todo/%s/tasks", todoID)
	bytes, _ := json.Marshal(task)
	request, _ := http.NewRequest("POST", url, strings.NewReader(string(bytes)))
	params := map[string]string{
		"id": todoID.String(),
	}
	request = mux.SetURLVars(request, params)
	responseRecorder := httptest.NewRecorder()

	h := NewTodoHandler(mockRepo)
	h.HandleAddTask(responseRecorder, request)

	assert.Equal(t, http.StatusCreated, responseRecorder.Code)

	val, _ = mockRepo.GetTodoByID(todoID)
	assert.Equal(t, 2, len(val.Tasks))

}

func TestTodoHandler_HandleAddTask_Duplicate(t *testing.T) {
	mockRepo := repositories.MockTodoRepository{}
	mockRepo.Clear()

	todoID := uuid.New()
	todo := newTodoID(todoID)
	mockRepo.AddTodo(todo)

	taskID := uuid.New()
	task := newTaskID(taskID)
	mockRepo.AddTask(todoID, task)

	url := fmt.Sprintf("/v1/todo/%s/tasks", todoID)
	// set request payload to the same task to simulate StatusConfict 409
	bytes, _ := json.Marshal(task)
	request, _ := http.NewRequest("POST", url, strings.NewReader(string(bytes)))
	params := map[string]string{
		"id": todoID.String(),
	}
	request = mux.SetURLVars(request, params)
	responseRecorder := httptest.NewRecorder()

	h := NewTodoHandler(mockRepo)
	h.HandleAddTask(responseRecorder, request)

	assert.Equal(t, http.StatusConflict, responseRecorder.Code)

}

func TestTodoHandler_HandleGetTodoList(t *testing.T) {
	mockRepo := repositories.MockTodoRepository{}
	mockRepo.Clear()

	var listID []uuid.UUID
	for i := 0; i < 10; i++ {
		listID = append(listID, uuid.New())
	}

	for _, id := range listID {
		todo := newTodoID(id)
		mockRepo.AddTodo(todo)
	}

	skip := 2
	limit := 5

	url := fmt.Sprintf("/v1/todo?skip=%v&limit=%v", skip, limit)

	request, _ := http.NewRequest("GET", url, nil)
	responseRecorder := httptest.NewRecorder()

	h := NewTodoHandler(mockRepo)
	h.HandleGetTodoList(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	// given list of 10 records todo, skip = 2, and limit 5
	// so response payload contains 5 records, and the list of ID
	// should be equal to list[2] - list[6]
	var todoList []models.Todo
	json.NewDecoder(responseRecorder.Body).Decode(&todoList)

	assert.Equal(t, limit, len(todoList))
	for i := 0; i < limit; i++ {
		assert.Equal(t, listID[skip+i], todoList[i].ID)
	}
}

func newTodo() models.Todo {
	return newTodoID(uuid.New())
}

func newTodoID(todoID uuid.UUID) models.Todo {
	taskID := uuid.New()
	return models.Todo{
		ID:          todoID,
		Name:        fmt.Sprintf("TODO:%v", todoID),
		Description: fmt.Sprintf("DESCRIPTION:%v", todoID),
		Tasks: []models.Task{
			{
				ID:   taskID,
				Name: fmt.Sprintf("TASK%v", taskID),
			},
		},
	}
}

func newTask() models.Task {
	return newTaskID(uuid.New())
}

func newTaskID(taskID uuid.UUID) models.Task {
	return models.Task{
		ID:   taskID,
		Name: fmt.Sprintf("TASK:%v", taskID),
	}
}
