package repository

import (
	"database/sql"
	"errors"
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

func (r *usr) GetPasswordHashByEmail(email string) (string, error) {
	var hash string

	err := r.db.QueryRow(`
		SELECT  password_hash
		FROM users WHERE email=$1`, email).Scan(&hash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user not found")
		}
		// FIX 3: Changed %d to %s because email is a string
		return "", fmt.Errorf("failed to get user by email %s: %w", email, err)
	}
	return hash, nil
}
