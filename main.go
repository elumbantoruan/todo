package main

import (
	"log"
	"net/http"

	"github.com/elumbantoruan/todo/handlers"
	"github.com/elumbantoruan/todo/repositories"
	"github.com/gorilla/mux"
)

func main() {

	m, err := registerHandlers()
	if err != nil {
		log.Panic(err)
	}
	http.Handle("/", m)

	err = http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Panic(err)
	}
}

func registerHandlers() (*mux.Router, error) {
	m := mux.NewRouter()

	// creating an instance of filerepository
	fr := repositories.NewFileStorageTodoRepository("data")

	// instance of handlers which requires a storage
	handle := handlers.NewTodoHandler(fr)

	// register the http handler for each operations
	m.HandleFunc("/v1/todo", handle.HandleAddTodo).Methods("POST")
	m.HandleFunc("/v1/todo/{id}/tasks", handle.HandleAddTask).Methods("POST")
	m.HandleFunc("/v1/todo/{id}/task/{taskID}/complete", handle.HandleUpdateTask).Methods("PUT")
	m.HandleFunc("/v1/todo", handle.HandleGetTodoList).Methods("GET") // may contains Queries("search", "{search}", "skip", "{skip}", "limit", "{limit}")
	m.HandleFunc("/v1/todo/{id}", handle.HandleGetTodoByID).Methods("GET")
	m.HandleFunc("/v1/todo/{id}", handle.HandleUpdateTodo).Methods("PUT")
	m.HandleFunc("/v1/todo/{id}", handle.HandleDeleteTodo).Methods("DELETE")
	m.HandleFunc("/v1/todo/{id}/task{taskID}", handle.HandleDeleteTask).Methods("DELETE")

	return m, nil
}
