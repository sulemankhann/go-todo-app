package cmd

import (
	"github.com/sulemankhann/go-todo-app/internal/task"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task to the todo list",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskStore := task.NewCSVStore("data.csv")
		tm := task.NewTaskManager(taskStore)
		tm.CreateTask(args[0])
	},
}
