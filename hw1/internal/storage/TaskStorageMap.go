package storage

import ( 
	"hw1/internal/model"
)

type TaskStorageMap struct {
	TaskMap map[model.TaskID]model.Task
	NextID model.TaskID
}

func (s *TaskStorageMap) Save(task model.Task) error {
	s.TaskMap[s.NextID] = task 
	return nil 
}

func (s *TaskStorageMap) GetAll() ([]model.Task, error) {
	tasks := []model.Task{}
	for _, task := range s.TaskMap {
		tasks = append(tasks, task)
	}
	return tasks, nil 
}

func (s *TaskStorageMap) GetByID(id model.TaskID) (model.Task, error) {
	return s.TaskMap[id], nil 
}

func (s *TaskStorageMap) UpdateStatus(id model.TaskID, newStatus model.TaskStatus) error {
	task := s.TaskMap[id]
	task.Status = newStatus
	s.TaskMap[id] = task
	return nil 
}

func (s *TaskStorageMap) EditTask(newTask model.Task) error {
	task := s.TaskMap[newTask.ID] 
	if newTask.Name != "" {
		task.Name = newTask.Name
	}
	if !newTask.Deadline.IsZero() {
		task.Deadline = newTask.Deadline
	}
	s.TaskMap[newTask.ID] = task
	return nil 
}

func (s *TaskStorageMap) Delete(id model.TaskID) error {
	delete(s.TaskMap, id)
	return nil 
}