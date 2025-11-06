package sqlite

import (
	"database/sql"
	"time"

	domain "github.com/hoyci/todo-ddd/pkg/domain/user"
	_ "modernc.org/sqlite"
)

type SQLiteUserRepository struct {
	db *sql.DB
	tx *sql.Tx
}

type SQLExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

func NewSQLiteUserRepositoryWithTx(tx *sql.Tx) *SQLiteUserRepository {
	return &SQLiteUserRepository{tx: tx}
}

func (r *SQLiteUserRepository) getExecutor() SQLExecutor {
	if r.tx != nil {
		return r.tx
	}

	return r.db
}

// ------------------- CREATE -------------------
func (r *SQLiteUserRepository) Save(user domain.User) error {
	_, err := r.getExecutor().Exec(`
		INSERT INTO users (id, name, email, created_at, updated_at, deleted_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt, user.DeletedAt)
	return err
}

// ------------------- READ -------------------
func (r *SQLiteUserRepository) FindByID(id string) (*domain.User, error) {
	row := r.getExecutor().QueryRow(`
		SELECT id, name, email, created_at, updated_at, deleted_at 
		FROM users WHERE id = ?`, id)

	u := &domain.User{}
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *SQLiteUserRepository) FindByEmail(email string) (*domain.User, error) {
	row := r.getExecutor().QueryRow(`
		SELECT id, name, email, created_at, updated_at, deleted_at 
		FROM users WHERE email = ?`, email)

	u := &domain.User{}
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// ------------------- LIST -------------------
func (r *SQLiteUserRepository) List() ([]*domain.User, error) {
	rows, err := r.getExecutor().Query(`
		SELECT id, name, email, created_at, updated_at, deleted_at 
		FROM users WHERE deleted_at IS NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u := &domain.User{}
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// ------------------- UPDATE -------------------
func (r *SQLiteUserRepository) Update(user domain.User) error {
	_, err := r.getExecutor().Exec(`
		UPDATE users 
		SET name = ?, email = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL`,
		user.Name, user.Email, user.UpdatedAt, user.ID)
	return err
}

// ------------------- DELETE (SOFT) -------------------
func (r *SQLiteUserRepository) Delete(id string, timestamp time.Time) error {
	_, err := r.getExecutor().Exec(`
		UPDATE users 
		SET deleted_at = ?
		WHERE id = ? AND deleted_at IS NULL`,
		timestamp, id)
	return err
}
