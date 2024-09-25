package types

import "time"

type Task struct {
	Id          int
	Description string
	Created     time.Time
	IsComplete  bool
}

type TaskStore interface {
	GetTaskList() ([]Task, error)
}
