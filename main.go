package main

import (
	"fmt"
	"log"
	"minhnghia2k3/task-tracker-cli/internal/task"
	"os"
	"strconv"
)

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
		task.AddTask(name)
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
		task.UpdateTask(taskId, name, "")
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
		task.DeleteTask(taskId)
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
		task.UpdateTask(taskId, "", task.TASK_STATUS_IN_PROGRESS)
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
		task.UpdateTask(taskId, "", task.TASK_STATUS_DONE)
	case "list":
		var filter = ""
		if len(os.Args) > 2 {
			filter = os.Args[2]
		}
		tasks := task.ListTasks(task.TaskStatus(filter))
		log.Println(tasks)
	default:
		log.Println("Invalid options: ", help())
		os.Exit(0)
	}
}
