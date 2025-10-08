package utils

import (
	"testing"
)

func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	password := "mySecret123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hash == "" {
		t.Fatal("expected a non-empty hash")
	}

	if !CheckPasswordHash(password, hash) {
		t.Error("expected password to match the hash")
	}

	if CheckPasswordHash("wrongPassword", hash) {
		t.Error("expected wrong password to not match the hash")
	}
}

// Test GenerateToken
func TestGenerateToken(t *testing.T) {
	token := GenerateToken()

	if len(token) != 32 {
		t.Errorf("expected token length 32, got %d", len(token))
	}

	token2 := GenerateToken()
	if token == token2 {
		t.Error("expected different tokens, but got the same")
	}
}
