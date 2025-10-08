package utils_test

import (
	"loopit/internal/utils"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndValidateJWT(t *testing.T) {
	userID := 123
	role := "admin"

	token, err := utils.GenerateJWT(userID, role)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected a non-empty token")
	}

	claims, err := utils.ValidateJWT(token)
	if err != nil {
		t.Fatalf("expected no error while validating, got %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected userID %d, got %d", userID, claims.UserID)
	}
	if claims.Role != role {
		t.Errorf("expected role %s, got %s", role, claims.Role)
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	invalidToken := "invalid.token.string"

	_, err := utils.ValidateJWT(invalidToken)
	if err == nil {
		t.Error("expected error for invalid token, got nil")
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {

	expiredClaims := &utils.Claims{
		UserID: 1,
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	tokenStr, err := token.SignedString([]byte("secret_key"))
	if err != nil {
		t.Fatalf("failed to sign expired token: %v", err)
	}

	_, err = utils.ValidateJWT(tokenStr)
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}
