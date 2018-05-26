package repositories

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/elumbantoruan/todo/models"
	"github.com/google/uuid"
	"github.com/peterbourgon/diskv"
	"github.com/pkg/errors"
)

// FileStorageTodoRepository represent a concerete implementation
// of TodoRepository
type FileStorageTodoRepository struct {
	disk *diskv.Diskv
}

// NewFileStorageTodoRepository creates an instance of
// FileStorageTodoRepository
func NewFileStorageTodoRepository(path string) *FileStorageTodoRepository {
	flatTransform := func(s string) []string { return []string{} }

	d := diskv.New(diskv.Options{
		BasePath:     path,
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})
	return &FileStorageTodoRepository{
		disk: d,
	}
}

// AddTodo adds new todo
// The record is stored in the disk folder defines in path
func (f *FileStorageTodoRepository) AddTodo(todo models.Todo) error {

	// check for dups
	var cancel = make(chan struct{})
	keys := f.disk.Keys(cancel)
	for key := range keys {
		if key == todo.ID.String() {
			return errors.New("duplicate todoID")
		}
	}

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	enc.Encode(todo)

	err := f.disk.Write(todo.ID.String(), buffer.Bytes())
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

// AddTask adds task to existing todo
// First, it needs to fetch existing todo, deserialize it, and append new task
// to list of task
func (f *FileStorageTodoRepository) AddTask(todoID uuid.UUID, task models.Task) error {
	value, _ := f.disk.Read(todoID.String())
	var (
		todo models.Todo
		err  error
	)

	dec := gob.NewDecoder(bytes.NewReader(value))

	err = dec.Decode(&todo)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	// check for dups id
	for _, t := range todo.Tasks {
		if t.ID == task.ID {
			return errors.New("duplicate taskId")
		}
	}

	todo.Tasks = append(todo.Tasks, task)

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err = enc.Encode(todo)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	err = f.disk.Write(todo.ID.String(), buffer.Bytes())
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

// GetTodo return list of todo
func (f *FileStorageTodoRepository) GetTodo() ([]models.Todo, error) {
	var (
		cancel   = make(chan struct{})
		todoList []models.Todo
	)

	keys := f.disk.Keys(cancel)
	for key := range keys {
		var todo models.Todo
		value, _ := f.disk.Read(key)
		dec := gob.NewDecoder(bytes.NewReader(value))
		err := dec.Decode(&todo)
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		todoList = append(todoList, todo)
	}
	return todoList, nil
}

// GetTodoByID return todo by id
func (f *FileStorageTodoRepository) GetTodoByID(todoID uuid.UUID) (*models.Todo, error) {
	var (
		todo models.Todo
		err  error
		bts  []byte
	)
	bts, err = f.disk.Read(todoID.String())
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}

	dec := gob.NewDecoder(bytes.NewReader(bts))

	err = dec.Decode(&todo)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	return &todo, nil
}

// UpdateTodo updates todo
func (f *FileStorageTodoRepository) UpdateTodo(todoID uuid.UUID, completed bool, dueDate *time.Time) error {
	todo, err := f.GetTodoByID(todoID)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	todo.Completed = completed
	todo.DueDate = dueDate

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err = enc.Encode(todo)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	err = f.disk.Write(todo.ID.String(), buffer.Bytes())
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

// UpdateTask updates task for a specific todo
func (f *FileStorageTodoRepository) UpdateTask(todoID uuid.UUID, taskID uuid.UUID, completed bool) error {
	value, _ := f.disk.Read(todoID.String())
	var (
		todo models.Todo
		err  error
	)

	dec := gob.NewDecoder(bytes.NewReader(value))

	err = dec.Decode(&todo)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	for i := 0; i < len(todo.Tasks); i++ {
		if todo.Tasks[i].ID == taskID {
			todo.Tasks[i].Completed = completed
			break
		}
	}
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err = enc.Encode(todo)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	err = f.disk.Write(todo.ID.String(), buffer.Bytes())
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

// DeleteTask deletes task
func (f *FileStorageTodoRepository) DeleteTask(todoID, taskID uuid.UUID) error {
	value, _ := f.disk.Read(todoID.String())
	var (
		todo models.Todo
		err  error
	)

	dec := gob.NewDecoder(bytes.NewReader(value))

	err = dec.Decode(&todo)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	for i := 0; i < len(todo.Tasks); i++ {
		if todo.Tasks[i].ID == taskID {
			// remove the task from the slice
			if i == 0 {
				todo.Tasks = todo.Tasks[1:]
			} else {
				// https://github.com/golang/go/wiki/SliceTricks#filtering-without-allocating
				todo.Tasks = append(todo.Tasks[:i], todo.Tasks[i+1:]...)
			}
			break
		}
	}
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err = enc.Encode(todo)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	err = f.disk.Write(todo.ID.String(), buffer.Bytes())
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

// DeleteTodo deletes todo
func (f *FileStorageTodoRepository) DeleteTodo(todoID uuid.UUID) error {
	err := f.disk.Erase(todoID.String())
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	return nil
}
