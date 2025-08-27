package usecase

import (
	"errors"
	"testing"

	"fiber-hello-world/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

// Mock repository for testing
type MockUserRepository struct {
	users  map[string]*entity.User
	nextID int
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:  make(map[string]*entity.User),
		nextID: 1,
	}
}

func (m *MockUserRepository) Create(user *entity.User) (*entity.User, error) {
	if _, exists := m.users[user.Email]; exists {
		return nil, errors.New("user already exists")
	}
	
	user.ID = m.nextID
	m.nextID++
	m.users[user.Email] = user
	return user, nil
}

func (m *MockUserRepository) GetByEmail(email string) (*entity.User, error) {
	if user, exists := m.users[email]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) GetByID(id int) (*entity.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) Update(user *entity.User) error {
	if _, exists := m.users[user.Email]; !exists {
		return errors.New("user not found")
	}
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) Delete(id int) error {
	for email, user := range m.users {
		if user.ID == id {
			delete(m.users, email)
			return nil
		}
	}
	return errors.New("user not found")
}

func TestNewUserUseCase(t *testing.T) {
	mockRepo := NewMockUserRepository()
	useCase := NewUserUseCase(mockRepo)

	if useCase == nil {
		t.Fatal("NewUserUseCase() returned nil")
	}
	if useCase.userRepo != mockRepo {
		t.Error("userRepo not set correctly")
	}
}

func TestUserUseCase_RegisterUser(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		password    string
		fullName    string
		phoneNumber string
		birthday    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid registration",
			email:       "test@example.com",
			password:    "password123",
			fullName:    "John Doe",
			phoneNumber: "0812345678",
			birthday:    "1990-01-15",
			expectError: false,
		},
		{
			name:        "invalid birthday format",
			email:       "test2@example.com",
			password:    "password123",
			fullName:    "Jane Doe",
			phoneNumber: "0812345678",
			birthday:    "1990/01/15", // wrong format
			expectError: true,
			errorMsg:    "invalid birthday format, should be YYYY-MM-DD",
		},
		{
			name:        "invalid birthday format - incomplete",
			email:       "test3@example.com",
			password:    "password123",
			fullName:    "Bob Smith",
			phoneNumber: "0812345678",
			birthday:    "1990-1-1", // wrong format
			expectError: true,
			errorMsg:    "invalid birthday format, should be YYYY-MM-DD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockUserRepository()
			useCase := NewUserUseCase(mockRepo)

			user, err := useCase.RegisterUser(tt.email, tt.password, tt.fullName, tt.phoneNumber, tt.birthday)

			if tt.expectError {
				if err == nil {
					t.Errorf("RegisterUser() should return error for %s", tt.name)
				}
				if err != nil && err.Error() != tt.errorMsg {
					t.Errorf("RegisterUser() error = %v, want %v", err.Error(), tt.errorMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("RegisterUser() error = %v", err)
				return
			}

			// Verify user data
			if user.Email != tt.email {
				t.Errorf("Email = %v, want %v", user.Email, tt.email)
			}
			if user.FullName != tt.fullName {
				t.Errorf("FullName = %v, want %v", user.FullName, tt.fullName)
			}
			if user.PhoneNumber != tt.phoneNumber {
				t.Errorf("PhoneNumber = %v, want %v", user.PhoneNumber, tt.phoneNumber)
			}
			if user.Birthday != tt.birthday {
				t.Errorf("Birthday = %v, want %v", user.Birthday, tt.birthday)
			}
			if user.Password != "" {
				t.Error("Password should be empty in returned user")
			}
			if user.ID == 0 {
				t.Error("ID should be set for registered user")
			}
		})
	}
}

func TestUserUseCase_RegisterUser_DuplicateEmail(t *testing.T) {
	mockRepo := NewMockUserRepository()
	useCase := NewUserUseCase(mockRepo)

	// Register first user
	_, err := useCase.RegisterUser("test@example.com", "password123", "John Doe", "0812345678", "1990-01-15")
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	// Try to register user with same email
	_, err = useCase.RegisterUser("test@example.com", "password456", "Jane Doe", "0898765432", "1985-05-20")
	if err == nil {
		t.Error("RegisterUser() should return error for duplicate email")
	}
	if err.Error() != "user with this email already exists" {
		t.Errorf("RegisterUser() error = %v, want 'user with this email already exists'", err.Error())
	}
}

func TestUserUseCase_AuthenticateUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	useCase := NewUserUseCase(mockRepo)

	// First, register a user
	email := "auth@example.com"
	password := "password123"
	_, err := useCase.RegisterUser(email, password, "Auth User", "0812345678", "1990-01-15")
	if err != nil {
		t.Fatalf("Failed to register user for auth test: %v", err)
	}

	// Test successful authentication
	user, err := useCase.AuthenticateUser(email, password)
	if err != nil {
		t.Errorf("AuthenticateUser() error = %v", err)
	}
	if user.Email != email {
		t.Errorf("Email = %v, want %v", user.Email, email)
	}

	// Test wrong password
	_, err = useCase.AuthenticateUser(email, "wrongpassword")
	if err == nil {
		t.Error("AuthenticateUser() should return error for wrong password")
	}
	if err.Error() != "invalid credentials" {
		t.Errorf("AuthenticateUser() error = %v, want 'invalid credentials'", err.Error())
	}

	// Test non-existent user
	_, err = useCase.AuthenticateUser("nonexistent@example.com", password)
	if err == nil {
		t.Error("AuthenticateUser() should return error for non-existent user")
	}
	if err.Error() != "invalid credentials" {
		t.Errorf("AuthenticateUser() error = %v, want 'invalid credentials'", err.Error())
	}
}

func TestUserUseCase_GetUserByID(t *testing.T) {
	mockRepo := NewMockUserRepository()
	useCase := NewUserUseCase(mockRepo)

	// Register a user first
	registeredUser, err := useCase.RegisterUser("get@example.com", "password123", "Get User", "0812345678", "1990-01-15")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Test getting existing user
	user, err := useCase.GetUserByID(registeredUser.ID)
	if err != nil {
		t.Errorf("GetUserByID() error = %v", err)
	}
	if user.ID != registeredUser.ID {
		t.Errorf("ID = %v, want %v", user.ID, registeredUser.ID)
	}
	if user.Email != registeredUser.Email {
		t.Errorf("Email = %v, want %v", user.Email, registeredUser.Email)
	}
	if user.Password != "" {
		t.Error("Password should be empty in returned user")
	}

	// Test getting non-existent user
	_, err = useCase.GetUserByID(999)
	if err == nil {
		t.Error("GetUserByID() should return error for non-existent user")
	}
	if err.Error() != "user not found" {
		t.Errorf("GetUserByID() error = %v, want 'user not found'", err.Error())
	}
}

func TestUserUseCase_RegisterUser_PasswordHashing(t *testing.T) {
	mockRepo := NewMockUserRepository()
	useCase := NewUserUseCase(mockRepo)

	email := "hash@example.com"
	password := "plaintextpassword"

	// Register user
	_, err := useCase.RegisterUser(email, password, "Hash User", "0812345678", "1990-01-15")
	if err != nil {
		t.Fatalf("RegisterUser() error = %v", err)
	}

	// Get user from repository directly to check password hashing
	savedUser, err := mockRepo.GetByEmail(email)
	if err != nil {
		t.Fatalf("Failed to get saved user: %v", err)
	}

	// Verify password is hashed (not plain text)
	if savedUser.Password == password {
		t.Error("Password should be hashed, not stored as plain text")
	}

	// Verify password can be verified with bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(savedUser.Password), []byte(password))
	if err != nil {
		t.Error("Saved password should be a valid bcrypt hash of the original password")
	}
}

// Mock repository that always fails on Create
type FailingMockUserRepository struct {
	*MockUserRepository
}

func (m *FailingMockUserRepository) Create(user *entity.User) (*entity.User, error) {
	return nil, errors.New("database error")
}

func TestUserUseCase_RegisterUser_RepositoryError(t *testing.T) {
	// Mock repository that always fails on Create
	mockRepo := &FailingMockUserRepository{
		MockUserRepository: NewMockUserRepository(),
	}
	
	useCase := NewUserUseCase(mockRepo)

	_, err := useCase.RegisterUser("repo@example.com", "password123", "Repo User", "0812345678", "1990-01-15")
	if err == nil {
		t.Error("RegisterUser() should return error when repository fails")
	}
	if err.Error() != "failed to save user" {
		t.Errorf("RegisterUser() error = %v, want 'failed to save user'", err.Error())
	}
}
