// I gonna build an application for Task Tracker CLI tool
// that using Golang without any 3rd party framwork or library.

// Task tracker is a project used to track and manage your tasks.
// In this task, you will build a simple command line interface (CLI) to track what you need to do, what you have done,
// and what you are currently working on. This project will help you practice your programming skills,
// including working with the filesystem, handling user inputs, and building a simple CLI application.

// Requirements:
// Add, Update, and Delete tasks
// Mark a task as in progress or done
// List all tasks
// List all tasks that are done
// List all tasks that are not done
// List all tasks that are in progress

// Implementation:
// Use a JSON file to store the tasks in the current directory.
// The JSON file should be created if it does not exist.
// Use the native file system module of your programming language to interact with the JSON file.
// Do not use any external libraries or frameworks to build this project.
// Ensure to handle errors and edge cases gracefully.

// Example:
// # Adding a new task
// task-cli add "Buy groceries"
// # Output: Task added successfully (ID: 1)

// # Updating and deleting tasks
// task-cli update 1 "Buy groceries and cook dinner"
// task-cli delete 1

// # Marking a task as in progress or done
// task-cli mark-in-progress 1
// task-cli mark-done 1

// # Listing all tasks
// task-cli list

// # Listing tasks by status
// task-cli list done
// task-cli list todo
// task-cli list in-progress.

// task:
// {"id", "name", "status", "created_at", "updated_at"}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"time"
)

type TaskStatus string

const (
	TODO        TaskStatus = "TODO"
	IN_PROGRESS TaskStatus = "IN_PROGRESS"
	DONE        TaskStatus = "DONE"
)

const fileName string = "task_list.json"

func (s TaskStatus) String() string {
	switch s {
	case TODO:
		return "todo"
	case IN_PROGRESS:
		return "in_progress"
	case DONE:
		return "done"
	default:
		return ""
	}
}

type Task struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Status    TaskStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func help() string {
	return `task-cli [OPTIONS]

options:
add 					add item to task list
update 					update name of an item by item id
delete 					delete item by id
mark-in-progress 		update item status to mark in progress
mark-done 				update item status to mark done
list 					list all item in the list
`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(help())
		return
	}

	inputAction := os.Args[1]

	switch inputAction {
	case "add":
		if len(os.Args) < 3 {
			log.Println("Missing task name")
			return
		}
		name := os.Args[2]
		addTask(name, fileName)
	case "update":
		if len(os.Args) < 4 {
			log.Println("Missing id or name")
			return
		}
		id := os.Args[2]
		name := os.Args[3]
		taskId, err := strconv.Atoi(id)
		if err != nil {
			log.Println("Invalid id:", err)
			return
		}
		updateTask(taskId, name, "", fileName)
	case "delete":
		if len(os.Args) < 3 {
			log.Println("Missing task id")
			return
		}
		id := os.Args[2]
		taskId, err := strconv.Atoi(id)
		if err != nil {
			log.Println("Invalid id:", err)
			return
		}
		deleteTask(taskId, fileName)
	case "mark-in-progress":
		if len(os.Args) < 3 {
			log.Println("Missing task id")
			return
		}
		id := os.Args[2]
		taskId, err := strconv.Atoi(id)
		if err != nil {
			log.Println("Invalid id:", err)
			return
		}
		updateTask(taskId, "", IN_PROGRESS, fileName)
	case "mark-done":
		if len(os.Args) < 3 {
			log.Println("Missing task id")
			return
		}
		id := os.Args[2]
		taskId, err := strconv.Atoi(id)
		if err != nil {
			log.Println("Invalid id:", err)
			return
		}
		updateTask(taskId, "", DONE, fileName)
	case "list":
		var filter = ""
		if len(os.Args) > 2 {
			filter = os.Args[2]
		}
		tasks := listTasks(TaskStatus(filter), fileName)
		log.Println(tasks)
	default:
		log.Println("Invalid options: ", help())
		os.Exit(0)
	}
}

func addTask(taskName, fileName string) {
	var tasks []Task

	// 1. Read file
	err := loadTasks(fileName, &tasks)
	if err != nil && !os.IsNotExist(err) {
		log.Println("Error load tasks:", err)
		return
	}

	task := Task{
		ID:        len(tasks) + 1,
		Name:      taskName,
		Status:    TODO,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 2. Append record into slice
	tasks = append(tasks, task)

	// 3. Marshalling & save data
	err = saveTasks(fileName, &tasks)
	if err != nil {
		log.Println("Error saving task:", err)
		return
	}

	log.Printf("Added %s successfully\n", taskName)
}

func listTasks(filter TaskStatus, fileName string) string {
	var tasks []Task

	// read file
	err := loadTasks(fileName, &tasks)
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

func updateTask(id int, name string, status TaskStatus, fileName string) {
	var tasks []Task

	// 1. Read file
	err := loadTasks(fileName, &tasks)
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

	err = saveTasks(fileName, &tasks)
	if err != nil {
		log.Println("Error update tasks:", err)
		return
	}
	log.Println("Updated task ID:", id)
}

func deleteTask(id int, fileName string) {
	var tasks []Task

	// 1. Read file
	err := loadTasks(fileName, &tasks)
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

	err = saveTasks(fileName, &tasks)
	if err != nil {
		log.Println("Error update tasks:", err)
		return
	}

	log.Println("Deleted task ID:", id)
}
