package entity

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	email := "test@example.com"
	password := "password123"
	fullName := "John Doe"
	phoneNumber := "0812345678"
	birthday := "1990-01-15"

	user := NewUser(email, password, fullName, phoneNumber, birthday)

	if user.Email != email {
		t.Errorf("Email = %v, want %v", user.Email, email)
	}
	if user.Password != password {
		t.Errorf("Password = %v, want %v", user.Password, password)
	}
	if user.FullName != fullName {
		t.Errorf("FullName = %v, want %v", user.FullName, fullName)
	}
	if user.PhoneNumber != phoneNumber {
		t.Errorf("PhoneNumber = %v, want %v", user.PhoneNumber, phoneNumber)
	}
	if user.Birthday != birthday {
		t.Errorf("Birthday = %v, want %v", user.Birthday, birthday)
	}
	if user.ID != 0 {
		t.Errorf("ID should be 0 for new user, got %v", user.ID)
	}
	if user.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
}

func TestUser_IsValidEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{
			name:     "valid email",
			email:    "test@example.com",
			expected: true,
		},
		{
			name:     "valid email with subdomain",
			email:    "user@mail.example.com",
			expected: true,
		},
		{
			name:     "empty email",
			email:    "",
			expected: false,
		},
		{
			name:     "short email",
			email:    "a@b.c",
			expected: false,
		},
		{
			name:     "email with 6 characters",
			email:    "a@b.co",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Email: tt.email}
			result := user.IsValidEmail()
			if result != tt.expected {
				t.Errorf("IsValidEmail() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUser_WithoutPassword(t *testing.T) {
	originalUser := &User{
		ID:          1,
		Email:       "test@example.com",
		Password:    "secretpassword",
		FullName:    "John Doe",
		PhoneNumber: "0812345678",
		Birthday:    "1990-01-15",
		CreatedAt:   time.Now(),
	}

	userWithoutPassword := originalUser.WithoutPassword()

	// Check that password is empty
	if userWithoutPassword.Password != "" {
		t.Errorf("Password should be empty, got %v", userWithoutPassword.Password)
	}

	// Check that other fields remain the same
	if userWithoutPassword.ID != originalUser.ID {
		t.Errorf("ID = %v, want %v", userWithoutPassword.ID, originalUser.ID)
	}
	if userWithoutPassword.Email != originalUser.Email {
		t.Errorf("Email = %v, want %v", userWithoutPassword.Email, originalUser.Email)
	}
	if userWithoutPassword.FullName != originalUser.FullName {
		t.Errorf("FullName = %v, want %v", userWithoutPassword.FullName, originalUser.FullName)
	}
	if userWithoutPassword.PhoneNumber != originalUser.PhoneNumber {
		t.Errorf("PhoneNumber = %v, want %v", userWithoutPassword.PhoneNumber, originalUser.PhoneNumber)
	}
	if userWithoutPassword.Birthday != originalUser.Birthday {
		t.Errorf("Birthday = %v, want %v", userWithoutPassword.Birthday, originalUser.Birthday)
	}
	if !userWithoutPassword.CreatedAt.Equal(originalUser.CreatedAt) {
		t.Errorf("CreatedAt = %v, want %v", userWithoutPassword.CreatedAt, originalUser.CreatedAt)
	}

	// Ensure original user is not modified
	if originalUser.Password == "" {
		t.Error("Original user password should not be modified")
	}
}

func TestUserStruct(t *testing.T) {
	now := time.Now()
	user := &User{
		ID:          123,
		Email:       "user@test.com",
		Password:    "hashedpassword",
		FullName:    "Test User",
		PhoneNumber: "0812345678",
		Birthday:    "1985-12-25",
		CreatedAt:   now,
	}

	// Test all fields are set correctly
	if user.ID != 123 {
		t.Errorf("ID = %v, want %v", user.ID, 123)
	}
	if user.Email != "user@test.com" {
		t.Errorf("Email = %v, want %v", user.Email, "user@test.com")
	}
	if user.Password != "hashedpassword" {
		t.Errorf("Password = %v, want %v", user.Password, "hashedpassword")
	}
	if user.FullName != "Test User" {
		t.Errorf("FullName = %v, want %v", user.FullName, "Test User")
	}
	if user.PhoneNumber != "0812345678" {
		t.Errorf("PhoneNumber = %v, want %v", user.PhoneNumber, "0812345678")
	}
	if user.Birthday != "1985-12-25" {
		t.Errorf("Birthday = %v, want %v", user.Birthday, "1985-12-25")
	}
	if !user.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", user.CreatedAt, now)
	}
}
