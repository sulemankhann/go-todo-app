package store

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sulemankhann/go-todo-app/todo"
)

type SqliteStore struct {
	db *sql.DB
}

func NewSqliteStore(dbPath string) (*SqliteStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Ensure tasks table exists
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS tasks (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            description TEXT,
            created_at TEXT,
            completed_at TEXT,
            due_date TEXT
        )
    `)
	if err != nil {
		return nil, err
	}

	return &SqliteStore{db: db}, nil
}

func (s *SqliteStore) CreateTask(
	description, dueDate string,
) (todo.Task, error) {
	task := todo.Task{
		Description: description,
		Created:     time.Now(),
		IsComplete:  time.Time{},
		DueDate:     dueDate,
	}

	res, err := s.db.Exec(
		"INSERT INTO tasks (description, created_at, completed_at, due_date) VALUES (?,?,?,?)",
		description,
		task.Created.Format(time.RFC3339),
		task.IsComplete.Format(time.RFC3339),
		dueDate,
	)
	if err != nil {
		return todo.Task{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return todo.Task{}, err
	}

	task.Id = int(id)

	return task, nil
}

func (s *SqliteStore) GetTaskList() ([]todo.Task, error) {
	rows, err := s.db.Query(
		"SELECT id, description, created_at, completed_at, due_date FROM tasks",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []todo.Task{}
	var createdAt, completedAt string

	for rows.Next() {
		var task todo.Task
		err = rows.Scan(
			&task.Id,
			&task.Description,
			&createdAt,
			&completedAt,
			&task.DueDate,
		)
		if err != nil {
			return nil, err
		}

		task.Created, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, err
		}

		task.IsComplete, err = time.Parse(time.RFC3339, completedAt)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *SqliteStore) MarkTaskCompleted(id int) (todo.Task, error) {
	task, err := s.GetTaskById(id)
	if err != nil {
		return todo.Task{}, err
	}

	completedAt := time.Now()
	_, err = s.db.Exec(
		"UPDATE tasks SET completed_at = ? WHERE id = ?",
		completedAt.Format(time.RFC3339),
		task.Id,
	)
	if err != nil {
		return todo.Task{}, err
	}

	task.IsComplete = completedAt

	return task, nil
}

func (s *SqliteStore) DeleteTask(id int) error {
	task, err := s.GetTaskById(id)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("DELETE FROM tasks WHERE id = ?", task.Id)
	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteStore) GetTaskById(id int) (todo.Task, error) {
	var task todo.Task
	var createdAt, completedAt string

	err := s.db.QueryRow(`
        SELECT id, description, created_at, completed_at, due_date 
        FROM tasks WHERE id = ?`, id).Scan(
		&task.Id, &task.Description, &createdAt, &completedAt, &task.DueDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return todo.Task{}, fmt.Errorf("Task with id %d not found", id)
		}
		return todo.Task{}, err
	}

	task.Created, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return todo.Task{}, err
	}

	task.IsComplete, err = time.Parse(time.RFC3339, completedAt)
	if err != nil {
		return todo.Task{}, err
	}

	return task, nil
}
