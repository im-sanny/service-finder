package repository

import (
	"database/sql"
	"fmt"

	"github.com/im-sanny/service-finder/model"
)

type UserRepository interface {
	Create(s *model.User) error
	GetPasswordHashByEmail(email string) (string, error)
}

type usr struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &usr{db: db}
}

func (r *usr) Create(u *model.User) error {
	err := r.db.QueryRow(`
	INSERT INTO users (name, email, password_hash)
	VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`,
		u.Name, u.Email, u.PasswordHash).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}
