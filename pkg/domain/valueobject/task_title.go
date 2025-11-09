package valueobject

import (
	"errors"
	"strings"
)

type TaskTitle struct {
	value string
}

var (
	ErrEmptyTitle   = errors.New("title cannot be empty")
	ErrTitleTooLong = errors.New("title exceeds 100 characters")
)

func NewTaskTitle(raw string) (TaskTitle, error) {
	title := strings.TrimSpace(raw)
	title = strings.Join(strings.Fields(title), " ")

	if len(title) == 0 {
		return TaskTitle{}, ErrEmptyTitle
	}
	if len([]rune(title)) > 100 {
		return TaskTitle{}, ErrTitleTooLong
	}

	return TaskTitle{value: title}, nil
}

func (t TaskTitle) String() string {
	return t.value
}
