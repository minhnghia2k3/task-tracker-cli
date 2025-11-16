package main

import (
	"minhnghia2k3/task-tracker-cli/internal/task"
	"os"
	"testing"
)

const filePath = "/Users/nghia.mle/Developer/pet-projects/task-tracker-cli/task_list.json"

func TestAddTask(t *testing.T) {
	tests := []struct {
		name         string
		taskName     string
		expectedID   int
		expectedName string
	}{
		{"Add single task", "Buy groceries", 1, "Buy groceries"},
		{"Add another task", "Clean house", 1, "Clean house"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Remove(filePath)       // Clean up before test
			defer os.Remove(filePath) // Clean up after test
			task.AddTask(tt.taskName)

			var tasks []task.Task
			task.LoadTasks(&tasks)

			if len(tasks) != 1 {
				t.Errorf("Expected 1 task, got %d", len(tasks))
			}
			if tasks[0].ID != tt.expectedID {
				t.Errorf("Expected ID %d, got %d", tt.expectedID, tasks[0].ID)
			}
			if tasks[0].Name != tt.expectedName {
				t.Errorf("Expected name '%s', got '%s'", tt.expectedName, tasks[0].Name)
			}
			if tasks[0].Status != task.TASK_STATUS_TODO {
				t.Errorf("Expected status TODO got %s", tasks[0].Status)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	tests := []struct {
		name           string
		taskID         int
		newName        string
		newStatus      task.TaskStatus
		expectedName   string
		expectedStatus task.TaskStatus
		shouldUpdate   bool
	}{
		{"Update task name", 1, "Buy coffee", "", "Buy coffee", task.TASK_STATUS_TODO, true},
		{"Update task status", 1, "", task.TASK_STATUS_DONE, "Test task", task.TASK_STATUS_DONE, true},
		{"Update non-existent task", 999, "New name", "", "Test task", task.TASK_STATUS_TODO, false},
		{"Update name and status", 1, "Complete project", task.TASK_STATUS_IN_PROGRESS, "Complete project", task.TASK_STATUS_IN_PROGRESS, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Remove(filePath)
			task.AddTask("Test task")
			task.UpdateTask(tt.taskID, tt.newName, tt.newStatus)

			var tasks []task.Task
			task.LoadTasks(&tasks)

			if tt.shouldUpdate {
				if tasks[0].Name != tt.expectedName {
					t.Errorf("Expected name '%s', got '%s'", tt.expectedName, tasks[0].Name)
				}
				if tasks[0].Status != tt.expectedStatus {
					t.Errorf("Expected status %s, got %s", tt.expectedStatus, tasks[0].Status)
				}
			} else {
				if tasks[0].Name == tt.newName {
					t.Error("Task should not have been updated")
				}
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name          string
		deleteID      int
		expectedCount int
		shouldDelete  bool
	}{
		{"Delete existing task", 1, 1, true},
		{"Delete non-existent task", 999, 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Remove(filePath)
			task.AddTask("Task 1")
			task.AddTask("Task 2")
			task.DeleteTask(tt.deleteID)

			var tasks []task.Task
			task.LoadTasks(&tasks)

			if len(tasks) != tt.expectedCount {
				t.Errorf("Expected %d tasks, got %d", tt.expectedCount, len(tasks))
			}
		})
	}
}

func TestListTask(t *testing.T) {
	tests := []struct {
		name              string
		filter            task.TaskStatus
		setupTasks        func()
		expectedDoneCount int
	}{
		{
			"List all tasks",
			"",
			func() {
				task.AddTask("Task 1")
				task.AddTask("Task 2")
			},
			0,
		},
		{
			"List DONE tasks",
			task.TASK_STATUS_DONE,
			func() {
				task.AddTask("Task 1")
				task.AddTask("Task 2")
				task.UpdateTask(1, "", task.TASK_STATUS_DONE)
			},
			1,
		},
		{
			"List TODO tasks",
			task.TASK_STATUS_TODO,
			func() {
				task.AddTask("Task 1")
				task.AddTask("Task 2")
				task.UpdateTask(1, "", task.TASK_STATUS_DONE)
			},
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Remove(filePath)
			tt.setupTasks()
			result := task.ListTasks(tt.filter)

			if result == "" {
				t.Error("Expected non-empty result")
			}

			var tasks []task.Task
			task.LoadTasks(&tasks)

			doneCount := 0
			if tt.filter == task.TASK_STATUS_DONE {
				for _, t := range tasks {
					if t.Status == task.TASK_STATUS_DONE {
						doneCount++
					}
				}
			}
			if doneCount != tt.expectedDoneCount {
				t.Errorf("Expected %d DONE tasks, got %d", tt.expectedDoneCount, doneCount)
			}
		})
	}
}
