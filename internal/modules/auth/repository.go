package auth

import (
	"context"
	"fmt"

	"vendix/internal/database"

	"github.com/google/uuid"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, schema string, user *User) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.users (id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, schema)

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Role,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (r *Repository) GetUserByEmail(ctx context.Context, schema, email string) (*User, error) {
	var user User
	query := fmt.Sprintf(`
		SELECT id, email, password_hash, first_name, last_name, role, is_active, last_login_at, created_at, updated_at
		FROM %s.users
		WHERE email = $1
	`, schema)

	err := r.db.GetContext(ctx, &user, query, email)
	return &user, err
}

func (r *Repository) GetUserByID(ctx context.Context, schema string, id uuid.UUID) (*User, error) {
	var user User
	query := fmt.Sprintf(`
		SELECT id, email, password_hash, first_name, last_name, role, is_active, last_login_at, created_at, updated_at
		FROM %s.users
		WHERE id = $1
	`, schema)

	err := r.db.GetContext(ctx, &user, query, id)
	return &user, err
}

func (r *Repository) UpdateLastLogin(ctx context.Context, schema string, userID uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s.users
		SET last_login_at = NOW()
		WHERE id = $1
	`, schema)

	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *Repository) CreateRefreshToken(ctx context.Context, schema string, token *RefreshToken) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.refresh_tokens (id, user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, schema)

	_, err := r.db.ExecContext(ctx, query,
		token.ID,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
	)

	return err
}

func (r *Repository) GetRefreshToken(ctx context.Context, schema, token string) (*RefreshToken, error) {
	var refreshToken RefreshToken
	query := fmt.Sprintf(`
		SELECT id, user_id, token, expires_at, created_at
		FROM %s.refresh_tokens
		WHERE token = $1
	`, schema)

	err := r.db.GetContext(ctx, &refreshToken, query, token)
	return &refreshToken, err
}

func (r *Repository) DeleteRefreshToken(ctx context.Context, schema, token string) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.refresh_tokens
		WHERE token = $1
	`, schema)

	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

func (r *Repository) ListUsers(ctx context.Context, schema string) ([]*User, error) {
	var users []*User
	query := fmt.Sprintf(`
		SELECT id, email, password_hash, first_name, last_name, role, is_active, last_login_at, created_at, updated_at
		FROM %s.users
		ORDER BY created_at DESC
	`, schema)

	err := r.db.SelectContext(ctx, &users, query)
	return users, err
}
