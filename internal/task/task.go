package task

import (
	"fmt"
	"sulemankhann/go-todo-app/types"
)

type TaskManager struct {
	store types.TaskStore
}

func NewTaskManager(store types.TaskStore) *TaskManager {
	return &TaskManager{store: store}
}

func (tm *TaskManager) ListTask() {
	tasks, err := tm.store.GetTaskList()
	if err != nil {
		panic(err)
	}

	fmt.Println(tasks)
}
