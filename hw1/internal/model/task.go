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
