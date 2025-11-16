package task

import (
	"encoding/json"
	"log"
	"slices"
	"time"
)

type TaskStatus string

const (
	TASK_STATUS_TODO        TaskStatus = "todo"
	TASK_STATUS_IN_PROGRESS TaskStatus = "in-progress"
	TASK_STATUS_DONE        TaskStatus = "done"
)

type Task struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Status    TaskStatus `json:"status"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

func AddTask(taskName string) {
	var tasks []Task

	// 1. Read file
	err := LoadTasks(&tasks)
	if err != nil {
		log.Println("Error load tasks:", err)
		return
	}

	task := Task{
		ID:        len(tasks) + 1,
		Name:      taskName,
		Status:    TASK_STATUS_TODO,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 2. Append record into slice
	tasks = append(tasks, task)

	// 3. Marshalling & save data
	err = SaveTasks(&tasks)
	if err != nil {
		log.Println("Error saving task:", err)
		return
	}

	log.Printf("Added %s successfully\n", taskName)
}

func ListTasks(filter TaskStatus) string {
	var tasks []Task

	// read file
	err := LoadTasks(&tasks)
	if err != nil {
		log.Println("Error load tasks:", err)
		return ""
	}

	if filter != "" {
		var filteredTask []Task
		for _, task := range tasks {
			if task.Status == filter {
				filteredTask = append(filteredTask, task)
			}
		}
		tasks = filteredTask
	}

	output, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		log.Printf("Error when marshalling the data: %v\n", err)
		return ""
	}

	return string(output)
}

func UpdateTask(id int, name string, status TaskStatus) {
	var tasks []Task

	// 1. Read file
	err := LoadTasks(&tasks)
	if err != nil {
		log.Println("Error load tasks:", err)
		return
	}

	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			if name != "" {
				tasks[i].Name = name
			}
			if status != "" {
				tasks[i].Status = status
			}
			tasks[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		log.Println("Task does not exist with id:", id)
		return
	}

	err = SaveTasks(&tasks)
	if err != nil {
		log.Println("Error update tasks:", err)
		return
	}
	log.Println("Updated task ID:", id)
}

func DeleteTask(id int) {
	var tasks []Task

	// 1. Read file
	err := LoadTasks(&tasks)
	if err != nil {
		log.Println("Error load tasks:", err)
		return
	}

	// Delete task by task id
	found := false
	for i, v := range tasks {
		if v.ID == id {
			tasks = slices.Delete(tasks, i, i+1)
			found = true
			break
		}
	}

	if !found {
		log.Println("Task does not exist:", id)
		return
	}

	err = SaveTasks(&tasks)
	if err != nil {
		log.Println("Error update tasks:", err)
		return
	}

	log.Println("Deleted task ID:", id)
}
