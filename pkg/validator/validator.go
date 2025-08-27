package validator

import (
	"github.com/go-playground/validator/v10"
)

// Service provides validation operations
type Service struct {
	validator *validator.Validate
}

// NewService creates a new validator service
func NewService() *Service {
	return &Service{
		validator: validator.New(),
	}
}

// Validate validates a struct based on validation tags
func (s *Service) Validate(data interface{}) error {
	return s.validator.Struct(data)
}
