package task

import (
	"encoding/csv"
	"os"
	"time"

	"github.com/sulemankhann/go-todo-app/types"
	"github.com/sulemankhann/go-todo-app/utils"
)

type CSVStore struct {
	filePath string
}

var header = []string{"ID", "Description", "CreatedAt", "IsComplete"}

func NewCSVStore(filePath string) *CSVStore {
	return &CSVStore{filePath: filePath}
}

func (s *CSVStore) GetTaskList() ([]types.Task, error) {
	records, err := getRawCSVRecords(s.filePath)
	if err != nil {
		return nil, err
	}

	tasks, err := utils.CsvRowsToTasks(records)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *CSVStore) CreateTask(description string) (types.Task, error) {
	records, err := getRawCSVRecords(s.filePath)
	if err != nil {
		return types.Task{}, err
	}

	id := utils.GetNextUniqueID(records)
	task := types.Task{
		Id:          id,
		Description: description,
		Created:     time.Now(),
		IsComplete:  false,
	}

	addHeader := len(records) == 0 // if no existing records, add csv header

	err = saveTaskToCSV(s.filePath, task, addHeader)
	if err != nil {
		return types.Task{}, err
	}

	return task, nil
}

func saveTaskToCSV(filePath string, task types.Task, addHeader bool) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

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

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}
