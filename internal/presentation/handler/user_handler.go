package handler

import (
	"fiber-hello-world/internal/presentation/dto"
	"fiber-hello-world/internal/usecase"
	"fiber-hello-world/pkg/jwt"
	"fiber-hello-world/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userUseCase *usecase.UserUseCase
	jwtService  *jwt.Service
	validator   *validator.Service
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUseCase *usecase.UserUseCase, jwtService *jwt.Service, validator *validator.Service) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		jwtService:  jwtService,
		validator:   validator,
	}
}

// @Summary Register a new user
// @Description Register a new user with email, password, full name, phone number, and birthday
// @Tags authentication
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequest true "User registration information"
// @Success 201 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	// Parse request body
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	// Validate input
	if err := h.validator.Validate(&req); err != nil {
		return c.Status(400).JSON(dto.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
	}

	// Register user
	user, err := h.userUseCase.RegisterUser(req.Email, req.Password, req.FullName, req.PhoneNumber, req.Birthday)
	if err != nil {
		status := 500
		if err.Error() == "user with this email already exists" {
			status = 409
		} else if err.Error() == "invalid birthday format, should be YYYY-MM-DD" {
			status = 400
		}

		return c.Status(status).JSON(dto.ErrorResponse{
			Error:   "Registration failed",
			Message: err.Error(),
		})
	}

	// Convert to response DTO
	userResponse := dto.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		Birthday:    user.Birthday,
		CreatedAt:   user.CreatedAt,
	}

	return c.Status(201).JSON(dto.SuccessResponse{
		Message: "User registered successfully",
		Data:    userResponse,
	})
}

// @Summary User login
// @Description Authenticate user with email and password, returns JWT token
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "User login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	// Parse request body
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	// Validate input
	if err := h.validator.Validate(&req); err != nil {
		return c.Status(400).JSON(dto.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
	}

	// Authenticate user
	user, err := h.userUseCase.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(dto.ErrorResponse{
			Error:   "Authentication failed",
			Message: err.Error(),
		})
	}

	// Generate JWT token
	token, expiresAt, err := h.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{
			Error:   "Token generation failed",
			Message: err.Error(),
		})
	}

	// Convert to response DTO
	userResponse := dto.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		Birthday:    user.Birthday,
		CreatedAt:   user.CreatedAt,
	}

	return c.JSON(dto.LoginResponse{
		Message:   "Login successful",
		Token:     token,
		User:      userResponse,
		ExpiresAt: expiresAt,
	})
}

// @Summary Get current user information
// @Description Get the current authenticated user's profile information using JWT token
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /me [get]
func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	// Get user claims from middleware
	claims, ok := c.Locals("user").(*jwt.Claims)
	if !ok {
		return c.Status(401).JSON(dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid token claims",
		})
	}

	// Get user by ID
	user, err := h.userUseCase.GetUserByID(claims.UserID)
	if err != nil {
		return c.Status(404).JSON(dto.ErrorResponse{
			Error:   "User not found",
			Message: err.Error(),
		})
	}

	// Convert to response DTO
	userResponse := dto.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		Birthday:    user.Birthday,
		CreatedAt:   user.CreatedAt,
	}

	return c.JSON(dto.SuccessResponse{
		Message: "User information retrieved successfully",
		Data:    userResponse,
	})
}
