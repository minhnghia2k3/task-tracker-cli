package main

import (
	"encoding/json"
	"os"
)

func loadTasks(fileName string, tasks *[]Task) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			// Create file with empty JSON array
			err = os.WriteFile(fileName, []byte("[]"), 0644)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	err = json.Unmarshal(data, tasks)
	if err != nil {
		return err
	}

	return nil
}

func saveTasks(fileName string, tasks *[]Task) error {
	data, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
