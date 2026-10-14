package app

import (
	"hw1/internal/storage"
)

type TaskHandler struct {
	storage TaskStorage
}

func (ts *TaskHandler) AddTask() {

}

func (ts *TaskHandler) UpdateTaskStatus() {

}

func (ts *TaskHandler) EditTask() {

}

func (ts *TaskHandler) DeleteTask() {

}

func MakeNewTaskHandler(storage storage.TaskStorageMap) TaskHandler {
	handler := TaskHandler{storage: &storage}
	return handler
}
