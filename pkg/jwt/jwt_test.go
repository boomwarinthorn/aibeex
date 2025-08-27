package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestNewService(t *testing.T) {
	secretKey := "test-secret-key"
	service := NewService(secretKey)

	if service == nil {
		t.Fatal("NewService() returned nil")
	}
	if string(service.secretKey) != secretKey {
		t.Errorf("secretKey = %v, want %v", string(service.secretKey), secretKey)
	}
}

func TestService_GenerateToken(t *testing.T) {
	service := NewService("test-secret")
	userID := 123
	email := "test@example.com"

	token, expirationTime, err := service.GenerateToken(userID, email)

	// Check no error
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	// Check token is not empty
	if token == "" {
		t.Error("GenerateToken() returned empty token")
	}

	// Check expiration time is in the future (approximately 24 hours)
	expectedExpiration := time.Now().Add(24 * time.Hour)
	if expirationTime.Before(time.Now()) {
		t.Error("ExpirationTime should be in the future")
	}
	if expirationTime.After(expectedExpiration.Add(time.Minute)) {
		t.Error("ExpirationTime should be approximately 24 hours from now")
	}

	// Parse token to verify claims
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return service.secretKey, nil
	})
	if err != nil {
		t.Fatalf("Failed to parse generated token: %v", err)
	}

	claims, ok := parsedToken.Claims.(*Claims)
	if !ok {
		t.Fatal("Failed to cast claims")
	}

	if claims.UserID != userID {
		t.Errorf("UserID = %v, want %v", claims.UserID, userID)
	}
	if claims.Email != email {
		t.Errorf("Email = %v, want %v", claims.Email, email)
	}
}

func TestService_ValidateToken(t *testing.T) {
	service := NewService("test-secret")
	userID := 456
	email := "validate@example.com"

	// Generate a token first
	token, _, err := service.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token for validation test: %v", err)
	}

	// Test valid token
	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("UserID = %v, want %v", claims.UserID, userID)
	}
	if claims.Email != email {
		t.Errorf("Email = %v, want %v", claims.Email, email)
	}

	// Test invalid token
	invalidToken := "invalid.jwt.token"
	_, err = service.ValidateToken(invalidToken)
	if err == nil {
		t.Error("ValidateToken() should return error for invalid token")
	}

	// Test token with wrong secret
	wrongSecretService := NewService("wrong-secret")
	_, err = wrongSecretService.ValidateToken(token)
	if err == nil {
		t.Error("ValidateToken() should return error for token with wrong secret")
	}

	// Test empty token
	_, err = service.ValidateToken("")
	if err == nil {
		t.Error("ValidateToken() should return error for empty token")
	}
}

func TestService_ValidateToken_ExpiredToken(t *testing.T) {
	service := NewService("test-secret")
	
	// Create an expired token manually
	pastTime := time.Now().Add(-time.Hour) // 1 hour ago
	claims := &Claims{
		UserID: 789,
		Email:  "expired@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(pastTime),
			IssuedAt:  jwt.NewNumericDate(pastTime.Add(-time.Hour)),
			Subject:   "789",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(service.secretKey)
	if err != nil {
		t.Fatalf("Failed to create expired token: %v", err)
	}

	// Try to validate expired token
	_, err = service.ValidateToken(tokenString)
	if err == nil {
		t.Error("ValidateToken() should return error for expired token")
	}
}

func TestClaims(t *testing.T) {
	userID := 101
	email := "claims@example.com"
	now := time.Now()

	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   "101",
		},
	}

	if claims.UserID != userID {
		t.Errorf("UserID = %v, want %v", claims.UserID, userID)
	}
	if claims.Email != email {
		t.Errorf("Email = %v, want %v", claims.Email, email)
	}
	if claims.Subject != "101" {
		t.Errorf("Subject = %v, want %v", claims.Subject, "101")
	}
}

func TestService_GenerateAndValidateRoundtrip(t *testing.T) {
	tests := []struct {
		name      string
		secretKey string
		userID    int
		email     string
	}{
		{
			name:      "normal case",
			secretKey: "normal-secret",
			userID:    1,
			email:     "normal@example.com",
		},
		{
			name:      "long secret key",
			secretKey: "very-long-secret-key-with-many-characters-for-security",
			userID:    999,
			email:     "long@example.com",
		},
		{
			name:      "special characters in email",
			secretKey: "special-secret",
			userID:    42,
			email:     "user+test@sub.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(tt.secretKey)

			// Generate token
			token, expTime, err := service.GenerateToken(tt.userID, tt.email)
			if err != nil {
				t.Fatalf("GenerateToken() error = %v", err)
			}

			// Validate token
			claims, err := service.ValidateToken(token)
			if err != nil {
				t.Fatalf("ValidateToken() error = %v", err)
			}

			// Verify claims
			if claims.UserID != tt.userID {
				t.Errorf("UserID = %v, want %v", claims.UserID, tt.userID)
			}
			if claims.Email != tt.email {
				t.Errorf("Email = %v, want %v", claims.Email, tt.email)
			}

			// Verify expiration time matches
			if !claims.ExpiresAt.Time.Equal(expTime) {
				t.Errorf("ExpiresAt = %v, want %v", claims.ExpiresAt.Time, expTime)
			}
		})
	}
}
