package sqlite

import (
	"database/sql"
	"time"

	domain "github.com/hoyci/todo-ddd/pkg/domain/task"
	_ "modernc.org/sqlite"
)

type SQLiteTaskRepository struct {
	db *sql.DB
	tx *sql.Tx
}

func NewSQLiteTaskRepository(db *sql.DB) *SQLiteTaskRepository {
	return &SQLiteTaskRepository{db: db}
}

func NewSQLiteTaskRepositoryWithTx(tx *sql.Tx) *SQLiteTaskRepository {
	return &SQLiteTaskRepository{tx: tx}
}

func (r *SQLiteTaskRepository) getExecutor() SQLExecutor {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *SQLiteTaskRepository) Save(task *domain.Task) error {
	_, err := r.getExecutor().Exec(`
		INSERT INTO tasks (id, title, description, priority, status, created_at, updated_at, deleted_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		task.ID, task.Title, task.Description, task.Priority, task.Status, task.CreatedAt, task.UpdatedAt, task.DeletedAt)
	return err
}

func (r *SQLiteTaskRepository) Update(task *domain.Task) error {
	_, err := r.getExecutor().Exec(`
		UPDATE tasks 
		SET title = ?, 
			description = ?, 
			priority = ?, 
			status = ?, 
			updated_at = ? 
		WHERE id = ? AND deleted_at IS NULL`,
		task.Title, task.Description, task.Priority, task.Status, task.UpdatedAt, task.ID)
	return err
}

func (r *SQLiteTaskRepository) FindByID(id, userID string) (*domain.Task, error) {
	row := r.getExecutor().QueryRow(`SELECT id, title, description, priority, status, created_at, updated_at, deleted_at FROM tasks t WHERE t.id = ? AND t.user_id = ?`, id, userID)
	task := &domain.Task{}
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.Status, &task.CreatedAt, &task.UpdatedAt, &task.DeletedAt)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *SQLiteTaskRepository) List(userID string) ([]*domain.Task, error) {
	rows, err := r.getExecutor().Query(`SELECT id, title, description, priority, status, created_at, updated_at FROM tasks t WHERE t.user_id = ? AND t.deleted_at is null`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		t := &domain.Task{}
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Priority, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *SQLiteTaskRepository) Delete(id string, timestamp time.Time) error {
	query := `
		UPDATE tasks 
		SET deleted_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`
	_, err := r.getExecutor().Exec(query, timestamp, id)
	return err
}
