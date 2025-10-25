package auth

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the tenant schema
type User struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	Email        string     `db:"email" json:"email"`
	PasswordHash string     `db:"password_hash" json:"-"`
	FirstName    *string    `db:"first_name" json:"first_name,omitempty"`
	LastName     *string    `db:"last_name" json:"last_name,omitempty"`
	Role         string     `db:"role" json:"role"`
	IsActive     bool       `db:"is_active" json:"is_active"`
	LastLoginAt  *time.Time `db:"last_login_at" json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}

// RefreshToken represents a refresh token
type RefreshToken struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	Token     string    `db:"token" json:"token"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=8"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Role      string  `json:"role,omitempty"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	TokenType    string        `json:"token_type"`
	ExpiresIn    int           `json:"expires_in"`
	User         *UserResponse `json:"user"`
}

// UserResponse represents a user response
type UserResponse struct {
	ID          string     `json:"id"`
	Email       string     `json:"email"`
	FirstName   *string    `json:"first_name,omitempty"`
	LastName    *string    `json:"last_name,omitempty"`
	Role        string     `json:"role"`
	IsActive    bool       `json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:          u.ID.String(),
		Email:       u.Email,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Role:        u.Role,
		IsActive:    u.IsActive,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
