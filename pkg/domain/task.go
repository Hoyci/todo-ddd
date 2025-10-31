package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/hoyci/todo-ddd/pkg/domain/valueobject"
)

type Task struct {
	ID          string
	Title       string
	Description string
	Priority    valueobject.Priority
	Status      valueobject.Status
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}

func NewTask(title, description string, priority valueobject.Priority) *Task {
	return &Task{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Priority:    priority,
		Status:      valueobject.StatusNew,
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
		DeletedAt:   nil,
	}
}

func (t *Task) setStatus(status valueobject.Status) {
	now := time.Now()
	t.Status = status
	t.UpdatedAt = &now
}

func (t *Task) SetAsNew()      { t.setStatus(valueobject.StatusNew) }
func (t *Task) SetInProgress() { t.setStatus(valueobject.StatusInProgress) }
func (t *Task) SetCompleted()  { t.setStatus(valueobject.StatusCompleted) }

func (t *Task) Update(title, description string, priority valueobject.Priority) {
	now := time.Now()
	t.Title = title
	t.Description = description
	t.Priority = priority
	t.UpdatedAt = &now
}

func (t *Task) Delete() {
	now := time.Now()
	t.UpdatedAt = &now
	t.DeletedAt = &now
}
