package cmd

import (
	"github.com/sulemankhann/go-todo-app/internal/task"

	"github.com/spf13/cobra"
)

var showAll bool

func init() {
	listCmd.Flags().
		BoolVarP(&showAll, "all", "a", false, "Show all of the tasks")
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of the tasks in your todo list",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		taskStore := task.NewCSVStore("data.csv")
		tm := task.NewTaskManager(taskStore)
		tm.ListTask(showAll)
	},
}
