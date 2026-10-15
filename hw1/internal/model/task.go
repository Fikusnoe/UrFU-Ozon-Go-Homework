package model

import (
	"time"
)

type TaskID int
type TaskStatus string

const (
	StatusPlanned    TaskStatus = "Planned"
	StatusInProgress TaskStatus = "In_progress"
	StatusCancelled  TaskStatus = "Cancelled"
	StatusDone       TaskStatus = "Done"
)

type Task struct {
	ID   TaskID
	Name string
	Status TaskStatus
	Deadline time.Time
}

func MakeNewTask(id TaskID, name string, status TaskStatus, deadline time.Time) Task {
	newTask := Task{
		ID: id,
		Name: name,
		Status: status,
		Deadline: deadline,
	}

	return newTask
}