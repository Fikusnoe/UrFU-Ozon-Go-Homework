package app

import (
	"hw1/internal/model"
	"hw1/internal/storage"
	"time"
	"fmt"
)

type TaskHandler struct {
	storage TaskStorage
}

func (ts *TaskHandler) AddTask(name string,deadline time.Time) {
	status := model.StatusPlanned
	id := ts.storage.GetNextID()
	newTask := model.MakeNewTask(id, name, status, deadline)
	ts.storage.Save(newTask)
}

func (ts *TaskHandler) GetTaskById(id model.TaskID) (model.Task, error) {
	task, err := ts.storage.GetByID(id)
	if err != nil {
		return model.Task{}, fmt.Errorf("ошибка получения задачи %d: %w", id, err)
	}
	return task, nil
}

func (ts *TaskHandler) GetAllTasks() ([]model.Task, error) {
	allTasks, _ := ts.storage.GetAll()
	return allTasks, nil
}

func (ts *TaskHandler) UpdateTaskStatus(id model.TaskID, status model.TaskStatus) error {
	err := ts.storage.UpdateStatus(id, status)
	return err
}

func (ts *TaskHandler) EditTask(id model.TaskID, name string, deadline time.Time) error {
	status := model.StatusPlanned
	newInfoTask := model.MakeNewTask(id, name, status, deadline)
	err := ts.storage.EditTask(newInfoTask)
	return err
}

func (ts *TaskHandler) DeleteTask(id model.TaskID) error {
	err := ts.storage.Delete(id)
	return err
}

func MakeNewTaskHandler(storage storage.TaskStorageMap) TaskHandler {
	handler := TaskHandler{storage: &storage}
	return handler
}

