package repository

import "fiber-hello-world/internal/domain/entity"

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Create saves a new user and returns the created user with ID
	Create(user *entity.User) (*entity.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(email string) (*entity.User, error)

	// GetByID retrieves a user by ID
	GetByID(id int) (*entity.User, error)

	// Update updates user information
	Update(user *entity.User) error

	// Delete removes a user by ID
	Delete(id int) error
}
