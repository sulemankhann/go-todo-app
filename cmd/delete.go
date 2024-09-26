package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/sulemankhann/go-todo-app/internal/task"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Removes a task from the todo list by it's id",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: argument must be an integer")
			return
		}

		taskStore := task.NewCSVStore("data.csv")
		tm := task.NewTaskManager(taskStore)
		tm.DeleteTask(id)
	},
}
