package storage

import (
	"encoding/json"
	"fmt"
	"hw1/internal/model"
	"math"
	"os"
)

type TaskStorageFile struct {
	filePath string
	NextID   model.TaskID
}

func (s *TaskStorageFile) GetNextID() (model.TaskID, error) {
	tasks, err := s.LoadTasks()
	if err != nil {
		return model.TaskID(0), err
	}

	currLastID := tasks[len(tasks)-1].ID

	if int(currLastID) >= math.MaxInt {
		for i := 1; i <= (len(tasks) - 1); i++ {
			if int(tasks[i].ID) != i {
				return model.TaskID(i), nil
			}
		}
		return 0, fmt.Errorf("Превышено количество записей")
	}

	return currLastID + 1, nil
}

func (s *TaskStorageFile) Save(task model.Task) error {
	tasks, err := s.LoadTasks()
	if err != nil {
		return err
	}

	tasks = append(tasks, task)

	err = s.SaveTasks(tasks)
	if err != nil {
		return err
	}

	return nil
}

func (s *TaskStorageFile) GetAll() ([]model.Task, error) {
	tasks, err := s.LoadTasks()
	if err != nil {
		return []model.Task{}, err
	}

	return tasks, nil
}

func (s *TaskStorageFile) GetByID(id model.TaskID) (model.Task, error) {
	tasks, err := s.LoadTasks()
	if err != nil {
		return model.Task{}, err
	}

	for _, task := range tasks {
		if task.ID == id {
			return task, nil
		}
	}

	return model.Task{}, fmt.Errorf("Такой задачи нет!")
}

func (s *TaskStorageFile) UpdateStatus(id model.TaskID, newStatus model.TaskStatus) error {
	tasks, err := s.LoadTasks()
	if err != nil {
		return err
	}

	for index, task := range tasks {
		if task.ID == id {
			if task.Status == newStatus {
				return fmt.Errorf("Текущий статус уже установлен.")
			}

			task.Status = newStatus
			tasks[index] = task

			err = s.SaveTasks(tasks)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("Такой задачи нет!")
}

func (s *TaskStorageFile) EditTask(newTask model.Task) error {
	tasks, err := s.LoadTasks()
	if err != nil {
		return err
	}

	for index, task := range tasks {
		if task.ID == newTask.ID {
			if newTask.Name != "" {
				task.Name = newTask.Name
			}
			task.Deadline = newTask.Deadline
			tasks[index] = task

			err = s.SaveTasks(tasks)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("Такой задачи нет!")
}

func (s *TaskStorageFile) Delete(id model.TaskID) error {
	tasks, err := s.LoadTasks()
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		return fmt.Errorf("Список задач пуст")
	}

	for index, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:index], tasks[index+1:]...)

			err = s.SaveTasks(tasks)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("Такой задачи нет!")
}

func (s *TaskStorageFile) LoadTasks() ([]model.Task, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return []model.Task{}, err
	}

	tasks := []model.Task{}

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return []model.Task{}, err
	}

	return tasks, nil
}

func (s *TaskStorageFile) SaveTasks(tasks []model.Task) error {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(s.filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func MakeNewStorageFile(filePath string) *TaskStorageFile {
	storage := TaskStorageFile{
		filePath: filePath,
		NextID:   1,
	}
	return &storage
}
