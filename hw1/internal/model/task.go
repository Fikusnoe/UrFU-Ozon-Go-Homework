package model

import (
	"errors"
	"strconv"
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
	ID   TaskID `json:"id"`
	Name string	`json:"title"`
	Status TaskStatus `json:"status"`
	Deadline time.Time `json:"deadline"`
}

func (s *Task) String() string {
	strID := strconv.Itoa(int(s.ID))
	strDeadline := s.Deadline.Format("2006-01-02")

	strTask := strID + " " + s.Name + " " + string(s.Status) + " " + strDeadline
	return strTask
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

func ParseID(strID string) (TaskID, error) {
	intID, _ := strconv.Atoi(strID)
	taskID := TaskID(intID)
	return taskID, nil
}

func ParseStatus(strStatus string) (TaskStatus, error) {
	status := TaskStatus("")
	
	switch strStatus{
	case "done":
		status = StatusDone
	case "in_progress":
		status = StatusInProgress
	case "cancelled":
		status = StatusCancelled
	case "planned":
		status = StatusPlanned
	default:
		return "", errors.New("Неизвестный статус")
	}

	return status, nil
}

