package task

import (
	"encoding/csv"
	"os"
	"strconv"
	"sulemankhann/go-todo-app/types"
	"time"
)

type CSVStore struct {
	filePath string
}

func NewCSVStore(filePath string) *CSVStore {
	return &CSVStore{filePath: filePath}
}

func (s *CSVStore) GetTaskList() ([]types.Task, error) {
	file, err := os.OpenFile(s.filePath, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	tasks, err := csvRowsToTasks(records)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func csvRowsToTasks(records [][]string) ([]types.Task, error) {
	tasks := []types.Task{}
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

		task := types.Task{
			Id:          id,
			Description: record[1],
			Created:     createdAt,
			IsComplete:  isComplete,
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
