package main

import (
	"os"
	"testing"
)

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
			fileName := "test_add_task.json"
			defer os.Remove(fileName)

			addTask(tt.taskName, fileName)

			var tasks []Task
			loadTasks(fileName, &tasks)

			if len(tasks) != 1 {
				t.Errorf("Expected 1 task, got %d", len(tasks))
			}
			if tasks[0].ID != tt.expectedID {
				t.Errorf("Expected ID %d, got %d", tt.expectedID, tasks[0].ID)
			}
			if tasks[0].Name != tt.expectedName {
				t.Errorf("Expected name '%s', got '%s'", tt.expectedName, tasks[0].Name)
			}
			if tasks[0].Status != TODO {
				t.Errorf("Expected status TODO, got %s", tasks[0].Status)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	tests := []struct {
		name           string
		taskID         int
		newName        string
		newStatus      TaskStatus
		expectedName   string
		expectedStatus TaskStatus
		shouldUpdate   bool
	}{
		{"Update task name", 1, "Buy coffee", "", "Buy coffee", TODO, true},
		{"Update task status", 1, "", DONE, "Test task", DONE, true},
		{"Update non-existent task", 999, "New name", "", "Test task", TODO, false},
		{"Update name and status", 1, "Complete project", IN_PROGRESS, "Complete project", IN_PROGRESS, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileName := "test_update_task.json"
			defer os.Remove(fileName)

			addTask("Test task", fileName)
			updateTask(tt.taskID, tt.newName, tt.newStatus, fileName)

			var tasks []Task
			loadTasks(fileName, &tasks)

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
			fileName := "test_delete_task.json"
			defer os.Remove(fileName)

			addTask("Task 1", fileName)
			addTask("Task 2", fileName)
			deleteTask(tt.deleteID, fileName)

			var tasks []Task
			loadTasks(fileName, &tasks)

			if len(tasks) != tt.expectedCount {
				t.Errorf("Expected %d tasks, got %d", tt.expectedCount, len(tasks))
			}
		})
	}
}

func TestListTask(t *testing.T) {
	tests := []struct {
		name              string
		filter            TaskStatus
		setupTasks        func(string)
		expectedDoneCount int
	}{
		{
			"List all tasks",
			"",
			func(f string) {
				addTask("Task 1", f)
				addTask("Task 2", f)
			},
			0,
		},
		{
			"List DONE tasks",
			DONE,
			func(f string) {
				addTask("Task 1", f)
				addTask("Task 2", f)
				updateTask(1, "", DONE, f)
			},
			1,
		},
		{
			"List TODO tasks",
			TODO,
			func(f string) {
				addTask("Task 1", f)
				addTask("Task 2", f)
				updateTask(1, "", DONE, f)
			},
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileName := "test_list_task.json"
			defer os.Remove(fileName)

			tt.setupTasks(fileName)
			result := listTasks(tt.filter, fileName)

			if result == "" {
				t.Error("Expected non-empty result")
			}

			var tasks []Task
			loadTasks(fileName, &tasks)

			doneCount := 0
			if tt.filter == DONE {
				for _, task := range tasks {
					if task.Status == DONE {
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
