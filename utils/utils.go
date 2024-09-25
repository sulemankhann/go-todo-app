package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/sulemankhann/go-todo-app/types"
)

// WrapText wraps the given text into lines of max width
func WrapText(text string, maxWidth int) (result string) {
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

// GetNextUniqueID takes a 2D slice of strings (CSV records)
// and returns the next unique ID based on existing IDs in the records
func GetNextUniqueID(records [][]string) (maxId int) {
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

// CsvRowsToTasks takes a 2D slice of strings (CSV records)
// and returns []Task
func CsvRowsToTasks(records [][]string) ([]types.Task, error) {
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
		if err != nil {
			return nil, err
		}

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
