package dto

import "time"

// RegisterRequest represents the request payload for user registration
type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6"`
	FullName    string `json:"fullName" validate:"required,min=2"`
	PhoneNumber string `json:"phoneNumber" validate:"required,min=10"`
	Birthday    string `json:"birthday" validate:"required"`
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse represents the response payload for user data
type UserResponse struct {
	ID          int       `json:"id"`
	Email       string    `json:"email"`
	FullName    string    `json:"fullName"`
	PhoneNumber string    `json:"phoneNumber"`
	Birthday    string    `json:"birthday"`
	CreatedAt   time.Time `json:"createdAt"`
}

// LoginResponse represents the response payload for login
type LoginResponse struct {
	Message   string       `json:"message"`
	Token     string       `json:"token"`
	User      UserResponse `json:"user"`
	ExpiresAt time.Time    `json:"expiresAt"`
}

// ErrorResponse represents the error response payload
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// SuccessResponse represents the success response payload
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
