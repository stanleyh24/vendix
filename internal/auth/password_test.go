package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hash == password {
		t.Error("Hash should not equal password")
	}

	if len(hash) == 0 {
		t.Error("Hash should not be empty")
	}
}

func TestComparePassword(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Test correct password
	err = ComparePassword(hash, password)
	if err != nil {
		t.Error("Password should match hash")
	}

	// Test incorrect password
	err = ComparePassword(hash, "wrongpassword")
	if err == nil {
		t.Error("Wrong password should not match hash")
	}
}
