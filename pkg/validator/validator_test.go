package validator

import (
	"testing"
)

type TestStruct struct {
	Email       string `validate:"required,email"`
	Password    string `validate:"required,min=6"`
	Name        string `validate:"required,min=2"`
	PhoneNumber string `validate:"required,min=10"`
	Age         int    `validate:"min=0,max=150"`
}

func TestNewService(t *testing.T) {
	service := NewService()
	if service == nil {
		t.Fatal("NewService() returned nil")
	}
	if service.validator == nil {
		t.Fatal("validator instance is nil")
	}
}

func TestService_Validate_ValidData(t *testing.T) {
	service := NewService()

	validData := &TestStruct{
		Email:       "test@example.com",
		Password:    "password123",
		Name:        "John Doe",
		PhoneNumber: "0812345678",
		Age:         25,
	}

	err := service.Validate(validData)
	if err != nil {
		t.Errorf("Validate() should not return error for valid data, got: %v", err)
	}
}

func TestService_Validate_InvalidData(t *testing.T) {
	service := NewService()

	tests := []struct {
		name        string
		data        *TestStruct
		expectError bool
		errorField  string
	}{
		{
			name: "missing email",
			data: &TestStruct{
				Email:       "",
				Password:    "password123",
				Name:        "John",
				PhoneNumber: "0812345678",
				Age:         25,
			},
			expectError: true,
			errorField:  "Email",
		},
		{
			name: "invalid email format",
			data: &TestStruct{
				Email:       "invalid-email",
				Password:    "password123",
				Name:        "John",
				PhoneNumber: "0812345678",
				Age:         25,
			},
			expectError: true,
			errorField:  "Email",
		},
		{
			name: "short password",
			data: &TestStruct{
				Email:       "test@example.com",
				Password:    "123",
				Name:        "John",
				PhoneNumber: "0812345678",
				Age:         25,
			},
			expectError: true,
			errorField:  "Password",
		},
		{
			name: "short name",
			data: &TestStruct{
				Email:       "test@example.com",
				Password:    "password123",
				Name:        "J",
				PhoneNumber: "0812345678",
				Age:         25,
			},
			expectError: true,
			errorField:  "Name",
		},
		{
			name: "short phone number",
			data: &TestStruct{
				Email:       "test@example.com",
				Password:    "password123",
				Name:        "John",
				PhoneNumber: "081234567",
				Age:         25,
			},
			expectError: true,
			errorField:  "PhoneNumber",
		},
		{
			name: "negative age",
			data: &TestStruct{
				Email:       "test@example.com",
				Password:    "password123",
				Name:        "John",
				PhoneNumber: "0812345678",
				Age:         -1,
			},
			expectError: true,
			errorField:  "Age",
		},
		{
			name: "age too high",
			data: &TestStruct{
				Email:       "test@example.com",
				Password:    "password123",
				Name:        "John",
				PhoneNumber: "0812345678",
				Age:         200,
			},
			expectError: true,
			errorField:  "Age",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Validate(tt.data)

			if tt.expectError && err == nil {
				t.Errorf("Validate() should return error for %s", tt.name)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Validate() should not return error for %s, got: %v", tt.name, err)
			}
		})
	}
}

func TestService_Validate_MultipleErrors(t *testing.T) {
	service := NewService()

	invalidData := &TestStruct{
		Email:       "",    // missing
		Password:    "123", // too short
		Name:        "",    // missing
		PhoneNumber: "123", // too short
		Age:         -5,    // negative
	}

	err := service.Validate(invalidData)
	if err == nil {
		t.Fatal("Validate() should return error for invalid data with multiple validation failures")
	}

	// Error message should be present
	errorMsg := err.Error()
	if errorMsg == "" {
		t.Error("Error message should not be empty")
	}
}

func TestService_Validate_NilData(t *testing.T) {
	service := NewService()

	err := service.Validate(nil)
	if err == nil {
		t.Error("Validate() should return error for nil data")
	}
}

type EmptyStruct struct{}

func TestService_Validate_EmptyStruct(t *testing.T) {
	service := NewService()

	emptyData := &EmptyStruct{}
	err := service.Validate(emptyData)
	if err != nil {
		t.Errorf("Validate() should not return error for empty struct, got: %v", err)
	}
}

type NestedStruct struct {
	User TestStruct `validate:"required"`
	ID   int        `validate:"min=1"`
}

func TestService_Validate_NestedStruct(t *testing.T) {
	service := NewService()

	validNested := &NestedStruct{
		User: TestStruct{
			Email:       "test@example.com",
			Password:    "password123",
			Name:        "John Doe",
			PhoneNumber: "0812345678",
			Age:         25,
		},
		ID: 1,
	}

	err := service.Validate(validNested)
	if err != nil {
		t.Errorf("Validate() should not return error for valid nested struct, got: %v", err)
	}

	invalidNested := &NestedStruct{
		User: TestStruct{
			Email:       "invalid-email", // invalid email
			Password:    "123",           // too short
			Name:        "John Doe",
			PhoneNumber: "0812345678",
			Age:         25,
		},
		ID: 0, // invalid ID
	}

	err = service.Validate(invalidNested)
	if err == nil {
		t.Error("Validate() should return error for invalid nested struct")
	}
}
