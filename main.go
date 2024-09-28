package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/sulemankhann/go-todo-app/store"
	"github.com/sulemankhann/go-todo-app/todo"
)

var rootCmd = &cobra.Command{
	Use:   "tasks",
	Short: "A todo list for the terminal",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	store := csv.NewStore("data.csv")
	tm := todo.NewTodoManager(store)

	var dueDate string
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new task to the todo list",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(dueDate)
			tm.CreateTask(args[0], dueDate)
		},
	}
	addCmd.Flags().
		StringVarP(&dueDate, "due", "d", "", "Due date of task")

	var showAll bool
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all of the tasks in your todo list",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			tm.ListTask(showAll)
		},
	}
	listCmd.Flags().
		BoolVarP(&showAll, "all", "a", false, "Show all of the tasks")

	completeCmd := &cobra.Command{
		Use:   "complete",
		Short: "Set a task as being completed",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error: argument must be an integer")
				return
			}
			tm.CompleteTask(id)
		},
	}

	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Removes a task from the todo list by it's id",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error: argument must be an integer")
				return
			}
			tm.DeleteTask(id)
		},
	}

	rootCmd.AddCommand(addCmd, listCmd, completeCmd, deleteCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
