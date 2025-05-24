package db

import (
	"github.com/alle/tasks/model"
)

// This service is not being used by the task manager. Instead we are using postgres for crus operations

type IDBService interface {
	Create(task *model.Task) error
	ReadById(id string) (*model.Task, error)
	// Update(task *model.Task) (error)
	// Delete(id string)
}

type DbService struct {
	storage map[string]*model.Task
}

func NewDBService () *DbService {
	mp := make(map[string]*model.Task)
	return &DbService{
		storage: mp,
	}
}

// Create functioncreates a task and store it in map
func (d *DbService) Create(task *model.Task) error {
	d.storage[task.Id.String()] = task
	return nil
}

// ReadById function return a tsk given it's id
func (d *DbService) ReadById(id string) (*model.Task, error) {
	return d.storage[id], nil
}