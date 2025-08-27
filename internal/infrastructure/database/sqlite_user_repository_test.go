package database

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"fiber-hello-world/internal/domain/entity"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	// Create temporary database file
	dbFile := "test_users.db"
	
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create users table
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		full_name TEXT NOT NULL,
		phone_number TEXT NOT NULL,
		birthday TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Return cleanup function
	cleanup := func() {
		db.Close()
		os.Remove(dbFile)
	}

	return db, cleanup
}

func TestInitDatabase(t *testing.T) {
	// Use a temporary database file
	dbFile := "test_init.db"
	defer os.Remove(dbFile)

	// Change working directory temporarily for the test
	originalDBPath := "users.db"
	
	// Test InitDatabase function
	db, err := InitDatabase()
	if err != nil {
		t.Fatalf("InitDatabase() error = %v", err)
	}
	defer db.Close()
	defer os.Remove(originalDBPath) // Clean up the created database

	// Verify database connection works
	err = db.Ping()
	if err != nil {
		t.Errorf("Database ping failed: %v", err)
	}

	// Verify table exists
	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='users'").Scan(&tableName)
	if err != nil {
		t.Errorf("Users table should exist: %v", err)
	}
	if tableName != "users" {
		t.Errorf("Table name = %v, want 'users'", tableName)
	}
}

func TestNewSQLiteUserRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewSQLiteUserRepository(db)
	if repo == nil {
		t.Fatal("NewSQLiteUserRepository() returned nil")
	}
	if repo.db != db {
		t.Error("Database connection not set correctly")
	}
}

func TestSQLiteUserRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewSQLiteUserRepository(db)
	now := time.Now()

	user := &entity.User{
		Email:       "test@example.com",
		Password:    "hashedpassword",
		FullName:    "Test User",
		PhoneNumber: "0812345678",
		Birthday:    "1990-01-15",
		CreatedAt:   now,
	}

	createdUser, err := repo.Create(user)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Verify ID is set
	if createdUser.ID == 0 {
		t.Error("ID should be set after creation")
	}

	// Verify all fields
	if createdUser.Email != user.Email {
		t.Errorf("Email = %v, want %v", createdUser.Email, user.Email)
	}
	if createdUser.Password != user.Password {
		t.Errorf("Password = %v, want %v", createdUser.Password, user.Password)
	}
	if createdUser.FullName != user.FullName {
		t.Errorf("FullName = %v, want %v", createdUser.FullName, user.FullName)
	}
	if createdUser.PhoneNumber != user.PhoneNumber {
		t.Errorf("PhoneNumber = %v, want %v", createdUser.PhoneNumber, user.PhoneNumber)
	}
	if createdUser.Birthday != user.Birthday {
		t.Errorf("Birthday = %v, want %v", createdUser.Birthday, user.Birthday)
	}
}

func TestSQLiteUserRepository_Create_DuplicateEmail(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewSQLiteUserRepository(db)

	user1 := &entity.User{
		Email:       "duplicate@example.com",
		Password:    "password1",
		FullName:    "User One",
		PhoneNumber: "0812345678",
		Birthday:    "1990-01-15",
		CreatedAt:   time.Now(),
	}

	user2 := &entity.User{
		Email:       "duplicate@example.com", // Same email
		Password:    "password2",
		FullName:    "User Two",
		PhoneNumber: "0898765432",
		Birthday:    "1985-05-20",
		CreatedAt:   time.Now(),
	}

	// Create first user
	_, err := repo.Create(user1)
	if err != nil {
		t.Fatalf("First user creation failed: %v", err)
	}

	// Try to create second user with duplicate email
	_, err = repo.Create(user2)
	if err == nil {
		t.Error("Create() should return error for duplicate email")
	}
}

func TestSQLiteUserRepository_GetByEmail(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewSQLiteUserRepository(db)

	// Create a user first
	user := &entity.User{
		Email:       "get@example.com",
		Password:    "hashedpassword",
		FullName:    "Get User",
		PhoneNumber: "0812345678",
		Birthday:    "1990-01-15",
		CreatedAt:   time.Now(),
	}

	createdUser, err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user for GetByEmail test: %v", err)
	}

	// Test getting existing user
	foundUser, err := repo.GetByEmail("get@example.com")
	if err != nil {
		t.Fatalf("GetByEmail() error = %v", err)
	}

	if foundUser.ID != createdUser.ID {
		t.Errorf("ID = %v, want %v", foundUser.ID, createdUser.ID)
	}
	if foundUser.Email != createdUser.Email {
		t.Errorf("Email = %v, want %v", foundUser.Email, createdUser.Email)
	}
	if foundUser.Password != createdUser.Password {
		t.Errorf("Password = %v, want %v", foundUser.Password, createdUser.Password)
	}

	// Test getting non-existent user
	_, err = repo.GetByEmail("nonexistent@example.com")
	if err == nil {
		t.Error("GetByEmail() should return error for non-existent user")
	}
}

func TestSQLiteUserRepository_GetByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewSQLiteUserRepository(db)

	// Create a user first
	user := &entity.User{
		Email:       "getbyid@example.com",
		Password:    "hashedpassword",
		FullName:    "GetByID User",
		PhoneNumber: "0812345678",
		Birthday:    "1990-01-15",
		CreatedAt:   time.Now(),
	}

	createdUser, err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user for GetByID test: %v", err)
	}

	// Test getting existing user by ID
	foundUser, err := repo.GetByID(createdUser.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if foundUser.ID != createdUser.ID {
		t.Errorf("ID = %v, want %v", foundUser.ID, createdUser.ID)
	}
	if foundUser.Email != createdUser.Email {
		t.Errorf("Email = %v, want %v", foundUser.Email, createdUser.Email)
	}

	// Test getting non-existent user
	_, err = repo.GetByID(999)
	if err == nil {
		t.Error("GetByID() should return error for non-existent user")
	}
}

func TestSQLiteUserRepository_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewSQLiteUserRepository(db)

	// Create a user first
	user := &entity.User{
		Email:       "update@example.com",
		Password:    "hashedpassword",
		FullName:    "Original Name",
		PhoneNumber: "0812345678",
		Birthday:    "1990-01-15",
		CreatedAt:   time.Now(),
	}

	createdUser, err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user for Update test: %v", err)
	}

	// Update user
	createdUser.FullName = "Updated Name"
	createdUser.PhoneNumber = "0898765432"
	createdUser.Birthday = "1985-05-20"

	err = repo.Update(createdUser)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	// Verify update
	updatedUser, err := repo.GetByID(createdUser.ID)
	if err != nil {
		t.Fatalf("Failed to get updated user: %v", err)
	}

	if updatedUser.FullName != "Updated Name" {
		t.Errorf("FullName = %v, want %v", updatedUser.FullName, "Updated Name")
	}
	if updatedUser.PhoneNumber != "0898765432" {
		t.Errorf("PhoneNumber = %v, want %v", updatedUser.PhoneNumber, "0898765432")
	}
	if updatedUser.Birthday != "1985-05-20" {
		t.Errorf("Birthday = %v, want %v", updatedUser.Birthday, "1985-05-20")
	}

	// Test updating non-existent user
	nonExistentUser := &entity.User{
		ID:          999,
		Email:       "nonexistent@example.com",
		FullName:    "Non Existent",
		PhoneNumber: "0812345678",
		Birthday:    "1990-01-15",
	}

	err = repo.Update(nonExistentUser)
	// SQLite UPDATE won't return error for non-existent ID, it just affects 0 rows
	// This is expected behavior
	if err != nil {
		t.Errorf("Update() should not return error for non-existent user in SQLite: %v", err)
	}
}

func TestSQLiteUserRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewSQLiteUserRepository(db)

	// Create a user first
	user := &entity.User{
		Email:       "delete@example.com",
		Password:    "hashedpassword",
		FullName:    "Delete User",
		PhoneNumber: "0812345678",
		Birthday:    "1990-01-15",
		CreatedAt:   time.Now(),
	}

	createdUser, err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user for Delete test: %v", err)
	}

	// Delete user
	err = repo.Delete(createdUser.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify user is deleted
	_, err = repo.GetByID(createdUser.ID)
	if err == nil {
		t.Error("GetByID() should return error for deleted user")
	}

	// Test deleting non-existent user
	err = repo.Delete(999)
	// SQLite DELETE won't return error for non-existent ID
	if err != nil {
		t.Errorf("Delete() should not return error for non-existent user in SQLite: %v", err)
	}
}

func TestSQLiteUserRepository_CRUD_Integration(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewSQLiteUserRepository(db)

	// Create
	user := &entity.User{
		Email:       "crud@example.com",
		Password:    "hashedpassword",
		FullName:    "CRUD User",
		PhoneNumber: "0812345678",
		Birthday:    "1990-01-15",
		CreatedAt:   time.Now(),
	}

	createdUser, err := repo.Create(user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Read by email
	foundUser, err := repo.GetByEmail("crud@example.com")
	if err != nil {
		t.Fatalf("GetByEmail failed: %v", err)
	}
	if foundUser.ID != createdUser.ID {
		t.Error("GetByEmail returned wrong user")
	}

	// Read by ID
	foundUser2, err := repo.GetByID(createdUser.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if foundUser2.Email != createdUser.Email {
		t.Error("GetByID returned wrong user")
	}

	// Update
	foundUser2.FullName = "Updated CRUD User"
	err = repo.Update(foundUser2)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify update
	updatedUser, err := repo.GetByID(createdUser.ID)
	if err != nil {
		t.Fatalf("GetByID after update failed: %v", err)
	}
	if updatedUser.FullName != "Updated CRUD User" {
		t.Error("Update was not persisted")
	}

	// Delete
	err = repo.Delete(createdUser.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID(createdUser.ID)
	if err == nil {
		t.Error("User should be deleted")
	}
}
