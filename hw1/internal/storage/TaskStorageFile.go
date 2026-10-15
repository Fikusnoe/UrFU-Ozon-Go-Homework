package storage

import (
	"encoding/json"
	"fmt"
	"hw1/internal/model"
	"os"
)

type TaskStorageFile struct {
	filePath string
	NextID   model.TaskID
}

func (s *TaskStorageFile) GetNextID() model.TaskID, error {
	tasks := s.LoadTasks()

	return model.TaskID(len(tasks) + 1)
}

func (s *TaskStorageFile) Save(task model.Task) error {
	tasks := s.LoadTasks()

	tasks = append(tasks, task)

	s.SaveTasks(tasks)

	return nil
}

func (s *TaskStorageFile) GetAll() ([]model.Task, error) {
	tasks := s.LoadTasks()

	return tasks, nil
}

func (s *TaskStorageFile) GetByID(id model.TaskID) (model.Task, error) {
	tasks := s.LoadTasks()

	for _, task := range tasks {
		if task.ID == id {
			return task, nil
		}
	}

	return model.Task{}, fmt.Errorf("Такой задачи нет!")
}

func (s *TaskStorageFile) UpdateStatus(id model.TaskID, newStatus model.TaskStatus) error {
	tasks := s.LoadTasks()

	for index, task := range tasks {
		if task.ID == id {
			task.Status = newStatus
			tasks[index] = task
			break
		}
	}

	s.SaveTasks(tasks)

	return nil
}

func (s *TaskStorageFile) EditTask(newTask model.Task) error {
	tasks := s.LoadTasks()

	for index, task := range tasks {
		if task.ID == newTask.ID {
			task.Name = newTask.Name
			task.Deadline = newTask.Deadline
			task = newTask
			tasks[index] = task
			break
		}
	}

	s.SaveTasks(tasks)

	return nil
}

func (s *TaskStorageFile) Delete(id model.TaskID) error {
	tasks := s.LoadTasks()

	for index, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:index], tasks[index+1:]...)
			break
		}
	}

	s.SaveTasks(tasks)

	return nil
}

func (s *TaskStorageFile) LoadTasks() []model.Task {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return []model.Task{}
	}

	tasks := []model.Task{}

	json.Unmarshal(data, &tasks)

	return tasks
}

func (s *TaskStorageFile) SaveTasks(tasks []model.Task) error {
	data, _ := json.MarshalIndent(tasks, "", " ")

	os.WriteFile(s.filePath, data, 0644)

	return nil
}

func MakeNewStorageFile(filepath string) *TaskStorageFile {
	storage := TaskStorageFile{
		filePath: filepath,
		NextID:   1,
	}
	return &storage
}
