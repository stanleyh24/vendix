package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/stanleyh24/vendix/internal/auth"
	"github.com/stanleyh24/vendix/internal/config"
	"github.com/stanleyh24/vendix/internal/database"
	"github.com/stanleyh24/vendix/internal/logger"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
	cfg  *config.Config
}

func NewService(db *database.DB, cfg *config.Config) *Service {
	return &Service{
		repo: NewRepository(db),
		cfg:  cfg,
	}
}

func (s *Service) Register(ctx context.Context, schema string, req *RegisterRequest) (*AuthResponse, error) {
	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         req.Role,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if user.Role == "" {
		user.Role = "viewer" // Default role
	}

	if err := s.repo.CreateUser(ctx, schema, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	logger.Info("User registered", "email", user.Email, "schema", schema)

	// Generate tokens
	return s.generateAuthResponse(ctx, schema, user)
}

func (s *Service) Login(ctx context.Context, schema string, req *LoginRequest) (*AuthResponse, error) {
	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, schema, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	// Verify password
	if err := auth.ComparePassword(user.PasswordHash, req.Password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Update last login
	user.LastLoginAt = timePtr(time.Now())
	if err := s.repo.UpdateLastLogin(ctx, schema, user.ID); err != nil {
		logger.Error("Failed to update last login", "error", err)
	}

	logger.Info("User logged in", "email", user.Email, "schema", schema)

	// Generate tokens
	return s.generateAuthResponse(ctx, schema, user)
}

func (s *Service) RefreshToken(ctx context.Context, schema string, refreshToken string) (*AuthResponse, error) {
	// Validate refresh token
	token, err := s.repo.GetRefreshToken(ctx, schema, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Check if token is expired
	if token.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("refresh token expired")
	}

	// Get user
	user, err := s.repo.GetUserByID(ctx, schema, token.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	// Delete old refresh token
	if err := s.repo.DeleteRefreshToken(ctx, schema, refreshToken); err != nil {
		logger.Error("Failed to delete old refresh token", "error", err)
	}

	// Generate new tokens
	return s.generateAuthResponse(ctx, schema, user)
}

func (s *Service) Logout(ctx context.Context, schema string, refreshToken string) error {
	return s.repo.DeleteRefreshToken(ctx, schema, refreshToken)
}

func (s *Service) GetUserByID(ctx context.Context, schema string, userID string) (*UserResponse, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := s.repo.GetUserByID(ctx, schema, id)
	if err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

func (s *Service) generateAuthResponse(ctx context.Context, schema string, user *User) (*AuthResponse, error) {
	// Generate access token
	accessToken, err := auth.GenerateAccessToken(
		user.ID.String(),
		user.Email,
		user.Role,
		s.cfg.JWTSecret,
		s.cfg.JWTAccessExpiry,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshTokenString := auth.GenerateRefreshToken()

	// Store refresh token
	refreshToken := &RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(s.cfg.JWTRefreshExpiry),
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateRefreshToken(ctx, schema, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.cfg.JWTAccessExpiry.Seconds()),
		User:         user.ToResponse(),
	}, nil
}

func timePtr(t time.Time) *time.Time {
	return &t
}
