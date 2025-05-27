package test

import (
	"testing"
	"user-service/internal/usecase"
)

func TestHashPasswordAndCheck(t *testing.T) {
	originalPassword := "securePa$$word123"

	hashedPassword, err := usecase.HashPassword(originalPassword)
	if err != nil {
		t.Fatalf("HashPassword() error: %v", err)
	}

	if hashedPassword == originalPassword {
		t.Error("Hashed password should not be equal to the original password")
	}

	if !usecase.CheckPasswordHash(originalPassword, hashedPassword) {
		t.Error("CheckPasswordHash() failed: expected password to match hash")
	}

	wrongPassword := "wrongPassword"
	if usecase.CheckPasswordHash(wrongPassword, hashedPassword) {
		t.Error("CheckPasswordHash() failed: expected password NOT to match hash")
	}
}
