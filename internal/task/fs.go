package task

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func taskFilePath() string {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting task file path:", err)
		return ""
	}

	return path.Join(cwd, "task_list.json")
}

func LoadTasks(tasks *[]Task) error {
	path := taskFilePath()

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.WriteFile(path, []byte("[]"), 0644)
			if err != nil {
				return err
			}
			// File created with empty array, return success
			*tasks = []Task{}
			return nil
		}
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&tasks)
	if err != nil {
		return err
	}

	return nil
}

func SaveTasks(tasks *[]Task) error {
	path := taskFilePath()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(&tasks)
	if err != nil {
		return err
	}

	return nil
}
