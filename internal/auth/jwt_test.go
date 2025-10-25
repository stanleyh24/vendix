package auth

import (
	"testing"
	"time"
)

func TestGenerateAccessToken(t *testing.T) {
	userID := "test-user-id"
	email := "test@example.com"
	role := "admin"
	secret := "test-secret"
	expiry := 15 * time.Minute

	token, err := GenerateAccessToken(userID, email, role, secret, expiry)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if len(token) == 0 {
		t.Error("Token should not be empty")
	}
}

func TestValidateAccessToken(t *testing.T) {
	userID := "test-user-id"
	email := "test@example.com"
	role := "admin"
	secret := "test-secret"
	expiry := 15 * time.Minute

	token, err := GenerateAccessToken(userID, email, role, secret, expiry)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Test valid token
	claims, err := ValidateAccessToken(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("Expected email %s, got %s", email, claims.Email)
	}

	if claims.Role != role {
		t.Errorf("Expected role %s, got %s", role, claims.Role)
	}

	// Test invalid secret
	_, err = ValidateAccessToken(token, "wrong-secret")
	if err == nil {
		t.Error("Token should not validate with wrong secret")
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	token := GenerateRefreshToken()

	if len(token) == 0 {
		t.Error("Refresh token should not be empty")
	}

	// Generate another token and ensure they're different
	token2 := GenerateRefreshToken()
	if token == token2 {
		t.Error("Refresh tokens should be unique")
	}
}
