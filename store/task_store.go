package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"stability-test-task-api/models"
)

const tasksFilePath = "data/tasks.json"

var (
	tasks = []models.Task{}
	mu    sync.RWMutex
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrNoChanges    = errors.New("no changes detected")
)

func init() {
	if err := loadTasks(); err != nil {
		panic(err)
	}
}

func loadTasks() error {
	if err := os.MkdirAll(filepath.Dir(tasksFilePath), 0o755); err != nil {
		return err
	}

	if _, err := os.Stat(tasksFilePath); os.IsNotExist(err) {
		mu.Lock()
		tasks = []models.Task{
			{ID: 1, Title: "Learn Go", Done: false},
			{ID: 2, Title: "Build API", Done: false},
		}
		err = saveTasksLocked()
		mu.Unlock()
		return err
	}

	data, err := os.ReadFile(tasksFilePath)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		mu.Lock()
		tasks = []models.Task{}
		mu.Unlock()
		return nil
	}

	loadedTasks := []models.Task{}
	if err := json.Unmarshal(data, &loadedTasks); err != nil {
		return err
	}

	mu.Lock()
	tasks = loadedTasks
	mu.Unlock()

	return nil
}

// digunakan untuk menyimpan tasks ke file, dipanggil setiap kali ada perubahan pada tasks
func saveTasksLocked() error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(tasksFilePath, data, 0o644)
}

func GetAllTasks() []models.Task {
	mu.RLock()
	defer mu.RUnlock()

	tasksCopy := make([]models.Task, len(tasks))
	copy(tasksCopy, tasks)

	return tasksCopy
}

func GetTaskByID(id int) *models.Task {
	mu.RLock()
	defer mu.RUnlock()

	for _, t := range tasks {
		if t.ID == id {
			taskCopy := t
			return &taskCopy
		}
	}
	return nil
}

func AddTask(task models.Task) {
	mu.Lock()
	defer mu.Unlock()

	tasks = append(tasks, task)
	_ = saveTasksLocked()
}

func DeleteTask(id int) {
	mu.Lock()
	defer mu.Unlock()

	for i, t := range tasks {
		if t.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			_ = saveTasksLocked()
			return
		}
	}
}

// endpoint baru untuk update
func UpdateTask(id int, title string, done bool) (*models.Task, error) {
	mu.Lock()
	defer mu.Unlock()

	for i, t := range tasks {
		if t.ID == id {
			if t.Title == title && t.Done == done {
				return nil, ErrNoChanges
			}
			tasks[i].Title = title
			tasks[i].Done = done
			_ = saveTasksLocked()
			return &tasks[i], nil
		}
	}
	return nil, ErrTaskNotFound
}
