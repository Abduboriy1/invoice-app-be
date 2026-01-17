// internal/infrastructure/database/postgres/user_repository.go
package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/invoice-app-be/internal/domain/user"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	query := `
        INSERT INTO users (id, email, password_hash, full_name, company_name, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := r.db.ExecContext(ctx, query, u.ID, u.Email, u.PasswordHash, u.FullName, u.CompanyName, u.CreatedAt, u.UpdatedAt)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var u user.User
	query := `SELECT id, email, password_hash, full_name, company_name, created_at, updated_at FROM users WHERE id = $1`
	if err := r.db.GetContext(ctx, &u, query, id); err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	query := `SELECT id, email, password_hash, full_name, company_name, created_at, updated_at FROM users WHERE email = $1`
	if err := r.db.GetContext(ctx, &u, query, email); err != nil {
		return nil, fmt.Errorf("getting user by email: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	query := `
        UPDATE users 
        SET email = $2, password_hash = $3, full_name = $4, company_name = $5, updated_at = $6
        WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, query, u.ID, u.Email, u.PasswordHash, u.FullName, u.CompanyName, u.UpdatedAt)
	return err
}
