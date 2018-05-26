package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/elumbantoruan/todo/models"

	"github.com/elumbantoruan/todo/repositories"
)

// TodoHandler handles Todo API operations
type TodoHandler struct {
	repo repositories.TodoRepository
}

// NewTodoHandler creates an instance of TodoHandler
func NewTodoHandler(repo repositories.TodoRepository) *TodoHandler {
	return &TodoHandler{
		repo: repo,
	}
}

// HandleAddTodo handles http POST action to add todo
func (t *TodoHandler) HandleAddTodo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var todo models.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// if empty guid
	if todo.ID == uuid.Nil {
		todo.ID = uuid.New()
	}
	for i := 0; i < len(todo.Tasks); i++ {
		if todo.Tasks[i].ID == uuid.Nil {
			todo.Tasks[i].ID = uuid.New()
		}
	}
	err = t.repo.AddTodo(todo)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate todoID") {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// HandleAddTask handles http POST action to add task
func (t *TodoHandler) HandleAddTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	if _, ok := vars["id"]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var task models.Task
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}
	if task.ID == uuid.Nil {
		task.ID = uuid.New()
	}
	err = t.repo.AddTask(id, task)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate taskId") {
			w.WriteHeader(http.StatusConflict) // 409
		} else {
			w.WriteHeader(http.StatusInternalServerError) // 500
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// HandleUpdateTask handles http PUT action
func (t *TodoHandler) HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	if _, ok := vars["id"]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, ok := vars["taskId"]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	taskID, err := uuid.Parse(vars["taskId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var tc models.CompletedTask
	err = json.NewDecoder(r.Body).Decode(&tc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = t.repo.UpdateTask(id, taskID, tc.Completed)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// HandleGetTodoList handles http GET action
func (t *TodoHandler) HandleGetTodoList(w http.ResponseWriter, r *http.Request) {

	var (
		todoList         []models.Todo
		filteredTodoList []models.Todo
		search           string
		skip             int
		limit            int
		err              error
	)
	todoList, err = t.repo.GetTodo()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(todoList) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	vars := r.URL.Query()
	if _, ok := vars["search"]; ok {
		search = vars["search"][0]
	}
	if _, ok := vars["skip"]; ok {
		skip, err = strconv.Atoi(vars["skip"][0])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if _, ok := vars["limit"]; ok {
		limit, err = strconv.Atoi(vars["limit"][0])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if len(search) > 0 {
		for _, t := range todoList {
			if strings.Contains(strings.ToLower(t.Name), strings.ToLower(search)) {
				filteredTodoList = append(filteredTodoList, t)
			}
		}
	} else {
		filteredTodoList = todoList
	}
	if skip > 0 {
		if len(filteredTodoList) > skip {
			filteredTodoList = filteredTodoList[skip:]
		}
	}
	if limit > 0 {
		if len(filteredTodoList) > limit {
			filteredTodoList = filteredTodoList[:limit]
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(filteredTodoList)
}

// HandleGetTodoByID handles http GET action for specific ToDoID
func (t *TodoHandler) HandleGetTodoByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if _, ok := vars["id"]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	todo, err := t.repo.GetTodoByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(todo)
}

// HandleDeleteTask handles http DELETE action for specific TaskID
func (t *TodoHandler) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if _, ok := vars["id"]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, ok := vars["taskId"]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	taskID, err := uuid.Parse(vars["taskId"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = t.repo.DeleteTask(id, taskID)
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// HandleUpdateTodo handles http PUT action for specific ToDoID
func (t *TodoHandler) HandleUpdateTodo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	if _, ok := vars["id"]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var ut models.UpdatedTodo
	err = json.NewDecoder(r.Body).Decode(&ut)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = t.repo.UpdateTodo(id, ut.Completed, ut.DueDate)
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// HandleDeleteTodo handles http DELETE action for specific ToDoID
func (t *TodoHandler) HandleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if _, ok := vars["id"]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = t.repo.DeleteTodo(id)
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
