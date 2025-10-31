package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/hoyci/todo-ddd/pkg/domain"
	_ "modernc.org/sqlite"
)

type SQLiteTaskRepository struct {
	db *sql.DB
}

func NewSQLiteTaskRepository(db *sql.DB) *SQLiteTaskRepository {
	return &SQLiteTaskRepository{db: db}
}

func (r *SQLiteTaskRepository) Save(task *domain.Task) error {
	_, err := r.db.Exec(`
		INSERT INTO tasks (id, title, description, priority, status, created_at, updated_at, deleted_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		task.ID, task.Title, task.Description, task.Priority, task.Status, task.CreatedAt, task.UpdatedAt, task.DeletedAt)
	return err
}

func (r *SQLiteTaskRepository) Update(task *domain.Task) error {
	_, err := r.db.Exec(`
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

func (r *SQLiteTaskRepository) FindByID(id string) (*domain.Task, error) {
	row := r.db.QueryRow(`SELECT id, title, description, priority, status, created_at, updated_at, deleted_at FROM tasks WHERE id = ?`, id)
	task := &domain.Task{}
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.Status, &task.CreatedAt, &task.UpdatedAt, &task.DeletedAt)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *SQLiteTaskRepository) List() ([]*domain.Task, error) {
	rows, err := r.db.Query(`SELECT id, title, description, priority, status, created_at, updated_at FROM tasks t WHERE t.deleted_at is null`)
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
	_, err := r.db.Exec(query, timestamp, id)
	return err
}

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./data/tasks.db")
	if err != nil {
		return nil, err
	}
	schema := `
	CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		priority INTEGER,
		status TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP,
		deleted_at TIMESTAMP
	);`
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("create schema: %w", err)
	}
	return db, nil
}
