package types

import (
	"strconv"
	"time"
)

type Task struct {
	Id          int
	Description string
	Created     time.Time
	IsComplete  bool
}

func (t Task) ToCSVRecord() []string {
	return []string{
		strconv.Itoa(t.Id),
		t.Description,
		t.Created.Format(time.RFC3339),
		strconv.FormatBool(t.IsComplete),
	}
}

type TaskStore interface {
	GetTaskList() ([]Task, error)
	CreateTask(description string) (Task, error)
}
