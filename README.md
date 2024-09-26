# Todo App

A simple terminal-based Todo application built with Go. This application allows users to manage their tasks efficiently through a command-line interface.

## Features

- **Add a Task**: Add a task to the todo list by running:
  ```bash
  ./bin/task add "<description>" example ./bin/task add "Tidy my desk"
  ```
  *Example*: 
  ```bash
  ./bin/task add "Tidy my desk"
  ```

- **Show List of Tasks**: Display the list of tasks by running:
  ```bash
  ./bin/task list
  ```
  To show all tasks, you can use a flag:
  ```bash
  ./bin/task list -a
  ```
  or
  ```bash
  ./bin/task list --all
  ```

- **Complete a Task**: Mark a task as done by running:
  ```bash
  ./bin/task complete <taskid>
  ```

- **Delete a Task**: Remove a task from the data store by running:
  ```bash
  ./bin/task delete <taskid>
  ```

## How to Build and Run the Project

1. Make sure you have the latest version of Go installed.
2. Clone the project repository:
   ```bash
   git clone <repository-url>
   ```
3. Navigate to the root of the project:
   ```bash
   cd <project-directory>
   ```
4. Run the following command to build the project:
   ```bash
   make build
   ```
5. Execute the application using:
   ```bash
   ./bin/task
   ```
