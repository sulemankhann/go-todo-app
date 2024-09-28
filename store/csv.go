package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/sulemankhann/go-todo-app/todo"
)

type Store struct {
	filePath string
}

var header = []string{"ID", "Description", "CreatedAt", "IsComplete"}

func NewStore(filePath string) *Store {
	return &Store{filePath: filePath}
}

func (s *Store) GetTaskList() ([]todo.Task, error) {
	records, err := getRawCSVRecords(s.filePath)
	if err != nil {
		return nil, err
	}

	tasks, err := csvRowsToTasks(records)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Store) CreateTask(description string) (todo.Task, error) {
	records, err := getRawCSVRecords(s.filePath)
	if err != nil {
		return todo.Task{}, err
	}

	id := getNextUniqueID(records)
	task := todo.Task{
		Id:          id,
		Description: description,
		Created:     time.Now(),
		IsComplete:  false,
	}

	addHeader := len(records) == 0 // if no existing records, add csv header

	err = saveTaskToCSV(s.filePath, task, addHeader)
	if err != nil {
		return todo.Task{}, err
	}

	return task, nil
}

func (s *Store) MarkTaskCompleted(id int) (todo.Task, error) {
	tasks, err := s.GetTaskList()
	if err != nil {
		return todo.Task{}, err
	}

	records := [][]string{header}
	task := todo.Task{}

	for _, t := range tasks {
		if t.Id == id {
			t.IsComplete = true
			task = t
		}
		records = append(records, t.ToCSVRecord())
	}

	if task == (todo.Task{}) {
		return todo.Task{}, fmt.Errorf("Task with id %d not found", id)
	}

	err = writeRecordsToCSV(s.filePath, records)
	if err != nil {
		return todo.Task{}, err
	}

	return task, nil
}

func (s *Store) DeleteTask(id int) error {
	tasks, err := s.GetTaskList()
	if err != nil {
		return err
	}

	records := [][]string{header}
	task := todo.Task{}

	for _, t := range tasks {
		if t.Id == id {
			task = t
		} else {
			records = append(records, t.ToCSVRecord())
		}
	}

	if task == (todo.Task{}) {
		return fmt.Errorf("Task with id %d not found", id)
	}

	err = writeRecordsToCSV(s.filePath, records)
	if err != nil {
		return err
	}

	return nil
}

func writeRecordsToCSV(filePath string, records [][]string) error {
	file, err := os.OpenFile(
		filePath,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0644,
	)
	if err != nil {
		return err
	}

	defer file.Close()

	// Lock the file for exclusive access
	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
	if err != nil {
		return err
	}
	defer syscall.Flock(int(file.Fd()), syscall.LOCK_UN) // Unlock when done

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(records)
	if err != nil {
		return err
	}

	return nil
}

func saveTaskToCSV(filePath string, task todo.Task, addHeader bool) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	// Lock the file for exclusive access
	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
	if err != nil {
		return err
	}
	defer syscall.Flock(int(file.Fd()), syscall.LOCK_UN) // Unlock when done

	w := csv.NewWriter(file)
	defer w.Flush()

	r := [][]string{}

	if addHeader {
		r = append(r, header)
	}

	r = append(r, task.ToCSVRecord())

	err = w.WriteAll(r)
	if err != nil {
		return err
	}

	return nil
}

func getRawCSVRecords(filePath string) ([][]string, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Apply a shared lock (multiple readers allowed, no writers)
	err = syscall.Flock(int(file.Fd()), syscall.LOCK_SH)
	if err != nil {
		return nil, err
	}
	defer syscall.Flock(int(file.Fd()), syscall.LOCK_UN) // Unlock after reading

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

// csvRowsToTasks takes a 2D slice of strings (CSV records)
// and returns []Task
func csvRowsToTasks(records [][]string) ([]todo.Task, error) {
	tasks := []todo.Task{}
	if len(records) == 0 {
		return tasks, nil
	}

	// Skip header
	for _, record := range records[1:] {
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}

		createdAt, err := time.Parse(time.RFC3339, record[2])
		if err != nil {
			return nil, err
		}

		isComplete, err := strconv.ParseBool(record[3])
		if err != nil {
			return nil, err
		}

		task := todo.Task{
			Id:          id,
			Description: record[1],
			Created:     createdAt,
			IsComplete:  isComplete,
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// getNextUniqueID takes a 2D slice of strings (CSV records)
// and returns the next unique ID based on existing IDs in the records
func getNextUniqueID(records [][]string) (maxId int) {
	if len(records) == 0 {
		return maxId + 1
	}

	for _, record := range records[1:] { // Skip header
		id, err := strconv.Atoi(record[0]) // Assuming ID is in the first column
		if err == nil && id > maxId {
			maxId = id
		}
	}

	return maxId + 1
}
