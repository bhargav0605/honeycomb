package tasks

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/bhargavparmar/hive-demo/pkg/metrics"
	"github.com/bhargavparmar/hive-demo/pkg/storage"
	"github.com/cilium/hive/cell"
)

// Cell provides task management business logic
var Cell = cell.Module(
	"tasks",
	"Task Management",

	cell.Provide(newTaskManager),
)

// Task represents a task in the system
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TaskManager manages tasks
type TaskManager interface {
	Create(title, description string) (*Task, error)
	Get(id string) (*Task, error)
	List() []*Task
	Update(id string, title, description, status string) (*Task, error)
	Delete(id string) error
	GetStats() map[string]interface{}
}

type taskManager struct {
	logger  *slog.Logger
	storage storage.Storage
	metrics metrics.Metrics
}

// newTaskManager creates a new task manager with dependencies
func newTaskManager(lc cell.Lifecycle, logger *slog.Logger, storage storage.Storage, metrics metrics.Metrics) TaskManager {
	tm := &taskManager{
		logger:  logger.With("component", "task-manager"),
		storage: storage,
		metrics: metrics,
	}

	lc.Append(cell.Hook{
		OnStart: func(ctx cell.HookContext) error {
			tm.logger.Info("Task manager started")
			return nil
		},
		OnStop: func(ctx cell.HookContext) error {
			count := tm.storage.Count()
			tm.logger.Info("Task manager stopping", "active_tasks", count)
			return nil
		},
	})

	return tm
}

func (tm *taskManager) Create(title, description string) (*Task, error) {
	if title == "" {
		tm.metrics.IncrementErrors()
		return nil, errors.New("title is required")
	}

	task := &Task{
		ID:          fmt.Sprintf("task-%d", time.Now().UnixNano()),
		Title:       title,
		Description: description,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tm.storage.Set(task.ID, task)
	tm.logger.Info("Task created", "id", task.ID, "title", task.Title)

	return task, nil
}

func (tm *taskManager) Get(id string) (*Task, error) {
	val, ok := tm.storage.Get(id)
	if !ok {
		tm.metrics.IncrementErrors()
		return nil, errors.New("task not found")
	}

	task, ok := val.(*Task)
	if !ok {
		tm.metrics.IncrementErrors()
		return nil, errors.New("invalid task data")
	}

	return task, nil
}

func (tm *taskManager) List() []*Task {
	all := tm.storage.List()
	tasks := make([]*Task, 0, len(all))

	for _, val := range all {
		if task, ok := val.(*Task); ok {
			tasks = append(tasks, task)
		}
	}

	return tasks
}

func (tm *taskManager) Update(id string, title, description, status string) (*Task, error) {
	task, err := tm.Get(id)
	if err != nil {
		return nil, err
	}

	if title != "" {
		task.Title = title
	}
	if description != "" {
		task.Description = description
	}
	if status != "" {
		task.Status = status
	}
	task.UpdatedAt = time.Now()

	tm.storage.Set(id, task)
	tm.logger.Info("Task updated", "id", task.ID)

	return task, nil
}

func (tm *taskManager) Delete(id string) error {
	_, err := tm.Get(id)
	if err != nil {
		return err
	}

	tm.storage.Delete(id)
	tm.logger.Info("Task deleted", "id", id)

	return nil
}

func (tm *taskManager) GetStats() map[string]interface{} {
	tasks := tm.List()
	stats := map[string]interface{}{
		"total_tasks":    len(tasks),
		"total_requests": tm.metrics.GetRequests(),
		"total_errors":   tm.metrics.GetErrors(),
	}

	// Count by status
	statusCount := make(map[string]int)
	for _, task := range tasks {
		statusCount[task.Status]++
	}
	stats["by_status"] = statusCount

	return stats
}
