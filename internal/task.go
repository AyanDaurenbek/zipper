package internal

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
)

type Task struct {
	ID      string
	Links   []string
	Status  TaskStatus
	Archive string
	Errors  []string
	mu      sync.Mutex
}

type TaskManager struct {
	tasks     map[string]*Task
	tasksMu   sync.Mutex
	activeSem chan struct{}
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks:     make(map[string]*Task),
		activeSem: make(chan struct{}, 3),
	}
}

func (m *TaskManager) CreateTask() (*Task, error) {
	m.tasksMu.Lock()
	defer m.tasksMu.Unlock()

	if len(m.activeSem) >= 3 {
		return nil, errors.New("сервер занят, попробуйте позже")
	}

	id := uuid.New().String()
	task := &Task{
		ID:     id,
		Status: StatusPending,
		Links:  []string{},
	}
	m.tasks[id] = task

	m.activeSem <- struct{}{}

	return task, nil
}

func (m *TaskManager) AddLink(taskID string, link string) error {
	task, ok := m.tasks[taskID]
	if !ok {
		return errors.New("задача не найдена")
	}

	task.mu.Lock()
	defer task.mu.Unlock()

	if task.Status != StatusPending {
		return errors.New("задача уже завершена")
	}
	if len(task.Links) >= 3 {
		return errors.New("нельзя добавить больше 3 ссылок")
	}

	task.Links = append(task.Links, link)

	if len(task.Links) == 3 {
		go m.buildArchive(task)
	}
	return nil
}

func (m *TaskManager) GetTask(id string) (*Task, bool) {
	m.tasksMu.Lock()
	defer m.tasksMu.Unlock()
	task, ok := m.tasks[id]
	return task, ok
}

func (m *TaskManager) buildArchive(task *Task) {
	task.mu.Lock()
	task.Status = StatusPending
	task.mu.Unlock()

	archivePath, badLinks, err := DownloadAndZip(task.ID, task.Links)
	task.mu.Lock()
	defer task.mu.Unlock()

	if err != nil {
		task.Status = StatusFailed
		task.Errors = append(task.Errors, err.Error())
	} else {
		task.Status = StatusCompleted
		task.Archive = archivePath
		task.Errors = badLinks
	}

	<-m.activeSem
}
