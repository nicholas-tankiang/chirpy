package auth_test

import (
	"chirpy/internal/auth"
	"testing"
)

func TestHashAndPassword(t *testing.T) {
	password := "VerySecurePassword"
	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	err = auth.CheckPasswordHash(hash, password)
	if err != nil {
		t.Fatalf("Password is incorrect: %v", err)
	}
}
