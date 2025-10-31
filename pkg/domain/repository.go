package domain

import "time"

type TaskRepository interface {
	Save(task *Task) error
	FindByID(id string) (*Task, error)
	List() ([]*Task, error)
	Update(task *Task) error
	Delete(id string, timestamp time.Time) error
}
