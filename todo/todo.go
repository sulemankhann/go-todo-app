package todo

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/mergestat/timediff"
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

type Store interface {
	GetTaskList() ([]Task, error)
	CreateTask(description string) (Task, error)
	MarkTaskCompleted(id int) (Task, error)
	DeleteTask(id int) error
}

type TodoManager struct {
	store Store
}

func NewTodoManager(store Store) *TodoManager {
	return &TodoManager{store: store}
}

func (tm *TodoManager) ListTask(showAll bool) {
	tasks, err := tm.store.GetTaskList()
	if err != nil {
		panic(err)
	}

	if !showAll {
		var incompleteTasks []Task
		for _, t := range tasks {
			if !t.IsComplete {
				incompleteTasks = append(incompleteTasks, t)
			}
		}
		tasks = incompleteTasks
	}

	printTasks(tasks)
}

func (tm *TodoManager) CreateTask(description string) {
	task, err := tm.store.CreateTask(description)
	if err != nil {
		panic(err)
	}
	printTasks([]Task{task})
}

func (tm *TodoManager) CompleteTask(id int) {
	task, err := tm.store.MarkTaskCompleted(id)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	printTasks([]Task{task})
}

func (tm *TodoManager) DeleteTask(id int) {
	err := tm.store.DeleteTask(id)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func printTasks(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println(
			"Yayyy!, You have nothing in todolist, Enjoy your free time :)",
		)
		return
	}

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
	wrappedTask := wrapText(task, maxWidth)

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

// wrapText wraps the given text into lines of max width
func wrapText(text string, maxWidth int) (result string) {
	words := strings.Split(text, " ")
	line := ""

	for _, word := range words {
		if len(line)+len(word)+1 > maxWidth {
			result += line + "\n"
			line = word
		} else {
			if len(line) > 0 {
				line += " "
			}
			line += word
		}
	}

	result += line

	return
}
