package app

import (
	"hw1/internal/model"
)

type TaskStorage interface {
    Save(task model.Task) error
	GetAll() ([]model.Task, error)
    GetByID(id model.TaskID) (model.Task, error)
    UpdateStatus(id model.TaskID, newStatus model.TaskStatus) error
	EditTask(newTask model.Task) error
    Delete(id model.TaskID) error
}
