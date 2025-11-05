package domain

import "time"

type TaskRepository interface {
	Save(task *Task) error
	FindByID(id, userID string) (*Task, error)
	List(userID string) ([]*Task, error)
	Update(task *Task) error
	Delete(id string, timestamp time.Time) error
}
