package task

import (
	"fmt"
	"os"
	"strings"
	"sulemankhann/go-todo-app/types"
	"sulemankhann/go-todo-app/utils"
	"text/tabwriter"

	"github.com/mergestat/timediff"
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

	printTasks(tasks)
}

func printTasks(tasks []types.Task) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	defer w.Flush()

	// Print the headers
	fmt.Fprintln(w, "ID\tTask\tCreated\tDone")

	// Print the rows
	for _, task := range tasks {
		printRow(
			w,
			task.Id,
			task.Description,
			timediff.TimeDiff(task.Created),
			task.IsComplete,
		)
	}
}

// printRow is a helper function to wrap text in the Task column if necessary
func printRow(w *tabwriter.Writer, id int, task, created string, status bool) {
	// Define max width for the Task column
	maxWidth := 70

	// Split the task into multiple lines if it's longer than maxWidth
	wrappedTask := utils.WrapText(task, maxWidth)

	// Split the task into lines and print each line with proper alignment
	taskLines := strings.Split(wrappedTask, "\n")
	for i, line := range taskLines {
		if i == 0 {
			// Print the first line with the all  column
			fmt.Fprintf(w, "%d\t%s\t%s\t%t\n", id, line, created, status)
		} else {
			// Print subsequent lines with blank ID and Created columns for alignment
			fmt.Fprintf(w, "\t%s\t\n", line)
		}
	}
}
