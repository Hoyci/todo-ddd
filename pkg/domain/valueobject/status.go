package valueobject

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
)
