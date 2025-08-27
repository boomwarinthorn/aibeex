package main

// @title Fiber Authentication API
// @version 2.0
// @description A Go Fiber API with JWT authentication, user registration, and login functionality
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://github.com/boomwarinthorn/aibeex/blob/main/LICENSE
// @host localhost:3000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"

	_ "fiber-hello-world/docs" // This line is needed for go-swagger to find your docs!
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

// LoginRequest struct for login validation
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// JWT Claims structure
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// Database connection
var db *sql.DB

// JWT secret key (in production, use environment variable)
var jwtSecret = []byte("your-secret-key")

// Initialize validator
var validate = validator.New()

// initDatabase initializes SQLite database and creates tables
func initDatabase() {
	var err error
	db, err = sql.Open("sqlite", "users.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Create users table if it doesn't exist
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
		log.Fatal("Failed to create table:", err)
	}

	log.Println("Database initialized successfully")
}

// createUser inserts a new user into the database
func createUser(user User) (*User, error) {
	query := `
	INSERT INTO users (email, password, full_name, phone_number, birthday, created_at)
	VALUES (?, ?, ?, ?, ?, ?)
	RETURNING id`

	var id int
	err := db.QueryRow(query, user.Email, user.Password, user.FullName, user.PhoneNumber, user.Birthday, user.CreatedAt).Scan(&id)
	if err != nil {
		return nil, err
	}

	user.ID = id
	return &user, nil
}

// getUserByEmail retrieves a user by email from the database
func getUserByEmail(email string) (*User, error) {
	query := `SELECT id, email, password, full_name, phone_number, birthday, created_at FROM users WHERE email = ?`

	var user User
	err := db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.FullName, &user.PhoneNumber, &user.Birthday, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// getUserByID retrieves a user by ID from the database
func getUserByID(id int) (*User, error) {
	query := `SELECT id, email, password, full_name, phone_number, birthday, created_at FROM users WHERE id = ?`

	var user User
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Password, &user.FullName, &user.PhoneNumber, &user.Birthday, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func main() {
	// Initialize database
	initDatabase()
	defer db.Close()

	// Create fiber app
	app := fiber.New()

	// Swagger documentation route
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	// @Summary Get hello world message
	// @Description Returns a simple hello world JSON response
	// @Tags general
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]string
	// @Router / [get]
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello World",
		})
	})

	// POST endpoint for user registration
	app.Post("/register", registerUser)

	// POST endpoint for user login
	app.Post("/login", loginUser)

	// GET endpoint for getting current user info (requires JWT token)
	app.Get("/me", getCurrentUser)

	// Start server on port 3000
	app.Listen(":3000")
}

// @Summary Register a new user
// @Description Register a new user with email, password, full name, phone number, and birthday
// @Tags authentication
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration information"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /register [post]
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
	existingUser, err := getUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		return c.Status(409).JSON(fiber.Map{
			"error":   "Email already exists",
			"message": "User with this email already registered",
		})
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
		Email:       req.Email,
		Password:    string(hashedPassword),
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Birthday:    birthday.Format("2006-01-02"),
		CreatedAt:   time.Now(),
	}

	// Save user to database
	savedUser, err := createUser(newUser)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to save user",
			"message": err.Error(),
		})
	}

	// Return success response (exclude password from response)
	savedUser.Password = ""
	return c.Status(201).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    savedUser,
	})
}

// @Summary User login
// @Description Authenticate user with email and password, returns JWT token
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "User login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /login [post]
func loginUser(c *fiber.Ctx) error {
	// Parse request body
	var req LoginRequest
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

	// Find user by email
	foundUser, err := getUserByEmail(req.Email)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error":   "Invalid credentials",
			"message": "Email or password is incorrect",
		})
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(req.Password))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error":   "Invalid credentials",
			"message": "Email or password is incorrect",
		})
	}

	// Create JWT token
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours
	claims := &Claims{
		UserID: foundUser.ID,
		Email:  foundUser.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   string(rune(foundUser.ID)),
		},
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to generate token",
			"message": err.Error(),
		})
	}

	// Return success response with token
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   tokenString,
		"user": fiber.Map{
			"id":          foundUser.ID,
			"email":       foundUser.Email,
			"fullName":    foundUser.FullName,
			"phoneNumber": foundUser.PhoneNumber,
			"birthday":    foundUser.Birthday,
			"createdAt":   foundUser.CreatedAt,
		},
		"expiresAt": expirationTime,
	})
}

// validateJWT validates JWT token from Authorization header
func validateJWT(c *fiber.Ctx) (*Claims, error) {
	// Get Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return nil, fiber.NewError(401, "Authorization header required")
	}

	// Check if it starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, fiber.NewError(401, "Bearer token required")
	}

	// Extract token from "Bearer <token>"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return nil, fiber.NewError(401, "Token required")
	}

	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Make sure the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(401, "Invalid token signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fiber.NewError(401, "Invalid token")
	}

	// Check if token is valid and extract claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fiber.NewError(401, "Invalid token claims")
}

// @Summary Get current user information
// @Description Get the current authenticated user's profile information using JWT token
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]interface{}
// @Router /me [get]
func getCurrentUser(c *fiber.Ctx) error {
	// Validate JWT token and get claims
	claims, err := validateJWT(c)
	if err != nil {
		return err
	}

	// Find user by ID from token claims
	foundUser, err := getUserByID(claims.UserID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "User not found",
			"message": "User associated with this token no longer exists",
		})
	}

	// Return user information (exclude password)
	return c.JSON(fiber.Map{
		"message": "User information retrieved successfully",
		"user": fiber.Map{
			"id":          foundUser.ID,
			"email":       foundUser.Email,
			"fullName":    foundUser.FullName,
			"phoneNumber": foundUser.PhoneNumber,
			"birthday":    foundUser.Birthday,
			"createdAt":   foundUser.CreatedAt,
		},
	})
}
