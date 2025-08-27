package usecase

import (
	"errors"
	"time"

	"fiber-hello-world/internal/domain/entity"
	"fiber-hello-world/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserUseCase handles user-related business logic
type UserUseCase struct {
	userRepo repository.UserRepository
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

// RegisterUser handles user registration logic
func (uc *UserUseCase) RegisterUser(email, password, fullName, phoneNumber, birthday string) (*entity.User, error) {
	// Check if user already exists
	existingUser, err := uc.userRepo.GetByEmail(email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Validate birthday format
	_, err = time.Parse("2006-01-02", birthday)
	if err != nil {
		return nil, errors.New("invalid birthday format, should be YYYY-MM-DD")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create new user entity
	user := entity.NewUser(email, string(hashedPassword), fullName, phoneNumber, birthday)

	// Save user to repository
	savedUser, err := uc.userRepo.Create(user)
	if err != nil {
		return nil, errors.New("failed to save user")
	}

	return savedUser.WithoutPassword(), nil
}

// AuthenticateUser handles user authentication
func (uc *UserUseCase) AuthenticateUser(email, password string) (*entity.User, error) {
	// Find user by email
	user, err := uc.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// GetUserByID retrieves user by ID
func (uc *UserUseCase) GetUserByID(id int) (*entity.User, error) {
	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user.WithoutPassword(), nil
}
