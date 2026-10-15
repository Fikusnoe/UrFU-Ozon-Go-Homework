package app

import (
	"fmt"
	"hw1/internal/model"
	"slices"
	"time"
)

type TaskHandler struct {
	storage TaskStorage
}

func MakeNewTaskHandler(storage TaskStorage) TaskHandler {
	handler := TaskHandler{storage: storage}
	return handler
}

func (ts *TaskHandler) AddTask(name string, deadline time.Time) error {
	status := model.StatusPlanned
	id, _ := ts.storage.GetNextID()
	newTask := model.MakeNewTask(id, name, status, deadline)
	err := ts.storage.Save(newTask)
	if err != nil{
		return fmt.Errorf("Не удалось сохранить задачу: %w", err)
	}
	
	return nil
}

func (ts *TaskHandler) GetTaskById(id model.TaskID) (model.Task, error) {
	task, err := ts.storage.GetByID(id)
	if err != nil {
		return model.Task{}, fmt.Errorf("Ошибка получения задачи %d: %w", id, err)
	}
	return task, nil
}

func (ts *TaskHandler) GetAllTasks() ([]model.Task, error) {
	allTasks, err := ts.storage.GetAll()
	if err != nil{
		return allTasks, fmt.Errorf("Не удалось получить список задач: %w", err)
	}
	allTasks = SortTasks(allTasks)
	return allTasks, nil
}

func (ts *TaskHandler) UpdateTaskStatus(id model.TaskID, status model.TaskStatus) error {
	err := ts.storage.UpdateStatus(id, status)
	if err != nil {
		return fmt.Errorf("Не удалось изменить статус задачи %d: %w", id, err)
	}
	return nil
}

func (ts *TaskHandler) EditTask(id model.TaskID, name string, deadline time.Time) error {
	status := model.StatusPlanned
	newInfoTask := model.MakeNewTask(id, name, status, deadline)
	err := ts.storage.EditTask(newInfoTask)
	if err != nil {
		return fmt.Errorf("Не удалось изменить задачу %d: %w", id, err)
	}
	return nil
}

func (ts *TaskHandler) DeleteTask(id model.TaskID) error {
	err := ts.storage.Delete(id)
	if err != nil {
		return fmt.Errorf("Не удалось удалить задачу %d: %w", id, err)
	}
	return nil
}

func SortTasks(tasks []model.Task) []model.Task {
	slices.SortFunc(tasks, func(a, b model.Task) int {
		if a.Deadline.Before(b.Deadline) {
			return 1
		}
		if a.Deadline.Before(b.Deadline) {
			return -1
		}
		if a.ID < b.ID {
			return 1
		}
		return -1
	})

	return tasks
}
