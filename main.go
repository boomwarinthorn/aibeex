package main

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// User struct for registration
type User struct {
	ID          int       `json:"id"`
	Email       string    `json:"email" validate:"required,email"`
	Password    string    `json:"password,omitempty" validate:"required,min=6"`
	FullName    string    `json:"fullName" validate:"required,min=2"`
	PhoneNumber string    `json:"phoneNumber" validate:"required,min=10"`
	Birthday    string    `json:"birthday" validate:"required"`
	CreatedAt   time.Time `json:"createdAt"`
}

// RegisterRequest struct for input validation
type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6"`
	FullName    string `json:"fullName" validate:"required,min=2"`
	PhoneNumber string `json:"phoneNumber" validate:"required,min=10"`
	Birthday    string `json:"birthday" validate:"required"`
}

// In-memory user storage (in production, use a database)
var users []User
var userIDCounter = 1

// Initialize validator
var validate = validator.New()

func main() {
	// Create fiber app
	app := fiber.New()

	// GET endpoint that returns JSON "Hello World"
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello World",
		})
	})

	// POST endpoint for user registration
	app.Post("/register", registerUser)

	// Start server on port 3000
	app.Listen(":3000")
}

// registerUser handles user registration
func registerUser(c *fiber.Ctx) error {
	// Parse request body
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
	}

	// Validate input
	if err := validate.Struct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Validation failed",
			"message": err.Error(),
		})
	}

	// Check if email already exists
	for _, user := range users {
		if user.Email == req.Email {
			return c.Status(409).JSON(fiber.Map{
				"error":   "Email already exists",
				"message": "User with this email already registered",
			})
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to hash password",
			"message": err.Error(),
		})
	}

	// Parse birthday
	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid birthday format",
			"message": "Birthday should be in YYYY-MM-DD format",
		})
	}

	// Create new user
	newUser := User{
		ID:          userIDCounter,
		Email:       req.Email,
		Password:    string(hashedPassword),
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Birthday:    birthday.Format("2006-01-02"),
		CreatedAt:   time.Now(),
	}

	// Add user to storage
	users = append(users, newUser)
	userIDCounter++

	// Return success response (exclude password from response)
	newUser.Password = ""
	return c.Status(201).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    newUser,
	})
}
